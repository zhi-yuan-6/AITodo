package dto

import "AITodo/models"

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha" binding:"required"`
}

type RegisterResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
