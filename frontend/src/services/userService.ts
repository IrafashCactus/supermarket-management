import api from './api'
import type { User, CreateUserRequest, UpdateUserRequest, ResetPasswordRequest } from '@/types/user'
import type { ListResponse, ApiResponse } from '@/types/api'

export const userService = {
  async getUsers(page = 1, perPage = 10): Promise<ListResponse<User>> {
    const response = await api.get<ApiResponse<ListResponse<User>>>('/users', {
      params: { page, per_page: perPage }
    })
    return response.data.data
  },

  async getUser(id: string): Promise<User> {
    const response = await api.get<ApiResponse<User>>(`/users/${id}`)
    return response.data.data
  },

  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await api.post<ApiResponse<User>>('/users', data)
    return response.data.data
  },

  async updateUser(id: string, data: UpdateUserRequest): Promise<User> {
    const response = await api.put<ApiResponse<User>>(`/users/${id}`, data)
    return response.data.data
  },

  async resetPassword(id: string, data: ResetPasswordRequest): Promise<void> {
    await api.put(`/users/${id}/password`, data)
  }
}
