<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-form">
        <h2 class="register-title">注册</h2>
        <t-form ref="form" :data="formData" @submit="handleSubmit">
          <t-form-item label="用户名" name="username">
            <t-input
              v-model="formData.username"
              placeholder="请输入用户名"
              size="large"
            />
          </t-form-item>
          <t-form-item label="邮箱" name="email">
            <t-input
              v-model="formData.email"
              placeholder="请输入邮箱"
              size="large"
            />
          </t-form-item>
          <t-form-item label="密码" name="password">
            <t-input
              v-model="formData.password"
              type="password"
              placeholder="请输入密码"
              size="large"
            />
          </t-form-item>
          <t-form-item label="确认密码" name="confirmPassword">
            <t-input
              v-model="formData.confirmPassword"
              type="password"
              placeholder="请确认密码"
              size="large"
            />
          </t-form-item>
          <t-form-item>
            <t-button 
              theme="primary" 
              type="submit" 
              size="large" 
              block
              :loading="loading"
            >
              注册
            </t-button>
          </t-form-item>
        </t-form>
        <div class="register-footer">
          <span>已有账号？</span>
          <t-link @click="goToLogin">立即登录</t-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { authApi } from '@/api/auth'

const router = useRouter()

const loading = ref(false)
const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const handleSubmit = async () => {
  if (formData.password !== formData.confirmPassword) {
    MessagePlugin.error('两次密码输入不一致')
    return
  }
  
  loading.value = true
  try {
    await authApi.register(formData)
    MessagePlugin.success('注册成功，请登录')
    router.push('/login')
  } catch (error) {
    MessagePlugin.error('注册失败，请重试')
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-container {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  width: 100%;
  max-width: 400px;
}

.register-title {
  text-align: center;
  margin-bottom: 30px;
  font-size: 24px;
  color: #1f2937;
}

.register-footer {
  text-align: center;
  margin-top: 20px;
  color: #6b7280;
}
</style>