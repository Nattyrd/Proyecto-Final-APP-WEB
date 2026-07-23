import { useState, useEffect, useCallback } from 'react';
import { productsAPI, receiptsAPI, usersAPI } from '../api/api';
import CreateProductModal  from '../components/CreateProductModal';
import EditProductModal    from '../components/EditProductModal';
import ReceiptDetailModal  from '../components/ReceiptDetailModal';

/* ── Utilidades ─────────────────────────────────── */
const fmt      = (v) => `$${parseFloat(v ?? 0).toFixed(2)}`;
const fmtDate  = (iso) =>
  iso ? new Date(iso).toLocaleString('es-EC', { dateStyle: 'medium', timeStyle: 'short' }) : '—';

const TABS = [
  { id: 'products', label: 'Productos' },
  { id: 'receipts', label: 'Recibos'   },
  { id: 'users',    label: 'Usuarios'  },
];

/* ═══════════════════════════════════════════════════
   Panel de Administración — Pestañas:
   Productos · Recibos · Usuarios
   ═══════════════════════════════════════════════════ */
export default function AdminPage() {
  const [activeTab, setActiveTab] = useState('products');

  return (
    <div className="page-wrapper">
      {/* Page header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-corporate-900">Panel de Administración</h1>
        <p className="text-sm text-corporate-500 mt-0.5">ShopCorp — Grupo 5</p>
      </div>

      {/* Tab bar */}
      <div className="flex gap-1 border-b border-corporate-200 mb-6">
        {TABS.map((t) => (
          <button
            key={t.id}
            id={`admin-tab-${t.id}`}
            onClick={() => setActiveTab(t.id)}
            className={`px-5 py-2.5 text-sm font-medium rounded-t-md transition-colors duration-150 ${
              activeTab === t.id
                ? 'bg-white border border-b-white border-corporate-200 text-primary-700 -mb-px'
                : 'text-corporate-500 hover:text-corporate-800 hover:bg-corporate-50'
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      {/* Tab content */}
      {activeTab === 'products' && <ProductsTab />}
      {activeTab === 'receipts' && <ReceiptsTab />}
      {activeTab === 'users'    && <UsersTab    />}
    </div>
  );
}

/* ══════════════════════════════════════════════════
   TAB: Productos
   ══════════════════════════════════════════════════ */
function ProductsTab() {
  const [products,   setProducts]   = useState([]);
  const [page,       setPage]       = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const [loading,    setLoading]    = useState(false);
  const [error,      setError]      = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  const [showCreate, setShowCreate] = useState(false);
  const [editTarget, setEditTarget] = useState(null);  // producto a editar
  const [deleting,   setDeleting]   = useState(null);

  const PAGE_SIZE = 15;

  const fetchProducts = useCallback(async (pg) => {
    setLoading(true);
    setError('');
    try {
      const { data } = await productsAPI.getAll(pg, PAGE_SIZE);
      setProducts(data.data ?? []);
      setPage(data.page);
      setTotalPages(data.totalPages);
      setTotalItems(data.totalItems);
    } catch {
      setError('No se pudieron cargar los productos.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchProducts(page); }, [page, fetchProducts]);

  const flash = (msg) => {
    setSuccessMsg(msg);
    setTimeout(() => setSuccessMsg(''), 4000);
  };

  const handleCreated = (p) => { flash(`Producto "${p.name}" creado.`); fetchProducts(1); };
  const handleUpdated = (p) => {
    setProducts((prev) => prev.map((x) => (x.id === p.id ? p : x)));
    flash(`Producto "${p.name}" actualizado.`);
  };

  const handleDelete = async (product) => {
    if (!window.confirm(`¿Eliminar "${product.name}"? Esta acción no se puede deshacer.`)) return;
    setDeleting(product.id);
    setError('');
    try {
      await productsAPI.remove(product.id);
      setProducts((prev) => prev.filter((p) => p.id !== product.id));
      setTotalItems((n) => n - 1);
      flash(`Producto "${product.name}" eliminado.`);
    } catch (err) {
      setError(err.response?.data?.error ?? 'Error al eliminar el producto.');
    } finally {
      setDeleting(null);
    }
  };

  return (
    <>
      {/* Sub-header */}
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-corporate-500">{totalItems} registros</p>
        <button
          id="admin-new-product-btn"
          onClick={() => setShowCreate(true)}
          className="btn btn-primary"
        >
          + Nuevo producto
        </button>
      </div>

      {successMsg && <div className="alert-success mb-4 animate-fade-in">{successMsg}</div>}
      {error      && <div className="alert-error  mb-4">{error}</div>}

      <div className="card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="table-base">
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                <th>Descripción</th>
                <th>Precio</th>
                <th>Stock</th>
                <th className="text-right">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {loading && (
                <tr><td colSpan={6} className="text-center py-10 text-corporate-400">Cargando…</td></tr>
              )}
              {!loading && products.length === 0 && (
                <tr><td colSpan={6} className="text-center py-10 text-corporate-400">No hay productos.</td></tr>
              )}
              {!loading && products.map((p) => (
                <tr key={p.id} className="animate-fade-in">
                  <td className="text-corporate-400 text-xs">#{p.id}</td>
                  <td className="font-medium text-corporate-900 max-w-[180px] truncate">{p.name}</td>
                  <td className="text-corporate-500 max-w-[220px] truncate text-xs">
                    {p.description || <span className="italic text-corporate-300">Sin descripción</span>}
                  </td>
                  <td className="font-semibold text-primary-700">{fmt(p.price)}</td>
                  <td>
                    <span className={`badge ${p.stock > 0 ? 'badge-success' : 'bg-red-100 text-red-700'}`}>
                      {p.stock}
                    </span>
                  </td>
                  <td className="text-right">
                    <div className="flex justify-end gap-2">
                      <button
                        id={`admin-edit-${p.id}`}
                        onClick={() => setEditTarget(p)}
                        className="btn btn-secondary btn-sm"
                      >
                        Editar
                      </button>
                      <button
                        id={`admin-delete-${p.id}`}
                        onClick={() => handleDelete(p)}
                        disabled={deleting === p.id}
                        className="btn btn-danger btn-sm"
                      >
                        {deleting === p.id ? '…' : 'Eliminar'}
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-corporate-100 bg-corporate-50">
            <span className="text-xs text-corporate-400">Página {page} de {totalPages}</span>
            <div className="flex gap-2">
              <button id="admin-prev-btn" disabled={page <= 1} onClick={() => setPage((p) => p - 1)} className="btn btn-secondary btn-sm">← Anterior</button>
              <button id="admin-next-btn" disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)} className="btn btn-secondary btn-sm">Siguiente →</button>
            </div>
          </div>
        )}
      </div>

      {showCreate && (
        <CreateProductModal onClose={() => setShowCreate(false)} onCreated={handleCreated} />
      )}
      {editTarget && (
        <EditProductModal
          product={editTarget}
          onClose={() => setEditTarget(null)}
          onUpdated={handleUpdated}
        />
      )}
    </>
  );
}

/* ══════════════════════════════════════════════════
   TAB: Recibos
   ══════════════════════════════════════════════════ */
function ReceiptsTab() {
  const [receipts,   setReceipts]   = useState([]);
  const [loading,    setLoading]    = useState(false);
  const [error,      setError]      = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  const [detailId,   setDetailId]   = useState(null);
  const [deleting,   setDeleting]   = useState(null);

  const fetchReceipts = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const { data } = await receiptsAPI.getAll();
      setReceipts(Array.isArray(data) ? data : (data.data ?? []));
    } catch {
      setError('No se pudieron cargar los recibos.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchReceipts(); }, [fetchReceipts]);

  const flash = (msg) => { setSuccessMsg(msg); setTimeout(() => setSuccessMsg(''), 4000); };

  const handleDelete = async (r) => {
    if (!window.confirm(`¿Eliminar el recibo #${r.id}? Esta acción no se puede deshacer.`)) return;
    setDeleting(r.id);
    setError('');
    try {
      await receiptsAPI.remove(r.id);
      setReceipts((prev) => prev.filter((x) => x.id !== r.id));
      flash(`Recibo #${r.id} eliminado.`);
    } catch (err) {
      setError(err.response?.data?.error ?? 'Error al eliminar el recibo.');
    } finally {
      setDeleting(null);
    }
  };

  return (
    <>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-corporate-500">{receipts.length} recibos</p>
        <button id="admin-refresh-receipts-btn" onClick={fetchReceipts} className="btn btn-secondary btn-sm">
          ↺ Actualizar
        </button>
      </div>

      {successMsg && <div className="alert-success mb-4 animate-fade-in">{successMsg}</div>}
      {error      && <div className="alert-error  mb-4">{error}</div>}

      <div className="card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="table-base">
            <thead>
              <tr>
                <th>ID</th>
                <th>Usuario</th>
                <th>Fecha</th>
                <th className="text-right">Total</th>
                <th className="text-right">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {loading && (
                <tr><td colSpan={5} className="text-center py-10 text-corporate-400">Cargando…</td></tr>
              )}
              {!loading && receipts.length === 0 && (
                <tr><td colSpan={5} className="text-center py-10 text-corporate-400">Sin recibos registrados.</td></tr>
              )}
              {!loading && receipts.map((r) => (
                <tr key={r.id} className="animate-fade-in">
                  <td className="text-corporate-400 text-xs font-mono">#{r.id}</td>
                  <td className="font-medium text-corporate-800">{r.username ?? r.userId}</td>
                  <td className="text-xs text-corporate-500">{fmtDate(r.createdAt)}</td>
                  <td className="text-right font-semibold text-primary-700">{fmt(r.total)}</td>
                  <td className="text-right">
                    <div className="flex justify-end gap-2">
                      <button
                        id={`receipt-detail-${r.id}`}
                        onClick={() => setDetailId(r.id)}
                        className="btn btn-secondary btn-sm"
                      >
                        Ver detalle
                      </button>
                      <button
                        id={`receipt-delete-${r.id}`}
                        onClick={() => handleDelete(r)}
                        disabled={deleting === r.id}
                        className="btn btn-danger btn-sm"
                      >
                        {deleting === r.id ? '…' : 'Eliminar'}
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {detailId && (
        <ReceiptDetailModal receiptId={detailId} onClose={() => setDetailId(null)} />
      )}
    </>
  );
}

/* ══════════════════════════════════════════════════
   TAB: Usuarios
   ══════════════════════════════════════════════════ */
function UsersTab() {
  const [users,      setUsers]      = useState([]);
  const [loading,    setLoading]    = useState(false);
  const [error,      setError]      = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  const [editTarget, setEditTarget] = useState(null);
  const [deleting,   setDeleting]   = useState(null);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const { data } = await usersAPI.getAll();
      setUsers(Array.isArray(data) ? data : (data.data ?? []));
    } catch {
      setError('No se pudieron cargar los usuarios.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchUsers(); }, [fetchUsers]);

  const flash = (msg) => { setSuccessMsg(msg); setTimeout(() => setSuccessMsg(''), 4000); };

  const handleDelete = async (u) => {
    if (!window.confirm(`¿Eliminar al usuario "${u.username}"? Esta acción no se puede deshacer.`)) return;
    setDeleting(u.id);
    setError('');
    try {
      await usersAPI.remove(u.id);
      setUsers((prev) => prev.filter((x) => x.id !== u.id));
      flash(`Usuario "${u.username}" eliminado.`);
    } catch (err) {
      setError(err.response?.data?.error ?? 'Error al eliminar el usuario.');
    } finally {
      setDeleting(null);
    }
  };

  const handleSaved = (updated) => {
    setUsers((prev) => prev.map((x) => (x.id === updated.id ? updated : x)));
    setEditTarget(null);
    flash(`Usuario "${updated.username}" actualizado.`);
  };

  return (
    <>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-corporate-500">{users.length} usuarios</p>
        <button id="admin-refresh-users-btn" onClick={fetchUsers} className="btn btn-secondary btn-sm">
          ↺ Actualizar
        </button>
      </div>

      {successMsg && <div className="alert-success mb-4 animate-fade-in">{successMsg}</div>}
      {error      && <div className="alert-error  mb-4">{error}</div>}

      <div className="card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="table-base">
            <thead>
              <tr>
                <th>ID</th>
                <th>Usuario</th>
                <th>Nombre</th>
                <th>Email</th>
                <th>Rol</th>
                <th className="text-right">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {loading && (
                <tr><td colSpan={6} className="text-center py-10 text-corporate-400">Cargando…</td></tr>
              )}
              {!loading && users.length === 0 && (
                <tr><td colSpan={6} className="text-center py-10 text-corporate-400">Sin usuarios.</td></tr>
              )}
              {!loading && users.map((u) => (
                <tr key={u.id} className="animate-fade-in">
                  <td className="text-corporate-400 text-xs">#{u.id}</td>
                  <td className="font-medium text-corporate-900">{u.username}</td>
                  <td className="text-corporate-600 text-sm">{u.firstName} {u.lastName}</td>
                  <td className="text-corporate-500 text-xs">{u.email}</td>
                  <td>
                    <span className={`badge ${u.role === 'ADMIN' ? 'badge-admin' : 'badge-client'}`}>
                      {u.role}
                    </span>
                  </td>
                  <td className="text-right">
                    <div className="flex justify-end gap-2">
                      <button
                        id={`user-edit-${u.id}`}
                        onClick={() => setEditTarget(u)}
                        className="btn btn-secondary btn-sm"
                      >
                        Editar
                      </button>
                      <button
                        id={`user-delete-${u.id}`}
                        onClick={() => handleDelete(u)}
                        disabled={deleting === u.id}
                        className="btn btn-danger btn-sm"
                      >
                        {deleting === u.id ? '…' : 'Eliminar'}
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {editTarget && (
        <EditUserModal
          user={editTarget}
          onClose={() => setEditTarget(null)}
          onSaved={handleSaved}
        />
      )}
    </>
  );
}

/* ── Modal edición de usuario (inline, solo para Admin) ── */
function EditUserModal({ user, onClose, onSaved }) {
  const [form,    setForm]    = useState({
    firstName: user.firstName ?? '',
    lastName:  user.lastName  ?? '',
    email:     user.email     ?? '',
    role:      user.role      ?? 'CLIENT',
  });
  const [apiErr,  setApiErr]  = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((p) => ({ ...p, [name]: value }));
    setApiErr('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const { data } = await usersAPI.update(user.id, form);
      onSaved(data);
    } catch (err) {
      setApiErr(err.response?.data?.error ?? 'Error al actualizar usuario.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-corporate-900/40 backdrop-blur-sm"
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div className="card animate-slide-in w-full max-w-md mx-4 p-6">
        <div className="flex items-center justify-between mb-5">
          <div>
            <h2 className="text-lg font-semibold text-corporate-900">Editar usuario</h2>
            <p className="text-xs text-corporate-400 mt-0.5">@{user.username} · ID #{user.id}</p>
          </div>
          <button onClick={onClose} className="btn btn-ghost btn-sm p-1 text-corporate-400 hover:text-corporate-700" aria-label="Cerrar">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {apiErr && <div className="alert-error mb-4">{apiErr}</div>}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label htmlFor="eu-firstName" className="label">Nombre</label>
              <input id="eu-firstName" name="firstName" value={form.firstName} onChange={handleChange} className="input" />
            </div>
            <div>
              <label htmlFor="eu-lastName" className="label">Apellido</label>
              <input id="eu-lastName" name="lastName" value={form.lastName} onChange={handleChange} className="input" />
            </div>
          </div>
          <div>
            <label htmlFor="eu-email" className="label">Email</label>
            <input id="eu-email" name="email" type="email" value={form.email} onChange={handleChange} className="input" />
          </div>
          <div>
            <label htmlFor="eu-role" className="label">Rol</label>
            <select id="eu-role" name="role" value={form.role} onChange={handleChange} className="input">
              <option value="CLIENT">CLIENT</option>
              <option value="ADMIN">ADMIN</option>
            </select>
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <button type="button" onClick={onClose} className="btn btn-secondary">Cancelar</button>
            <button id="eu-submit-btn" type="submit" disabled={loading} className="btn btn-primary">
              {loading ? 'Guardando…' : 'Guardar cambios'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
