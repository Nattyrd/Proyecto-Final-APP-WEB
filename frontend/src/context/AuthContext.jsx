import { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { authAPI } from '../api/api';

/* ── Tipos de datos devueltos por la API ───────── *
 * AuthResponse  →  { token, user: UserResponse }
 * UserResponse  →  { id, username, email, firstName, lastName, role }
 * ──────────────────────────────────────────────── */

const AuthContext = createContext(null);

/**
 * AuthProvider
 * Gestiona el usuario autenticado (token + datos).
 * Persiste la sesión en localStorage para sobrevivir recargas.
 */
export function AuthProvider({ children }) {
  const [user,  setUser]  = useState(() => {
    try {
      const raw = localStorage.getItem('user');
      return raw ? JSON.parse(raw) : null;
    } catch { return null; }
  });
  const [token, setToken] = useState(() => localStorage.getItem('token') ?? null);
  const [loading, setLoading] = useState(false);
  const [error,   setError]   = useState(null);

  /* ── helpers ─────────────────────────────────── */
  const persist = useCallback((tok, usr) => {
    localStorage.setItem('token', tok);
    localStorage.setItem('user',  JSON.stringify(usr));
    setToken(tok);
    setUser(usr);
  }, []);

  const clearSession = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setToken(null);
    setUser(null);
  }, []);

  /* ── login ───────────────────────────────────── */
  const login = useCallback(async ({ username, password }) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await authAPI.login({ username, password });
      // data: { token, user }
      persist(data.token, data.user);
      return { ok: true };
    } catch (err) {
      const msg = err.response?.data?.error ?? err.response?.data?.message ?? 'Error al iniciar sesión';
      setError(msg);
      return { ok: false, error: msg };
    } finally {
      setLoading(false);
    }
  }, [persist]);

  /* ── register ────────────────────────────────── */
  const register = useCallback(async (formData) => {
    setLoading(true);
    setError(null);
    try {
      const payload = {
        username:    formData.username,
        email:       formData.email,
        password:    formData.password,
        firstName:   formData.firstName,
        lastName:    formData.lastName,
      };
      if (formData.adminSecret?.trim()) {
        payload.adminSecret = formData.adminSecret.trim();
      }
      const { data } = await authAPI.register(payload);
      persist(data.token, data.user);
      return { ok: true };
    } catch (err) {
      const msg = err.response?.data?.error ?? err.response?.data?.message ?? 'Error al registrarse';
      setError(msg);
      return { ok: false, error: msg };
    } finally {
      setLoading(false);
    }
  }, [persist]);

  /* ── logout ──────────────────────────────────── */
  const logout = useCallback(() => clearSession(), [clearSession]);

  /* ── computed ────────────────────────────────── */
  const isAuthenticated = !!token;
  const isAdmin = user?.role === 'ADMIN';

  const value = {
    user, token, loading, error,
    isAuthenticated, isAdmin,
    login, register, logout,
    setError,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/** Hook para consumir el contexto desde cualquier componente */
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth debe usarse dentro de <AuthProvider>');
  return ctx;
}
