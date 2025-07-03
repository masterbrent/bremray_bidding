<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  
  export let isOpen: boolean = false;
  export let title: string = '';
  export let size: 'small' | 'medium' | 'large' = 'medium';
  
  const dispatch = createEventDispatcher<{
    close: void;
  }>();
  
  let dialog: HTMLDialogElement;
  
  $: if (dialog && isOpen) {
    dialog.showModal();
  } else if (dialog && !isOpen) {
    dialog.close();
  }
  
  function handleClose() {
    isOpen = false;
    dispatch('close');
  }
  
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === dialog) {
      handleClose();
    }
  }
  
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      handleClose();
    }
  }
</script>

<dialog
  bind:this={dialog}
  class="modal {size}"
  on:click={handleBackdropClick}
  on:keydown={handleKeydown}
>
  <div class="modal-content">
    <div class="modal-header">
      <h2>{title}</h2>
      <button
        type="button"
        class="close-button"
        on:click={handleClose}
        aria-label="Close modal"
      >
        ×
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
</dialog>

<style>
  dialog {
    padding: 0;
    border: none;
    border-radius: 8px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
    overflow: visible;
  }
  
  dialog::backdrop {
    background-color: rgba(0, 0, 0, 0.5);
  }
  
  .modal-content {
    background: white;
    border-radius: 8px;
    overflow: hidden;
  }
  
  .modal.small .modal-content {
    width: 400px;
    max-width: 90vw;
  }
  
  .modal.medium .modal-content {
    width: 600px;
    max-width: 90vw;
  }
  
  .modal.large .modal-content {
    width: 800px;
    max-width: 90vw;
  }
  
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e5e7eb;
  }
  
  .modal-header h2 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
    color: #111827;
  }
  
  .close-button {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    border: none;
    border-radius: 4px;
    font-size: 1.5rem;
    color: #6b7280;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .close-button:hover {
    background-color: #f3f4f6;
  }
  
  .modal-body {
    padding: 1.5rem;
    max-height: 70vh;
    overflow-y: auto;
  }
  
  .modal-footer {
    padding: 1rem 1.5rem;
    border-top: 1px solid #e5e7eb;
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }
</style>