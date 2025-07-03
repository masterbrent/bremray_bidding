# Bremray Electrical Backend

## Overview
Go backend API for the Bremray Electrical bidding application.

## Architecture
- **Clean Architecture**: Separation of concerns with layers
- **Repository Pattern**: Database abstraction
- **Service Layer**: Business logic
- **TDD**: Test-Driven Development approach

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### Development Setup

1. Start the database:
```bash
make db-up
```

2. Run tests:
```bash
make test
```

3. Start the server:
```bash
make dev
```

The API will be available at `http://localhost:8080`

### API Endpoints

#### Items
- `GET /api/items` - List all items
- `POST /api/items` - Create new item
- `GET /api/items/:id` - Get item by ID
- `PUT /api/items/:id` - Update item
- `DELETE /api/items/:id` - Delete item

#### Health Check
- `GET /api/health` - Service health status

## Testing

Run all tests:
```bash
make test
```

Run with verbose output:
```bash
make test-verbose
```

Run with coverage:
```bash
make test-coverage
```

## Project Structure
```
backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── models/          # Domain models
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── repository/      # Data access
│   └── middleware/      # HTTP middleware
├── migrations/          # Database migrations
├── tests/              # Integration tests
└── pkg/                # Shared utilities
```

## Environment Variables
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string

## Database Migrations
Migrations are automatically applied when using docker-compose.

## Deployment
Build for production:
```bash
make build
```

The binary will be in `bin/server`.