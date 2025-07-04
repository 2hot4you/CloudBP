<template>
  <div class="home">
    <t-layout>
      <!-- 头部导航 -->
      <t-header class="header">
        <t-head-menu value="home" layout="horizontal">
          <template #logo>
            <div class="logo">
              <span class="logo-text">云服务器销售平台</span>
            </div>
          </template>
          <t-menu-item value="home">首页</t-menu-item>
          <t-menu-item value="products">产品</t-menu-item>
          <t-menu-item value="pricing">价格</t-menu-item>
          <t-menu-item value="about">关于我们</t-menu-item>
          <template #operations>
            <t-button theme="primary" @click="goToLogin">登录</t-button>
            <t-button variant="outline" @click="goToRegister">注册</t-button>
          </template>
        </t-head-menu>
      </t-header>

      <!-- 主要内容 -->
      <t-content class="content">
        <!-- 首页横幅 -->
        <div class="hero-section">
          <div class="hero-content">
            <h1 class="hero-title">专业的云服务器销售平台</h1>
            <p class="hero-subtitle">多厂商产品对比，一站式云服务器采购解决方案</p>
            <div class="hero-buttons">
              <t-button theme="primary" size="large" @click="goToPurchase">立即购买</t-button>
              <t-button variant="outline" size="large">了解更多</t-button>
            </div>
          </div>
        </div>

        <!-- 产品特色 -->
        <div class="features-section">
          <h2 class="section-title">为什么选择我们</h2>
          <div class="features-grid">
            <div class="feature-card">
              <div class="feature-icon">🚀</div>
              <h3>高性能</h3>
              <p>采用最新硬件配置，提供卓越性能体验</p>
            </div>
            <div class="feature-card">
              <div class="feature-icon">🔒</div>
              <h3>安全可靠</h3>
              <p>多重安全防护，保障数据安全</p>
            </div>
            <div class="feature-card">
              <div class="feature-icon">💰</div>
              <h3>价格优势</h3>
              <p>多厂商价格对比，找到最优性价比</p>
            </div>
            <div class="feature-card">
              <div class="feature-icon">🛠️</div>
              <h3>技术支持</h3>
              <p>7×24小时技术支持，快速响应</p>
            </div>
          </div>
        </div>
      </t-content>
    </t-layout>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goToLogin = () => {
  router.push('/login')
}

const goToRegister = () => {
  router.push('/register')
}

const goToPurchase = () => {
  router.push('/purchase')
}
</script>

<style scoped>
.home {
  min-height: 100vh;
}

.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  color: #0052d9;
}

.content {
  padding: 0;
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 120px 20px;
  text-align: center;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
}

.hero-title {
  font-size: 48px;
  margin-bottom: 20px;
  font-weight: bold;
}

.hero-subtitle {
  font-size: 20px;
  margin-bottom: 40px;
  opacity: 0.9;
}

.hero-buttons {
  display: flex;
  gap: 20px;
  justify-content: center;
}

.features-section {
  padding: 80px 20px;
  background: #f8f9ff;
}

.section-title {
  text-align: center;
  font-size: 32px;
  margin-bottom: 60px;
  color: #1f2937;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
  max-width: 1200px;
  margin: 0 auto;
}

.feature-card {
  background: white;
  padding: 40px 30px;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
}

.feature-icon {
  font-size: 48px;
  margin-bottom: 20px;
}

.feature-card h3 {
  font-size: 24px;
  margin-bottom: 15px;
  color: #1f2937;
}

.feature-card p {
  color: #6b7280;
  line-height: 1.6;
}
</style>