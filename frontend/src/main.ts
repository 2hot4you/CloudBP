import { createApp } from 'vue'
import { createPinia } from 'pinia'
import TDesign from 'tdesign-vue-next'
import App from './App.vue'
import router from './router'

// 引入TDesign样式
import 'tdesign-vue-next/es/style/index.css'

// 引入全局样式
import '@/assets/css/global.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(TDesign)

app.mount('#app')