<template>
  <div>
    <div v-if="!formData.wechat_openid" id="login_container"></div>
    <!-- 其他原有内容 -->
  </div>
</template>

<script>
export default {
  data() {
    return {
      formData: {
        wechat_openid: '', // 添加wechat_openid字段
        // 其他原有字段
      }
    }
  },
  created() {
    // 检查是否已经授权过
    const openid = localStorage.getItem('wechat_openid')
    if (openid) {
      this.formData.wechat_openid = openid
    } else {
      // 引入微信JS SDK
      const script = document.createElement('script')
      script.src = 'https://res.wx.qq.com/connect/zh_CN/htmledition/js/wxLogin.js'
      script.onload = () => {
        this.initWechatLogin()
      }
      document.head.appendChild(script)
    }
  },
  methods: {
    initWechatLogin() {
      // 获取当前页面URL作为回调地址
      const redirectUri = encodeURIComponent(window.location.origin + '/process/create-ticket')
      const appid = 'YOUR_APPID' // 需要替换为实际的微信公众号APPID
      const scope = 'snsapi_login' // 网页应用授权作用域

      // 初始化微信登录
      const obj = new WxLogin({
        self_redirect: true,
        id: "login_container", 
        appid: appid,
        scope: scope,
        redirect_uri: redirectUri,
        state: Math.random().toString(36).substr(2),
        style: "black",
        href: "",
        onReady: (isReady) => {
          console.log('微信登录组件准备就绪:', isReady)
        }
      })
    },
    // 处理微信回调
    handleWechatCallback() {
      const urlParams = new URLSearchParams(window.location.search)
      const code = urlParams.get('code')
      const state = urlParams.get('state')
      
      if (code) {
        this.getWechatOpenid(code)
      }
    },
    async getWechatOpenid(code) {
      try {
        const response = await this.$http.get(`/api/v1/wechat/openid?code=${code}`)
        if (response.data.code === 0) {
          const openid = response.data.data.openid
          localStorage.setItem('wechat_openid', openid)
          this.formData.wechat_openid = openid
          // 清除URL中的code参数
          window.history.replaceState({}, document.title, window.location.pathname)
        }
      } catch (error) {
        console.error('获取微信openid失败:', error)
      }
    },
    // 提交工单前检查
    async submitWorkOrder() {
      if (!this.formData.wechat_openid) {
        this.$message.error('请先进行微信授权登录')
        return
      }
      // 其他提交逻辑
    }
  }
}
</script>

<style>
#login_container {
  width: 300px;
  height: 300px;
  margin: 0 auto;
}
</style> 