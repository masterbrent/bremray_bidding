import { writable, derived } from 'svelte/store';
import type { JobTemplate, JobTemplateItem } from '../types/models';
import { api } from '../api/client';

interface TemplatesState {
  templates: JobTemplate[];
  loading: boolean;
  error: string | null;
}

// Template items need to match backend format
interface TemplateItemRequest {
  itemId: string;
  defaultQuantity: number;
}

function createJobTemplatesStore() {
  const { subscribe, set, update } = writable<TemplatesState>({
    templates: [],
    loading: false,
    error: null
  });

  return {
    subscribe,
    
    // Load all templates
    async load() {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const templates = await api.get<JobTemplate[]>('/templates');
        set({ templates, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load templates'
        }));
      }
    },

    // Load only active templates
    async loadActive() {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const templates = await api.get<JobTemplate[]>('/templates?active=true');
        set({ templates, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load templates'
        }));
      }
    },

    // Create a new template
    async add(name: string, description: string, items: TemplateItemRequest[], phases: { name: string; order: number; description?: string }[] = []) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const template = await api.post<JobTemplate>('/templates', {
          name,
          description,
          items,
          phases
        });
        
        update(state => ({
          ...state,
          templates: [...state.templates, template],
          loading: false,
          error: null
        }));
        
        return template;
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to create template'
        }));
        throw error;
      }
    },

    // Update template
    async update(id: string, updates: { name?: string; description?: string; isActive?: boolean }) {
      try {
        const template = await api.put<JobTemplate>(`/templates/${id}`, updates);
        update(state => ({
          ...state,
          templates: state.templates.map(t => t.id === id ? template : t)
        }));
        return template;
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to update template'
        }));
        throw error;
      }
    },

    // Add item to template
    async addItemToTemplate(templateId: string, itemId: string, defaultQuantity: number) {
      try {
        await api.post(`/templates/${templateId}/items`, {
          itemId,
          defaultQuantity
        });
        
        // Reload the template to get updated data
        const template = await api.get<JobTemplate>(`/templates/${templateId}`);
        update(state => ({
          ...state,
          templates: state.templates.map(t => t.id === templateId ? template : t)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to add item to template'
        }));
        throw error;
      }
    },

    // Update template item
    async updateItem(templateId: string, itemId: string, defaultQuantity: number) {
      try {
        await api.put(`/templates/${templateId}/items/${itemId}`, {
          defaultQuantity
        });
        
        // Reload the template to get updated data
        const template = await api.get<JobTemplate>(`/templates/${templateId}`);
        update(state => ({
          ...state,
          templates: state.templates.map(t => t.id === templateId ? template : t)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to update template item'
        }));
        throw error;
      }
    },

    // Remove item from template
    async removeItemFromTemplate(templateId: string, itemId: string) {
      try {
        await api.delete(`/templates/${templateId}/items/${itemId}`);
        
        // Reload the template to get updated data
        const template = await api.get<JobTemplate>(`/templates/${templateId}`);
        update(state => ({
          ...state,
          templates: state.templates.map(t => t.id === templateId ? template : t)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to remove item from template'
        }));
        throw error;
      }
    },

    // Delete template
    async remove(id: string) {
      try {
        await api.delete(`/templates/${id}`);
        update(state => ({
          ...state,
          templates: state.templates.filter(t => t.id !== id)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to delete template'
        }));
        throw error;
      }
    }
  };
}

export const jobTemplatesStore = createJobTemplatesStore();

// Derived store for easy access to just the templates array
export const templates = derived(jobTemplatesStore, $store => $store.templates);