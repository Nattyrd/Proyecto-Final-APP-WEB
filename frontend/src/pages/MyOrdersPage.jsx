import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { receiptsAPI } from '../api/api';
import ReceiptDetailModal from '../components/ReceiptDetailModal';

const fmt     = (v)   => `$${parseFloat(v ?? 0).toFixed(2)}`;
const fmtDate = (iso) =>
  iso ? new Date(iso).toLocaleString('es-EC', { dateStyle: 'medium', timeStyle: 'short' }) : '—';

/**
 * Vista "Mis Compras" — solo para rol CLIENT.
 * Lista todos los recibos del usuario autenticado.
 */
export default function MyOrdersPage() {
  const { user } = useAuth();
  const [receipts, setReceipts] = useState([]);
  const [loading,  setLoading]  = useState(false);
  const [error,    setError]    = useState('');
  const [detailId, setDetailId] = useState(null);

  const fetchOrders = useCallback(async () => {
    if (!user?.id) return;
    setLoading(true);
    setError('');
    try {
      const { data } = await receiptsAPI.getByUser(user.id);
      setReceipts(Array.isArray(data) ? data : (data.data ?? []));
    } catch {
      setError('No se pudo cargar tu historial de compras.');
    } finally {
      setLoading(false);
    }
  }, [user?.id]);

  useEffect(() => { fetchOrders(); }, [fetchOrders]);

  const totalGastado = receipts.reduce((acc, r) => acc + parseFloat(r.total ?? 0), 0);

  return (
    <div className="page-wrapper">
      {/* Header */}
      <div className="page-header">
        <div>
          <h1 className="text-2xl font-bold text-corporate-900">Mis Compras</h1>
          <p className="text-sm text-corporate-500 mt-0.5">
            Hola, <span className="font-medium text-corporate-700">{user?.firstName}</span> — historial completo de tus pedidos
          </p>
        </div>
        <button id="orders-refresh-btn" onClick={fetchOrders} className="btn btn-secondary btn-sm">
          ↺ Actualizar
        </button>
      </div>

      {/* Summary cards */}
      {!loading && receipts.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
          <StatCard label="Total de pedidos"   value={receipts.length} unit="pedidos" color="primary" />
          <StatCard label="Total gastado"      value={fmt(totalGastado)} color="emerald" />
          <StatCard label="Último pedido"      value={fmtDate(receipts[0]?.createdAt)} color="slate" />
        </div>
      )}

      {error && <div className="alert-error mb-4">{error}</div>}

      {/* Table */}
      <div className="card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="table-base">
            <thead>
              <tr>
                <th>Pedido</th>
                <th>Fecha</th>
                <th className="text-right">Total</th>
                <th className="text-right">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {loading && (
                <tr>
                  <td colSpan={4} className="text-center py-14 text-corporate-400">
                    <svg className="animate-spin w-6 h-6 mx-auto mb-2" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
                    </svg>
                    Cargando tus compras…
                  </td>
                </tr>
              )}
              {!loading && receipts.length === 0 && (
                <tr>
                  <td colSpan={4} className="text-center py-16">
                    <div className="flex flex-col items-center gap-3 text-corporate-400">
                      <svg className="w-12 h-12 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                          d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2 9m12-9l2 9M9 21h6" />
                      </svg>
                      <p className="text-sm font-medium">Aún no tienes compras</p>
                      <a href="/catalog" className="btn btn-primary btn-sm">Explorar catálogo</a>
                    </div>
                  </td>
                </tr>
              )}
              {!loading && receipts.map((r) => (
                <tr key={r.id} className="animate-fade-in">
                  <td>
                    <span className="font-mono text-xs text-corporate-400">#{r.id}</span>
                  </td>
                  <td className="text-sm text-corporate-600">{fmtDate(r.createdAt)}</td>
                  <td className="text-right font-bold text-primary-700 text-base">{fmt(r.total)}</td>
                  <td className="text-right">
                    <button
                      id={`my-order-detail-${r.id}`}
                      onClick={() => setDetailId(r.id)}
                      className="btn btn-secondary btn-sm"
                    >
                      Ver detalle
                    </button>
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
    </div>
  );
}

/* ── Tarjeta de estadística ── */
function StatCard({ label, value, unit = '', color = 'primary' }) {
  const colors = {
    primary: 'bg-primary-50 border-primary-200 text-primary-800',
    emerald: 'bg-emerald-50 border-emerald-200 text-emerald-800',
    slate:   'bg-corporate-50 border-corporate-200 text-corporate-700',
  };
  return (
    <div className={`rounded-xl border px-5 py-4 animate-fade-in ${colors[color]}`}>
      <p className="text-xs font-semibold uppercase tracking-wider opacity-70">{label}</p>
      <p className="text-2xl font-bold mt-1">{value} <span className="text-sm font-normal opacity-60">{unit}</span></p>
    </div>
  );
}
