package services

import (
	"AITodo/dto"
	"AITodo/models"
	"AITodo/util"
	"fmt"
)

func LoginService(req dto.LoginRequest) (*models.User, error) {
	var identifier = req.Identifier
	//根据判断用户使用手机号还是者邮箱
	if util.IsPhone(identifier) {
		//查询验证
		user, err := models.SearchByPhone(identifier)
		if err != nil {
			return nil, fmt.Errorf("该手机号未注册:%w", err)
		}
		if user.Status != models.UserStatusActive {
			return nil, fmt.Errorf("账户未激活或已被禁用")
		}
		//验证密码
		err = util.VerifyPassword(req.Password, user.Password)
		if err != nil {
			return nil, fmt.Errorf("密码错误:%w", err)
		}
		return user, nil
	} else if util.IsEmaill(identifier) {
		//查询验证
		user, err := models.SearchByEmail(identifier)
		if err != nil {
			return nil, fmt.Errorf("该邮箱未注册:%w", err)
		}
		if user.Status != models.UserStatusActive {
			return nil, fmt.Errorf("账户未激活或已被禁用")
		}
		err = util.VerifyPassword(user.Password, req.Password)
		if err != nil {
			return nil, fmt.Errorf("密码错误:%w", err)
		}
		return user, nil
	} else {
		return nil, fmt.Errorf("格式错误")
	}
}
