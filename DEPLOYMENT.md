# Deployment Configuration

## Environment Variables

### Backend Environment Variables

Create a `.env` file in the backend directory with these variables:

```bash
# Server Configuration
PORT=8080

# Database Configuration
DATABASE_URL=postgres://username:password@host:port/database?sslmode=require

# CORS Configuration
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
# For development: ALLOWED_ORIGINS=*

# Optional: Connection Pool Settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

### Frontend Environment Variables

Create a `.env` file in the frontend directory:

```bash
# API Configuration
VITE_API_URL=https://api.yourdomain.com/api
```

## Digital Ocean Deployment

### App Platform Configuration

1. **Backend Service**:
   - Type: Web Service
   - Environment: Go
   - Build Command: `go build -o bin/server ./cmd/server`
   - Run Command: `./bin/server`
   - HTTP Port: 8080
   - Health Check Path: `/api/health`

2. **Frontend Service**:
   - Type: Static Site
   - Environment: Node.js
   - Build Command: `npm install && npm run build`
   - Output Directory: `dist`

3. **Database**:
   - Add a PostgreSQL database
   - Run migrations after deployment

### Database Migrations

Run migrations in order:
```bash
psql $DATABASE_URL -f migrations/001_create_items_table.sql
psql $DATABASE_URL -f migrations/002_create_customers_table.sql
psql $DATABASE_URL -f migrations/003_create_job_templates_table.sql
psql $DATABASE_URL -f migrations/004_create_jobs_table.sql
psql $DATABASE_URL -f migrations/005_create_company_table.sql
psql $DATABASE_URL -f migrations/005_create_template_phases_table.sql
psql $DATABASE_URL -f migrations/006_add_item_nickname.sql
```

## Security Considerations

1. **CORS**: Update the backend to use environment-based CORS configuration
2. **Database**: Use SSL connections in production
3. **API Keys**: If adding external services, use environment variables
4. **HTTPS**: Ensure all traffic uses HTTPS in production

## Monitoring

- Backend health check: `/api/health`
- Consider adding application monitoring (e.g., Sentry, DataDog)
- Set up database backups
