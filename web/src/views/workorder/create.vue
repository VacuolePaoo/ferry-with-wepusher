<template>
  <div class="app-container">
    <el-card>
      <div slot="header">
        <span>创建工单</span>
        <el-button type="primary" size="small" style="float: right" @click="handleWechatAuth" v-if="!hasOpenid">
          微信授权
        </el-button>
      </div>
      <el-form ref="form" :model="form" :rules="rules" label-width="80px">
        <!-- 其他表单字段 -->
        <el-form-item>
          <el-button type="primary" @click="submitForm">提交</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { getAuthUrl } from '@/api/wechat'

export default {
  name: 'WorkOrderCreate',
  data() {
    return {
      hasOpenid: false,
      form: {
        openid: '',
        // 其他表单字段
      },
      rules: {
        // 表单验证规则
      }
    }
  },
  created() {
    this.checkOpenid()
  },
  methods: {
    checkOpenid() {
      // 检查是否有openid
      const openid = this.$cookies.get('openid')
      if (openid) {
        this.hasOpenid = true
        this.form.openid = openid
      }
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
      this.$refs.form.validate(valid => {
        if (valid) {
          // 提交表单
          this.$store.dispatch('workorder/create', this.form).then(() => {
            this.$message.success('创建成功')
            this.$router.push('/workorder/list')
          })
        }
      })
    }
  }
}
</script> 