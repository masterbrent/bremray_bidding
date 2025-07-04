import { writable } from 'svelte/store';

export type Route = 'jobs' | 'items' | 'templates' | 'job-detail' | 'settings';

interface RouterState {
  route: Route;
  params?: Record<string, string>;
}

function createRouter() {
  const { subscribe, set, update } = writable<RouterState>({
    route: 'jobs'
  });

  return {
    subscribe,
    navigate: (route: Route, params?: Record<string, string>) => {
      set({ route, params });
    },
    getParam: (state: RouterState, key: string) => {
      return state.params?.[key];
    }
  };
}

export const router = createRouter();