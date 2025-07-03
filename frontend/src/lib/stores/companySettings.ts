import { writable } from 'svelte/store';

export interface CompanySettings {
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
}

function createCompanySettingsStore() {
  const defaultSettings: CompanySettings = {
    name: 'Bremray Electrical',
    logo: null,
    address: '',
    city: '',
    state: '',
    zip: '',
    phone: '',
    email: '',
    license: '',
    website: ''
  };

  const { subscribe, set, update } = writable<CompanySettings>(defaultSettings);

  return {
    subscribe,
    updateSettings: (settings: Partial<CompanySettings>) => {
      update(current => ({ ...current, ...settings }));
    },
    uploadLogo: (logoUrl: string) => {
      update(current => ({ ...current, logo: logoUrl }));
    },
    removeLogo: () => {
      update(current => ({ ...current, logo: null }));
    },
    reset: () => {
      set(defaultSettings);
    }
  };
}

export const companySettingsStore = createCompanySettingsStore();