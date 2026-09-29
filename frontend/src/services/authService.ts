import api from './api'
import type { LoginRequest, LoginResponse, User } from '@/types/auth'
import type { ApiResponse } from '@/types/api'

export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', credentials)
    return response.data.data
  },

  async refreshToken(): Promise<{ access_token: string }> {
    const response = await api.post<ApiResponse<{ access_token: string }>>('/auth/refresh')
    return response.data.data
  },

  async logout(): Promise<void> {
    await api.post('/auth/logout')
  },

  async getMe(): Promise<User> {
    const response = await api.get<ApiResponse<User>>('/auth/me')
    return response.data.data
  }
}
