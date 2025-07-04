<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Cloud, Zap, CloudOff, ZapOff, Loader } from 'lucide-svelte';
  import { serviceHealth, type ServiceStatus } from '../services/serviceHealth';
  
  let status: ServiceStatus = {
    wave: 'checking',
    cloudflare: 'checking'
  };
  
  let unsubscribe: (() => void) | null = null;
  
  onMount(() => {
    // Subscribe to service status updates
    unsubscribe = serviceHealth.subscribe((newStatus) => {
      status = newStatus;
    });
  });
  
  onDestroy(() => {
    if (unsubscribe) {
      unsubscribe();
    }
  });
  
  function getStatusColor(state: 'connected' | 'disconnected' | 'checking'): string {
    switch (state) {
      case 'connected': return '#10b981'; // green
      case 'disconnected': return '#6b7280'; // gray
      case 'checking': return '#f59e0b'; // amber
    }
  }
  
  function getStatusTitle(service: string, state: 'connected' | 'disconnected' | 'checking'): string {
    switch (state) {
      case 'connected': return `${service} connected`;
      case 'disconnected': return `${service} disconnected`;
      case 'checking': return `Checking ${service} connection...`;
    }
  }
</script>

<div class="status-indicators">
  <!-- Cloudflare R2 Status -->
  <div 
    class="status-icon"
    title={getStatusTitle('Cloudflare R2', status.cloudflare)}
  >
    {#if status.cloudflare === 'checking'}
      <Loader size={16} class="spinning" color={getStatusColor('checking')} />
    {:else if status.cloudflare === 'connected'}
      <Cloud size={16} color={getStatusColor('connected')} />
    {:else}
      <CloudOff size={16} color={getStatusColor('disconnected')} />
    {/if}
  </div>
  
  <!-- Wave Status -->
  <div 
    class="status-icon"
    title={getStatusTitle('Wave', status.wave)}
  >
    {#if status.wave === 'checking'}
      <Loader size={16} class="spinning" color={getStatusColor('checking')} />
    {:else if status.wave === 'connected'}
      <Zap size={16} color={getStatusColor('connected')} />
    {:else}
      <ZapOff size={16} color={getStatusColor('disconnected')} />
    {/if}
  </div>
</div>

<style>
  .status-indicators {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .status-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    background: #f3f4f6;
    cursor: help;
    transition: all 0.2s ease;
  }
  
  .status-icon:hover {
    background: #e5e7eb;
  }
  
  :global(.spinning) {
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
