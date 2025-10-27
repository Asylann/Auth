package model

import "time"

type User struct {
	ID       int       `json:"id"`
	Email    string    `json:"email" binding:"min=5,max=30,required" validate:"required,email"`
	Password string    `json:"password" binding:"required"`
	CreateAt time.Time `json:"createAt"`
	Role     string    `json:"role"`
}
