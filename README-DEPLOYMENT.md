# Bremray Electrical Bidding App - Production Deployment

## Pre-deployment Checklist

### 1. Environment Configuration

#### Backend (.env)
```bash
# Required
PORT=8080
DATABASE_URL=postgres://user:pass@host:port/db?sslmode=require
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Optional
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

#### Frontend (.env)
```bash
VITE_API_URL=https://api.yourdomain.com/api
```

### 2. Security Checklist

- [ ] Change default database credentials
- [ ] Set strong database password
- [ ] Configure ALLOWED_ORIGINS for production domains only
- [ ] Enable SSL/TLS for database connections
- [ ] Use HTTPS for all traffic
- [ ] Review and remove any debug/test code
- [ ] Ensure .env files are not committed to git

### 3. Database Setup

Run migrations in order:
```bash
for file in backend/migrations/*.sql; do
  psql $DATABASE_URL -f "$file"
done
```

## Deployment Options

### Option 1: Digital Ocean App Platform

1. **Create Database**
   - Add a managed PostgreSQL database
   - Note the connection string

2. **Deploy Backend**
   - Source: GitHub repository
   - Branch: main
   - Source Directory: `/backend`
   - Build Command: `go build -o bin/server ./cmd/server`
   - Run Command: `./bin/server`
   - Environment Variables: Add all from backend .env
   - HTTP Port: 8080
   - Health Check Path: `/api/health`

3. **Deploy Frontend**
   - Source: Same GitHub repository
   - Branch: main
   - Source Directory: `/frontend`
   - Build Command: `npm install && npm run build`
   - Output Directory: `dist`
   - Environment Variables: Add VITE_API_URL

### Option 2: Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   # Create .env file with production values
   cp .env.example .env
   # Edit .env with production values
   
   # Deploy
   docker-compose -f docker-compose.prod.yml up -d
   ```

2. **With external database**
   ```bash
   # Only run the services, not the database
   docker-compose -f docker-compose.prod.yml up -d backend frontend
   ```

### Option 3: Traditional VPS

1. **Install dependencies**
   ```bash
   # Install Go, Node.js, PostgreSQL, Nginx
   ```

2. **Build backend**
   ```bash
   cd backend
   go build -o bin/server ./cmd/server
   ```

3. **Build frontend**
   ```bash
   cd frontend
   npm install
   npm run build
   ```

4. **Configure Nginx**
   ```nginx
   server {
       server_name yourdomain.com;
       
       location / {
           root /path/to/frontend/dist;
           try_files $uri $uri/ /index.html;
       }
   }
   
   server {
       server_name api.yourdomain.com;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

5. **Setup systemd service for backend**
   ```ini
   [Unit]
   Description=Bremray Backend
   After=network.target
   
   [Service]
   Type=simple
   User=appuser
   WorkingDirectory=/path/to/backend
   ExecStart=/path/to/backend/bin/server
   Restart=on-failure
   Environment="DATABASE_URL=..."
   
   [Install]
   WantedBy=multi-user.target
   ```

## Post-deployment

1. **Verify deployment**
   - Check health endpoint: `curl https://api.yourdomain.com/api/health`
   - Test frontend loads
   - Create a test job

2. **Setup monitoring**
   - Configure uptime monitoring
   - Setup error tracking (e.g., Sentry)
   - Configure database backups

3. **SSL/TLS**
   - Use Let's Encrypt for free SSL certificates
   - Configure auto-renewal

## Troubleshooting

### Backend won't start
- Check database connectivity
- Verify environment variables are set
- Check logs: `docker logs <container>` or systemd logs

### Frontend can't reach backend
- Verify VITE_API_URL is correct
- Check CORS configuration
- Ensure backend is accessible from frontend

### Database connection issues
- Verify DATABASE_URL format
- Check firewall rules
- Ensure SSL mode matches database configuration
