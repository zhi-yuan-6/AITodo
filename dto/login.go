package dto

import (
	"AITodo/models"
	"AITodo/util"
)

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` //手机号或邮箱
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	TokenPair *util.TokenPair `json:"token_pair"`
	User      models.User     `json:"user"`
}
