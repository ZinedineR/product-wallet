package route

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "product-wallet/docs"
)

func (h *Router) SwaggerRouter() {
	h.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
