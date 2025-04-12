package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	GetUserInfoUrl = "https://api.weixin.qq.com/sns/userinfo"
)

type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}

// 获取微信授权URL
func GetAuthUrl(c *gin.Context) {
	appid := viper.GetString("settings.wechat.appid")
	redirectUri := viper.GetString("settings.wechat.redirect_uri")
	scope := "snsapi_base" // 只获取openid
	state := "STATE"

	authUrl := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		appid, redirectUri, scope, state)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": authUrl,
	})
}

// 微信授权回调
func Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "code不能为空",
		})
		return
	}

	appid := viper.GetString("settings.wechat.appid")
	appsecret := viper.GetString("settings.wechat.appsecret")

	// 获取access_token
	accessTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appid, appsecret, code)

	resp, err := http.Get(accessTokenUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取access_token失败",
		})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "解析access_token失败",
		})
		return
	}

	if _, ok := result["errcode"]; ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  result["errmsg"],
		})
		return
	}

	openid := result["openid"].(string)
	if openid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取openid失败",
		})
		return
	}

	// 将openid写入cookie
	c.SetCookie("openid", openid, 3600*24*30, "/", "", false, true)

	// 重定向回工单填写页面
	c.Redirect(http.StatusFound, "/#/process/create-ticket")
}

// 更新用户OpenID
func updateUserOpenID(c *gin.Context, openid string) error {
	// TODO: 实现更新用户OpenID的逻辑
	return nil
} 