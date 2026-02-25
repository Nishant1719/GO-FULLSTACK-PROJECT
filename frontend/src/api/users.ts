export interface User {
  id: string
  username: string
  email: string
  first_name?: string
  last_name?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface UserListResponse {
  data: User[]
  pagination: { limit: number; offset: number; total: number }
}

export interface CreateUserPayload {
  username: string
  email: string
  password: string
  first_name?: string
  last_name?: string
}

export interface UpdateUserPayload {
  username?: string
  email?: string
  first_name?: string
  last_name?: string
  is_active?: boolean
}

import { api } from './client'

export const usersApi = {
  list: (limit = 10, offset = 0) =>
    api.get<UserListResponse>(`/api/v1/users?limit=${limit}&offset=${offset}`),
  get: (id: string) => api.get<User>(`/api/v1/users/${id}`),
  create: (payload: CreateUserPayload) =>
    api.post<User>('/api/v1/users', payload),
  update: (id: string, payload: UpdateUserPayload) =>
    api.patch<User>(`/api/v1/users/${id}`, payload),
  delete: (id: string) => api.delete(`/api/v1/users/${id}`),
}
