// ============================================================
// 【学生自己编写的代码】Vite构建工具配置
// 作用：配置前端开发服务器和API代理
// ============================================================
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,  // 前端运行在3000端口
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})


