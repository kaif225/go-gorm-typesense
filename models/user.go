package models

import "database/sql"

type Users struct {
	ID                   int            `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	FirstName            string         `gorm:"column:first_name;size:100;not null" json:"first_name,omitempty" validate:"required"`
	LastName             string         `gorm:"column:last_name;size:100;not null" json:"last_name,omitempty" validate:"required"`
	Email                string         `gorm:"column:email;size:150;unique;not null" json:"email,omitempty" validate:"required,email"`
	Username             string         `gorm:"column:username;size:100;unique;not null" json:"username,omitempty" validate:"required"`
	Password             string         `gorm:"column:password;size:255;not null" json:"password,omitempty" validate:"required,min=8,max=16"`
	PasswordChangedAt    sql.NullString `gorm:"column:password_changed_at" json:"password_changed_at,omitempty"`
	UserCreatedAt        string         `gorm:"column:user_created_at;autoCreateTime" json:"user_created_at,omitempty"`
	PasswordResetToken   sql.NullString `gorm:"column:password_reset_token" json:"password_reset_token,omitempty"`
	PasswordTokenExpired sql.NullString `gorm:"column:password_token_expired" json:"password_token_expired,omitempty"`
	InactiveStatus       bool           `gorm:"column:inactive_status;default:false" json:"inactive_status,omitempty"`
	Role                 string         `gorm:"column:role;size:50;default:user" json:"role,omitempty"`
	TypesenseSynced      bool           `gorm:"column:typesense_synced;default:false" json:"typesense_synced,omitempty"`
}

type TypesenseUser struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	UserCreatedAt  string `json:"user_created_at"`
	InactiveStatus bool   `json:"inactive_status"`
	Role           string `json:"role"`
}
