import { writable, derived } from 'svelte/store';
import type { Item } from '../types/models';
import { itemsApi } from '../api/items';

interface ItemsState {
  items: Item[];
  loading: boolean;
  error: string | null;
}

function createItemsStore() {
  const { subscribe, set, update } = writable<ItemsState>({
    items: [],
    loading: false,
    error: null,
  });

  return {
    subscribe,
    
    // Load all items from API
    async load() {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const items = await itemsApi.getAll();
        set({ items, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load items',
        }));
      }
    },

    // Add new item
    async add(item: Omit<Item, 'id' | 'createdAt' | 'updatedAt'>) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const newItem = await itemsApi.create(item);
        update(state => ({
          items: [...state.items, newItem],
          loading: false,
          error: null,
        }));
        return newItem;
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to create item',
        }));
        throw error;
      }
    },

    // Update existing item
    async update(id: string, updates: Partial<Item>) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const updatedItem = await itemsApi.update(id, updates);
        update(state => ({
          items: state.items.map(item => 
            item.id === id ? updatedItem : item
          ),
          loading: false,
          error: null,
        }));
        return updatedItem;
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to update item',
        }));
        throw error;
      }
    },

    // Remove item
    async remove(id: string) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        await itemsApi.delete(id);
        update(state => ({
          items: state.items.filter(item => item.id !== id),
          loading: false,
          error: null,
        }));
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to delete item',
        }));
        throw error;
      }
    },

    // Clear error
    clearError() {
      update(state => ({ ...state, error: null }));
    },
  };
}

export const itemsStore = createItemsStore();

// Derived store for just the items array (for backward compatibility)
export const items = derived(itemsStore, $itemsStore => $itemsStore.items);