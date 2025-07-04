<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-form">
        <h2 class="login-title">登录</h2>
        <t-form ref="form" :data="formData" @submit="handleSubmit">
          <t-form-item label="用户名" name="username">
            <t-input
              v-model="formData.username"
              placeholder="请输入用户名"
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
          <t-form-item>
            <t-button 
              theme="primary" 
              type="submit" 
              size="large" 
              block
              :loading="loading"
            >
              登录
            </t-button>
          </t-form-item>
        </t-form>
        <div class="login-footer">
          <span>还没有账号？</span>
          <t-link @click="goToRegister">立即注册</t-link>
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
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const formData = reactive({
  username: '',
  password: ''
})

const handleSubmit = async () => {
  loading.value = true
  try {
    const response = await authApi.login(formData)
    authStore.setToken(response.data.token)
    authStore.setUser(response.data.user)
    
    MessagePlugin.success('登录成功')
    router.push('/dashboard')
  } catch (error) {
    MessagePlugin.error('登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-container {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  width: 100%;
  max-width: 400px;
}

.login-title {
  text-align: center;
  margin-bottom: 30px;
  font-size: 24px;
  color: #1f2937;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  color: #6b7280;
}
</style>