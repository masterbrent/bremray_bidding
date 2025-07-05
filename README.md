# Bremray Electrical Bidding App

A web application for managing electrical job bidding, tracking, and invoicing.

## Tech Stack

- **Frontend**: SvelteKit, TypeScript, Tailwind CSS
- **Backend**: Go, PostgreSQL
- **Storage**: Cloudflare R2
- **Invoicing**: Wave API Integration

## Setup

1. Clone the repository
2. Copy `.env.example` files to `.env` in both `backend/` and `frontend/` directories
3. Configure environment variables
4. Install dependencies:
   ```bash
   # Backend
   cd backend
   go mod download
   
   # Frontend
   cd frontend
   npm install
   ```
5. Run database migrations
6. Start the application:
   ```bash
   ./start-all.sh
   ```

## Features

- Job management and tracking
- Photo upload and gallery
- Customer management
- Invoice generation via Wave
- Role-based access control
- Template-based job creation

## License

Private repository - All rights reserved
