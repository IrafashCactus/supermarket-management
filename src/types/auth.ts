export interface User {
  id: string
  name: string
  email: string
  role: 'admin' | 'manager' | 'accountant' | 'storekeeper' | 'cashier'
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  user: User
}

export type UserRole = User['role']
