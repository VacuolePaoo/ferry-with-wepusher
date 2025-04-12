package wechat

import (
    "encoding/json"
    "ferry/pkg/notify/wechat"
    "ferry/tools/app"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
)

// 获取微信授权URL
func GetAuthUrl(c *gin.Context) {
    state := c.Query("state")  // state可以用来传递其他参数
    url := wechat.GetOAuthUrl(state)
    app.OK(c, map[string]string{"url": url}, "获取授权URL成功")
}

// 微信授权回调
func Callback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        app.Error(c, -1, fmt.Errorf("未获取到授权码"), "")
        return
    }

    // 获取access_token
    appId := viper.GetString("settings.wechat.appid")
    appSecret := viper.GetString("settings.wechat.appsecret")
    url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
        appId, appSecret, code)
    
    resp, err := http.Get(url)
    if err != nil {
        app.Error(c, -1, err, "获取access_token失败")
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var tokenResp struct {
        AccessToken  string `json:"access_token"`
        ExpiresIn   int    `json:"expires_in"`
        RefreshToken string `json:"refresh_token"`
        OpenId      string `json:"openid"`
        Scope       string `json:"scope"`
        ErrCode     int    `json:"errcode"`
        ErrMsg      string `json:"errmsg"`
    }
    if err := json.Unmarshal(body, &tokenResp); err != nil {
        app.Error(c, -1, err, "解析access_token响应失败")
        return
    }

    if tokenResp.ErrCode != 0 {
        app.Error(c, -1, fmt.Errorf(tokenResp.ErrMsg), "获取access_token失败")
        return
    }

    // 获取用户信息
    userInfoUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
        tokenResp.AccessToken, tokenResp.OpenId)
    resp, err = http.Get(userInfoUrl)
    if err != nil {
        app.Error(c, -1, err, "获取用户信息失败")
        return
    }
    defer resp.Body.Close()

    body, _ = ioutil.ReadAll(resp.Body)
    var userInfo struct {
        OpenId     string   `json:"openid"`
        Nickname   string   `json:"nickname"`
        Sex        int      `json:"sex"`
        Province   string   `json:"province"`
        City       string   `json:"city"`
        Country    string   `json:"country"`
        HeadImgUrl string   `json:"headimgurl"`
        Privilege  []string `json:"privilege"`
        UnionId    string   `json:"unionid"`
        ErrCode    int      `json:"errcode"`
        ErrMsg     string   `json:"errmsg"`
    }
    if err := json.Unmarshal(body, &userInfo); err != nil {
        app.Error(c, -1, err, "解析用户信息失败")
        return
    }

    if userInfo.ErrCode != 0 {
        app.Error(c, -1, fmt.Errorf(userInfo.ErrMsg), "获取用户信息失败")
        return
    }

    // 返回用户信息
    app.OK(c, userInfo, "获取用户信息成功")
} 