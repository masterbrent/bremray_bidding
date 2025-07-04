import { writable, derived } from 'svelte/store';
import type { Job, Customer } from '../types/models';
import { api } from '../api/client';

interface JobsState {
  jobs: Job[];
  loading: boolean;
  error: string | null;
}

function createJobsStore() {
  const { subscribe, set, update } = writable<JobsState>({
    jobs: [],
    loading: false,
    error: null
  });

  return {
    subscribe,
    
    // Load all jobs
    async load() {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const jobs = await api.get<Job[]>('/jobs');
        set({ jobs, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load jobs'
        }));
      }
    },

    // Get jobs by status
    async loadByStatus(status: string) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const jobs = await api.get<Job[]>(`/jobs?status=${status}`);
        set({ jobs, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load jobs'
        }));
      }
    },

    // Get jobs by customer
    async loadByCustomer(customerId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const jobs = await api.get<Job[]>(`/jobs?customerId=${customerId}`);
        set({ jobs, loading: false, error: null });
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load jobs'
        }));
      }
    },

    // Create a new job from template
    async createFromTemplate(
      customerId: string,
      address: string,
      templateId: string,
      scheduledDate?: Date,
      notes?: string
    ) {
      update(state => ({ ...state, loading: true, error: null }));
      try {
        const payload = {
          customerId,
          templateId,
          address,
          scheduledDate: scheduledDate ? scheduledDate.toISOString() : undefined,
          notes
        };
        console.log('Creating job with payload:', payload);
        
        const job = await api.post<Job>('/jobs', payload);
        
        update(state => ({
          ...state,
          jobs: [...state.jobs, job],
          loading: false,
          error: null
        }));
        
        return job.id;
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to create job'
        }));
        throw error;
      }
    },

    // Update job status
    async updateStatus(jobId: string, status: string) {
      try {
        const job = await api.put<Job>(`/jobs/${jobId}`, { status });
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to update job status'
        }));
        throw error;
      }
    },
    
    // Generic update method
    async update(jobId: string, updates: Partial<Job>) {
      try {
        const job = await api.put<Job>(`/jobs/${jobId}`, updates);
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
        return job;
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to update job'
        }));
        throw error;
      }
    },

    // Toggle permit requirement
    async togglePermit(jobId: string) {
      update(state => {
        const job = state.jobs.find(j => j.id === jobId);
        if (!job) return state;
        
        const updatedJobs = state.jobs.map(j => 
          j.id === jobId ? { ...j, permitRequired: !j.permitRequired } : j
        );
        
        // Update backend
        api.put(`/jobs/${jobId}`, { permitRequired: !job.permitRequired })
          .catch(error => {
            // Revert on error
            update(s => ({
              ...s,
              jobs: state.jobs,
              error: 'Failed to update permit status'
            }));
          });
        
        return { ...state, jobs: updatedJobs };
      });
    },

    // Update job dates
    async updateDates(jobId: string, scheduledDate?: Date, startDate?: Date, endDate?: Date) {
      try {
        const updates: any = {};
        if (scheduledDate) updates.scheduledDate = scheduledDate.toISOString();
        if (startDate) updates.startDate = startDate.toISOString();
        if (endDate) updates.endDate = endDate.toISOString();
        
        const job = await api.put<Job>(`/jobs/${jobId}`, updates);
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to update job dates'
        }));
        throw error;
      }
    },

    // Add item to job
    async addItem(jobId: string, itemId: string, quantity: number) {
      try {
        await api.post(`/jobs/${jobId}/items`, { itemId, quantity });
        // Reload the job to get updated data
        const job = await api.get<Job>(`/jobs/${jobId}`);
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to add item'
        }));
        throw error;
      }
    },

    // Update job item quantity
    async updateItemQuantity(jobId: string, itemId: string, quantity: number) {
      try {
        await api.put(`/jobs/${jobId}/items/${itemId}`, { quantity });
        // Reload the job to get updated data
        const job = await api.get<Job>(`/jobs/${jobId}`);
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to update item quantity'
        }));
        throw error;
      }
    },
    
    // Alias for updateItemQuantity for compatibility
    async updateJobItem(jobId: string, itemId: string, quantity: number) {
      return this.updateItemQuantity(jobId, itemId, quantity);
    },

    // Remove item from job
    async removeItem(jobId: string, itemId: string) {
      try {
        await api.delete(`/jobs/${jobId}/items/${itemId}`);
        // Reload the job to get updated data
        const job = await api.get<Job>(`/jobs/${jobId}`);
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to remove item'
        }));
        throw error;
      }
    },

    // Add photo to job
    async addPhoto(jobId: string, url: string, caption?: string) {
      try {
        await api.post(`/jobs/${jobId}/photos`, { url, caption });
        // Reload the job to get updated data
        const job = await api.get<Job>(`/jobs/${jobId}`);
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to add photo'
        }));
        throw error;
      }
    },

    // Remove photo from job
    async removePhoto(jobId: string, photoId: string) {
      try {
        await api.delete(`/jobs/${jobId}/photos/${photoId}`);
        // Reload the job to get updated data
        const job = await api.get<Job>(`/jobs/${jobId}`);
        update(state => ({
          ...state,
          jobs: state.jobs.map(j => j.id === jobId ? job : j)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to remove photo'
        }));
        throw error;
      }
    },

    // Load a single job by ID
    async loadById(jobId: string) {
      try {
        const job = await api.get<Job>(`/jobs/${jobId}`);
        update(state => {
          const existingIndex = state.jobs.findIndex(j => j.id === jobId);
          const jobs = existingIndex >= 0
            ? state.jobs.map(j => j.id === jobId ? job : j)
            : [...state.jobs, job];
          return { ...state, jobs };
        });
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to load job'
        }));
        throw error;
      }
    },

    // Delete job
    async remove(id: string) {
      try {
        await api.delete(`/jobs/${id}`);
        update(state => ({
          ...state,
          jobs: state.jobs.filter(j => j.id !== id)
        }));
      } catch (error) {
        update(state => ({
          ...state,
          error: error instanceof Error ? error.message : 'Failed to delete job'
        }));
        throw error;
      }
    },

    // Calculate job total
    calculateJobTotal(job: Job): number {
      if (!job.items) return 0;
      return job.totalAmount || 0;
    }
  };
}

export const jobsStore = createJobsStore();

// Derived store for easy access to just the jobs array
export const jobs = derived(jobsStore, $jobsStore => $jobsStore.jobs);