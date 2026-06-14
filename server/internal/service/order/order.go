package order

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"gosh/internal/database"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	addressRepo "gosh/internal/repository/address"
	cartRepo "gosh/internal/repository/cart"
	orderRepo "gosh/internal/repository/order"
	productRepo "gosh/internal/repository/product"
)

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrInvalidOrderStatus     = errors.New("invalid order status for this operation")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrCartEmpty              = errors.New("cart is empty")
	ErrNoDefaultAddress       = errors.New("no default address")
	ErrMissingIdempotentKey   = errors.New("missing Idempotent-Key")
	ErrOrderNotBelongToUser   = errors.New("order does not belong to user")
	ErrInsufficientPoints      = errors.New("insufficient points")
)

type Service interface {
	Create(userID uint, req *request.CreateOrderRequest, idempotentKey string) (*response.OrderResponse, error)
	List(userID uint, req *request.ListOrderRequest) ([]response.OrderResponse, int64, error)
	GetByID(userID, orderID uint) (*response.OrderResponse, error)
	Cancel(userID, orderID uint, req *request.CancelOrderRequest) error
	Pay(userID, orderID uint) error
	Ship(orderID uint) error
	Confirm(userID, orderID uint) error
	Rebuy(userID, orderID uint) error
	ApplyPoints(userID, orderID uint, req *request.ApplyPointsRequest) error
}

type service struct {
	orderRepo   orderRepo.Repository
	cartRepo    cartRepo.Repository
	productRepo productRepo.Repository
	addressRepo addressRepo.Repository
}

func New() Service {
	return &service{
		orderRepo:   orderRepo.New(),
		cartRepo:    cartRepo.New(),
		productRepo: productRepo.New(),
		addressRepo: addressRepo.New(),
	}
}

func (s *service) Create(userID uint, req *request.CreateOrderRequest, idempotentKey string) (*response.OrderResponse, error) {
	if idempotentKey == "" {
		return nil, ErrMissingIdempotentKey
	}

	// 幂等检查
	existing, err := s.orderRepo.FindIdempotency(idempotentKey)
	if err == nil && existing != nil {
		var resp response.OrderResponse
		return &resp, nil
	}

	carts, err := s.cartRepo.FindSelectedByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(carts) == 0 {
		return nil, ErrCartEmpty
	}

	// 获取默认地址
	addr, err := s.addressRepo.FindDefaultByUserID(userID)
	if err != nil {
		return nil, ErrNoDefaultAddress
	}

	var order *model.Order

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 记录幂等键
		if err := tx.Create(&model.IdempotencyRecord{Key: idempotentKey}).Error; err != nil {
			return err
		}

		orderNo := generateOrderNo()

		totalAmount := int64(0)
		var items []model.OrderItem

		for _, cart := range carts {
			// 读取 SKU（从 DB，绝不信任客户端）
			var sku model.ProductSKU
			if err := tx.First(&sku, cart.SKUID).Error; err != nil {
				return fmt.Errorf("SKU %d not found: %w", cart.SKUID, err)
			}

			// 验证商品上架
			var product model.Product
			if err := tx.First(&product, sku.ProductID).Error; err != nil {
				return fmt.Errorf("product %d not found: %w", sku.ProductID, err)
			}
			if product.Status != model.ProductStatusOn {
				return fmt.Errorf("product %s is off shelf", product.Name)
			}

			// 乐观锁扣库存
			res := tx.Model(&model.ProductSKU{}).
				Where("id = ? AND stock >= ? AND version = ?", sku.ID, cart.Quantity, sku.Version).
				Updates(map[string]interface{}{
					"stock":   gorm.Expr("stock - ?", cart.Quantity),
					"version": gorm.Expr("version + 1"),
				})
			if res.RowsAffected == 0 {
				return fmt.Errorf("%w: SKU %d", ErrInsufficientStock, sku.ID)
			}

			// 创建订单项快照
			subtotal := sku.Price * int64(cart.Quantity)
			items = append(items, model.OrderItem{
				SKUID:       sku.ID,
				ProductName: product.Name,
				SKUName:     sku.Name,
				Image:       product.Images,
				Price:       sku.Price,
				Quantity:    cart.Quantity,
				Subtotal:    subtotal,
			})
			totalAmount += subtotal
		}

		shippingFee := int64(0)
		payAmount := totalAmount + shippingFee

		order = &model.Order{
			OrderNo:        orderNo,
			UserID:         userID,
			Status:         model.OrderStatusPendingPayment,
			TotalAmount:    totalAmount,
			ShippingFee:    shippingFee,
			PayAmount:      payAmount,
			Remark:         req.Remark,
			DeliveryMethod: req.DeliveryMethod,
			AddressName:    addr.Name,
			AddressPhone:   addr.Phone,
			AddressDetail:  response.AddressToDetail(addr),
			Version:        0,
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		// 清空已购购物车
		if err := tx.Where("user_id = ? AND selected = ?", userID, true).Delete(&model.Cart{}).Error; err != nil {
			return err
		}

		// 审计日志
		if err := tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: "",
			ToStatus:   model.OrderStatusPendingPayment,
			Operator:   fmt.Sprintf("user:%d", userID),
			Note:       "订单创建",
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 重新查询获取完整的 Order + Items
	fullOrder, err := s.orderRepo.FindByID(order.ID)
	if err != nil {
		return nil, err
	}
	resp := response.ToOrderResponse(fullOrder)
	return &resp, nil
}

func (s *service) List(userID uint, req *request.ListOrderRequest) ([]response.OrderResponse, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 || size > 50 {
		size = 10
	}
	orders, total, err := s.orderRepo.ListByUserID(userID, req.Status, page, size)
	if err != nil {
		return nil, 0, err
	}
	return response.ToOrderList(orders), total, nil
}

func (s *service) GetByID(userID, orderID uint) (*response.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	if order.UserID != userID {
		return nil, ErrOrderNotBelongToUser
	}
	resp := response.ToOrderResponse(order)
	return &resp, nil
}

func (s *service) Cancel(userID, orderID uint, req *request.CancelOrderRequest) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusPendingPayment {
		return ErrInvalidOrderStatus
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"status":        model.OrderStatusCancelled,
				"cancelled_at":  time.Now(),
				"cancel_reason": req.Reason,
				"version":       order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}

		// 恢复库存
		var items []model.OrderItem
		if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			if err := tx.Model(&model.ProductSKU{}).
				Where("id = ?", item.SKUID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// 审计日志
		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: model.OrderStatusPendingPayment,
			ToStatus:   model.OrderStatusCancelled,
			Operator:   fmt.Sprintf("user:%d", userID),
			Note:       req.Reason,
		}).Error
	})
}

func (s *service) Pay(userID, orderID uint) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusPendingPayment {
		return ErrInvalidOrderStatus
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ? AND status = ?", order.ID, order.Version, model.OrderStatusPendingPayment).
			Updates(map[string]interface{}{
				"status":   model.OrderStatusPendingDelivery,
				"paid_at":  now,
				"version":  order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}
		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: model.OrderStatusPendingPayment,
			ToStatus:   model.OrderStatusPendingDelivery,
			Operator:   fmt.Sprintf("user:%d", userID),
			Note:       "支付成功",
		}).Error
	})
}

func (s *service) Ship(orderID uint) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.Status != model.OrderStatusPendingDelivery {
		return ErrInvalidOrderStatus
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"status":     model.OrderStatusPendingReceipt,
				"shipped_at": now,
				"version":    order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}
		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: model.OrderStatusPendingDelivery,
			ToStatus:   model.OrderStatusPendingReceipt,
			Operator:   "admin",
			Note:       "卖家已发货",
		}).Error
	})
}

func (s *service) Confirm(userID, orderID uint) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusPendingReceipt {
		return ErrInvalidOrderStatus
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"status":       model.OrderStatusCompleted,
				"completed_at": now,
				"version":      order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}
		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: model.OrderStatusPendingReceipt,
			ToStatus:   model.OrderStatusCompleted,
			Operator:   fmt.Sprintf("user:%d", userID),
			Note:       "确认收货",
		}).Error
	})
}

func (s *service) Rebuy(userID, orderID uint) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusCompleted {
		return ErrInvalidOrderStatus
	}

	var items []model.OrderItem
	if err := database.DB.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
		return err
	}

	var skipped []string
	for _, item := range items {
		var sku model.ProductSKU
		if err := database.DB.First(&sku, item.SKUID).Error; err != nil {
			skipped = append(skipped, item.ProductName)
			continue
		}
		var product model.Product
		if err := database.DB.First(&product, sku.ProductID).Error; err != nil || product.Status != model.ProductStatusOn {
			skipped = append(skipped, item.ProductName)
			continue
		}

		var existingCart model.Cart
		err := database.DB.Where("user_id = ? AND sku_id = ?", userID, item.SKUID).First(&existingCart).Error
		if err == nil {
			database.DB.Model(&existingCart).Update("quantity", item.Quantity)
		} else {
			database.DB.Create(&model.Cart{
				UserID:   userID,
				SKUID:    item.SKUID,
				Quantity: item.Quantity,
				Selected: true,
			})
		}
	}

	if len(skipped) > 0 {
		return fmt.Errorf("部分商品已下架，已跳过: %v", skipped)
	}
	return nil
}

func (s *service) ApplyPoints(userID, orderID uint, req *request.ApplyPointsRequest) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusPendingPayment {
		return ErrInvalidOrderStatus
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Select("points").First(&user, userID).Error; err != nil {
			return err
		}
		if user.Points < req.Points {
			return ErrInsufficientPoints
		}

		pointsAmount := int64(req.Points)
		if pointsAmount > order.PayAmount {
			pointsAmount = order.PayAmount
		}

		res := tx.Model(&model.User{}).
			Where("id = ? AND points >= ?", userID, req.Points).
			Update("points", gorm.Expr("points - ?", req.Points))
		if res.RowsAffected == 0 {
			return ErrInsufficientPoints
		}

		newDiscount := order.DiscountAmount + pointsAmount
		newPayAmount := order.TotalAmount + order.ShippingFee - newDiscount
		if newPayAmount < 0 {
			newPayAmount = 0
		}

		if err := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"discount_amount": newDiscount,
				"pay_amount":      newPayAmount,
				"points_deducted": order.PointsDeducted + req.Points,
				"version":         order.Version + 1,
			}).Error; err != nil {
			return err
		}

		var updatedUser model.User
		tx.Select("points").First(&updatedUser, userID)
		return tx.Create(&model.PointLog{
			UserID:  userID,
			Type:    model.PointTypeSpend,
			Amount:  req.Points,
			Balance: updatedUser.Points,
			OrderID: &orderID,
			Note:    fmt.Sprintf("订单抵扣 %d 积分，减免 %.2f 元", req.Points, float64(pointsAmount)/100),
		}).Error
	})
}

// orderNo 生成: YYYYMMDDHHmmss + 4随机 + 2校验 = 20位
func generateOrderNo() string {
	now := time.Now().Format("20060102150405")
	randPart := fmt.Sprintf("%04d", rand.Intn(10000))
	checkPart := checksum(now + randPart)
	return now + randPart + checkPart
}

func checksum(s string) string {
	sum := 0
	for _, c := range s {
		sum += int(c)
	}
	return fmt.Sprintf("%02d", sum%100)
}
