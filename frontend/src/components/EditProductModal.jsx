import { useState, useEffect } from 'react';
import { productsAPI } from '../api/api';

/**
 * Modal para editar un producto existente.
 * Props:
 *   product   – objeto producto a editar
 *   onClose   () => void
 *   onUpdated (product) => void
 */
export default function EditProductModal({ product, onClose, onUpdated }) {
  const [form,    setForm]    = useState({ name: '', description: '', price: '', stock: '' });
  const [errors,  setErrors]  = useState({});
  const [apiErr,  setApiErr]  = useState('');
  const [loading, setLoading] = useState(false);

  /* Pre-rellenar con datos actuales */
  useEffect(() => {
    if (product) {
      setForm({
        name:        product.name        ?? '',
        description: product.description ?? '',
        price:       String(product.price ?? ''),
        stock:       String(product.stock ?? ''),
      });
    }
  }, [product]);

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
    setForm((p)   => ({ ...p, [name]: value }));
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
      const { data } = await productsAPI.update(product.id, payload);
      onUpdated(data);
      onClose();
    } catch (err) {
      setApiErr(
        err.response?.data?.error ?? err.response?.data?.message ?? 'Error al actualizar el producto.',
      );
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
        {/* Header */}
        <div className="flex items-center justify-between mb-5">
          <div>
            <h2 className="text-lg font-semibold text-corporate-900">Editar producto</h2>
            <p className="text-xs text-corporate-400 mt-0.5">ID #{product?.id}</p>
          </div>
          <button
            id="edit-modal-close-btn"
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
            <label htmlFor="edit-prod-name" className="label">Nombre <span className="text-red-400">*</span></label>
            <input
              id="edit-prod-name" name="name" type="text"
              value={form.name} onChange={handleChange}
              className={`input ${errors.name ? 'input-error' : ''}`}
            />
            {errors.name && <p className="field-error">{errors.name}</p>}
          </div>

          {/* Descripción */}
          <div>
            <label htmlFor="edit-prod-description" className="label">
              Descripción <span className="text-corporate-400 font-normal">(opcional)</span>
            </label>
            <textarea
              id="edit-prod-description" name="description" rows={3}
              value={form.description} onChange={handleChange}
              className="input resize-none"
            />
          </div>

          {/* Precio + Stock */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label htmlFor="edit-prod-price" className="label">Precio (USD) <span className="text-red-400">*</span></label>
              <input
                id="edit-prod-price" name="price" type="number" min="0.01" step="0.01"
                value={form.price} onChange={handleChange}
                className={`input ${errors.price ? 'input-error' : ''}`}
              />
              {errors.price && <p className="field-error">{errors.price}</p>}
            </div>
            <div>
              <label htmlFor="edit-prod-stock" className="label">Stock <span className="text-red-400">*</span></label>
              <input
                id="edit-prod-stock" name="stock" type="number" min="0" step="1"
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
              id="edit-product-submit-btn"
              type="submit"
              disabled={loading}
              className="btn btn-primary"
            >
              {loading ? 'Guardando…' : 'Guardar cambios'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
