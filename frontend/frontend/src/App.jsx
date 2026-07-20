import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { CartProvider } from './context/CartContext';
import Navbar from './components/Navbar';
import { RequireAuth, RequireAdmin } from './components/ProtectedRoute';

import HomePage     from './pages/HomePage';
import LoginPage    from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import CatalogPage  from './pages/CatalogPage';
import CartPage     from './pages/CartPage';
import AdminPage    from './pages/AdminPage';

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <CartProvider>
          <div className="flex flex-col min-h-screen">
            <Navbar />

            <main className="flex-1">
              <Routes>
                {/* Public */}
                <Route path="/"         element={<HomePage />} />
                <Route path="/login"    element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />
                <Route path="/catalog"  element={<CatalogPage />} />

                {/* Auth required */}
                <Route path="/cart" element={
                  <RequireAuth><CartPage /></RequireAuth>
                } />

                {/* Admin only */}
                <Route path="/admin" element={
                  <RequireAdmin><AdminPage /></RequireAdmin>
                } />

                {/* Fallback */}
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </main>

            {/* Footer */}
            <footer className="border-t border-corporate-200 bg-white text-center py-4 text-xs text-corporate-400">
              © {new Date().getFullYear()} ShopCorp · E-Commerce Grupo 5 ·{' '}
              <span className="text-primary-600">API:</span> http://localhost:8080/api
            </footer>
          </div>
        </CartProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}
