<script lang="ts">
  import { X } from 'lucide-svelte';
  
  export let isOpen = false;
  export let title = '';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      isOpen = false;
    }
  }
  
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      isOpen = false;
    }
  }
</script>

{#if isOpen}
  <div class="modal-backdrop" on:click={handleBackdropClick} on:keydown={handleKeydown}>
    <div class="modal modal-{size}">
      <div class="modal-header">
        <h2>{title}</h2>
        <button class="close-btn" on:click={() => isOpen = false}>
          <X size={20} />
        </button>
      </div>
      
      <div class="modal-body">
        <slot />
      </div>
      
      {#if $$slots.footer}
        <div class="modal-footer">
          <slot name="footer" />
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
    animation: fadeIn 0.2s ease-out;
  }
  
  .modal {
    background: white;
    border-radius: var(--radius-xl);
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: var(--shadow-2xl);
    animation: slideUp 0.3s ease-out;
  }
  
  .modal-sm {
    width: 100%;
    max-width: 400px;
  }
  
  .modal-md {
    width: 100%;
    max-width: 600px;
  }
  
  .modal-lg {
    width: 100%;
    max-width: 800px;
  }
  
  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem;
    border-bottom: 1px solid var(--gray-100);
  }
  
  .modal-header h2 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
  }
  
  .close-btn {
    background: none;
    border: none;
    padding: 0.5rem;
    cursor: pointer;
    border-radius: var(--radius-md);
    color: var(--text-tertiary);
    transition: all var(--transition-base);
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .close-btn:hover {
    background: var(--gray-100);
    color: var(--text-primary);
  }
  
  .modal-body {
    flex: 1;
    padding: 1.5rem;
    overflow-y: auto;
  }
  
  .modal-footer {
    padding: 1.5rem;
    border-top: 1px solid var(--gray-100);
    background: var(--gray-50);
    border-radius: 0 0 var(--radius-xl) var(--radius-xl);
  }
  
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
  
  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
  
  @media (max-width: 640px) {
    .modal-backdrop {
      padding: 0;
      align-items: flex-end;
    }
    
    .modal {
      max-height: 80vh;
      width: 100%;
      max-width: 100%;
      border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    }
  }
</style>