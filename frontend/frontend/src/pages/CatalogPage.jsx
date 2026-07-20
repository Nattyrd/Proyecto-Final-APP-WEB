import { useState, useEffect, useCallback } from 'react';
import { productsAPI } from '../api/api';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';

/* ── Product Card ────────────────────────────────── */
function ProductCard({ product, onAdd }) {
  const [added, setAdded] = useState(false);

  const handleAdd = () => {
    onAdd(product);
    setAdded(true);
    setTimeout(() => setAdded(false), 1500);
  };

  return (
    <article className="card flex flex-col animate-fade-in overflow-hidden">
      {/* Imagen placeholder */}
      <div className="bg-gradient-to-br from-primary-50 to-corporate-100 h-40 flex items-center justify-center">
        <svg className="w-16 h-16 text-primary-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1}
            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10" />
        </svg>
      </div>

      <div className="flex flex-col flex-1 p-4 gap-2">
        <h3 className="text-sm font-semibold text-corporate-900 line-clamp-2">{product.name}</h3>
        {product.description && (
          <p className="text-xs text-corporate-500 line-clamp-2 flex-1">{product.description}</p>
        )}

        <div className="flex items-end justify-between mt-auto pt-2">
          <div>
            <p className="text-lg font-bold text-primary-700">
              ${parseFloat(product.price).toFixed(2)}
            </p>
            <p className={`text-[11px] font-medium ${product.stock > 0 ? 'text-emerald-600' : 'text-red-500'}`}>
              {product.stock > 0 ? `${product.stock} disponibles` : 'Sin stock'}
            </p>
          </div>

          <button
            id={`add-to-cart-${product.id}`}
            onClick={handleAdd}
            disabled={product.stock === 0 || added}
            className={`btn btn-sm transition-all duration-300 ${
              added ? 'btn-secondary text-emerald-600 border-emerald-300' : 'btn-primary'
            }`}
          >
            {added ? '✓ Añadido' : '+ Carrito'}
          </button>
        </div>
      </div>
    </article>
  );
}

/* ── Pagination ──────────────────────────────────── */
function Pagination({ page, totalPages, onPrev, onNext }) {
  return (
    <div className="flex items-center justify-center gap-4 mt-8">
      <button
        id="catalog-prev-btn"
        onClick={onPrev}
        disabled={page <= 1}
        className="btn btn-secondary"
      >
        ← Anterior
      </button>
      <span className="text-sm text-corporate-500">
        Página <strong className="text-corporate-800">{page}</strong> de{' '}
        <strong className="text-corporate-800">{totalPages}</strong>
      </span>
      <button
        id="catalog-next-btn"
        onClick={onNext}
        disabled={page >= totalPages}
        className="btn btn-secondary"
      >
        Siguiente →
      </button>
    </div>
  );
}

/* ── Catalog Page ────────────────────────────────── */
export default function CatalogPage() {
  const { addItem } = useCart();
  const { isAuthenticated } = useAuth();

  const [products,    setProducts]    = useState([]);
  const [page,        setPage]        = useState(1);
  const [totalPages,  setTotalPages]  = useState(1);
  const [totalItems,  setTotalItems]  = useState(0);
  const [loading,     setLoading]     = useState(false);
  const [error,       setError]       = useState('');
  const [search,      setSearch]      = useState('');

  const PAGE_SIZE = 12;

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
      setError('No se pudieron cargar los productos. Verifica que el servidor esté activo.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchProducts(page); }, [page, fetchProducts]);

  const filtered = search.trim()
    ? products.filter((p) =>
        p.name.toLowerCase().includes(search.toLowerCase()) ||
        p.description?.toLowerCase().includes(search.toLowerCase()),
      )
    : products;

  return (
    <div className="page-wrapper">
      {/* Header */}
      <div className="page-header">
        <div>
          <h1>Catálogo de Productos</h1>
          {!loading && (
            <p className="text-sm text-corporate-500 mt-0.5">
              {totalItems} productos encontrados
            </p>
          )}
        </div>

        {/* Búsqueda local */}
        <div className="relative w-full max-w-xs">
          <input
            id="catalog-search"
            type="search"
            placeholder="Buscar producto…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="input pl-9"
          />
          <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-corporate-400"
            fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
              d="M21 21l-4.35-4.35M17 11A6 6 0 1 1 5 11a6 6 0 0 1 12 0z" />
          </svg>
        </div>
      </div>

      {/* Auth hint */}
      {!isAuthenticated && (
        <div className="alert-info mb-6">
          Inicia sesión para poder añadir productos al carrito.
        </div>
      )}

      {/* Error */}
      {error && <div className="alert-error mb-6">{error}</div>}

      {/* Loading skeleton */}
      {loading && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="card h-60 animate-pulse bg-corporate-100" />
          ))}
        </div>
      )}

      {/* Grid */}
      {!loading && filtered.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
          {filtered.map((p) => (
            <ProductCard key={p.id} product={p} onAdd={addItem} />
          ))}
        </div>
      )}

      {/* Empty */}
      {!loading && filtered.length === 0 && !error && (
        <div className="text-center py-20 text-corporate-400">
          <svg className="w-12 h-12 mx-auto mb-3 text-corporate-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
              d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0H4" />
          </svg>
          <p className="text-sm">No se encontraron productos.</p>
        </div>
      )}

      {/* Pagination */}
      {!loading && totalPages > 1 && !search && (
        <Pagination
          page={page}
          totalPages={totalPages}
          onPrev={() => setPage((p) => p - 1)}
          onNext={() => setPage((p) => p + 1)}
        />
      )}
    </div>
  );
}
