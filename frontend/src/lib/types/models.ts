export interface Item {
  id: string;
  name: string;
  nickname?: string;  // Display name for job cards
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
  id: string;
  templateId: string;
  itemId: string;
  defaultQuantity: number;
}

export interface TemplatePhase {
  id: string;
  templateId: string;
  name: string;
  order: number;
  description?: string;
}

export interface JobTemplate {
  id: string;
  name: string;
  description: string;
  items: JobTemplateItem[];
  phases: TemplatePhase[];
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface JobItem {
  id: string;
  jobId: string;
  itemId: string;
  name: string;
  nickname?: string;
  quantity: number;
  price: number;
  total: number;
}

export interface Customer {
  id: string;
  name: string;
  phone?: string;
  email?: string;
}

export interface JobPhoto {
  id: string;
  jobId: string;
  url: string;
  caption?: string;
  uploadedAt: Date;
}

export interface Job {
  id: string;
  customerId: string;
  templateId: string;
  address: string;
  status: 'scheduled' | 'in_progress' | 'completed' | 'cancelled';
  currentPhaseId?: string;
  scheduledDate: Date;
  startDate?: Date;
  endDate?: Date;
  permitRequired: boolean;
  permitNumber?: string;
  totalAmount: number;
  items: JobItem[];
  photos: JobPhoto[];
  notes?: string;
  waveInvoiceId?: string;
  waveInvoiceUrl?: string;
  createdAt: Date;
  updatedAt: Date;
  // Frontend only - these will be loaded separately
  customer?: Customer;
  template?: JobTemplate;
  currentPhase?: TemplatePhase;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  zip: string;
  formatted: string;
}