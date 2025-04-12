package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateWorkOrder(c *gin.Context) {
	var (
		workOrder work_order.WorkOrderInfo
		err       error
	)

	if err = c.ShouldBind(&workOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// 检查是否包含微信openid
	if workOrder.WechatOpenid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "请先进行微信授权登录",
		})
		return
	}

	// ... 其他原有逻辑 ...

	// 保存工单
	if err = models.DB.Create(&workOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建工单成功",
		"data": workOrder,
	})
} 