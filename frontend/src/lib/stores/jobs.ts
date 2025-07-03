import { writable, get } from 'svelte/store';
import type { Job, JobItem, Customer, JobPhase } from '../types/models';
import { itemsStore } from './items';
import { jobTemplatesStore } from './jobTemplates';

function createJobsStore() {
  const { subscribe, set, update } = writable<Job[]>([]);

  return {
    subscribe,
    createFromTemplate: (
      customer: Customer,
      address: string,
      templateId: string,
      requiresPermit: boolean,
      startDate?: Date,
      endDate?: Date
    ) => {
      const templates = get(jobTemplatesStore);
      const items = get(itemsStore);
      const template = templates.find(t => t.id === templateId);
      
      if (!template) {
        throw new Error('Template not found');
      }

      const jobItems: JobItem[] = template.items.map(templateItem => {
        const item = items.find(i => i.id === templateItem.itemId);
        if (!item) throw new Error(`Item ${templateItem.itemId} not found`);
        
        return {
          id: crypto.randomUUID(),
          itemId: templateItem.itemId,
          item,
          quantity: templateItem.defaultQuantity,
          installedQuantity: 0
        };
      });

      const jobPhases: JobPhase[] = template.phases.map(phase => ({
        ...phase,
        id: crypto.randomUUID(),
        isCompleted: false,
        completedAt: undefined
      }));

      const newJob: Job = {
        id: crypto.randomUUID(),
        customer,
        address,
        templateId,
        template,
        items: jobItems,
        phases: jobPhases,
        requiresPermit,
        startDate,
        endDate,
        status: 'pending',
        photos: [
          'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800',
          'https://images.unsplash.com/photo-1565608087341-404b25492fee?w=800',
          'https://images.unsplash.com/photo-1621905251189-08b45d6a269e?w=800'
        ],
        createdAt: new Date(),
        updatedAt: new Date()
      };

      update(jobs => [...jobs, newJob]);
      return newJob.id;
    },
    updateItemQuantity: (jobId: string, jobItemId: string, installedQuantity: number) => {
      update(jobs => jobs.map(job => {
        if (job.id === jobId) {
          return {
            ...job,
            items: job.items.map(item => 
              item.id === jobItemId 
                ? { ...item, installedQuantity }
                : item
            ),
            updatedAt: new Date()
          };
        }
        return job;
      }));
    },
    updatePhase: (jobId: string, phaseId: string, isCompleted: boolean) => {
      update(jobs => jobs.map(job => {
        if (job.id === jobId) {
          return {
            ...job,
            phases: job.phases.map(phase => 
              phase.id === phaseId 
                ? { 
                    ...phase, 
                    isCompleted,
                    completedAt: isCompleted ? new Date() : undefined
                  }
                : phase
            ),
            updatedAt: new Date()
          };
        }
        return job;
      }));
    },
    updateStatus: (jobId: string, status: Job['status']) => {
      update(jobs => jobs.map(job => 
        job.id === jobId 
          ? { ...job, status, updatedAt: new Date() }
          : job
      ));
    },
    togglePermit: (jobId: string) => {
      update(jobs => jobs.map(job => 
        job.id === jobId 
          ? { ...job, requiresPermit: !job.requiresPermit, updatedAt: new Date() }
          : job
      ));
    },
    updateDates: (jobId: string, startDate?: Date, endDate?: Date) => {
      update(jobs => jobs.map(job => 
        job.id === jobId 
          ? { ...job, startDate, endDate, updatedAt: new Date() }
          : job
      ));
    },
    setWaveInvoice: (jobId: string, invoiceId: string, invoiceUrl: string) => {
      update(jobs => jobs.map(job => 
        job.id === jobId 
          ? { ...job, waveInvoiceId: invoiceId, waveInvoiceUrl: invoiceUrl, updatedAt: new Date() }
          : job
      ));
    },
    calculateJobTotal: (job: Job): number => {
      const subtotal = job.items.reduce((sum, item) => {
        return sum + (item.installedQuantity * item.item.unitPrice);
      }, 0);
      const tax = subtotal * 0.08; // 8% tax
      return subtotal + tax;
    },
    addPhotos: (jobId: string, photoUrls: string[]) => {
      update(jobs => jobs.map(job => 
        job.id === jobId 
          ? { ...job, photos: [...(job.photos || []), ...photoUrls], updatedAt: new Date() }
          : job
      ));
    },
    removePhotos: (jobId: string, photoUrls: string[]) => {
      update(jobs => jobs.map(job => 
        job.id === jobId 
          ? { ...job, photos: (job.photos || []).filter(p => !photoUrls.includes(p)), updatedAt: new Date() }
          : job
      ));
    },
    loadJobs: (jobs: Job[]) => {
      set(jobs);
    }
  };
}

export const jobsStore = createJobsStore();