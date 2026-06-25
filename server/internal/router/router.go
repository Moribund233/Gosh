package router

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gosh/internal/handler/user"
	"gosh/internal/handler/address"
	"gosh/internal/handler/favorite"
	browseHistory "gosh/internal/handler/browse_history"
	"gosh/internal/handler/merchant"
	"gosh/internal/handler/category"
	"gosh/internal/handler/coupon"
	"gosh/internal/handler/flash_sale"
	"gosh/internal/handler/payment"
	"gosh/internal/handler/point"
	"gosh/internal/handler/product"
	"gosh/internal/handler/review"
	"gosh/internal/handler/banner"
	"gosh/internal/handler/upload"
	cartHandler "gosh/internal/handler/cart"
	orderHandler "gosh/internal/handler/order"
	"gosh/internal/middleware"
	"gosh/internal/model"
	"gosh/pkg/response"

	_ "gosh/docs"
)

func New(log *zap.Logger) *gin.Engine {
	r := gin.New()
	r.MaxMultipartMemory = 10 << 20 // 10MB
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.SanitizeInput())
	r.Use(middleware.RateLimitByRole(200, 100, 200, time.Minute))

	// Static files for uploads
	r.Static("/uploads", "./storage/upload")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	userH := user.NewHandler()
	addressH := address.NewHandler()
	favoriteH := favorite.NewHandler()
	browseHistoryH := browseHistory.NewHandler()
	merchantH := merchant.NewHandler()
	categoryH := category.NewHandler()
	productH := product.NewHandler()
	reviewH := review.NewHandler()
	bannerH := banner.NewHandler()
	cartH := cartHandler.NewHandler()
	orderH := orderHandler.NewHandler()
	paymentH := payment.NewHandler()
	couponH := coupon.NewHandler()
	flashSaleH := flash_sale.NewHandler()
	pointH := point.NewHandler()

	api := r.Group("/api/v1")
	{
		api.POST("/user/register", userH.Register)
		api.POST("/user/login", userH.Login)

		// Public routes
		api.GET("/categories", categoryH.Tree)
		api.GET("/categories/:id", categoryH.GetByID)
		api.GET("/products", productH.List)
		api.GET("/products/:id", productH.GetByID)
		api.GET("/products/search", productH.Search)
		api.GET("/products/hot-search", productH.HotSearch)
		api.GET("/reviews", reviewH.List)
		api.GET("/banners", bannerH.GetActive)
		api.GET("/brand-story", bannerH.GetBrandStory)
		api.GET("/payment/methods", paymentH.GetMethods)
		api.POST("/payment/callback/:method", paymentH.Callback)
		api.GET("/flash-sales", flashSaleH.ListActive)

		auth := api.Group("")
		auth.Use(middleware.Auth())
		{
			// User profile
			auth.GET("/user/profile", userH.GetProfile)
			auth.PUT("/user/profile", userH.UpdateProfile)

			// Addresses
			auth.POST("/addresses", addressH.Create)
			auth.GET("/addresses", addressH.List)
			auth.PUT("/addresses/:id", addressH.Update)
			auth.DELETE("/addresses/:id", addressH.Delete)

			// Favorites
			auth.POST("/favorites", favoriteH.Add)
			auth.POST("/favorites/remove", favoriteH.Remove)
			auth.GET("/favorites", favoriteH.List)

			// Browse history
			auth.POST("/browse-history", browseHistoryH.Add)
			auth.GET("/browse-history", browseHistoryH.List)

			// Merchant
			auth.POST("/merchant/apply", merchantH.Apply)
			auth.GET("/merchant/application", merchantH.MyApplication)

			// Reviews
			auth.POST("/reviews", reviewH.Create)

			// Search history
			auth.GET("/products/search-history", productH.SearchHistory)
			auth.POST("/products/search-history/clear", productH.ClearSearchHistory)

			// Cart
			auth.GET("/cart", cartH.List)
			auth.POST("/cart", cartH.Add)
			auth.PUT("/cart/:id", cartH.Update)
			auth.DELETE("/cart/:id", cartH.Delete)
			auth.POST("/cart/select", cartH.Select)
			auth.POST("/cart/merge", cartH.Merge)
			auth.GET("/cart/count", cartH.Count)

			// Orders
			auth.POST("/orders", orderH.Create)
			auth.GET("/orders", orderH.List)
			auth.GET("/orders/:id", orderH.GetByID)
			auth.POST("/orders/:id/cancel", orderH.Cancel)
			auth.POST("/orders/:id/pay", orderH.Pay)
			auth.POST("/orders/:id/confirm", orderH.Confirm)
			auth.POST("/orders/:id/rebuy", orderH.Rebuy)
			auth.POST("/orders/:id/apply-points", orderH.ApplyPoints)

			// Payment
			auth.POST("/payment/pay", paymentH.Pay)
			auth.GET("/payment/status/:order_no", paymentH.GetStatus)

			// Points
			auth.GET("/points", pointH.GetBalance)
			auth.GET("/points/logs", pointH.ListLogs)

			// Coupons
			auth.GET("/coupons/available", couponH.GetAvailable)
			auth.POST("/coupons/calculate", couponH.Calculate)
			auth.POST("/coupons/:id/receive", couponH.Receive)

			// Upload
			auth.POST("/upload", upload.Upload)
			auth.POST("/upload/base64", upload.UploadBase64)

			// Admin routes
			admin := auth.Group("/admin")
			admin.Use(middleware.RequireRole(model.RoleSuperAdmin, model.RoleOperator, model.RoleMerchant))
			{
				// Super admin only
				super := admin.Group("")
				super.Use(middleware.RequireRole(model.RoleSuperAdmin))
				{
					super.GET("/users", userH.ListUsers)
					super.PUT("/users/role", userH.UpdateRole)
					super.POST("/merchant/review", merchantH.Review)
					super.GET("/merchant/applications", merchantH.List)
				}

				// Category management
				admin.POST("/categories", categoryH.Create)
				admin.PUT("/categories/:id", categoryH.Update)
				admin.DELETE("/categories/:id", categoryH.Delete)

				// Product management
				admin.POST("/products", productH.Create)
				admin.PUT("/products/:id", productH.Update)
				admin.DELETE("/products/:id", productH.Delete)
				admin.PUT("/products/:id/status/:status", productH.UpdateStatus)

				// Banner management
				admin.POST("/banners", bannerH.Create)
				admin.PUT("/banners/:id", bannerH.Update)
				admin.DELETE("/banners/:id", bannerH.Delete)
				admin.GET("/banners", bannerH.List)

				// Brand story
				admin.PUT("/brand-story", bannerH.UpdateBrandStory)

				// Order management
				admin.POST("/orders/:id/ship", orderH.Ship)

				// Coupon management
				admin.POST("/coupons", couponH.Create)

				// Payment refund
				admin.POST("/payment/refund", paymentH.Refund)
			}
		}
	}

	return r
}
