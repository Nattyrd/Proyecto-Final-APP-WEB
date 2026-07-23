import { useState } from 'react';
import { productsAPI } from '../api/api';

const INITIAL = { name: '', description: '', price: '', stock: '' };

/**
 * Modal con formulario para crear un nuevo producto.
 * Props:
 *   onClose   () => void
 *   onCreated (product) => void
 */
export default function CreateProductModal({ onClose, onCreated }) {
  const [form,    setForm]    = useState(INITIAL);
  const [errors,  setErrors]  = useState({});
  const [apiErr,  setApiErr]  = useState('');
  const [loading, setLoading] = useState(false);

  const validate = () => {
    const e = {};
    if (!form.name.trim() || form.name.length < 2) e.name  = 'Mínimo 2 caracteres';
    const price = parseFloat(form.price);
    if (isNaN(price) || price <= 0)                e.price = 'Precio debe ser un número positivo';
    const stock = parseInt(form.stock, 10);
    if (isNaN(stock) || stock < 0)                 e.stock = 'Stock debe ser ≥ 0';
    return e;
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((p) => ({ ...p, [name]: value }));
    setErrors((p) => ({ ...p, [name]: '' }));
    setApiErr('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length) { setErrors(errs); return; }

    setLoading(true);
    try {
      const payload = {
        name:        form.name.trim(),
        description: form.description.trim(),
        price:       parseFloat(form.price),
        stock:       parseInt(form.stock, 10),
      };
      const { data } = await productsAPI.create(payload);
      onCreated(data);
      onClose();
    } catch (err) {
      setApiErr(
        err.response?.data?.error ?? err.response?.data?.message ?? 'Error al crear el producto.',
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    /* Backdrop */
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-corporate-900/40 backdrop-blur-sm"
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div className="card animate-slide-in w-full max-w-md mx-4 p-6">
        {/* Header */}
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-lg font-semibold text-corporate-900">Nuevo producto</h2>
          <button
            id="modal-close-btn"
            onClick={onClose}
            className="btn btn-ghost btn-sm p-1 text-corporate-400 hover:text-corporate-700"
            aria-label="Cerrar"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {apiErr && <div className="alert-error mb-4">{apiErr}</div>}

        <form onSubmit={handleSubmit} noValidate className="space-y-4">
          {/* Nombre */}
          <div>
            <label htmlFor="prod-name" className="label">Nombre <span className="text-red-400">*</span></label>
            <input
              id="prod-name" name="name" type="text"
              placeholder="Ej. Laptop Corporativa Pro"
              value={form.name} onChange={handleChange}
              className={`input ${errors.name ? 'input-error' : ''}`}
            />
            {errors.name && <p className="field-error">{errors.name}</p>}
          </div>

          {/* Descripción */}
          <div>
            <label htmlFor="prod-description" className="label">
              Descripción <span className="text-corporate-400 font-normal">(opcional)</span>
            </label>
            <textarea
              id="prod-description" name="description" rows={3}
              placeholder="Describe brevemente el producto…"
              value={form.description} onChange={handleChange}
              className="input resize-none"
            />
          </div>

          {/* Precio + Stock */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label htmlFor="prod-price" className="label">Precio (USD) <span className="text-red-400">*</span></label>
              <input
                id="prod-price" name="price" type="number" min="0.01" step="0.01"
                placeholder="0.00"
                value={form.price} onChange={handleChange}
                className={`input ${errors.price ? 'input-error' : ''}`}
              />
              {errors.price && <p className="field-error">{errors.price}</p>}
            </div>
            <div>
              <label htmlFor="prod-stock" className="label">Stock <span className="text-red-400">*</span></label>
              <input
                id="prod-stock" name="stock" type="number" min="0" step="1"
                placeholder="0"
                value={form.stock} onChange={handleChange}
                className={`input ${errors.stock ? 'input-error' : ''}`}
              />
              {errors.stock && <p className="field-error">{errors.stock}</p>}
            </div>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3 pt-2">
            <button type="button" onClick={onClose} className="btn btn-secondary">
              Cancelar
            </button>
            <button
              id="create-product-submit-btn"
              type="submit"
              disabled={loading}
              className="btn btn-primary"
            >
              {loading ? 'Guardando…' : 'Crear producto'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
