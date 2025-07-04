import { api } from './client';
import type { Item } from '../types/models';

// Convert between backend snake_case and frontend camelCase
function toFrontendItem(backendItem: any): Item {
  return {
    id: backendItem.id,
    name: backendItem.name,
    nickname: backendItem.nickname,
    description: backendItem.description,
    unit: backendItem.unit,
    unitPrice: backendItem.unitPrice,
    category: backendItem.category,
    createdAt: new Date(backendItem.createdAt),
    updatedAt: new Date(backendItem.updatedAt),
  };
}

function toBackendItem(item: Partial<Item>): any {
  const backendItem: any = {};
  
  if (item.name !== undefined) backendItem.name = item.name;
  if (item.nickname !== undefined) backendItem.nickname = item.nickname;
  if (item.description !== undefined) backendItem.description = item.description;
  if (item.unit !== undefined) backendItem.unit = item.unit;
  if (item.unitPrice !== undefined) backendItem.unitPrice = item.unitPrice;
  if (item.category !== undefined) backendItem.category = item.category;
  
  return backendItem;
}

export const itemsApi = {
  async getAll(): Promise<Item[]> {
    const items = await api.get<any[]>('/items');
    return items.map(toFrontendItem);
  },

  async getById(id: string): Promise<Item> {
    const item = await api.get<any>(`/items/${id}`);
    return toFrontendItem(item);
  },

  async create(item: Omit<Item, 'id' | 'createdAt' | 'updatedAt'>): Promise<Item> {
    const backendItem = toBackendItem(item);
    const created = await api.post<any>('/items', backendItem);
    return toFrontendItem(created);
  },

  async update(id: string, updates: Partial<Item>): Promise<Item> {
    const backendUpdates = toBackendItem(updates);
    const updated = await api.put<any>(`/items/${id}`, backendUpdates);
    return toFrontendItem(updated);
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/items/${id}`);
  },
};