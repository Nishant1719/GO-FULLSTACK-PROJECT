import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { usersApi, User, CreateUserPayload, UpdateUserPayload } from '../api/users'
import { UserForm } from '../components/UserForm'

export function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [total, setTotal] = useState(0)
  const [offset, setOffset] = useState(0)
  const [showAddForm, setShowAddForm] = useState(false)
  const [copyFromUser, setCopyFromUser] = useState<User | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)
  const limit = 10

  const fetchUsers = () => {
    setError(null)
    usersApi
      .list(limit, offset)
      .then((res) => {
        setUsers(res.data)
        setTotal(res.pagination.total)
      })
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    setLoading(true)
    fetchUsers()
  }, [offset])

  const handleAdd = async (payload: CreateUserPayload | UpdateUserPayload) => {
    const p = payload as CreateUserPayload
    if (!p.password) throw new Error('Password required')
    await usersApi.create(p)
    setShowAddForm(false)
    setCopyFromUser(null)
    fetchUsers()
  }

  const handleEdit = async (id: string, payload: Parameters<typeof usersApi.update>[1]) => {
    await usersApi.update(id, payload)
    setEditingId(null)
    fetchUsers()
  }

  const handleDelete = async (id: string, username: string) => {
    if (!confirm(`Delete user "${username}"?`)) return
    await usersApi.delete(id)
    fetchUsers()
  }

  if (loading) return <p>Loading users...</p>
  if (error) {
    return (
      <div>
        <p style={{ color: '#dc2626', marginBottom: '1rem' }}>Error: {error}</p>
        <button onClick={() => { setLoading(true); fetchUsers() }}>Retry</button>
      </div>
    )
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem', flexWrap: 'wrap', gap: '0.5rem' }}>
        <h1 style={{ margin: 0 }}>Users</h1>
        <button onClick={() => { setCopyFromUser(null); setShowAddForm(true) }} disabled={showAddForm && !copyFromUser}>Add User</button>
      </div>

      {(showAddForm || copyFromUser) && (
        <div style={{ padding: '1rem', marginBottom: '1rem', border: '1px solid #e2e8f0', borderRadius: '8px', background: '#f8fafc' }}>
          <h2 style={{ marginBottom: '1rem', fontSize: '1.125rem' }}>{copyFromUser ? `Copy user: ${copyFromUser.username}` : 'New User'}</h2>
          <UserForm
            copyFrom={copyFromUser ?? undefined}
            onSave={handleAdd}
            onCancel={() => { setShowAddForm(false); setCopyFromUser(null) }}
          />
        </div>
      )}

      {total > 0 && (
        <div style={{ marginBottom: '1rem' }}>
          Showing {offset + 1}–{Math.min(offset + limit, total)} of {total}
        </div>
      )}

      {total === 0 ? (
        <p style={{ color: '#64748b', padding: '2rem', textAlign: 'center' }}>No users available</p>
      ) : (
      <ul style={{ listStyle: 'none' }}>
        {users.map((u) => (
          <li
            key={u.id}
            style={{
              padding: '0.75rem',
              marginBottom: '0.5rem',
              border: '1px solid #e2e8f0',
              borderRadius: '6px',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              flexWrap: 'wrap',
              gap: '0.5rem',
            }}
          >
            {editingId === u.id ? (
              <div style={{ flex: '1 1 100%' }}>
                <UserForm
                  user={u}
                  onSave={(payload) => handleEdit(u.id, payload)}
                  onCancel={() => setEditingId(null)}
                />
              </div>
            ) : (
              <>
                <span>
                  <Link to={`/users/${u.id}`} style={{ color: '#2563eb', textDecoration: 'none', fontWeight: 500 }}>
                    {u.username}
                  </Link>
                  <span style={{ color: '#64748b', marginLeft: '0.5rem' }}>— {u.email}</span>
                </span>
                <div style={{ display: 'flex', gap: '0.5rem' }}>
                  <button onClick={() => setEditingId(u.id)}>Edit</button>
                  <button onClick={() => setCopyFromUser(u)}>Copy</button>
                  <button onClick={() => handleDelete(u.id, u.username)} style={{ color: '#dc2626' }}>Delete</button>
                </div>
              </>
            )}
          </li>
        ))}
      </ul>
      )}

      <div style={{ display: 'flex', gap: '0.5rem', marginTop: '1rem' }}>
        <button
          disabled={offset === 0}
          onClick={() => setOffset((o) => Math.max(0, o - limit))}
        >
          Previous
        </button>
        <button
          disabled={offset + limit >= total}
          onClick={() => setOffset((o) => o + limit)}
        >
          Next
        </button>
      </div>
    </div>
  )
}
