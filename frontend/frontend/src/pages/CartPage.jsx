import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';
import { receiptsAPI } from '../api/api';

export default function CartPage() {
  const { items, totalItems, totalPrice, updateQty, removeItem, clearCart } = useCart();
  const { user, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  const [loading,  setLoading]  = useState(false);
  const [success,  setSuccess]  = useState(null);
  const [apiError, setApiError] = useState('');

  const handleCheckout = async () => {
    if (!isAuthenticated) { navigate('/login'); return; }
    if (items.length === 0) return;

    setLoading(true);
    setApiError('');
    try {
      const payload = {
        userId: user.id,
        items: items.map((i) => ({
          productId: i.product.id,
          quantity:  i.quantity,
        })),
      };
      const { data } = await receiptsAPI.create(payload);
      setSuccess(data);
      clearCart();
    } catch (err) {
      setApiError(
        err.response?.data?.error ?? err.response?.data?.message ?? 'Error al procesar la compra.',
      );
    } finally {
      setLoading(false);
    }
  };

  /* ── Order success screen ─────────────────────── */
  if (success) {
    return (
      <div className="page-wrapper max-w-2xl">
        <div className="card p-8 text-center animate-fade-in">
          <div className="w-16 h-16 bg-emerald-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 className="text-emerald-700 mb-1">¡Compra realizada!</h1>
          <p className="text-sm text-corporate-500 mb-6">
            Recibo #{success.id} — Total:{' '}
            <strong>${parseFloat(success.total).toFixed(2)}</strong>
          </p>

          {/* Receipt items */}
          <div className="text-left border border-corporate-200 rounded-lg overflow-hidden mb-6">
            <table className="table-base">
              <thead>
                <tr>
                  <th>Producto ID</th>
                  <th>Cant.</th>
                  <th>Precio unit.</th>
                  <th>Subtotal</th>
                </tr>
              </thead>
              <tbody>
                {success.items?.map((it) => (
                  <tr key={it.id}>
                    <td>#{it.productId}</td>
                    <td>{it.quantity}</td>
                    <td>${parseFloat(it.unitPrice).toFixed(2)}</td>
                    <td>${parseFloat(it.subtotal).toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <button className="btn btn-primary" onClick={() => navigate('/')}>
            Volver al catálogo
          </button>
        </div>
      </div>
    );
  }

  /* ── Empty cart ───────────────────────────────── */
  if (totalItems === 0) {
    return (
      <div className="page-wrapper max-w-xl">
        <div className="card p-12 text-center animate-fade-in">
          <svg className="w-14 h-14 mx-auto mb-4 text-corporate-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
              d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2 9m12-9l2 9M9 21h6" />
          </svg>
          <h2 className="text-corporate-700 mb-2">Tu carrito está vacío</h2>
          <p className="text-sm text-corporate-400 mb-6">Añade productos desde el catálogo.</p>
          <button className="btn btn-primary" onClick={() => navigate('/catalog')}>
            Ir al catálogo
          </button>
        </div>
      </div>
    );
  }

  /* ── Cart view ────────────────────────────────── */
  return (
    <div className="page-wrapper">
      <div className="page-header">
        <h1>Carrito de compras</h1>
        <span className="text-sm text-corporate-500">{totalItems} artículo(s)</span>
      </div>

      {apiError && <div className="alert-error mb-4">{apiError}</div>}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Items */}
        <div className="lg:col-span-2 space-y-3">
          {items.map(({ product, quantity }) => (
            <div key={product.id} className="card p-4 flex gap-4 items-center animate-fade-in">
              {/* Icon */}
              <div className="w-12 h-12 rounded-lg bg-primary-50 flex items-center justify-center flex-shrink-0">
                <svg className="w-6 h-6 text-primary-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10" />
                </svg>
              </div>

              <div className="flex-1 min-w-0">
                <p className="text-sm font-semibold text-corporate-900 truncate">{product.name}</p>
                <p className="text-xs text-primary-700 font-medium">
                  ${parseFloat(product.price).toFixed(2)} / unidad
                </p>
              </div>

              {/* Quantity control */}
              <div className="flex items-center gap-1">
                <button
                  id={`cart-dec-${product.id}`}
                  onClick={() => updateQty(product.id, quantity - 1)}
                  className="btn btn-ghost btn-sm rounded-full w-7 h-7 p-0 justify-center"
                >−</button>
                <span className="w-8 text-center text-sm font-semibold text-corporate-800">
                  {quantity}
                </span>
                <button
                  id={`cart-inc-${product.id}`}
                  onClick={() => updateQty(product.id, quantity + 1)}
                  disabled={quantity >= product.stock}
                  className="btn btn-ghost btn-sm rounded-full w-7 h-7 p-0 justify-center"
                >+</button>
              </div>

              {/* Subtotal */}
              <p className="text-sm font-bold text-corporate-800 w-20 text-right">
                ${(parseFloat(product.price) * quantity).toFixed(2)}
              </p>

              {/* Remove */}
              <button
                id={`cart-remove-${product.id}`}
                onClick={() => removeItem(product.id)}
                className="btn btn-ghost btn-sm text-red-400 hover:text-red-600 p-1"
                title="Eliminar"
              >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          ))}
        </div>

        {/* Summary */}
        <div className="card p-6 h-fit space-y-4">
          <h2 className="text-base font-semibold text-corporate-800 border-b border-corporate-100 pb-3">
            Resumen del pedido
          </h2>

          <div className="space-y-2 text-sm">
            {items.map(({ product, quantity }) => (
              <div key={product.id} className="flex justify-between text-corporate-600">
                <span className="truncate max-w-[60%]">{product.name} × {quantity}</span>
                <span>${(parseFloat(product.price) * quantity).toFixed(2)}</span>
              </div>
            ))}
          </div>

          <div className="border-t border-corporate-100 pt-3 flex justify-between font-bold text-corporate-900">
            <span>Total</span>
            <span className="text-primary-700">${totalPrice.toFixed(2)}</span>
          </div>

          <button
            id="checkout-btn"
            onClick={handleCheckout}
            disabled={loading}
            className="btn btn-primary w-full justify-center"
          >
            {loading ? 'Procesando…' : 'Finalizar compra'}
          </button>

          <button
            id="cart-clear-btn"
            onClick={clearCart}
            className="btn btn-ghost w-full justify-center text-red-400 text-xs"
          >
            Vaciar carrito
          </button>
        </div>
      </div>
    </div>
  );
}
