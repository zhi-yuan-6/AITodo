package dto

import "AITodo/models"

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` //手机号或邮箱
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
