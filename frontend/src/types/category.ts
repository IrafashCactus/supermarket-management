export interface Category {
  id: string
  name: string
  description: string
  product_count: number
  created_at: string
  updated_at: string
}

export interface CreateCategoryRequest {
  name: string
  description: string
}

export interface UpdateCategoryRequest {
  name: string
  description: string
}
