<template>
  <div class="app-container">
    <el-button type="primary" size="small" style="float: right; margin-bottom: 10px;" @click="handleWechatAuth" v-if="!hasOpenid">
      微信授权
    </el-button>
    <!-- 原有的表单内容 -->
  </div>
</template>

<script>
import { getAuthUrl } from '@/api/wechat'
// ... 其他已有的imports ...

export default {
  name: 'CreateTicket',
  data() {
    return {
      hasOpenid: false,
      // ... 其他已有的data ...
    }
  },
  created() {
    this.checkOpenid()
    // ... 其他已有的created逻辑 ...
  },
  methods: {
    checkOpenid() {
      const openid = this.$cookies.get('openid')
      this.hasOpenid = !!openid
    },
    async handleWechatAuth() {
      try {
        const res = await getAuthUrl()
        if (res.code === 200) {
          window.location.href = res.data
        }
      } catch (error) {
        this.$message.error('获取授权链接失败')
      }
    },
    submitForm() {
      // 在原有的提交方法中添加openid
      const openid = this.$cookies.get('openid')
      if (!openid) {
        this.$message.warning('请先进行微信授权')
        return
      }
      this.formData.openid = openid
      // ... 原有的提交逻辑 ...
    }
    // ... 其他已有的methods ...
  }
}
</script>
// ... existing code ... 