import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import 'ant-design-vue/dist/reset.css'
import './assets/styles.css'
import 'xgplayer/dist/index.min.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
// ant-design-vue 组件由 unplugin-vue-components 按需自动引入（见 vite.config.ts），
// 这里不再 app.use(Antd) 整包注册

app.mount('#app')
