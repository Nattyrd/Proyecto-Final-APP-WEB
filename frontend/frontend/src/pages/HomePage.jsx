import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function HomePage() {
  const { isAuthenticated, isAdmin } = useAuth();

  return (
    <div className="page-wrapper">
      {/* Hero */}
      <section className="text-center py-16 animate-fade-in">
        <h1 className="text-3xl sm:text-4xl font-bold text-corporate-900 mb-4">
          Bienvenido a <span className="text-primary-700">ShopCorp</span>
        </h1>
        <p className="text-corporate-500 max-w-md mx-auto mb-8">
          Explora nuestro catálogo, gestiona tu carrito y realiza tus pedidos de
          forma rápida y segura.
        </p>
        <div className="flex flex-wrap justify-center gap-3">
          <Link to="/catalog" className="btn btn-primary btn-lg">
            Ver catálogo
          </Link>
          {!isAuthenticated && (
            <Link to="/register" className="btn btn-secondary btn-lg">
              Crear cuenta
            </Link>
          )}
          {isAdmin && (
            <Link to="/admin" className="btn btn-secondary btn-lg">
              Panel Admin
            </Link>
          )}
        </div>
      </section>

      {/* Feature cards */}
      <section className="grid grid-cols-1 sm:grid-cols-3 gap-5 mt-4">
        {[
          {
            title: "Catálogo completo",
            desc: "Navega nuestra selección de productos con filtros y paginación.",
          },
          {
            title: "Compra segura",
            desc: "Autenticación JWT y manejo de roles para máxima seguridad.",
          },
          {
            title: "Recibos instantáneos",
            desc: "Al finalizar tu compra recibes tu comprobante de forma inmediata.",
          },
        ].map((f) => (
          <div key={f.title} className="card p-6 text-center animate-fade-in">
            <div className="text-3xl mb-3">{f.icon}</div>
            <h3 className="text-corporate-900 mb-1">{f.title}</h3>
            <p className="text-sm text-corporate-500">{f.desc}</p>
          </div>
        ))}
      </section>
    </div>
  );
}
