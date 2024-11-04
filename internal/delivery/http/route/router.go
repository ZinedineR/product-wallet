package route

import (
	"github.com/gin-gonic/gin"
	"product-wallet/internal/delivery/http"
	api "product-wallet/internal/delivery/http/middleware"
)

type Router struct {
	App                *gin.Engine
	UserHandler        *http.UserHTTPHandler
	ProductHandler     *http.ProductHTTPHandler
	WalletHandler      *http.WalletHTTPHandler
	TransactionHandler *http.TransactionHTTPHandler
	AuthMiddleware     *api.AuthMiddleware
}

func (h *Router) Setup() {
	h.App.Use(h.AuthMiddleware.ErrorHandler)

	// Guest routes for unauthenticated users
	guestApi := h.App.Group("/auth")
	{
		guestApi.POST("/register", h.UserHandler.Register)
		guestApi.POST("/login", h.UserHandler.Login)
	}

	// Private routes for authenticated users
	privateApi := h.App.Group("")
	privateApi.Use(h.AuthMiddleware.JWTAuthentication)
	{
		// Product Routes
		productApi := privateApi.Group("/products")
		{
			productApi.POST("", h.ProductHandler.Create)
			productApi.PUT("/:id", h.ProductHandler.Update)
			productApi.GET("", h.ProductHandler.Find)
			productApi.GET("/:id", h.ProductHandler.Detail)
			productApi.DELETE("/:id", h.ProductHandler.Delete)
		}

		// Wallet Routes
		walletApi := privateApi.Group("/wallets")
		{
			walletApi.POST("", h.WalletHandler.Create)
			walletApi.PUT("/:id", h.WalletHandler.Update)
			walletApi.GET("", h.WalletHandler.Find)
			walletApi.GET("/:id", h.WalletHandler.Detail)
			walletApi.GET("/transaction/:id", h.WalletHandler.DetailWalletTransaction)
			walletApi.DELETE("/:id", h.WalletHandler.Delete)
		}

		// Transaction Routes
		transactionApi := privateApi.Group("/transactions")
		{
			transactionApi.POST("", h.TransactionHandler.Create)
			transactionApi.GET("/:id", h.TransactionHandler.Detail)
			transactionApi.GET("", h.TransactionHandler.Find)
			transactionApi.POST("/credit", h.TransactionHandler.Credit)
			transactionApi.POST("/transfer", h.TransactionHandler.Transfer)
			transactionApi.DELETE("/:id", h.TransactionHandler.Delete)
		}
	}
}
