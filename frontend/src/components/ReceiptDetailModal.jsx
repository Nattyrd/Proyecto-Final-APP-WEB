import { useState, useEffect } from 'react';
import { receiptsAPI } from '../api/api';

/**
 * Modal que muestra el detalle completo de un recibo.
 * Props:
 *   receiptId  – ID del recibo a cargar
 *   onClose    () => void
 */
export default function ReceiptDetailModal({ receiptId, onClose }) {
  const [receipt, setReceipt] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error,   setError]   = useState('');

  useEffect(() => {
    if (!receiptId) return;
    setLoading(true);
    setError('');
    receiptsAPI.getById(receiptId)
      .then(({ data }) => setReceipt(data))
      .catch(() => setError('No se pudo cargar el detalle del recibo.'))
      .finally(() => setLoading(false));
  }, [receiptId]);

  const fmt = (v) => `$${parseFloat(v ?? 0).toFixed(2)}`;
  const fmtDate = (iso) =>
    iso ? new Date(iso).toLocaleString('es-EC', { dateStyle: 'medium', timeStyle: 'short' }) : '—';

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-corporate-900/40 backdrop-blur-sm"
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div className="card animate-slide-in w-full max-w-2xl mx-4 p-6 max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between mb-5">
          <div>
            <h2 className="text-lg font-semibold text-corporate-900">
              Detalle del Recibo {receipt ? `#${receipt.id}` : ''}
            </h2>
            {receipt && (
              <p className="text-xs text-corporate-400 mt-0.5">
                {fmtDate(receipt.createdAt)} · {receipt.username ?? receipt.userId}
              </p>
            )}
          </div>
          <button
            id="receipt-detail-close-btn"
            onClick={onClose}
            className="btn btn-ghost btn-sm p-1 text-corporate-400 hover:text-corporate-700"
            aria-label="Cerrar"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Body */}
        <div className="overflow-y-auto flex-1">
          {loading && (
            <div className="flex items-center justify-center py-16 text-corporate-400">
              <svg className="animate-spin w-6 h-6 mr-2" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
              </svg>
              Cargando…
            </div>
          )}
          {error && <div className="alert-error">{error}</div>}

          {!loading && receipt && (
            <>
              {/* Items table */}
              <div className="overflow-x-auto rounded-lg border border-corporate-200">
                <table className="table-base">
                  <thead>
                    <tr>
                      <th>Producto</th>
                      <th className="text-right">Precio unit.</th>
                      <th className="text-right">Cant.</th>
                      <th className="text-right">Subtotal</th>
                    </tr>
                  </thead>
                  <tbody>
                    {(receipt.items ?? []).map((item, i) => (
                      <tr key={item.id ?? i}>
                        <td className="font-medium text-corporate-800">{item.productName ?? `#${item.productId}`}</td>
                        <td className="text-right text-corporate-600">{fmt(item.unitPrice)}</td>
                        <td className="text-right">
                          <span className="badge badge-client">{item.quantity}</span>
                        </td>
                        <td className="text-right font-semibold text-primary-700">
                          {fmt(parseFloat(item.unitPrice ?? 0) * (item.quantity ?? 1))}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Total */}
              <div className="flex justify-end mt-4">
                <div className="bg-primary-50 border border-primary-200 rounded-lg px-6 py-3 text-right">
                  <p className="text-xs text-primary-600 font-medium uppercase tracking-wider">Total</p>
                  <p className="text-2xl font-bold text-primary-800">{fmt(receipt.total)}</p>
                </div>
              </div>
            </>
          )}
        </div>

        {/* Footer */}
        <div className="flex justify-end mt-5 pt-4 border-t border-corporate-100">
          <button id="receipt-detail-ok-btn" onClick={onClose} className="btn btn-secondary">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  );
}
