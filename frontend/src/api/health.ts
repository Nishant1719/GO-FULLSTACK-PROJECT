interface HealthStatus {
  status: 'healthy' | 'unhealthy'
  bff?: { status: string }
  goDomain?: { status: string }
}

import { api } from './client'

export const healthApi = {
  check: () => api.get<HealthStatus>('/health'),
}
