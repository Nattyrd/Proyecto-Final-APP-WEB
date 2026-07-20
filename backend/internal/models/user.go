package models

import "time"

// RoleClient es el rol por defecto asignado al registrarse.
const RoleClient = "CLIENT"

// RoleAdmin es el rol con permisos administrativos.
const RoleAdmin = "ADMIN"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"size:120;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	FirstName string    `gorm:"size:80;not null" json:"firstName"`
	LastName  string    `gorm:"size:80;not null" json:"lastName"`
	Role      string    `gorm:"size:20;not null;default:'CLIENT'" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
