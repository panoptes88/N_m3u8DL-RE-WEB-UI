import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'
import { resolve } from 'path'

export default defineConfig({
  plugins: [
    vue(),
    // 按需引入 ant-design-vue：此前在 main.js 里 app.use(Antd) 会整包注册，
    // 无法被 tree-shaking，主 chunk 约 1.55MB。改为模板中用到才引入。
    Components({
      dts: false,
      resolvers: [
        AntDesignVueResolver({
          // v4 使用 CSS-in-JS，仅需 main.js 里的 reset.css，无需逐组件引样式
          importStyle: false
        })
      ]
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: 'dist',
    assetsDir: 'static'
  }
})
