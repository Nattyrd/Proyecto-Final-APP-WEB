package dto

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email,max=120"`
	Password  string `json:"password" binding:"required,min=6,max=72"`
	FirstName string `json:"firstName" binding:"required,min=2,max=80"`
	LastName  string `json:"lastName" binding:"required,min=2,max=80"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email,max=120"`
	FirstName string `json:"firstName" binding:"omitempty,min=2,max=80"`
	LastName  string `json:"lastName" binding:"omitempty,min=2,max=80"`
	Password  string `json:"password" binding:"omitempty,min=6,max=72"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
