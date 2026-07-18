# E-commerce API - Grupo 5

API RESTful de e-commerce desarrollada para la tarea integradora de **Aplicaciones Web (EPN)**.

**Stack:** Go + Gin + GORM + PostgreSQL + JWT + Swagger

## Descripcion

Este backend implementa el mismo dominio funcional del laboratorio de referencia:

- **User:** registro, login, consulta, actualizacion y eliminacion
- **Product:** CRUD completo con manejo de stock
- **Receipt / ReceiptItem:** compras con calculo de total en backend y descuento automatico de stock

## Arquitectura

```
cmd/api/              -> Punto de entrada
internal/
  config/             -> Variables de entorno
  database/           -> Conexion y migraciones GORM
  models/             -> Entidades de dominio
  dto/                -> Objetos de entrada/salida
  repository/         -> Acceso a datos
  service/            -> Reglas de negocio
  handlers/           -> Controladores HTTP
  middleware/         -> JWT y CORS
  auth/               -> Generacion y validacion JWT
  apperrors/          -> Manejo centralizado de errores
  router/             -> Rutas REST
docs/                 -> Documentacion Swagger generada
postman/              -> Coleccion de pruebas
scripts/              -> Script inicial de base de datos
```

## Requisitos

- Go 1.22+
- Docker (opcional, para PostgreSQL)
- PostgreSQL 16+

## Configuracion

1. Copiar variables de entorno:

```bash
cp .env.example .env
```

2. Levantar PostgreSQL con Docker:

```bash
docker compose up -d
```

3. Instalar dependencias:

```bash
go mod tidy
```

4. Generar documentacion Swagger:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs
```

5. Ejecutar la API:

```bash
go run ./cmd/api
```

La API quedara disponible en `http://localhost:8080`.

## Variables de entorno

| Variable | Descripcion | Valor por defecto |
|---|---|---|
| `APP_PORT` | Puerto HTTP | `8080` |
| `APP_ENV` | Entorno (`development` / `production`) | `development` |
| `DB_HOST` | Host PostgreSQL | `localhost` |
| `DB_PORT` | Puerto PostgreSQL | `5432` |
| `DB_USER` | Usuario PostgreSQL | `postgres` |
| `DB_PASSWORD` | Contrasena PostgreSQL | `postgres` |
| `DB_NAME` | Nombre de base de datos | `ecommerce` |
| `DB_SSLMODE` | Modo SSL PostgreSQL | `disable` |
| `JWT_SECRET` | Clave secreta JWT | valor de desarrollo |
| `JWT_EXPIRATION_HOURS` | Expiracion del token | `24` |

## Endpoints

### Users

| Metodo | Ruta | Auth | Descripcion |
|---|---|---|---|
| POST | `/api/users/register` | No | Registrar usuario |
| POST | `/api/users/login` | No | Iniciar sesion |
| GET | `/api/users/{id}` | No | Consultar usuario |
| PUT | `/api/users/{id}` | JWT | Actualizar usuario |
| DELETE | `/api/users/{id}` | JWT | Eliminar usuario |

### Products

| Metodo | Ruta | Auth | Descripcion |
|---|---|---|---|
| POST | `/api/products` | JWT | Crear producto |
| GET | `/api/products` | No | Listar productos |
| GET | `/api/products/{id}` | No | Consultar producto |
| PUT | `/api/products/{id}` | JWT | Actualizar producto |
| DELETE | `/api/products/{id}` | JWT | Eliminar producto |

### Receipts

| Metodo | Ruta | Auth | Descripcion |
|---|---|---|---|
| POST | `/api/receipts` | JWT | Crear recibo/compra |
| GET | `/api/receipts` | JWT | Listar recibos |
| GET | `/api/receipts/{id}` | JWT | Consultar recibo |
| GET | `/api/receipts/user/{userId}` | JWT | Listar recibos por usuario |
| DELETE | `/api/receipts/{id}` | JWT | Eliminar recibo |

## Reglas de negocio implementadas

- El total de la compra se calcula en backend usando el precio real del producto en base de datos
- Se valida stock disponible antes de confirmar un recibo
- Al crear un recibo, el stock se descuenta automaticamente dentro de una transaccion
- Las contrasenas se almacenan cifradas con BCrypt y nunca se devuelven en respuestas HTTP
- Los montos monetarios usan `decimal.Decimal` para evitar errores de precision
- Los errores se manejan de forma centralizada con codigos HTTP y mensajes consistentes

## Documentacion Swagger

- URL: `http://localhost:8080/swagger/index.html`
- Para endpoints protegidos, usar el boton **Authorize** con: `Bearer <token>`

## Pruebas con Postman

Importar la coleccion:

`postman/Ecommerce-Grupo5.postman_collection.json`

Flujo sugerido:

1. Register o Login (guarda token automaticamente)
2. Create Product
3. Create Receipt
4. List Receipts / List Receipts By User

## Health check

`GET http://localhost:8080/health`

## Integrantes

Grupo 5 - Aplicaciones Web - EPN
