package dto

type RegisterDTO struct {
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	Email     string `json:"email" form:"email" binding:"required,email" `
	Password  string `json:"password" form:"password" binding:"required"`
}

type LoginDTO struct {
	Email    string `json:"email" form:"email" binding:"required,email" `
	Password string `json:"password" form:"password" binding:"required"`
}

type UserUpdateDTO struct {
	ID        uint64 `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	Email     string `json:"email" form:"email" binding:"required,email" `
	Password  string `json:"password" form:"password" binding:"required"`
}