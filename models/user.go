package models

import (
	"AITodo/db"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	UserStatusInactiv = 0
	UserStatusActive  = 1
	UserStatusBanned  = 2
)

type User struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"user_name"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"` //密码哈希不序列化到json
	Phone    *string   `gorm:"type:varchar(20);uniqueIndex;default:NUll" json:"phone"`
	Email    *string   `gorm:"type:varchar(255);uniqueIndex;default:NULL" json:"email"` //带*表示允许为NULL
	Status   int       `gorm:"type:tinyint(1);default:0" json:"status"`
	CreateAt time.Time `gorm:"column:create_at;type:datetime;autoCreateTime" json:"create_at"`
	UpdateAt time.Time `gorm:"column:create_at;type:datetime;autoUpdateTime" json:"update_at"`
}

func (User) TableName() string {
	return "users"
}

func SearchByPhone(phone string) (*User, error) {
	var user User
	result := db.DB.Where("phone=?", phone).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户名或密码错误")
		} else {
			return nil, fmt.Errorf("查询失败:%v", result.Error)
		}
	}
	return &user, nil
}

func SearchByEmail(email string) (*User, error) {
	var user User
	result := db.DB.Where("email=?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户名或密码错误")
		} else {
			return nil, fmt.Errorf("查询失败:%v", result.Error)
		}
	}
	return &user, nil
}

func PhoneExists(phone string) (bool, error) {
	var user User
	result := db.DB.Where("phone=?", phone).First(&user)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func CreateUser(user *User) error {
	return db.DB.Create(user).Error
}
