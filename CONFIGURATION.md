# Configuration Improvements

## No More Hardcoded URLs

1. **Frontend Configuration**:
   - All API calls now use `VITE_API_URL` environment variable
   - Default is `/api` which uses Vite's proxy in development
   - Single source of truth: `src/lib/config.ts`
   - No more hardcoded `localhost` anywhere

2. **Vite Proxy Configuration**:
   - Added proxy in `vite.config.ts` to forward `/api` requests to backend
   - This means frontend can use relative URLs like `/api/jobs`
   - In production, nginx or similar would handle this routing

3. **Environment Variables**:
   - Frontend: `.env` file with `VITE_API_URL=/api`
   - Backend: `.env` file with all service configurations
   - Both have `.env.example` files for documentation

## DRY Improvements

1. **Single API Configuration**:
   - Removed duplicate `API_BASE_URL` definitions
   - All services import from `lib/config.ts`
   - API client, photo service, and health service all use same config

2. **User Initialization**:
   - Moved `userStore.init()` to App.svelte
   - Only initialized once when app loads
   - No more duplicate initialization in JobsPage

3. **Consistent API Client**:
   - All API calls go through `lib/api/client.ts`
   - Centralized error handling
   - Consistent response handling

## Production Ready

1. **No Hardcoded Values**:
   - All configuration comes from environment variables
   - Frontend uses relative URLs
   - Backend reads all secrets from .env

2. **Flexible Deployment**:
   - Frontend can be served from same domain as API
   - Or can be configured to use separate API domain
   - Just change `VITE_API_URL` for different environments

3. **Security**:
   - Secrets only in backend .env
   - Frontend only knows API endpoint
   - CORS configured through environment variables
