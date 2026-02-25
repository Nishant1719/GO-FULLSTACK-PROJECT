import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { usersApi, User } from '../api/users'
import { UserForm } from '../components/UserForm'

export function UserDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editing, setEditing] = useState(false)

  const fetchUser = () => {
    if (!id) return
    setError(null)
    usersApi
      .get(id)
      .then(setUser)
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    setLoading(true)
    fetchUser()
  }, [id])

  const handleEdit = async (payload: Parameters<typeof usersApi.update>[1]) => {
    if (!id) return
    await usersApi.update(id, payload)
    setEditing(false)
    fetchUser()
  }

  const handleDelete = async () => {
    if (!id || !user) return
    if (!confirm(`Delete user "${user.username}"?`)) return
    await usersApi.delete(id)
    navigate('/users')
  }

  if (loading) return <p>Loading...</p>
  if (error) {
    return (
      <div>
        <p style={{ color: '#dc2626', marginBottom: '1rem' }}>Error: {error}</p>
        <button onClick={() => { setLoading(true); fetchUser() }}>Retry</button>
      </div>
    )
  }
  if (!user) return <p>User not found</p>

  if (editing) {
    return (
      <div>
        <h2 style={{ marginBottom: '1rem' }}>Edit User</h2>
        <UserForm user={user} onSave={handleEdit} onCancel={() => setEditing(false)} />
      </div>
    )
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1rem', flexWrap: 'wrap', gap: '0.5rem' }}>
        <h1 style={{ margin: 0 }}>{user.username}</h1>
        <div style={{ display: 'flex', gap: '0.5rem' }}>
          <button onClick={() => setEditing(true)}>Edit</button>
          <button onClick={handleDelete} style={{ color: '#dc2626' }}>Remove</button>
        </div>
      </div>
      <dl style={{ display: 'grid', gap: '0.5rem' }}>
        <div>
          <dt style={{ color: '#64748b', fontSize: '0.875rem' }}>Email</dt>
          <dd>{user.email}</dd>
        </div>
        {user.first_name && (
          <div>
            <dt style={{ color: '#64748b', fontSize: '0.875rem' }}>First name</dt>
            <dd>{user.first_name}</dd>
          </div>
        )}
        {user.last_name && (
          <div>
            <dt style={{ color: '#64748b', fontSize: '0.875rem' }}>Last name</dt>
            <dd>{user.last_name}</dd>
          </div>
        )}
        <div>
          <dt style={{ color: '#64748b', fontSize: '0.875rem' }}>Active</dt>
          <dd>{user.is_active ? 'Yes' : 'No'}</dd>
        </div>
        <div>
          <dt style={{ color: '#64748b', fontSize: '0.875rem' }}>Created</dt>
          <dd>{new Date(user.created_at).toLocaleString()}</dd>
        </div>
      </dl>
    </div>
  )
}
