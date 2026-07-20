package service

import (
	"context"
	"errors"

	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/auth"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/models"
	"github.com/grupo5/ecommerce-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo     *repository.UserRepository
	tokenService *auth.TokenService
	adminSecret  string // leído desde ADMIN_SECRET en .env
}

func NewUserService(userRepo *repository.UserRepository, tokenService *auth.TokenService, adminSecret string) *UserService {
	return &UserService{
		userRepo:     userRepo,
		tokenService: tokenService,
		adminSecret:  adminSecret,
	}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.Internal("No se pudo validar el usuario")
	}
	if exists {
		return nil, apperrors.Conflict("El nombre de usuario ya esta registrado")
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.Internal("No se pudo validar el correo")
	}
	if exists {
		return nil, apperrors.Conflict("El correo ya esta registrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal("No se pudo cifrar la contrasena")
	}

	// Asignar rol ADMIN solo si el secret coincide con el de .env y no está vacío.
	role := models.RoleClient
	if s.adminSecret != "" && req.AdminSecret == s.adminSecret {
		role = models.RoleAdmin
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.Internal("No se pudo registrar el usuario")
	}

	token, err := s.tokenService.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal("No se pudo generar el token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Unauthorized("Credenciales invalidas")
		}
		return nil, apperrors.Internal("No se pudo autenticar al usuario")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.Unauthorized("Credenciales invalidas")
	}

	token, err := s.tokenService.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal("No se pudo generar el token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Usuario no encontrado")
		}
		return nil, apperrors.Internal("No se pudo consultar el usuario")
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *UserService) Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Usuario no encontrado")
		}
		return nil, apperrors.Internal("No se pudo consultar el usuario")
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, apperrors.Internal("No se pudo validar el correo")
		}
		if exists {
			return nil, apperrors.Conflict("El correo ya esta registrado")
		}
		user.Email = req.Email
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, apperrors.Internal("No se pudo cifrar la contrasena")
		}
		user.Password = string(hashedPassword)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("No se pudo actualizar el usuario")
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	if _, err := s.userRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound("Usuario no encontrado")
		}
		return apperrors.Internal("No se pudo consultar el usuario")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return apperrors.Internal("No se pudo eliminar el usuario")
	}

	return nil
}

func toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}
}
