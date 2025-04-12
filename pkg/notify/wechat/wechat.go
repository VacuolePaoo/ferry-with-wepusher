package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"ferry/pkg/logger"
	"github.com/spf13/viper"
)

const (
	GetAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
	SendTemplateMsgUrl = "https://api.weixin.qq.com/cgi-bin/message/template/send"
)

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// 模板消息数据项
type TemplateDataItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// 模板消息主体结构
type TemplateMsg struct {
	ToUser     string                    `json:"touser"`
	TemplateID string                    `json:"template_id"`
	URL        string                    `json:"url,omitempty"` // 可选的跳转URL
	Data       map[string]TemplateDataItem `json:"data"`
}

type SendTemplateMsgResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgID   int64  `json:"msgid"`
}

// 获取access_token
func GetAccessToken() (string, error) {
	appid := viper.GetString("settings.wechat.appid")
	appsecret := viper.GetString("settings.wechat.appsecret")
	
	if appid == "" || appsecret == "" {
		return "", fmt.Errorf("微信配置不完整，请检查配置文件中的appid和appsecret")
	}

	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", 
		GetAccessTokenUrl, appid, appsecret)
	
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("获取access_token请求失败：", err)
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取access_token响应失败：", err)
		return "", err
	}
	
	var tokenResp AccessTokenResp
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		logger.Error("解析access_token响应失败：", err)
		return "", err
	}
	
	if tokenResp.ErrCode != 0 {
		logger.Error("获取access_token失败：", tokenResp.ErrMsg)
		return "", fmt.Errorf("获取access_token失败: %s", tokenResp.ErrMsg)
	}
	
	return tokenResp.AccessToken, nil
}

// 发送工单处理结果通知
func SendWorkOrderResult(openid, title, result, remarks, time string) error {
	if openid == "" {
		return fmt.Errorf("openid不能为空")
	}

	// 检查微信是否启用
	if !viper.GetBool("settings.wechat.enable") {
		logger.Warn("微信通知未启用")
		return nil // 或者返回特定错误，取决于业务逻辑
	}

	templateId := viper.GetString("settings.wechat.template_id")
	if templateId == "" {
		return fmt.Errorf("未配置模板ID，请检查配置文件")
	}

	// 获取模板格式配置
	templateFormatConfig := viper.GetStringMap("settings.wechat.template_format")
	if len(templateFormatConfig) == 0 {
		return fmt.Errorf("未配置模板格式 template_format，请检查配置文件")
	}

	accessToken, err := GetAccessToken()
	if err != nil {
		return err
	}

	// 构建模板消息数据
	msgData := make(map[string]TemplateDataItem)
	placeholders := map[string]string{
		"{title}":   title,
		"{result}":  result,
		"{remarks}": remarks,
		"{time}":    time,
	}

	for key, formatItem := range templateFormatConfig {
		itemConfig, ok := formatItem.(map[string]interface{})
		if !ok {
			logger.Warnf("无效的模板格式项配置：%s", key)
			continue
		}
		value := fmt.Sprintf("%v", itemConfig["value"])
		color := fmt.Sprintf("%v", itemConfig["color"])

		// 替换占位符
		for placeholder, actualValue := range placeholders {
			value = strings.ReplaceAll(value, placeholder, actualValue)
		}

		msgData[key] = TemplateDataItem{
			Value: value,
			Color: color,
		}
	}

	msg := TemplateMsg{
		ToUser:     openid,
		TemplateID: templateId,
		Data:       msgData,
		// URL: "http://your-detail-page.com/" + workOrderId, // 可以选择添加跳转URL
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error("序列化消息失败：", err)
		return err
	}

	url := fmt.Sprintf("%s?access_token=%s", SendTemplateMsgUrl, accessToken)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		logger.Error("发送模板消息请求失败：", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取模板消息响应失败：", err)
		return err
	}

	var resultResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal(body, &resultResp)
	if err != nil {
		logger.Error("解析模板消息响应失败：", err)
		return err
	}

	if resultResp.ErrCode != 0 {
		logger.Errorf("发送模板消息失败：%d - %s", resultResp.ErrCode, resultResp.ErrMsg)
		return fmt.Errorf("发送模板消息失败: %s", resultResp.ErrMsg)
	}

	logger.Infof("成功向 %s 发送模板消息", openid)
	return nil
} 