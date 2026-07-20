import { Link, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';

export default function Navbar() {
  const { isAuthenticated, isAdmin, user, logout } = useAuth();
  const { totalItems } = useCart();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const linkClass = ({ isActive }) =>
    `text-sm font-medium transition-colors duration-150 ${
      isActive
        ? 'text-primary-700 border-b-2 border-primary-700 pb-0.5'
        : 'text-corporate-600 hover:text-corporate-900'
    }`;

  return (
    <header className="bg-white border-b border-corporate-200 shadow-sm sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-14">
          {/* Logo */}
          <Link
            to="/"
            className="flex items-center gap-2 text-primary-800 font-bold text-lg tracking-tight"
          >
            <svg
              className="w-7 h-7 text-primary-700"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2 9m12-9l2 9M9 21h6"
              />
            </svg>
            ShopCorp
          </Link>

          {/* Nav links */}
          <nav className="hidden sm:flex items-center gap-6">
            <NavLink to="/" end className={linkClass}>
              Inicio
            </NavLink>
            <NavLink to="/catalog" className={linkClass}>
              Catálogo
            </NavLink>
          </nav>

          {/* Right side */}
          <div className="flex items-center gap-2">
            {isAuthenticated ? (
              <>
                {/* Carrito */}
                <Link
                  to="/cart"
                  id="nav-cart-btn"
                  className="relative btn-ghost btn rounded-full p-2"
                  title="Carrito"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
                    <path strokeLinecap="round" strokeLinejoin="round"
                      d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2 9m12-9l2 9M9 21h6" />
                  </svg>
                  {totalItems > 0 && (
                    <span className="absolute -top-1 -right-1 bg-primary-700 text-white text-[10px] font-bold
                                     rounded-full w-4 h-4 flex items-center justify-center">
                      {totalItems > 99 ? '99+' : totalItems}
                    </span>
                  )}
                </Link>

                {/* Mis Compras (Client) */}
                {!isAdmin && (
                  <NavLink to="/my-orders" className={linkClass}>
                    Mis Compras
                  </NavLink>
                )}

                {/* Admin panel */}
                {isAdmin && (
                  <NavLink to="/admin" className={linkClass}>
                    Panel Admin
                  </NavLink>
                )}

                {/* User info + logout */}
                <div className="flex items-center gap-2 pl-2 border-l border-corporate-200">
                  <div className="hidden sm:block text-right">
                    <p className="text-xs font-semibold text-corporate-800 leading-none">
                      {user?.firstName} {user?.lastName}
                    </p>
                    <span className={`text-[10px] font-medium ${isAdmin ? 'text-primary-700' : 'text-corporate-400'}`}>
                      {user?.role}
                    </span>
                  </div>
                  <button
                    id="nav-logout-btn"
                    onClick={handleLogout}
                    className="btn btn-secondary btn-sm"
                  >
                    Salir
                  </button>
                </div>
              </>
            ) : (
              <>
                <Link to="/login" id="nav-login-btn" className="btn btn-ghost btn-sm">
                  Iniciar sesión
                </Link>
                <Link to="/register" id="nav-register-btn" className="btn btn-primary btn-sm">
                  Registrarse
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
