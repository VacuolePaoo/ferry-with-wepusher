package apis

import (
	"encoding/json"
	"ferry/pkg/logger"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetWechatOpenid 获取微信openid
func GetWechatOpenid(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "code不能为空",
		})
		return
	}

	// 从配置文件中获取微信配置
	appid := "YOUR_APPID"     // 需要替换为实际的微信公众号APPID
	secret := "YOUR_SECRET"   // 需要替换为实际的微信公众号SECRET

	// 调用微信接口获取access_token和openid
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appid + "&secret=" + secret + "&code=" + code + "&grant_type=authorization_code"
	resp, err := http.Get(url)
	if err != nil {
		logger.Errorf("调用微信接口失败，%v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取openid失败",
		})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("读取微信接口响应失败，%v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取openid失败",
		})
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Errorf("解析微信接口响应失败，%v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取openid失败",
		})
		return
	}

	// 检查是否有错误
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		logger.Errorf("微信接口返回错误，%v", result["errmsg"])
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  result["errmsg"],
		})
		return
	}

	// 获取access_token和openid
	accessToken := result["access_token"].(string)
	openid := result["openid"].(string)
	refreshToken := result["refresh_token"].(string)
	expiresIn := int(result["expires_in"].(float64))

	// 将refresh_token保存到数据库或缓存中，用于后续刷新access_token
	// TODO: 实现refresh_token的存储逻辑

	// 设置access_token的过期时间
	expireTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取openid成功",
		"data": gin.H{
			"openid": openid,
			"access_token": accessToken,
			"expires_in": expiresIn,
			"expire_time": expireTime,
		},
	})
} 