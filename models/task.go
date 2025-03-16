package models

import (
	"AITodo/db"
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `json:"user_id"`
	Title       string    `gorm:"size:255;not null" json:"title" binding:"required"`
	Category    string    `gorm:"size:100;default:'其他' json:"category" binding:"required"` //工作、学习、生活、健身
	Location    string    `gorm:"size:255" json:"location"`
	Description string    `gorm:"type:text" json:"description"`
	StartDate   time.Time `gorm:"index" json:"start_date" binding:"required"`
	DueDate     time.Time `gorm:"index" json:"due_date" binding:"required"`
	Status      string    `gorm:"size:50;default:'pending';index" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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

func CountAllTasks() (int64, error) {
	var count int64
	err := db.DB.Table("tasks").Count(&count).Error
	return count, err
}

type Result struct {
	Category string
	Total    int
}

func CountTaskByCategory() *[]Result {
	var rets *[]Result
	db.DB.Table("tasks").Select("category,sum(category)").Group("category").Scan(rets)
	return rets
}
