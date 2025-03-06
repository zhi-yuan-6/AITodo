package services

import (
	"AITodo/models"
	"fmt"
)

// GetAllTasks 获取所有任务
func GetAllTasks() (*[]models.Task, error) {
	tasks, err := models.GetAllTasks()
	if err != nil {
		return &[]models.Task{}, fmt.Errorf("无法获取任务列表:%w", err)
	}
	return tasks, nil
}

// CreateTask 创建新任务
func CreateTask(task models.Task) error {
	if err := task.Create(); err != nil {
		return fmt.Errorf("创建任务失败:%w", err)
	}
	return nil
}

// UpdateTask 更新任务
func UpdateTask(id uint, req models.Task) (*models.Task, error) {
	task, err := models.GetTaskById(id)
	if err != nil {
		return &models.Task{}, fmt.Errorf("获取任务失败:%w", err)
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.StartDate = req.StartDate
	task.DueDate = req.DueDate

	if err := task.Update(); err != nil {
		return &models.Task{}, fmt.Errorf("更新任务失败:%w", err)
	}

	return task, nil
}

// DeleteTask 删除任务
func DeleteTask(id uint) error {
	if err := models.DeleteTask(id); err != nil {
		return fmt.Errorf("删除任务失败:%w", err)
	}
	return nil
}
