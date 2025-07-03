<script lang="ts">
  import Modal from './Modal.svelte';
  import Button from './Button.svelte';
  
  export let isOpen: boolean = false;
  export let title: string = 'Confirm Action';
  export let message: string = 'Are you sure you want to proceed?';
  export let confirmText: string = 'Confirm';
  export let cancelText: string = 'Cancel';
  export let variant: 'danger' | 'warning' | 'info' = 'danger';
  
  export let onConfirm: () => void = () => {};
  export let onCancel: () => void = () => {};
  
  function handleConfirm() {
    onConfirm();
    isOpen = false;
  }
  
  function handleCancel() {
    onCancel();
    isOpen = false;
  }
</script>

<Modal bind:isOpen {title} size="small">
  <p class="message">{message}</p>
  
  <div slot="footer">
    <Button variant="secondary" on:click={handleCancel}>
      {cancelText}
    </Button>
    <Button variant={variant === 'danger' ? 'danger' : 'primary'} on:click={handleConfirm}>
      {confirmText}
    </Button>
  </div>
</Modal>

<style>
  .message {
    margin: 0;
    color: #374151;
    line-height: 1.5;
  }
</style>