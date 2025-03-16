package dto

import (
	"AITodo/models"
	"AITodo/util"
)

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha" binding:"required"`
}

type RegisterResponse struct {
	TokenPair *util.TokenPair `json:"token_pair"`
	User      models.User     `json:"user"`
}
