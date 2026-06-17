// ============================================================
// 【学生自己编写的代码】Axios配置 + 请求/响应拦截器
// 作用：配置HTTP请求的基地址、超时时间，统一处理Token和错误
// ============================================================

import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建Axios实例（所有API请求的基配置）
const request = axios.create({
  baseURL: '/api',              // 请求自动加上/api前缀，Vite代理转发到后端8000端口
  timeout: 15000                // 超时时间15秒
})

// 【请求拦截器】每次发请求前，自动从localStorage取出JWT Token附加到请求头
request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// 【响应拦截器】收到响应后统一处理错误
request.interceptors.response.use(
  response => {
    const res = response.data
    // GoFrame框架返回 code=0 表示业务成功
    if (res.code !== 0 && res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message))
    }
    return res
  },
  error => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        // Token过期或无效，清除本地Token并跳回登录页
        localStorage.removeItem('token')
        window.location.hash = '#/login'
        ElMessage.error('登录已过期，请重新登录')
      } else {
        ElMessage.error(data?.message || `请求错误 (${status})`)
      }
    } else {
      ElMessage.error('网络错误，请检查后端是否启动')
    }
    return Promise.reject(error)
  }
)

export default request
