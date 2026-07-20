import { useState, useEffect, useCallback } from 'react';
import { productsAPI } from '../api/api';
import CreateProductModal from '../components/CreateProductModal';

export default function AdminPage() {
  const [products,    setProducts]    = useState([]);
  const [page,        setPage]        = useState(1);
  const [totalPages,  setTotalPages]  = useState(1);
  const [totalItems,  setTotalItems]  = useState(0);
  const [loading,     setLoading]     = useState(false);
  const [error,       setError]       = useState('');
  const [showModal,   setShowModal]   = useState(false);
  const [deleting,    setDeleting]    = useState(null);
  const [successMsg,  setSuccessMsg]  = useState('');

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

  const handleCreated = (newProduct) => {
    setSuccessMsg(`Producto "${newProduct.name}" creado correctamente.`);
    fetchProducts(1);
    setTimeout(() => setSuccessMsg(''), 4000);
  };

  const handleDelete = async (product) => {
    if (!window.confirm(`¿Eliminar "${product.name}"? Esta acción no se puede deshacer.`)) return;
    setDeleting(product.id);
    try {
      await productsAPI.remove(product.id);
      setProducts((prev) => prev.filter((p) => p.id !== product.id));
      setTotalItems((n) => n - 1);
      setSuccessMsg(`Producto "${product.name}" eliminado.`);
      setTimeout(() => setSuccessMsg(''), 4000);
    } catch (err) {
      setError(err.response?.data?.error ?? 'Error al eliminar el producto.');
    } finally {
      setDeleting(null);
    }
  };

  return (
    <div className="page-wrapper">
      {/* Page header */}
      <div className="page-header">
        <div>
          <h1>Panel de Administración</h1>
          <p className="text-sm text-corporate-500 mt-0.5">
            Gestión de productos — {totalItems} registros
          </p>
        </div>
        <button
          id="admin-new-product-btn"
          onClick={() => setShowModal(true)}
          className="btn btn-primary"
        >
          + Nuevo producto
        </button>
      </div>

      {/* Alerts */}
      {successMsg && <div className="alert-success mb-4 animate-fade-in">{successMsg}</div>}
      {error      && <div className="alert-error  mb-4">{error}</div>}

      {/* Table */}
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
                <tr>
                  <td colSpan={6} className="text-center py-10 text-corporate-400">
                    Cargando…
                  </td>
                </tr>
              )}
              {!loading && products.length === 0 && (
                <tr>
                  <td colSpan={6} className="text-center py-10 text-corporate-400">
                    No hay productos registrados.
                  </td>
                </tr>
              )}
              {!loading && products.map((p) => (
                <tr key={p.id} className="animate-fade-in">
                  <td className="text-corporate-400 text-xs">#{p.id}</td>
                  <td className="font-medium text-corporate-900 max-w-[180px] truncate">
                    {p.name}
                  </td>
                  <td className="text-corporate-500 max-w-[220px] truncate text-xs">
                    {p.description || <span className="italic text-corporate-300">Sin descripción</span>}
                  </td>
                  <td className="font-semibold text-primary-700">
                    ${parseFloat(p.price).toFixed(2)}
                  </td>
                  <td>
                    <span className={`badge ${p.stock > 0 ? 'badge-success' : 'bg-red-100 text-red-700'}`}>
                      {p.stock}
                    </span>
                  </td>
                  <td className="text-right">
                    <button
                      id={`admin-delete-${p.id}`}
                      onClick={() => handleDelete(p)}
                      disabled={deleting === p.id}
                      className="btn btn-danger btn-sm"
                    >
                      {deleting === p.id ? '…' : 'Eliminar'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-corporate-100 bg-corporate-50">
            <span className="text-xs text-corporate-400">
              Página {page} de {totalPages}
            </span>
            <div className="flex gap-2">
              <button
                id="admin-prev-btn"
                disabled={page <= 1}
                onClick={() => setPage((p) => p - 1)}
                className="btn btn-secondary btn-sm"
              >
                ← Anterior
              </button>
              <button
                id="admin-next-btn"
                disabled={page >= totalPages}
                onClick={() => setPage((p) => p + 1)}
                className="btn btn-secondary btn-sm"
              >
                Siguiente →
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Modal */}
      {showModal && (
        <CreateProductModal
          onClose={() => setShowModal(false)}
          onCreated={handleCreated}
        />
      )}
    </div>
  );
}
