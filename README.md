# Autohost Cloud API

A REST API for managing distributed nodes and agents, built with Go and PostgreSQL.

## Features

- 🔐 User authentication with JWT tokens
- 🖥️ Node/server registration and management
- 🎫 Enrollment token system for secure agent onboarding
- 💓 Heartbeat monitoring
- 🔄 Refresh token rotation
- 📊 PostgreSQL database with migrations

## Tech Stack

- **Language:** Go 1.21+
- **Database:** PostgreSQL 16
- **Router:** Chi v5
- **Migrations:** golang-migrate
- **Live reload:** Air

## Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Make (optional, but recommended)

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/mazapanuwu13/autohost-cloud-api.git
cd autohost-cloud-api
```

### 2. Set up environment variables

Create a `.env` file in the root directory:

```bash
PORT=8080
DATABASE_URL=postgres://autohost:autohost@localhost:5432/autohost?sslmode=disable
JWT_SECRET=your_super_secret_key_change_this_in_production
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=720h
FRONTEND_URL=http://localhost:3000
ENV=development
```

### 3. Start PostgreSQL with Docker Compose

```bash
docker compose up -d
```

This will start:
- PostgreSQL on port `5432`
- pgAdmin on port `5050` (optional)

### 4. Run database migrations

```bash
make migrate-up
```

### 5. Install Go dependencies

```bash
go mod download
```

### 6. Run the application

**Development mode with live reload:**
```bash
make run-air
```

**Production mode:**
```bash
make run
```

The API will be available at `http://localhost:8080`

## Project Structure

```
.
├── cmd/
│   ├── api/              # Main application entry point
│   └── migrate/          # Migration runner
├── internal/
│   ├── domain/           # Business logic & entities
│   │   ├── auth/
│   │   ├── node/
│   │   └── enrollment/
│   ├── repository/       # Data persistence layer
│   │   └── postgres/
│   ├── handler/          # HTTP handlers
│   │   └── middleware/
│   └── platform/         # Infrastructure utilities
├── migrations/           # Database migrations
├── request/              # HTTP request examples (.http files)
├── compose.yml           # Docker Compose configuration
└── Makefile             # Common commands
```

## API Endpoints

### Authentication

- `POST /v1/auth/register` - Register new user
- `POST /v1/auth/login` - Login user
- `POST /v1/auth/refresh` - Refresh access token
- `POST /v1/auth/logout` - Logout user
- `GET /v1/auth/me` - Get current user (requires auth)

### Nodes

- `POST /v1/nodes/register` - Register a new node (requires auth)
- `GET /v1/nodes` - List user's nodes (requires auth)

### Enrollment

- `POST /v1/enrollments/generate` - Generate enrollment token (requires auth)
- `POST /v1/enrollments/enroll` - Enroll node with token

## Testing API Endpoints

### Using REST Client (VS Code)

1. Install the REST Client extension:
```bash
code --install-extension humao.rest-client
```

2. Open `request/user.http` and click "Send Request" above each endpoint

### Using curl

See examples in `request/register_user.txt` or use the `.http` files

## Database Management

### pgAdmin

Access pgAdmin at `http://localhost:5050`:
- **Email:** `admin@localhost.com`
- **Password:** `devpass`

Connection settings:
- **Host:** `postgres`
- **Port:** `5432`
- **Database:** `autohost`
- **Username:** `autohost`
- **Password:** `autohost`

### Migrations

```bash
# Create new migration
make migrate-create name=your_migration_name

# Apply migrations
make migrate-up

# Rollback one migration
make migrate-down

# Check migration status
make migrate-version
```

### Direct database access

```bash
docker exec -it ahcl-postgres psql -U autohost -d autohost
```

## Available Make Commands

```bash
make run          # Run the application
make run-air      # Run with live reload
make build        # Build the binary
make test         # Run tests
make migrate-up   # Apply database migrations
make migrate-down # Rollback last migration
make migrate-create name=xyz  # Create new migration
```

## Development

### Project follows Clean Architecture:

1. **Domain Layer** (`internal/domain/`): Business logic, entities, and repository interfaces
2. **Repository Layer** (`internal/repository/`): Database implementations
3. **Handler Layer** (`internal/handler/`): HTTP request handling
4. **Platform Layer** (`internal/platform/`): Infrastructure utilities (JWT, password hashing, etc.)

### Adding a new feature:

1. Define domain entity and service in `internal/domain/`
2. Implement repository in `internal/repository/postgres/`
3. Create HTTP handler in `internal/handler/`
4. Register routes in `internal/handler/router.go`
5. Add migration if database changes needed

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
