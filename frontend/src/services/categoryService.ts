import api from './api'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/types/category'
import type { ListResponse, ApiResponse } from '@/types/api'

export const categoryService = {
  async getCategories(page = 1, perPage = 10): Promise<ListResponse<Category>> {
    const response = await api.get<ApiResponse<ListResponse<Category>>>('/categories', {
      params: { page, per_page: perPage }
    })
    return response.data.data
  },

  async getAllCategories(): Promise<Category[]> {
    const categories: Category[] = []
    let page = 1
    let totalPages = 1
    do {
      const response = await this.getCategories(page, 100)
      categories.push(...response.data)
      totalPages = response.meta.total_pages
      page += 1
    } while (page <= totalPages)
    return categories
  },

  async createCategory(data: CreateCategoryRequest): Promise<Category> {
    const response = await api.post<ApiResponse<Category>>('/categories', data)
    return response.data.data
  },

  async updateCategory(id: string, data: UpdateCategoryRequest): Promise<Category> {
    const response = await api.put<ApiResponse<Category>>(`/categories/${id}`, data)
    return response.data.data
  },

  async deleteCategory(id: string): Promise<void> {
    await api.delete(`/categories/${id}`)
  }
}
