import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const INITIAL = {
  username: '', email: '', password: '', confirmPassword: '',
  firstName: '', lastName: '', adminSecret: '',
};

/* ── Field definido FUERA del componente padre ──────────────────────────────
   Si estuviera dentro, React lo trataría como un tipo nuevo en cada render,
   desmontando y remontando el <input> (perdiendo el foco al escribir).
   ──────────────────────────────────────────────────────────────────────── */
function Field({ id, name, label, type = 'text', placeholder, autoComplete, form, errors, onChange }) {
  return (
    <div>
      <label htmlFor={id} className="label">{label}</label>
      <input
        id={id}
        name={name}
        type={type}
        autoComplete={autoComplete}
        placeholder={placeholder}
        value={form[name]}
        onChange={onChange}
        className={`input ${errors[name] ? 'input-error' : ''}`}
      />
      {errors[name] && <p className="field-error">{errors[name]}</p>}
    </div>
  );
}

export default function RegisterPage() {
  const { register, loading } = useAuth();
  const navigate = useNavigate();
  const [form,   setForm]   = useState(INITIAL);
  const [errors, setErrors] = useState({});
  const [apiErr, setApiErr] = useState('');

  const validate = () => {
    const e = {};
    if (!form.username.trim() || form.username.length < 3) e.username        = 'Mínimo 3 caracteres';
    if (!form.email.match(/^[^\s@]+@[^\s@]+\.[^\s@]+$/))  e.email           = 'Email inválido';
    if (form.password.length < 6)                          e.password        = 'Mínimo 6 caracteres';
    if (form.password !== form.confirmPassword)            e.confirmPassword = 'Las contraseñas no coinciden';
    if (!form.firstName.trim() || form.firstName.length < 2) e.firstName     = 'Mínimo 2 caracteres';
    if (!form.lastName.trim()  || form.lastName.length  < 2) e.lastName      = 'Mínimo 2 caracteres';
    return e;
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev)   => ({ ...prev, [name]: value }));
    setErrors((prev) => ({ ...prev, [name]: '' }));
    setApiErr('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length) { setErrors(errs); return; }
    const res = await register(form);
    if (res.ok) navigate('/');
    else setApiErr(res.error);
  };

  /* Props comunes que Field necesita */
  const fieldProps = { form, errors, onChange: handleChange };

  return (
    <div className="min-h-[calc(100vh-56px)] flex items-center justify-center bg-corporate-50 px-4 py-12">
      <div className="card animate-fade-in w-full max-w-lg p-8">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-bold text-corporate-900">Crear cuenta</h1>
          <p className="text-sm text-corporate-500 mt-1">Completa el formulario para registrarte</p>
        </div>

        {apiErr && <div className="alert-error mb-4">{apiErr}</div>}

        <form onSubmit={handleSubmit} noValidate className="space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Field id="reg-firstName" name="firstName" label="Nombre"   placeholder="Juan"  {...fieldProps} />
            <Field id="reg-lastName"  name="lastName"  label="Apellido" placeholder="Pérez" {...fieldProps} />
          </div>

          <Field id="reg-username" name="username" label="Usuario"
            placeholder="juan_perez" autoComplete="username" {...fieldProps} />

          <Field id="reg-email" name="email" label="Correo electrónico"
            type="email" placeholder="juan@empresa.com" autoComplete="email" {...fieldProps} />

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Field id="reg-password" name="password" type="password"
              label="Contraseña" placeholder="••••••" autoComplete="new-password" {...fieldProps} />
            <Field id="reg-confirmPassword" name="confirmPassword" type="password"
              label="Confirmar contraseña" placeholder="••••••" autoComplete="new-password" {...fieldProps} />
          </div>

          {/* Admin secret (opcional) */}
          <div className="border-t border-corporate-100 pt-4">
            <label htmlFor="reg-adminSecret" className="label">
              Código de administrador{' '}
              <span className="text-corporate-400 font-normal">(opcional)</span>
            </label>
            <input
              id="reg-adminSecret"
              name="adminSecret"
              type="password"
              placeholder="Déjalo vacío si eres cliente"
              value={form.adminSecret}
              onChange={handleChange}
              className="input"
            />
            <p className="text-[11px] text-corporate-400 mt-1">
              Solo si tienes el código de acceso administrativo.
            </p>
          </div>

          <button
            id="reg-submit-btn"
            type="submit"
            disabled={loading}
            className="btn btn-primary w-full justify-center mt-2"
          >
            {loading ? 'Creando cuenta…' : 'Registrarse'}
          </button>
        </form>

        <p className="text-center text-sm text-corporate-500 mt-6">
          ¿Ya tienes cuenta?{' '}
          <Link to="/login" className="text-primary-700 font-medium hover:underline">
            Iniciar sesión
          </Link>
        </p>
      </div>
    </div>
  );
}
