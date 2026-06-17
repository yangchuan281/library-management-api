// ============================================================
// 【学生自己编写的代码】登录/注册/重置密码的API调用函数
// ============================================================

import request from './request'

export const sendCodeApi = (email, type) =>
  request.post('/auth/verification-codes', { email, type })

export const registerApi = (data) =>
  request.post('/auth/register', data)

export const loginApi = (data) =>
  request.post('/auth/login', data)

export const resetPasswordApi = (data) =>
  request.put('/auth/password', data)

export const getProfileApi = () =>
  request.get('/users/me')


