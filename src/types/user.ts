import type { User } from './auth'
export type { User }

export interface CreateUserRequest {
  name: string
  email: string
  password: string
  role: User['role']
}

export interface UpdateUserRequest {
  name: string
  email: string
  role: User['role']
  is_active: boolean
}

export interface ResetPasswordRequest {
  password: string
}
