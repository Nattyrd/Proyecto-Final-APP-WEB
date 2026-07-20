import axios from 'axios';

const BASE_URL = 'http://localhost:8080/api';

/**
 * Instancia principal de Axios.
 * El interceptor de request adjunta automáticamente el JWT
 * almacenado en localStorage si existe.
 */
const api = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 10_000,
});

/* ── Request interceptor: adjunta Bearer token ── */
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

/* ── Response interceptor: manejo global de 401 ─ */
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expirado o inválido → limpiar sesión
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  },
);

/* ──────────────────────────────────────────────── *
 *  Auth endpoints
 * ──────────────────────────────────────────────── */
export const authAPI = {
  /** POST /api/users/register */
  register: (data) => api.post('/users/register', data),
  /** POST /api/users/login */
  login:    (data) => api.post('/users/login',    data),
};

/* ──────────────────────────────────────────────── *
 *  Products endpoints
 * ──────────────────────────────────────────────── */
export const productsAPI = {
  /** GET /api/products?page=1&pageSize=12 */
  getAll:    (page = 1, pageSize = 12) =>
    api.get('/products', { params: { page, pageSize } }),
  /** GET /api/products/:id */
  getById:   (id)  => api.get(`/products/${id}`),
  /** POST /api/products  (ADMIN only) */
  create:    (data) => api.post('/products', data),
  /** PUT /api/products/:id  (ADMIN only) */
  update:    (id, data) => api.put(`/products/${id}`, data),
  /** DELETE /api/products/:id  (ADMIN only) */
  remove:    (id)  => api.delete(`/products/${id}`),
};

/* ──────────────────────────────────────────────── *
 *  Receipts endpoints  (require auth)
 * ──────────────────────────────────────────────── */
export const receiptsAPI = {
  /** POST /api/receipts */
  create:       (data)   => api.post('/receipts', data),
  /** GET  /api/receipts */
  getAll:       ()       => api.get('/receipts'),
  /** GET  /api/receipts/user/:userId */
  getByUser:    (userId) => api.get(`/receipts/user/${userId}`),
  /** GET  /api/receipts/:id */
  getById:      (id)     => api.get(`/receipts/${id}`),
  /** DELETE /api/receipts/:id */
  remove:       (id)     => api.delete(`/receipts/${id}`),
};

export default api;
