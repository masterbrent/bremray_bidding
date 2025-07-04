import { writable, get } from 'svelte/store';
import type { Customer } from '../types/models';
import { api } from '../api/client';

interface CustomersState {
  customers: Customer[];
  loading: boolean;
  error: string | null;
}

function createCustomersStore() {
  const { subscribe, set, update } = writable<CustomersState>({
    customers: [],
    loading: false,
    error: null
  });

  return {
    subscribe,
    
    // Load all customers
    async load() {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const customers = await api.get<Customer[]>('/customers');
        set({ customers, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load customers'
        }));
      }
    },

    // Get customer by ID
    async getById(id: string): Promise<Customer | undefined> {
      const state = get(customersStore);
      let customer = state.customers.find(c => c.id === id);
      
      if (!customer) {
        try {
          customer = await api.get<Customer>(`/customers/${id}`);
          update(state => ({
            ...state,
            customers: [...state.customers, customer as Customer]
          }));
        } catch (error) {
          console.error('Failed to load customer:', error);
        }
      }
      
      return customer;
    },
    
    // Load a single customer by ID (alias for getById for consistency)
    async loadById(id: string) {
      return this.getById(id);
    },

    // Create a new customer
    async create(name: string, email: string, phone?: string) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const customer = await api.post<Customer>('/customers', { name, email, phone });
        update(state => ({
          ...state,
          customers: [...state.customers, customer],
          loading: false,
          error: null
        }));
        return customer;
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to create customer'
        }));
        throw error;
      }
    },

    // Update customer
    async update(id: string, updates: Partial<Customer>) {
      try {
        const customer = await api.put<Customer>(`/customers/${id}`, updates);
        update(state => ({
          ...state,
          customers: state.customers.map(c => c.id === id ? customer : c)
        }));
        return customer;
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to update customer'
        }));
        throw error;
      }
    },

    // Delete customer
    async remove(id: string) {
      try {
        await api.delete(`/customers/${id}`);
        update(state => ({
          ...state,
          customers: state.customers.filter(c => c.id !== id)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to delete customer'
        }));
        throw error;
      }
    }
  };
}

export const customersStore = createCustomersStore();