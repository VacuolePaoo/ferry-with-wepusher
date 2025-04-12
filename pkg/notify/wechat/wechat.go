package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"ferry/pkg/logger"
	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
  @Desc : 微信服务号通知
*/

const (
	GetAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
	SendTemplateMsgUrl = "https://api.weixin.qq.com/cgi-bin/message/template/send"
)

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

type TemplateMsg struct {
	Touser     string                 `json:"touser"`
	TemplateId string                 `json:"template_id"`
	Url        string                 `json:"url"`
	Data       map[string]interface{} `json:"data"`
}

type SendTemplateMsgResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Msgid   int64  `json:"msgid"`
}

// 获取微信服务号的access_token
func GetAccessToken(appId, appSecret string) string {
	res, err := http.Get(fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", 
		GetAccessTokenUrl, appId, appSecret))
	if err != nil {
		logger.Error(err)
		return ""
	}
	defer res.Body.Close()
	
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return ""
	}
	
	var token AccessTokenResp
	err = json.Unmarshal(body, &token)
	if err != nil {
		logger.Error(err)
		return ""
	}
	
	if token.Errcode != 0 {
		logger.Error(token.Errmsg)
		return ""
	}
	
	return token.AccessToken
}

// 发送模板消息
func SendTemplateMsg(token string, msg TemplateMsg) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(fmt.Sprintf("%s?access_token=%s", SendTemplateMsgUrl, token), 
		"application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	var result SendTemplateMsgResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	
	if result.Errcode != 0 {
		return fmt.Errorf("发送模板消息失败: %s", result.Errmsg)
	}
	
	return nil
}

// 发送工单通知到微信服务号
func SendWechatMsg(openid string, url string, title string, creator string, priority string, createdAt string) error {
	appId := viper.GetString("settings.wechat.appid")
	appSecret := viper.GetString("settings.wechat.appsecret")
	templateId := viper.GetString("settings.wechat.template_id")
	
	// 获取access_token
	token := GetAccessToken(appId, appSecret)
	if token == "" {
		return fmt.Errorf("获取access_token失败")
	}
	
	// 构建模板消息
	msg := TemplateMsg{
		Touser:     openid,
		TemplateId: templateId,
		Url:        url,
		Data: map[string]interface{}{
			"first": map[string]string{
				"value": "您有一个工单审批已完成",
				"color": "#173177",
			},
			"keyword1": map[string]string{
				"value": title,
				"color": "#173177",
			},
			"keyword2": map[string]string{
				"value": creator,
				"color": "#173177",
			},
			"keyword3": map[string]string{
				"value": priority,
				"color": "#173177",
			},
			"keyword4": map[string]string{
				"value": createdAt,
				"color": "#173177",
			},
			"remark": map[string]string{
				"value": "请点击查看详情",
				"color": "#173177",
			},
		},
	}
	
	// 发送模板消息
	return SendTemplateMsg(token, msg)
} 