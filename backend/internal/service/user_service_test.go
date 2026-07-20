package service_test

import (
	"context"
	"testing"

	"github.com/grupo5/ecommerce-api/internal/auth"
	"github.com/grupo5/ecommerce-api/internal/config"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/models"
	"github.com/grupo5/ecommerce-api/internal/repository"
	"github.com/grupo5/ecommerce-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// testAdminSecret es la clave que se inyecta al servicio durante los tests.
// En producción este valor viene de ADMIN_SECRET en el .env.
const testAdminSecret = "test-admin-secret-123"

// setupTestDB crea una base de datos SQLite en memoria y migra el modelo User.
// Cada test que llame a esta función obtiene una BD limpia e independiente.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&mode=memory"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "no se pudo abrir la BD en memoria")

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err, "no se pudo migrar el modelo User")

	return db
}

// newTestUserService construye las dependencias reales con la BD de test.
func newTestUserService(t *testing.T) (*service.UserService, *gorm.DB) {
	t.Helper()
	db := setupTestDB(t)

	cfg := &config.Config{
		JWTSecret:          "test-secret-key",
		JWTExpirationHours: 1,
	}
	tokenSvc := auth.NewTokenService(cfg)
	userRepo := repository.NewUserRepository(db)
	// Pasamos el secret de test: en producción viene de cfg.AdminSecret (ADMIN_SECRET en .env)
	userSvc := service.NewUserService(userRepo, tokenSvc, testAdminSecret)

	return userSvc, db
}

// ---------------------------------------------------------------------------
// Tests de Register
// ---------------------------------------------------------------------------

func TestRegister_Exitoso_RolClientePorDefecto(t *testing.T) {
	svc, _ := newTestUserService(t)
	ctx := context.Background()

	req := dto.RegisterRequest{
		Username:  "juanito",
		Email:     "juan@test.com",
		Password:  "password123",
		FirstName: "Juan",
		LastName:  "Perez",
		// No se envía adminSecret → rol CLIENT
	}

	resp, err := svc.Register(ctx, req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token, "debe devolver un token JWT")
	assert.Equal(t, "juanito", resp.User.Username)
	assert.Equal(t, "juan@test.com", resp.User.Email)
	assert.Equal(t, models.RoleClient, resp.User.Role, "el rol debe ser CLIENT por defecto")
}

func TestRegister_ConSecretoAdmin_RolAdmin(t *testing.T) {
	svc, _ := newTestUserService(t)
	ctx := context.Background()

	req := dto.RegisterRequest{
		Username:    "adminuser",
		Email:       "admin@test.com",
		Password:    "adminpass123",
		FirstName:   "Admin",
		LastName:    "User",
		AdminSecret: testAdminSecret, // mismo valor inyectado al servicio
	}

	resp, err := svc.Register(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, resp.User.Role, "el rol debe ser ADMIN con el código correcto")
}

func TestRegister_UsernameYaExiste_DevuelveConflict(t *testing.T) {
	svc, _ := newTestUserService(t)
	ctx := context.Background()

	req := dto.RegisterRequest{
		Username:  "duplicado",
		Email:     "original@test.com",
		Password:  "pass1234",
		FirstName: "A",
		LastName:  "B",
	}
	_, err := svc.Register(ctx, req)
	require.NoError(t, err, "primer registro debe ser exitoso")

	// Segundo intento con el mismo username
	req2 := dto.RegisterRequest{
		Username:  "duplicado",
		Email:     "otro@test.com",
		Password:  "pass1234",
		FirstName: "C",
		LastName:  "D",
	}
	_, err2 := svc.Register(ctx, req2)

	require.Error(t, err2)
	assert.Contains(t, err2.Error(), "usuario ya esta registrado")
}

// ---------------------------------------------------------------------------
// Tests de Login
// ---------------------------------------------------------------------------

func TestLogin_Exitoso(t *testing.T) {
	svc, _ := newTestUserService(t)
	ctx := context.Background()

	// Primero registramos el usuario
	_, err := svc.Register(ctx, dto.RegisterRequest{
		Username:  "loginuser",
		Email:     "login@test.com",
		Password:  "secret123",
		FirstName: "Login",
		LastName:  "User",
	})
	require.NoError(t, err)

	// Ahora hacemos login
	resp, err := svc.Login(ctx, dto.LoginRequest{
		Username: "loginuser",
		Password: "secret123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token, "debe devolver un token JWT")
	assert.Equal(t, "loginuser", resp.User.Username)
	assert.Equal(t, models.RoleClient, resp.User.Role)
}

func TestLogin_ContrasenaIncorrecta_DevuelveUnauthorized(t *testing.T) {
	svc, _ := newTestUserService(t)
	ctx := context.Background()

	_, _ = svc.Register(ctx, dto.RegisterRequest{
		Username:  "userpass",
		Email:     "pass@test.com",
		Password:  "correctpass",
		FirstName: "U",
		LastName:  "P",
	})

	_, err := svc.Login(ctx, dto.LoginRequest{
		Username: "userpass",
		Password: "wrongpass",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Credenciales invalidas")
}

func TestLogin_UsuarioNoExiste_DevuelveUnauthorized(t *testing.T) {
	svc, _ := newTestUserService(t)
	ctx := context.Background()

	_, err := svc.Login(ctx, dto.LoginRequest{
		Username: "noexiste",
		Password: "anypass",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Credenciales invalidas")
}
