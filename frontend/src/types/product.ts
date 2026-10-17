export interface Product {
  id: string
  name: string
  barcode: string
  sku: string
  category_id: string
  category_name: string
  unit_price: number
  cost_price: number
  tax_rate: number
  unit_of_measure: string
  reorder_level: number
  image_url: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateProductRequest {
  name: string
  barcode: string
  sku: string
  category_id: string
  unit_price: number
  cost_price: number
  tax_rate: number
  unit_of_measure: string
  reorder_level: number
  image_url: string
}

export interface UpdateProductRequest extends CreateProductRequest {
  is_active: boolean
}

export interface ProductListParams {
  page?: number
  per_page?: number
  search?: string
  category_id?: string
  is_active?: boolean | null
}
