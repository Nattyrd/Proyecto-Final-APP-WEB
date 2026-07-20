// @title           E-commerce API - Grupo 5
// @version         1.0
// @description     API RESTful de e-commerce desarrollada con Go, Gin, GORM, PostgreSQL y JWT.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Grupo 5 - Aplicaciones Web
// @contact.email  grupo5@epn.edu.ec

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ingrese el token JWT con el prefijo Bearer. Ejemplo: Bearer eyJhbGciOiJIUzI1NiIs...
package main

import (
	"log"

	"github.com/grupo5/ecommerce-api/internal/auth"
	"github.com/grupo5/ecommerce-api/internal/config"
	"github.com/grupo5/ecommerce-api/internal/database"
	"github.com/grupo5/ecommerce-api/internal/handlers"
	"github.com/grupo5/ecommerce-api/internal/repository"
	"github.com/grupo5/ecommerce-api/internal/router"
	"github.com/grupo5/ecommerce-api/internal/service"

	_ "github.com/grupo5/ecommerce-api/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuracion: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	receiptRepo := repository.NewReceiptRepository(db)

	tokenService := auth.NewTokenService(cfg)
	userService := service.NewUserService(userRepo, tokenService, cfg.AdminSecret)
	productService := service.NewProductService(productRepo)
	receiptService := service.NewReceiptService(receiptRepo, productRepo, userRepo)

	engine := router.Setup(router.Dependencies{
		TokenService:   tokenService,
		UserHandler:    handlers.NewUserHandler(userService),
		ProductHandler: handlers.NewProductHandler(productService),
		ReceiptHandler: handlers.NewReceiptHandler(receiptService),
	})

	log.Printf("Servidor iniciado en http://localhost:%s", cfg.AppPort)
	log.Printf("Swagger disponible en http://localhost:%s/swagger/index.html", cfg.AppPort)

	if err := engine.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
