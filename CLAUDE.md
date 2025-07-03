# Bremray Electrical Bidding App - Developer Guide for Claude

## Project Overview
This is an electrical contractor bidding application built for Bremray Electrical. The app helps field electricians create job estimates, track work progress, and manage invoicing.

## Technology Stack
- **Frontend**: Svelte + TypeScript + Vite
- **Backend**: Go (Golang)
- **Database**: PostgreSQL (planned)
- **Deployment**: Digital Ocean
- **Storage**: Cloudflare R2 (for photos)
- **Payment Integration**: Wave (planned)

## Architecture Principles
1. **Complete separation of concerns**: Frontend is purely for rendering, backend handles all business logic
2. **RESTful API design**: Backend provides JSON APIs consumed by frontend
3. **Stateless architecture**: Use JWT tokens for authentication (no sessions)
4. **Mock-first development**: Frontend uses mock services until backend is ready

## Current State (as of last update)

### Frontend Status
- ✅ Complete UI implementation with finance dashboard theme (blue accents)
- ✅ Job management with cards showing all required info
- ✅ Photo capture and gallery functionality
- ✅ Company settings page
- ✅ Mock data and services for all features
- ⚠️ Contains hardcoded mock data that needs backend integration

### Backend Status
- 🚧 Not yet implemented
- 📋 Planned structure follows Go best practices

### Key Features Implemented (Frontend Only)
1. **Jobs Page**: List view with job cards showing template, permit status, dates, amount
2. **Job Details**: "What we did" section with item tracking and photo gallery
3. **Items Management**: CRUD for inventory items
4. **Job Templates**: Reusable job configurations
5. **Photo Management**: Camera capture, upload, gallery with multi-select
6. **Company Settings**: Logo, contact info, license details

## Design Decisions Made

### UI/UX
- Light theme with blue primary colors (#3b82f6)
- Modern finance dashboard aesthetic
- Rounded cards and minimal shadows
- Company name: "Bremray" (not "ElectriBid")
- Lightning bolt favicon

### Frontend Architecture
```
src/
├── lib/
│   ├── stores/          # Svelte stores for state management
│   ├── components/      # Reusable UI components
│   ├── services/        # Mock services (to be replaced with API calls)
│   ├── types/           # TypeScript interfaces
│   └── router.ts        # Simple routing solution
└── pages/              # Page components
```

### Data Models (TypeScript interfaces to guide backend)
```typescript
interface Job {
  id: string;
  customer: Customer;
  address: string;
  templateId: string;
  template: JobTemplate;
  items: JobItem[];
  phases: JobPhase[];
  requiresPermit: boolean;
  startDate?: Date;
  endDate?: Date;
  status: 'pending' | 'in-progress' | 'completed';
  photos: string[];
  waveInvoiceId?: string;
  waveInvoiceUrl?: string;
  createdAt: Date;
  updatedAt: Date;
}

interface Item {
  id: string;
  name: string;
  unit: string;
  unitPrice: number;
  category: string;
}

interface JobTemplate {
  id: string;
  name: string;
  description: string;
  items: TemplateItem[];
  phases: TemplatePhase[];
}
```

## Backend Implementation Guide

### Recommended Go Project Structure
```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── models/              # Database models
│   ├── handlers/            # HTTP handlers (controllers)
│   ├── services/            # Business logic
│   ├── repository/          # Database access layer
│   └── middleware/          # Auth, CORS, logging
├── pkg/
│   └── utils/               # Shared utilities
├── migrations/              # SQL migrations
└── config/                  # Configuration files
```

### API Endpoints to Implement
```
Authentication:
POST   /api/auth/login
POST   /api/auth/refresh
POST   /api/auth/logout

Jobs:
GET    /api/jobs              # List all jobs
POST   /api/jobs              # Create job from template
GET    /api/jobs/:id          # Get job details
PUT    /api/jobs/:id          # Update job
DELETE /api/jobs/:id          # Delete job
POST   /api/jobs/:id/photos   # Upload photos
DELETE /api/jobs/:id/photos   # Delete photos

Items:
GET    /api/items             # List all items
POST   /api/items             # Create item
PUT    /api/items/:id         # Update item
DELETE /api/items/:id         # Delete item

Templates:
GET    /api/templates         # List all templates
POST   /api/templates         # Create template
PUT    /api/templates/:id     # Update template
DELETE /api/templates/:id     # Delete template

Company:
GET    /api/company/settings  # Get company settings
PUT    /api/company/settings  # Update company settings
POST   /api/company/logo      # Upload logo
```

### Database Schema (PostgreSQL)
```sql
-- Core tables needed
companies
users
items
job_templates
template_items
template_phases
jobs
job_items
job_phases
job_photos
customers
```

## Migration Strategy (Frontend → Backend)

### Phase 1: Backend Setup
1. Initialize Go module and project structure
2. Set up PostgreSQL database
3. Create models matching TypeScript interfaces
4. Implement basic CRUD endpoints
5. Add CORS middleware for local development

### Phase 2: Frontend Integration
1. Create API service layer in frontend
2. Replace mock store functions with API calls
3. Add error handling and loading states
4. Implement authentication flow

### Phase 3: Deployment Prep
1. Environment configuration
2. Docker containerization
3. Database migrations setup
4. CI/CD pipeline for Digital Ocean

## Important Implementation Notes

### Frontend Cleanup Required
- Remove all mock data from stores
- Convert store functions to API calls
- Add proper error handling
- Implement loading states
- Add authentication checks

### Security Considerations
- JWT token authentication
- CORS configuration for production
- Input validation on all endpoints
- SQL injection prevention
- XSS protection in frontend

### Performance Considerations
- Implement pagination for lists
- Add caching where appropriate
- Optimize photo uploads/storage
- Use connection pooling for database

## Common Commands

### Frontend Development
```bash
cd frontend
npm run dev          # Start dev server
npm run build        # Build for production
npm run preview      # Preview production build
```

### Backend Development (once implemented)
```bash
cd backend
go mod download      # Install dependencies
go run cmd/server/main.go  # Run server
go test ./...        # Run tests
```

## Known Issues / TODOs
1. Photo service uses blob URLs (needs R2 integration)
2. Wave invoice integration not implemented
3. No authentication system yet
4. All data is client-side only
5. No data validation on forms

## Questions for Product Owner
1. User roles and permissions model?
2. Multi-company support needed?
3. Offline functionality requirements?
4. Backup and data retention policies?
5. Invoice numbering scheme?

## Progress Documentation Requirements

### IMPORTANT: Update PROGRESS.md Regularly
Every Claude instance working on this project MUST update the PROGRESS.md file:

#### When to Update:
1. **Start of Session**: Review current progress and plan tasks
2. **During Session**: Mark tasks as completed/in-progress as you work
3. **End of Session**: Update metrics, add session notes, plan next steps

#### What to Update:
- Task completion status (mark with ✅ when done)
- Phase progress percentages
- New issues discovered
- Code statistics/metrics
- Session log with work completed
- Next steps for the following session

#### How to Update:
```bash
# At start of session
1. Read PROGRESS.md first
2. Review what was done last session
3. Plan current session tasks

# During work
4. Mark tasks as you complete them
5. Add new tasks if discovered
6. Update completion percentages

# Before ending session
7. Add session summary
8. Update metrics
9. Document any blockers
10. Set clear next steps
```

### Example Update Pattern:
```markdown
### Session 3 (Date)
- Completed backend structure setup ✅
- Implemented Jobs API endpoints ✅  
- Started authentication system 🚧
- **Blocker**: Need clarification on user roles
- **Next**: Complete auth and start frontend integration
```

## Final Notes
- Always maintain separation between frontend and backend
- Follow RESTful conventions for API design
- Keep TypeScript interfaces in sync with Go models
- Test API endpoints before integrating with frontend
- Document any deviations from this plan
- **UPDATE PROGRESS.MD AFTER EVERY WORK SESSION**

This project is designed to be maintainable and scalable. Stick to the established patterns and conventions for consistency.