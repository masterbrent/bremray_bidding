export interface Item {
  id: string;
  name: string;
  description?: string;
  unit: 'each' | 'ft' | 'hr' | 'lot';
  unitPrice: number;
  category?: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface JobPhase {
  id: string;
  name: string;
  order: number;
  isCompleted: boolean;
  completedAt?: Date;
}

export interface JobTemplateItem {
  itemId: string;
  defaultQuantity: number;
}

export interface JobTemplate {
  id: string;
  name: string;
  description?: string;
  items: JobTemplateItem[];
  phases: JobPhase[];
  createdAt: Date;
  updatedAt: Date;
}

export interface JobItem {
  id: string;
  itemId: string;
  item: Item;
  quantity: number;
  installedQuantity: number;
  notes?: string;
}

export interface Customer {
  id: string;
  name: string;
  phone?: string;
  email?: string;
}

export interface Job {
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
  status: 'pending' | 'in-progress' | 'completed' | 'cancelled';
  photos?: string[]; // Array of photo URLs
  waveInvoiceId?: string; // Wave invoice ID if sent
  waveInvoiceUrl?: string; // Wave invoice URL
  createdAt: Date;
  updatedAt: Date;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  zip: string;
  formatted: string;
}