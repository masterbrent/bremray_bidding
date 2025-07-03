<script lang="ts">
  export let clickable: boolean = false;
  export let padding: 'none' | 'small' | 'medium' | 'large' = 'medium';
  export let variant: 'default' | 'elevated' | 'glass' = 'default';
</script>

{#if clickable}
  <button 
    class="card {variant}"
    class:clickable
    class:padding-none={padding === 'none'}
    class:padding-small={padding === 'small'}
    class:padding-medium={padding === 'medium'}
    class:padding-large={padding === 'large'}
    on:click
  >
    <slot />
  </button>
{:else}
  <div 
    class="card {variant}"
    class:padding-none={padding === 'none'}
    class:padding-small={padding === 'small'}
    class:padding-medium={padding === 'medium'}
    class:padding-large={padding === 'large'}
  >
    <slot />
  </div>
{/if}

<style>
  .card {
    border-radius: 24px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
  }
  
  .card.default {
    background: var(--card-bg);
    border: 1px solid var(--gray-200);
  }
  
  .card.elevated {
    background: var(--card-bg);
    border: none;
    box-shadow: var(--shadow-card);
  }
  
  .card.glass {
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: var(--shadow-lg);
  }
  
  button.card {
    width: 100%;
    text-align: left;
    font-family: inherit;
    cursor: pointer;
  }
  
  .card.clickable:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-2xl);
  }
  
  .card.clickable:active {
    transform: translateY(0);
  }
  
  .card.padding-none {
    padding: 0;
  }
  
  .card.padding-small {
    padding: 1.25rem;
  }
  
  .card.padding-medium {
    padding: 1.75rem;
  }
  
  .card.padding-large {
    padding: 2.25rem;
  }
  
</style>