package cart

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/cart"
	productRepo "gosh/internal/repository/product"
	"gosh/pkg/cache"
)

const CartMaxQuantity = 99

var (
	ErrCartNotFound    = errors.New("cart item not found")
	ErrProductOffShelf = errors.New("product is off shelf")
	ErrSKUNotFound     = errors.New("SKU not found")
)

type Service interface {
	Add(userID uint, req *request.AddCartRequest) (*response.CartItemResponse, error)
	List(userID uint) (*response.CartSummaryResponse, error)
	Update(id, userID uint, req *request.UpdateCartRequest) (*response.CartItemResponse, error)
	Delete(id, userID uint) error
	Select(userID uint, req *request.SelectRequest) error
	Merge(userID uint, req *request.MergeCartRequest) (*response.CartSummaryResponse, error)
	Count(userID uint) (int64, error)
}

type service struct {
	repo        repo.Repository
	productRepo productRepo.Repository
}

func New() Service {
	return &service{
		repo:        repo.New(),
		productRepo: productRepo.New(),
	}
}

func cartCacheKey(userID uint) string {
	return fmt.Sprintf("cache:cart:user:%d", userID)
}

func invalidateCartCache(userID uint) {
	if c := cache.Default(); c != nil {
		c.Del(context.Background(), cartCacheKey(userID))
	}
}

func (s *service) Add(userID uint, req *request.AddCartRequest) (*response.CartItemResponse, error) {
	sku, err := s.productRepo.FindSKUByID(req.SKUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSKUNotFound
		}
		return nil, err
	}
	if !s.isProductOnline(sku.ProductID) {
		return nil, ErrProductOffShelf
	}

	existing, err := s.repo.FindByUserAndSKU(userID, req.SKUID)
	if err == nil {
		existing.Quantity += req.Quantity
		if existing.Quantity > sku.Stock {
			existing.Quantity = sku.Stock
		}
		if err := s.repo.Update(existing); err != nil {
			return nil, err
		}
		invalidateCartCache(userID)
		return s.toItemResponse(existing, sku, nil), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	qty := req.Quantity
	if qty > sku.Stock {
		qty = sku.Stock
	}
	cart := &model.Cart{
		UserID:   userID,
		SKUID:    req.SKUID,
		Quantity: qty,
		Selected: true,
	}
	if err := s.repo.Create(cart); err != nil {
		return nil, err
	}
	invalidateCartCache(userID)
	return s.toItemResponse(cart, sku, nil), nil
}

func (s *service) List(userID uint) (*response.CartSummaryResponse, error) {
	c := cache.Default()
	if c != nil {
		var result response.CartSummaryResponse
		err := c.Remember(context.Background(), cartCacheKey(userID), 5*time.Minute, func() (interface{}, error) {
			return s.buildCartSummary(userID)
		}, &result)
		if err == nil {
			return &result, nil
		}
	}

	return s.buildCartSummary(userID)
}

func (s *service) buildCartSummary(userID uint) (*response.CartSummaryResponse, error) {
	carts, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(carts) == 0 {
		return &response.CartSummaryResponse{
			Items: []response.CartItemResponse{},
		}, nil
	}

	skuIDs := make([]uint, len(carts))
	for i, c := range carts {
		skuIDs[i] = c.SKUID
	}

	skus, err := s.productRepo.FindSKUsByIDs(skuIDs)
	if err != nil {
		return nil, err
	}
	skuMap := make(map[uint]model.ProductSKU, len(skus))
	for _, sku := range skus {
		skuMap[sku.ID] = sku
	}

	productIDs := make([]uint, 0)
	productIDSet := make(map[uint]bool)
	for _, sku := range skus {
		if !productIDSet[sku.ProductID] {
			productIDs = append(productIDs, sku.ProductID)
			productIDSet[sku.ProductID] = true
		}
	}
	products, err := s.productRepo.FindByIDs(productIDs)
	if err != nil {
		return nil, err
	}
	productMap := make(map[uint]model.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	items := make([]response.CartItemResponse, 0)
	removed := make([]response.RemovedItem, 0)
	var selectedIDs []uint
	var totalAmount int64
	allSelected := len(carts) > 0

	for _, c := range carts {
		sku, ok := skuMap[c.SKUID]
		if !ok {
			removed = append(removed, response.RemovedItem{
				SKUID:  c.SKUID,
				Reason: "商品已下架或不存在",
			})
			continue
		}

		product, ok := productMap[sku.ProductID]
		if !ok || product.Status != model.ProductStatusOn {
			removed = append(removed, response.RemovedItem{
				SKUID:  c.SKUID,
				Reason: "商品已下架",
			})
			continue
		}

		item := s.toItemResponse(&c, &sku, &product)
		items = append(items, *item)

		if c.Selected {
			selectedIDs = append(selectedIDs, c.SKUID)
			totalAmount += sku.Price * int64(c.Quantity)
		} else {
			allSelected = false
		}
	}

	if len(removed) > 0 {
		var removedIDs []uint
		for _, r := range removed {
			removedIDs = append(removedIDs, r.SKUID)
		}
		s.repo.DeleteBySKUIDs(userID, removedIDs)
	}

	return &response.CartSummaryResponse{
		Items:        items,
		SelectedAll:  allSelected && len(items) > 0,
		SelectedIDs:  selectedIDs,
		TotalAmount:  totalAmount,
		TotalCount:   len(items),
		RemovedItems: removed,
	}, nil
}

func (s *service) Update(id, userID uint, req *request.UpdateCartRequest) (*response.CartItemResponse, error) {
	cart, err := s.repo.FindByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	if req.Quantity > 0 {
		cart.Quantity = req.Quantity
	}
	if req.Selected != nil {
		cart.Selected = *req.Selected
	}
	if err := s.repo.Update(cart); err != nil {
		return nil, err
	}

	invalidateCartCache(userID)

	sku, _ := s.productRepo.FindSKUByID(cart.SKUID)
	if sku != nil {
		return s.toItemResponse(cart, sku, nil), nil
	}
	resp := response.ToCartItemResponse(cart)
	return &resp, nil
}

func (s *service) Delete(id, userID uint) error {
	if _, err := s.repo.FindByID(id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartNotFound
		}
		return err
	}
	if err := s.repo.Delete(id, userID); err != nil {
		return err
	}
	invalidateCartCache(userID)
	return nil
}

func (s *service) Select(userID uint, req *request.SelectRequest) error {
	if req.All != nil {
		if err := s.repo.UpdateAllSelected(userID, *req.All); err != nil {
			return err
		}
		invalidateCartCache(userID)
		return nil
	}
	if len(req.SKUIDs) > 0 {
		for _, skuID := range req.SKUIDs {
			if err := s.repo.UpdateSelected(userID, skuID, req.Select); err != nil {
				return err
			}
		}
		invalidateCartCache(userID)
		return nil
	}
	if err := s.repo.UpdateAllSelected(userID, req.Select); err != nil {
		return err
	}
	invalidateCartCache(userID)
	return nil
}

func (s *service) Merge(userID uint, req *request.MergeCartRequest) (*response.CartSummaryResponse, error) {
	for _, item := range req.Items {
		existing, err := s.repo.FindByUserAndSKU(userID, item.SKUID)
		if err == nil {
			if item.Quantity > existing.Quantity {
				existing.Quantity = item.Quantity
				s.repo.Update(existing)
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			sku, err := s.productRepo.FindSKUByID(item.SKUID)
			if err != nil {
				continue
			}
			if !s.isProductOnline(sku.ProductID) {
				continue
			}
			s.repo.Create(&model.Cart{
				UserID:   userID,
				SKUID:    item.SKUID,
				Quantity: item.Quantity,
				Selected: true,
			})
		}
	}
	invalidateCartCache(userID)
	return s.List(userID)
}

func (s *service) Count(userID uint) (int64, error) {
	c := cache.Default()
	if c != nil {
		var summary response.CartSummaryResponse
		err := c.Get(context.Background(), cartCacheKey(userID), &summary)
		if err == nil {
			return int64(summary.TotalCount), nil
		}
	}

	return s.repo.Count(userID)
}

// internal helpers

func (s *service) isProductOnline(productID uint) bool {
	online, err := s.productRepo.IsProductOnline(productID)
	if err != nil {
		return false
	}
	return online
}

func (s *service) toItemResponse(cart *model.Cart, sku *model.ProductSKU, product *model.Product) *response.CartItemResponse {
	image := ""
	productName := ""
	online := true
	if product != nil {
		image = firstImage(product.Images)
		productName = product.Name
		online = product.Status == model.ProductStatusOn
	} else {
		p, err := s.productRepo.FindByID(sku.ProductID)
		if err == nil {
			image = firstImage(p.Images)
			productName = p.Name
			online = p.Status == model.ProductStatusOn
		}
	}

	maxBuyable := sku.Stock
	if maxBuyable > CartMaxQuantity {
		maxBuyable = CartMaxQuantity
	}

	return &response.CartItemResponse{
		ID:          cart.ID,
		UserID:      cart.UserID,
		SKUID:       cart.SKUID,
		Quantity:    cart.Quantity,
		Selected:    cart.Selected,
		ProductName: productName,
		SKUName:     sku.Name,
		Image:       image,
		Price:       sku.Price,
		Stock:       sku.Stock,
		MaxBuyable:  maxBuyable,
		ProductID:   sku.ProductID,
		ProductOn:   online,
		CreatedAt:   cart.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func firstImage(images []string) string {
	if len(images) > 0 {
		return images[0]
	}
	return ""
}
