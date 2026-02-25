import { useEffect, useState } from 'react'
import { healthApi } from '../api/health'

interface HealthState {
  status: string
  bff?: string
  goDomain?: string
}

export function HealthPage() {
  const [health, setHealth] = useState<HealthState | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    healthApi
      .check()
      .then((res) => setHealth(res as HealthState))
      .catch((e) => setError(e instanceof Error ? e.message : 'Check failed'))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <p>Checking health...</p>
  if (error) return <p style={{ color: '#dc2626' }}>Error: {error}</p>
  if (!health) return null

  const isHealthy = health.status === 'healthy'

  return (
    <div>
      <h1 style={{ marginBottom: '1rem' }}>Health Check</h1>
      <p style={{ marginBottom: '1rem', color: isHealthy ? '#16a34a' : '#dc2626' }}>
        Overall: <strong>{health.status}</strong>
      </p>
      <ul style={{ listStyle: 'none' }}>
        {health.bff != null && (
          <li style={{ marginBottom: '0.5rem' }}>BFF: {health.bff}</li>
        )}
        {health.goDomain != null && (
          <li style={{ marginBottom: '0.5rem' }}>Go Domain: {health.goDomain}</li>
        )}
      </ul>
    </div>
  )
}
