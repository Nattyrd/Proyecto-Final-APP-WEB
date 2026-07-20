package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/auth"
	"github.com/grupo5/ecommerce-api/internal/config"
	"github.com/grupo5/ecommerce-api/internal/handlers"
	"github.com/grupo5/ecommerce-api/internal/middleware"
	"github.com/grupo5/ecommerce-api/internal/models"
	"github.com/grupo5/ecommerce-api/internal/repository"
	"github.com/grupo5/ecommerce-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupProductTestEnv crea toda la cadena de dependencias con SQLite en memoria.
// IMPORTANTE: incluye apperrors.ErrorHandler() igual que el router de producción,
// para que los c.Error()+c.Abort() del middleware se serialicen con el código HTTP correcto.
func setupProductTestEnv(t *testing.T) (*gin.Engine, *auth.TokenService) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&mode=memory"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Product{}))

	cfg := &config.Config{
		JWTSecret:          "test-secret",
		JWTExpirationHours: 1,
	}
	tokenSvc := auth.NewTokenService(cfg)
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	authMiddleware := middleware.JWTAuth(tokenSvc)
	adminMiddleware := middleware.RequireAdmin()

	router := gin.New()
	// ErrorHandler convierte los errores adjuntados por apperrors.Abort en respuestas JSON
	router.Use(apperrors.ErrorHandler())
	router.POST("/api/products", authMiddleware, adminMiddleware, productHandler.Create)
	router.GET("/api/products", productHandler.GetAll)

	return router, tokenSvc
}

// makeAuthHeader genera el header "Authorization: Bearer <token>" para el usuario dado.
func makeAuthHeader(t *testing.T, tokenSvc *auth.TokenService, userID uint, username, role string) string {
	t.Helper()
	token, err := tokenSvc.Generate(userID, username, role)
	require.NoError(t, err)
	return "Bearer " + token
}

// ---------------------------------------------------------------------------
// Test 1: Usuario SIN token → debe recibir 401
// ---------------------------------------------------------------------------

func TestCreateProduct_SinToken_Devuelve401(t *testing.T) {
	router, _ := setupProductTestEnv(t)

	body := map[string]interface{}{
		"name":        "Laptop",
		"description": "Una laptop potente",
		"price":       "999.99",
		"stock":       10,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// No se agrega Authorization header

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "sin token debe devolver 401")
}

// ---------------------------------------------------------------------------
// Test 2: Usuario con rol CLIENT → debe recibir 403
// ---------------------------------------------------------------------------

func TestCreateProduct_ConTokenCliente_Devuelve403(t *testing.T) {
	router, tokenSvc := setupProductTestEnv(t)

	authHeader := makeAuthHeader(t, tokenSvc, 1, "clientuser", models.RoleClient)

	body := map[string]interface{}{
		"name":        "Laptop",
		"description": "Una laptop potente",
		"price":       "999.99",
		"stock":       10,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code, "token de CLIENT debe devolver 403")
}

// ---------------------------------------------------------------------------
// Test 3: Usuario con rol ADMIN → debe recibir 201
// ---------------------------------------------------------------------------

func TestCreateProduct_ConTokenAdmin_Devuelve201(t *testing.T) {
	router, tokenSvc := setupProductTestEnv(t)

	authHeader := makeAuthHeader(t, tokenSvc, 99, "adminuser", models.RoleAdmin)

	body := map[string]interface{}{
		"name":        "Laptop Gamer",
		"description": "Alta potencia para gaming",
		"price":       "1299.99",
		"stock":       5,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "token de ADMIN debe devolver 201")

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Laptop Gamer", resp["name"])
}

// ---------------------------------------------------------------------------
// Test 4: Listar productos es público (sin token → 200)
// ---------------------------------------------------------------------------

func TestGetProducts_SinToken_Devuelve200(t *testing.T) {
	router, _ := setupProductTestEnv(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "GET /products es público y debe devolver 200")
}
