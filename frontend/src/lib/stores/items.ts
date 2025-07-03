import { writable } from 'svelte/store';
import type { Item } from '../types/models';

function createItemsStore() {
  const { subscribe, set, update } = writable<Item[]>([]);

  return {
    subscribe,
    add: (item: Omit<Item, 'id' | 'createdAt' | 'updatedAt'>) => {
      update(items => {
        const newItem: Item = {
          ...item,
          id: crypto.randomUUID(),
          createdAt: new Date(),
          updatedAt: new Date()
        };
        return [...items, newItem];
      });
    },
    update: (id: string, updates: Partial<Omit<Item, 'id' | 'createdAt'>>) => {
      update(items => items.map(item => 
        item.id === id 
          ? { ...item, ...updates, updatedAt: new Date() }
          : item
      ));
    },
    remove: (id: string) => {
      update(items => items.filter(item => item.id !== id));
    },
    loadItems: (items: Item[]) => {
      set(items);
    }
  };
}

export const itemsStore = createItemsStore();