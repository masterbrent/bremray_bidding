import { writable, derived } from 'svelte/store';

type UserRole = 'admin' | 'tech';

interface User {
  email: string;
  role: UserRole;
  isViewingAsTech?: boolean; // For admins to toggle tech view
}

// Create the user store
function createUserStore() {
  const { subscribe, set, update } = writable<User | null>(null);

  // Check localStorage for saved view preference
  const savedViewMode = typeof window !== 'undefined' 
    ? localStorage.getItem('adminViewMode') 
    : null;

  return {
    subscribe,
    
    // Initialize user (this would come from your auth system)
    init: (email: string) => {
      // Master admin
      if (email === 'brenthall@gmail.com') {
        set({
          email,
          role: 'admin',
          isViewingAsTech: savedViewMode === 'tech'
        });
      } else {
        // Default to tech for now - you'll assign roles separately
        set({
          email,
          role: 'tech'
        });
      }
    },

    // Toggle between admin and tech view (only for admins)
    toggleViewMode: () => {
      update(user => {
        if (!user || user.role !== 'admin') return user;
        
        const newViewMode = !user.isViewingAsTech;
        
        // Save preference to localStorage
        if (typeof window !== 'undefined') {
          localStorage.setItem('adminViewMode', newViewMode ? 'tech' : 'admin');
        }
        
        return {
          ...user,
          isViewingAsTech: newViewMode
        };
      });
    },

    // Set user role (admin only function)
    setUserRole: (email: string, role: UserRole) => {
      // This would typically update the backend
      console.log(`Setting ${email} role to ${role}`);
    },

    // Logout
    logout: () => {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('adminViewMode');
      }
      set(null);
    }
  };
}

export const userStore = createUserStore();

// Derived store for effective role (considering view mode)
export const effectiveRole = derived(
  userStore,
  $user => {
    if (!$user) return null;
    
    // If admin is viewing as tech, return tech role
    if ($user.role === 'admin' && $user.isViewingAsTech) {
      return 'tech';
    }
    
    return $user.role;
  }
);

// Helper functions for permission checks
export const permissions = {
  canSeePrices: (role: UserRole | null) => role === 'admin',
  canEditJobs: (role: UserRole | null) => role === 'admin',
  canDeleteJobs: (role: UserRole | null) => role === 'admin',
  canCreateJobs: (role: UserRole | null) => role === 'admin',
  canEditQuantities: (role: UserRole | null) => role === 'admin' || role === 'tech',
  canTakePhotos: (role: UserRole | null) => role === 'admin' || role === 'tech',
};
