// ============================================================
// 【学生自己编写的代码】图书增删改查的API调用函数
// ============================================================

import request from './request'

export const getBooksApi = (params) =>
  request.get('/books', { params })

export const getBookApi = (id) =>
  request.get(`/books/${id}`)

export const createBookApi = (data) =>
  request.post('/books', data)

export const updateBookApi = (id, data) =>
  request.put(`/books/${id}`, data)

export const deleteBookApi = (id) =>
  request.delete(`/books/${id}`)


