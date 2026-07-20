package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/auth"
	"github.com/grupo5/ecommerce-api/internal/handlers"
	"github.com/grupo5/ecommerce-api/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Dependencies struct {
	TokenService   *auth.TokenService
	UserHandler    *handlers.UserHandler
	ProductHandler *handlers.ProductHandler
	ReceiptHandler *handlers.ReceiptHandler
}

func Setup(deps Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(), apperrors.ErrorHandler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ecommerce-api-grupo5",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := middleware.JWTAuth(deps.TokenService)
	adminMiddleware := middleware.RequireAdmin() // sólo ADMIN puede escribir productos

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", deps.UserHandler.Register)
			users.POST("/login", deps.UserHandler.Login)
			users.GET("/:id", deps.UserHandler.GetByID)
			users.PUT("/:id", authMiddleware, deps.UserHandler.Update)
			users.DELETE("/:id", authMiddleware, deps.UserHandler.Delete)
		}

		products := api.Group("/products")
		{
			// Lectura pública (sin token)
			products.GET("", deps.ProductHandler.GetAll)
			products.GET("/:id", deps.ProductHandler.GetByID)

			// Escritura: primero autenticar (JWT), luego verificar rol ADMIN
			products.POST("", authMiddleware, adminMiddleware, deps.ProductHandler.Create)
			products.PUT("/:id", authMiddleware, adminMiddleware, deps.ProductHandler.Update)
			products.DELETE("/:id", authMiddleware, adminMiddleware, deps.ProductHandler.Delete)
		}

		receipts := api.Group("/receipts")
		receipts.Use(authMiddleware)
		{
			receipts.POST("", deps.ReceiptHandler.Create)
			receipts.GET("", deps.ReceiptHandler.GetAll)
			receipts.GET("/user/:userId", deps.ReceiptHandler.GetByUserID)
			receipts.GET("/:id", deps.ReceiptHandler.GetByID)
			receipts.DELETE("/:id", deps.ReceiptHandler.Delete)
		}
	}

	return router
}
