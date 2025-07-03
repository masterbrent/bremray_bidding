<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' | 'ghost' = 'primary';
  export let size: 'small' | 'medium' | 'large' = 'medium';
  export let disabled: boolean = false;
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let fullWidth: boolean = false;
  export let loading: boolean = false;
</script>

<button
  {type}
  {disabled}
  class="button {variant} {size}"
  class:full-width={fullWidth}
  class:loading
  on:click
>
  {#if loading}
    <span class="spinner"></span>
  {/if}
  <slot />
</button>

<style>
  .button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    border: none;
    border-radius: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-family: inherit;
    position: relative;
    overflow: hidden;
    letter-spacing: -0.01em;
  }
  
  .button::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.3);
    transform: translate(-50%, -50%);
    transition: width 0.6s, height 0.6s;
  }
  
  .button:active::before {
    width: 300px;
    height: 300px;
  }
  
  .button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .button.loading {
    color: transparent;
  }
  
  .spinner {
    position: absolute;
    width: 20px;
    height: 20px;
    border: 2px solid currentColor;
    border-right-color: transparent;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  .button.full-width {
    width: 100%;
  }
  
  /* Sizes */
  .button.small {
    padding: 0.625rem 1.25rem;
    font-size: 0.875rem;
    min-height: 40px;
  }
  
  .button.medium {
    padding: 0.875rem 1.75rem;
    font-size: 0.9375rem;
    min-height: 48px;
  }
  
  .button.large {
    padding: 1.125rem 2.25rem;
    font-size: 1.0625rem;
    min-height: 56px;
  }
  
  /* Variants */
  .button.primary {
    background: var(--primary-500);
    color: white;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  }
  
  .button.primary:hover:not(:disabled) {
    background: var(--primary-600);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px 0 rgba(59, 130, 246, 0.25);
  }
  
  .button.primary:active:not(:disabled) {
    transform: translateY(0);
  }
  
  .button.secondary {
    background: var(--card-bg);
    color: var(--gray-700);
    border: 1px solid var(--gray-300);
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  }
  
  .button.secondary:hover:not(:disabled) {
    background: var(--gray-50);
    border-color: var(--gray-400);
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }
  
  .button.secondary:active:not(:disabled) {
    transform: translateY(0);
  }
  
  .button.danger {
    background: var(--danger-500);
    color: white;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  }
  
  .button.danger:hover:not(:disabled) {
    background: var(--danger-600);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px 0 rgba(239, 68, 68, 0.25);
  }
  
  .button.danger:active:not(:disabled) {
    transform: translateY(0);
  }
  
  .button.ghost {
    background: transparent;
    color: var(--gray-600);
    box-shadow: none;
  }
  
  .button.ghost:hover:not(:disabled) {
    background: var(--gray-100);
    color: var(--gray-900);
  }
  
  .button.ghost:active:not(:disabled) {
    background: var(--gray-200);
  }
</style>