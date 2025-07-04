<template>
  <div class="dashboard">
    <t-layout>
      <!-- 侧边栏 -->
      <t-aside class="sidebar">
        <div class="sidebar-header">
          <h3>云服务器平台</h3>
        </div>
        <t-menu
          v-model:value="activeMenu"
          :default-value="['dashboard']"
          theme="light"
        >
          <t-menu-item value="dashboard">
            <template #icon><div>📊</div></template>
            仪表板
          </t-menu-item>
          <t-menu-item value="servers">
            <template #icon><div>🖥️</div></template>
            我的服务器
          </t-menu-item>
          <t-menu-item value="purchase">
            <template #icon><div>🛒</div></template>
            购买服务器
          </t-menu-item>
          <t-menu-item value="billing">
            <template #icon><div>💳</div></template>
            费用中心
          </t-menu-item>
          <t-menu-item value="profile">
            <template #icon><div>👤</div></template>
            个人资料
          </t-menu-item>
        </t-menu>
      </t-aside>

      <!-- 主内容区 -->
      <t-layout>
        <t-header class="header">
          <div class="header-content">
            <h2>控制台</h2>
            <div class="header-actions">
              <t-button variant="text" @click="logout">退出登录</t-button>
            </div>
          </div>
        </t-header>

        <t-content class="content">
          <!-- 统计卡片 -->
          <div class="stats-cards">
            <div class="stat-card">
              <div class="stat-icon">🖥️</div>
              <div class="stat-info">
                <h3>{{ stats.totalServers }}</h3>
                <p>服务器总数</p>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">✅</div>
              <div class="stat-info">
                <h3>{{ stats.runningServers }}</h3>
                <p>运行中</p>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">💰</div>
              <div class="stat-info">
                <h3>¥{{ stats.monthlySpend }}</h3>
                <p>本月消费</p>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">📈</div>
              <div class="stat-info">
                <h3>99.9%</h3>
                <p>可用性</p>
              </div>
            </div>
          </div>

          <!-- 最近活动 -->
          <div class="activity-section">
            <h3>最近活动</h3>
            <t-table
              :data="recentActivity"
              :columns="activityColumns"
              :pagination="false"
              size="small"
            />
          </div>
        </t-content>
      </t-layout>
    </t-layout>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const activeMenu = ref(['dashboard'])

const stats = reactive({
  totalServers: 8,
  runningServers: 6,
  monthlySpend: 2580,
  uptime: 99.9
})

const recentActivity = [
  {
    id: 1,
    action: '服务器启动',
    server: 'web-server-01',
    time: '2024-01-15 10:30:00',
    status: '成功'
  },
  {
    id: 2,
    action: '购买服务器',
    server: 'db-server-02',
    time: '2024-01-15 09:15:00',
    status: '成功'
  },
  {
    id: 3,
    action: '配置更新',
    server: 'api-server-01',
    time: '2024-01-14 16:45:00',
    status: '成功'
  }
]

const activityColumns = [
  {
    colKey: 'action',
    title: '操作',
    width: 120
  },
  {
    colKey: 'server',
    title: '服务器',
    width: 150
  },
  {
    colKey: 'time',
    title: '时间',
    width: 180
  },
  {
    colKey: 'status',
    title: '状态',
    width: 100
  }
]

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.dashboard {
  height: 100vh;
}

.sidebar {
  width: 250px;
  background: #f8f9ff;
  border-right: 1px solid #e7e7e7;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e7e7e7;
}

.sidebar-header h3 {
  margin: 0;
  color: #1f2937;
}

.header {
  background: white;
  border-bottom: 1px solid #e7e7e7;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.content {
  padding: 20px;
  background: #f8f9ff;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stat-icon {
  font-size: 32px;
  margin-right: 15px;
}

.stat-info h3 {
  margin: 0 0 5px 0;
  font-size: 24px;
  color: #1f2937;
}

.stat-info p {
  margin: 0;
  color: #6b7280;
}

.activity-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.activity-section h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #1f2937;
}
</style>