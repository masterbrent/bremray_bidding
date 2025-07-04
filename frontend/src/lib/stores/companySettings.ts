import { writable, get } from 'svelte/store';
import { api } from '../api/client';

export interface CompanySettings {
  id: string;
  name: string;
  logo: string | null;
  address: string;
  city: string;
  state: string;
  zip: string;
  phone: string;
  email: string;
  license: string;
  website: string;
  createdAt: Date;
  updatedAt: Date;
}

interface CompanyState {
  settings: CompanySettings | null;
  loading: boolean;
  error: string | null;
}

function createCompanySettingsStore() {
  const { subscribe, set, update } = writable<CompanyState>({
    settings: null,
    loading: false,
    error: null
  });

  return {
    subscribe,
    
    // Load company settings from backend
    async load() {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const settings = await api.get<CompanySettings>('/company');
        set({ settings, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load company settings'
        }));
      }
    },
    
    // Update company settings
    async updateSettings(updates: Partial<CompanySettings>) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const settings = await api.put<CompanySettings>('/company', updates);
        set({ settings, loading: false, error: null });
        return settings;
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to update company settings'
        }));
        throw error;
      }
    },
    
    // Upload logo (placeholder - will be implemented with R2)
    async uploadLogo(logoUrl: string) {
      const state = get(companySettingsStore);
      if (!state.settings) return;
      
      return this.updateSettings({ ...state.settings, logo: logoUrl });
    },
    
    // Remove logo
    async removeLogo() {
      const state = get(companySettingsStore);
      if (!state.settings) return;
      
      return this.updateSettings({ ...state.settings, logo: '' });
    },
    
    // Get current settings (for backward compatibility)
    getSettings(): CompanySettings | null {
      const state = get(companySettingsStore);
      return state.settings;
    }
  };
}

export const companySettingsStore = createCompanySettingsStore();