import { ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'

interface LayoutProps {
  children: ReactNode
}

export function Layout({ children }: LayoutProps) {
  const location = useLocation()

  const navLinks = [
    { path: '/', label: 'Home' },
    { path: '/users', label: 'Users' },
    { path: '/health', label: 'Health' },
  ]

  return (
    <div>
      <nav style={{
        display: 'flex',
        gap: '1.5rem',
        padding: '1rem 0',
        marginBottom: '2rem',
        borderBottom: '1px solid #e2e8f0',
      }}>
        {navLinks.map(({ path, label }) => (
          <Link
            key={path}
            to={path}
            style={{
              color: location.pathname === path ? '#2563eb' : '#64748b',
              textDecoration: 'none',
              fontWeight: location.pathname === path ? 600 : 400,
            }}
          >
            {label}
          </Link>
        ))}
      </nav>
      <main>{children}</main>
    </div>
  )
}
