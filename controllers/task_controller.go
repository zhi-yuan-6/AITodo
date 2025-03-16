package controllers

import (
	"AITodo/models"
	"AITodo/services"
	"AITodo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetAllTasks 获取所有任务
func GetAllTasks(c *gin.Context) {
	tasks, err := services.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// CreateTask 创建新任务
func CreateTask(c *gin.Context) {
	uid, err := util.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"解析task失败": err.Error()})
		return
	}
	task.UserID = uid

	if err := services.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"创建任务失败": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": "创建成功"})
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "无效的任务ID"})
		return
	}

	var req models.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := services.UpdateTask(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := services.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
