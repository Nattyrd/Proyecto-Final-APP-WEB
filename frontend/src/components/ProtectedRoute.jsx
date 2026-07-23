import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

/** Redirige a /login si no hay sesión activa */
export function RequireAuth({ children }) {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? children : <Navigate to="/login" replace />;
}

/** Redirige a / si el usuario no es ADMIN */
export function RequireAdmin({ children }) {
  const { isAuthenticated, isAdmin } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  if (!isAdmin)         return <Navigate to="/"      replace />;
  return children;
}
