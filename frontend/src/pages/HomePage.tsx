import { Link } from 'react-router-dom'

export function HomePage() {
  return (
    <div>
      <h1 style={{ marginBottom: '1rem' }}>GO Fullstack Project</h1>
      <p style={{ color: '#64748b', marginBottom: '1.5rem' }}>
        Layered architecture: React → Node.js (BFF) → Go → PostgreSQL
      </p>
      <ul style={{ listStyle: 'none' }}>
        <li style={{ marginBottom: '0.5rem' }}>
          <Link to="/users" style={{ color: '#2563eb' }}>Browse Users</Link>
        </li>
        <li style={{ marginBottom: '0.5rem' }}>
          <Link to="/health" style={{ color: '#2563eb' }}>Health Check</Link>
        </li>
      </ul>
    </div>
  )
}
