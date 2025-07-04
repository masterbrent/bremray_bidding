# Role-Based Access Control

The Bremray Electrical Bidding App now has a role-based access control system with two roles:

## Roles

### Admin
- Can see all dollar values and prices
- Can create, edit, and delete jobs
- Can manage items and templates
- Can assign roles to other users
- Can take photos
- Can edit quantities
- Has access to all pages

### Tech
- **Cannot** see any dollar values or prices
- **Cannot** create, edit, or delete jobs
- **Can** view job cards and job details
- **Can** edit quantities on jobs
- **Can** take photos
- Limited navigation (only Jobs, Stats, and Settings)

## Master Admin

The master admin is hardcoded as `brenthall@gmail.com`. This account:
- Always has admin privileges
- Cannot be demoted to tech role
- Cannot be removed from the system

## Admin Testing Feature

Admins have a hidden toggle button to switch to "Tech View" for testing purposes:
- Located in the bottom-right corner of the Jobs page
- Semi-transparent icon that becomes visible on hover
- Shield icon = currently in tech view
- Users icon = currently in admin view
- The view preference persists across sessions

## User Management

Admins can manage users through the Users page:
- Add new users with email and role
- Change user roles (except master admin)
- Remove users (except master admin)

## Implementation Details

### Hidden UI Elements for Techs:
- "New Job" button
- Delete buttons on job cards
- All price/amount displays
- Items, Templates, and Users navigation items

### Visible to Techs:
- Job cards (without prices)
- Job details (with quantity editing)
- Photo capture and gallery
- Basic job information (customer, address, status)

### View Mode Banner
When an admin switches to tech view, a yellow banner appears at the top of the Jobs page indicating "Viewing as Tech" mode.

## Future Enhancements

Currently, users are managed in the frontend. In production:
1. User authentication would be handled by a proper auth system
2. Roles would be stored in the backend database
3. API endpoints would check permissions before allowing actions
4. The user's role would be determined during login
