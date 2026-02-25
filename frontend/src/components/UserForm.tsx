import { useState } from 'react'
import type { CreateUserPayload, UpdateUserPayload, User } from '../api/users'

interface UserFormProps {
  user?: User
  /** Pre-fill form for copy flow (create mode with initial values) */
  copyFrom?: User
  onSave: (payload: CreateUserPayload | UpdateUserPayload) => Promise<void>
  onCancel: () => void
}

export function UserForm({ user, copyFrom, onSave, onCancel }: UserFormProps) {
  const isEdit = !!user
  const source = user ?? copyFrom
  const [username, setUsername] = useState(source?.username ?? '')
  const [email, setEmail] = useState(source?.email ?? '')
  const [password, setPassword] = useState('')
  const [firstName, setFirstName] = useState(source?.first_name ?? '')
  const [lastName, setLastName] = useState(source?.last_name ?? '')
  const [isActive, setIsActive] = useState(source?.is_active ?? true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setSaving(true)
    try {
      if (isEdit) {
        await onSave({ username, email, first_name: firstName || undefined, last_name: lastName || undefined, is_active: isActive })
      } else {
        if (!password || password.length < 8) {
          throw new Error('Password must be at least 8 characters')
        }
        await onSave({ username, email, password, first_name: firstName || undefined, last_name: lastName || undefined })
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save')
    } finally {
      setSaving(false)
    }
  }

  const inputStyle = { width: '100%', padding: '0.5rem', marginBottom: '0.75rem', borderRadius: '4px', border: '1px solid #e2e8f0' }
  const labelStyle = { display: 'block', marginBottom: '0.25rem', fontSize: '0.875rem', color: '#64748b' }

  return (
    <form onSubmit={handleSubmit} style={{ maxWidth: '400px' }}>
      {error && <p style={{ color: '#dc2626', marginBottom: '1rem' }}>{error}</p>}
      <div>
        <label style={labelStyle}>Username</label>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          minLength={3}
          style={inputStyle}
        />
      </div>
      <div>
        <label style={labelStyle}>Email</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          style={inputStyle}
        />
      </div>
      {!isEdit && (
        <div>
          <label style={labelStyle}>Password</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
            style={inputStyle}
          />
        </div>
      )}
      <div>
        <label style={labelStyle}>First name</label>
        <input type="text" value={firstName} onChange={(e) => setFirstName(e.target.value)} style={inputStyle} />
      </div>
      <div>
        <label style={labelStyle}>Last name</label>
        <input type="text" value={lastName} onChange={(e) => setLastName(e.target.value)} style={inputStyle} />
      </div>
      {isEdit && (
        <div style={{ marginBottom: '1rem' }}>
          <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
            <input type="checkbox" checked={isActive} onChange={(e) => setIsActive(e.target.checked)} />
            Active
          </label>
        </div>
      )}
      <div style={{ display: 'flex', gap: '0.5rem', marginTop: '1rem' }}>
        <button type="submit" disabled={saving}>
          {saving ? 'Saving...' : isEdit ? 'Save' : copyFrom ? 'Copy User' : 'Add User'}
        </button>
        <button type="button" onClick={onCancel}>
          Cancel
        </button>
      </div>
    </form>
  )
}
