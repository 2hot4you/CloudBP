import api from './request'
import type { LoginRequest, RegisterRequest, AuthResponse } from '@/types/user'

export const authApi = {
  // 登录
  login: (data: LoginRequest) => 
    api.post<AuthResponse>('/auth/login', data),
  
  // 注册
  register: (data: RegisterRequest) => 
    api.post<AuthResponse>('/auth/register', data),
  
  // 登出
  logout: () => 
    api.post('/auth/logout'),
  
  // 刷新令牌
  refresh: () => 
    api.post('/auth/refresh'),
  
  // 获取用户信息
  getUserInfo: () => 
    api.get('/user/profile')
}