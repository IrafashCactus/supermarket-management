import api from './api'
import type { Product, CreateProductRequest, UpdateProductRequest, ProductListParams } from '@/types/product'
import type { ListResponse, ApiResponse } from '@/types/api'

export const productService = {
  async getProducts(params: ProductListParams = {}): Promise<ListResponse<Product>> {
    const queryParams: Record<string, string | number | boolean> = {
      page: params.page ?? 1,
      per_page: params.per_page ?? 10
    }
    if (params.search) queryParams.search = params.search
    if (params.category_id) queryParams.category_id = params.category_id
    if (params.is_active !== null && params.is_active !== undefined) queryParams.is_active = params.is_active

    const response = await api.get<ApiResponse<ListResponse<Product>>>('/products', {
      params: queryParams
    })
    return response.data.data
  },

  async getProduct(id: string): Promise<Product> {
    const response = await api.get<ApiResponse<Product>>(`/products/${id}`)
    return response.data.data
  },

  async getProductByBarcode(barcode: string): Promise<Product> {
    const response = await api.get<ApiResponse<Product>>(`/products/barcode/${barcode}`)
    return response.data.data
  },

  async createProduct(data: CreateProductRequest): Promise<Product> {
    const response = await api.post<ApiResponse<Product>>('/products', data)
    return response.data.data
  },

  async updateProduct(id: string, data: UpdateProductRequest): Promise<Product> {
    const response = await api.put<ApiResponse<Product>>(`/products/${id}`, data)
    return response.data.data
  }
}
