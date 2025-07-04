<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' | 'ghost' | 'success' = 'primary';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let fullWidth = false;
  export let disabled = false;
  export let loading = false;
  export let icon = false;
  export let type: 'button' | 'submit' | 'reset' = 'button';
</script>

<button
  {type}
  class="btn btn-{variant} btn-{size}"
  class:btn-full={fullWidth}
  class:btn-icon={icon}
  class:btn-loading={loading}
  disabled={disabled || loading}
  on:click
>
  {#if loading}
    <span class="spinner"></span>
  {/if}
  <slot />
</button>

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0 1.25rem;
    height: 40px;
    font-family: inherit;
    font-weight: 500;
    font-size: 0.875rem;
    line-height: 1;
    border: none;
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-base);
    position: relative;
    white-space: nowrap;
  }

  /* Size variants */
  .btn-sm {
    height: 32px;
    padding: 0 1rem;
    font-size: 0.8125rem;
  }

  .btn-lg {
    height: 48px;
    padding: 0 1.5rem;
    font-size: 0.9375rem;
  }

  /* Icon button */
  .btn-icon {
    width: 40px;
    padding: 0;
  }

  .btn-icon.btn-sm {
    width: 32px;
  }

  .btn-icon.btn-lg {
    width: 48px;
  }

  /* Color variants */
  .btn-primary {
    background: var(--primary-500);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--primary-600);
  }

  .btn-primary:active:not(:disabled) {
    transform: scale(0.98);
  }

  .btn-secondary {
    background: var(--gray-100);
    color: var(--text-primary);
  }

  .btn-secondary:hover:not(:disabled) {
    background: var(--gray-200);
  }

  .btn-success {
    background: var(--success-500);
    color: white;
  }

  .btn-success:hover:not(:disabled) {
    background: var(--success-600);
  }

  .btn-danger {
    background: var(--danger-500);
    color: white;
  }

  .btn-danger:hover:not(:disabled) {
    background: var(--danger-600);
  }

  .btn-ghost {
    background: transparent;
    color: var(--text-secondary);
  }

  .btn-ghost:hover:not(:disabled) {
    background: var(--gray-100);
    color: var(--text-primary);
  }

  /* States */
  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-full {
    width: 100%;
  }

  .btn:focus-visible {
    outline: 2px solid var(--primary-500);
    outline-offset: 2px;
  }

  /* Loading spinner */
  .spinner {
    width: 14px;
    height: 14px;
    border: 2px solid transparent;
    border-top-color: currentColor;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  .btn-loading {
    color: transparent;
  }

  .btn-loading .spinner {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>