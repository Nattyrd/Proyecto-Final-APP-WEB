import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const { login, loading } = useAuth();
  const navigate = useNavigate();

  const [form,   setForm]   = useState({ username: '', password: '' });
  const [errors, setErrors] = useState({});
  const [apiErr, setApiErr] = useState('');

  const validate = () => {
    const e = {};
    if (!form.username.trim()) e.username = 'El usuario es obligatorio';
    if (!form.password)        e.password = 'La contraseña es obligatoria';
    return e;
  };

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
    setErrors((prev) => ({ ...prev, [e.target.name]: '' }));
    setApiErr('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length) { setErrors(errs); return; }
    const res = await login(form);
    if (res.ok) navigate('/');
    else setApiErr(res.error);
  };

  return (
    <div className="min-h-[calc(100vh-56px)] flex items-center justify-center bg-corporate-50 px-4 py-12">
      <div className="card animate-fade-in w-full max-w-sm p-8">
        {/* Header */}
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-bold text-corporate-900">Iniciar sesión</h1>
          <p className="text-sm text-corporate-500 mt-1">Accede a tu cuenta corporativa</p>
        </div>

        {/* API error */}
        {apiErr && <div className="alert-error mb-4">{apiErr}</div>}

        <form onSubmit={handleSubmit} noValidate className="space-y-4">
          {/* Usuario */}
          <div>
            <label htmlFor="login-username" className="label">Usuario</label>
            <input
              id="login-username"
              name="username"
              type="text"
              autoComplete="username"
              className={`input ${errors.username ? 'input-error' : ''}`}
              placeholder="tu_usuario"
              value={form.username}
              onChange={handleChange}
            />
            {errors.username && <p className="field-error">{errors.username}</p>}
          </div>

          {/* Contraseña */}
          <div>
            <label htmlFor="login-password" className="label">Contraseña</label>
            <input
              id="login-password"
              name="password"
              type="password"
              autoComplete="current-password"
              className={`input ${errors.password ? 'input-error' : ''}`}
              placeholder="••••••"
              value={form.password}
              onChange={handleChange}
            />
            {errors.password && <p className="field-error">{errors.password}</p>}
          </div>

          <button
            id="login-submit-btn"
            type="submit"
            disabled={loading}
            className="btn btn-primary w-full justify-center mt-2"
          >
            {loading ? 'Ingresando…' : 'Iniciar sesión'}
          </button>
        </form>

        <p className="text-center text-sm text-corporate-500 mt-6">
          ¿No tienes cuenta?{' '}
          <Link to="/register" className="text-primary-700 font-medium hover:underline">
            Regístrate
          </Link>
        </p>
      </div>
    </div>
  );
}
