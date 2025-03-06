package models

import (
	"AITodo/db"
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Title       string    `gorm:"size:255;not null" json:"title" binding:"required"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:50;default:'pending';index" json:"status"`
	StartDate   time.Time `gorm:"index" json:"start_date" binding:"required"`
	DueDate     time.Time `gorm:"index" json:"due_date" binding:"required"`
}

// 数据库操作封装
func GetAllTasks() (*[]Task, error) {
	var tasks []Task
	err := db.DB.Find(&tasks).Error
	return &tasks, err
}

func GetTaskById(id uint) (*Task, error) {
	var task Task
	err := db.DB.First(&task, id).Error
	return &task, err
}

func (t *Task) Create() error {
	return db.DB.Create(t).Error
}

func (t *Task) Update() error {
	return db.DB.Save(t).Error
}

func DeleteTask(id uint) error {
	return db.DB.Delete(&Task{}, id).Error
}
