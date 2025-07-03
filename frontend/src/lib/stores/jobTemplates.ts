import { writable } from 'svelte/store';
import type { JobTemplate, JobPhase } from '../types/models';

function createJobTemplatesStore() {
  const { subscribe, set, update } = writable<JobTemplate[]>([]);

  return {
    subscribe,
    add: (template: Omit<JobTemplate, 'id' | 'createdAt' | 'updatedAt'>) => {
      update(templates => {
        const newTemplate: JobTemplate = {
          ...template,
          id: crypto.randomUUID(),
          createdAt: new Date(),
          updatedAt: new Date()
        };
        return [...templates, newTemplate];
      });
    },
    update: (id: string, updates: Partial<Omit<JobTemplate, 'id' | 'createdAt'>>) => {
      update(templates => templates.map(template => 
        template.id === id 
          ? { ...template, ...updates, updatedAt: new Date() }
          : template
      ));
    },
    remove: (id: string) => {
      update(templates => templates.filter(template => template.id !== id));
    },
    addItemToTemplate: (templateId: string, itemId: string, defaultQuantity: number = 1) => {
      update(templates => templates.map(template => {
        if (template.id === templateId) {
          const exists = template.items.some(item => item.itemId === itemId);
          if (!exists) {
            return {
              ...template,
              items: [...template.items, { itemId, defaultQuantity }],
              updatedAt: new Date()
            };
          }
        }
        return template;
      }));
    },
    removeItemFromTemplate: (templateId: string, itemId: string) => {
      update(templates => templates.map(template => 
        template.id === templateId 
          ? {
              ...template,
              items: template.items.filter(item => item.itemId !== itemId),
              updatedAt: new Date()
            }
          : template
      ));
    },
    loadTemplates: (templates: JobTemplate[]) => {
      set(templates);
    }
  };
}

export const jobTemplatesStore = createJobTemplatesStore();