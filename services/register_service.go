package services

import (
	"AITodo/dto"
	"AITodo/models"
	"AITodo/util"
	"fmt"
)

// 验证验证码是否正确
func RegisterService(req dto.RegisterRequest) (*models.User, error) {
	//验证验证码否符合要求
	if valid := util.VerifySMS(req.Phone, req.Captcha); !valid {
		return nil, fmt.Errorf("验证码错误")
	}
	//验证手机号格式
	if !util.IsPhone(req.Phone) {
		return nil, fmt.Errorf("手机号格式不正确")
	}
	//验证手机号唯一性
	if req.Phone != "" {
		if exists, _ := models.PhoneExists(req.Phone); exists {
			return nil, fmt.Errorf("手机号%s已被注册", req.Phone)
		}
	}

	//验证密码是否合规
	if err := util.ValidatePasswordPolicy(req.Password); err != nil {
		return nil, fmt.Errorf("密码不符合要求")
	}

	//对密码进行加密
	password, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("加密失败:%w", err)
	}

	//创建用户对象
	user := models.User{
		Phone:    &req.Phone,
		Password: password,
		Status:   models.UserStatusActive,
	}
	//保存到数据库
	if err = models.CreateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
