package cart

import (
	"errors"

	"gorm.io/gorm"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/cart"
	productRepo "gosh/internal/repository/product"
)

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
		return s.toItemResponse(existing, sku), nil
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
	return s.toItemResponse(cart, sku), nil
}

func (s *service) List(userID uint) (*response.CartSummaryResponse, error) {
	carts, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	items := make([]response.CartItemResponse, 0)
	removed := make([]response.RemovedItem, 0)
	var selectedIDs []uint
	var totalAmount int64
	allSelected := len(carts) > 0

	for _, c := range carts {
		sku, skuErr := s.productRepo.FindSKUByID(c.SKUID)
		if skuErr != nil {
			removed = append(removed, response.RemovedItem{
				SKUID:  c.SKUID,
				Reason: "商品已下架或不存在",
			})
			continue
		}
		if !s.isProductOnline(sku.ProductID) {
			removed = append(removed, response.RemovedItem{
				SKUID:  c.SKUID,
				Reason: "商品已下架",
			})
			continue
		}

		item := s.toItemResponse(&c, sku)
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

	sku, _ := s.productRepo.FindSKUByID(cart.SKUID)
	if sku != nil {
		return s.toItemResponse(cart, sku), nil
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
	return s.repo.Delete(id, userID)
}

func (s *service) Select(userID uint, req *request.SelectRequest) error {
	if req.All != nil {
		return s.repo.UpdateAllSelected(userID, *req.All)
	}
	if len(req.SKUIDs) > 0 {
		for _, skuID := range req.SKUIDs {
			if err := s.repo.UpdateSelected(userID, skuID, req.Select); err != nil {
				return err
			}
		}
		return nil
	}
	return s.repo.UpdateAllSelected(userID, req.Select)
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
	return s.List(userID)
}

func (s *service) Count(userID uint) (int64, error) {
	return s.repo.Count(userID)
}

// internal helpers

func (s *service) isProductOnline(productID uint) bool {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return false
	}
	return product.Status == model.ProductStatusOn
}

func (s *service) toItemResponse(cart *model.Cart, sku *model.ProductSKU) *response.CartItemResponse {
	product, _ := s.productRepo.FindByID(sku.ProductID)
	image := ""
	productName := ""
	online := true
	if product != nil {
		image = product.Images
		productName = product.Name
		online = product.Status == model.ProductStatusOn
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
		ProductID:   sku.ProductID,
		ProductOn:   online,
		CreatedAt:   cart.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
