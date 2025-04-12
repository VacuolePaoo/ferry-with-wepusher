package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateWorkOrder(c *gin.Context) {
	var workOrder models.WorkOrderInfo
	if err := c.ShouldBindJSON(&workOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 获取openid
	openid, err := c.Cookie("openid")
	if err == nil && openid != "" {
		workOrder.CreatorOpenID = openid
	}

	// 创建工单
	if err := models.CreateWorkOrder(&workOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建工单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建工单成功",
		"data": workOrder,
	})
} 