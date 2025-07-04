<script lang="ts">
  export let hoverable = false;
  export let clickable = false;
  export let selected = false;
  export let glass = false;
  export let noPadding = false;
  export let variant: 'default' | 'bordered' | 'elevated' = 'default';
</script>

<div 
  class="card card-{variant}"
  class:card-hoverable={hoverable}
  class:card-clickable={clickable}
  class:card-selected={selected}
  class:card-glass={glass}
  class:card-no-padding={noPadding}
  on:click
  role={clickable ? 'button' : undefined}
  tabindex={clickable ? 0 : undefined}
>
  <slot />
</div>

<style>
  .card {
    background: var(--bg-primary);
    border-radius: var(--radius-2xl);
    padding: var(--space-lg);
    transition: all var(--transition-base);
    position: relative;
    overflow: hidden;
  }

  /* Variants */
  .card-default {
    background: var(--bg-primary);
  }

  .card-bordered {
    background: var(--bg-primary);
    border: 1px solid var(--gray-200);
  }

  .card-elevated {
    background: var(--bg-primary);
    box-shadow: var(--shadow-md);
  }

  /* Glass morphism */
  .card-glass {
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.2);
  }

  /* Hoverable */
  .card-hoverable:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-lg);
  }

  /* Clickable */
  .card-clickable {
    cursor: pointer;
    user-select: none;
  }

  .card-clickable:active {
    transform: translateY(0);
  }

  /* Selected */
  .card-selected {
    border: 2px solid var(--primary-500);
    box-shadow: 0 0 0 4px rgba(0, 102, 255, 0.1);
  }

  /* No padding */
  .card-no-padding {
    padding: 0;
  }

  /* Focus styles for accessibility */
  .card-clickable:focus-visible {
    outline: 2px solid var(--primary-500);
    outline-offset: 2px;
  }

  /* Subtle inner glow on hover */
  .card::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: radial-gradient(
      circle at var(--mouse-x, 50%) var(--mouse-y, 50%),
      rgba(0, 102, 255, 0.06) 0%,
      transparent 50%
    );
    opacity: 0;
    transition: opacity var(--transition-slow);
    pointer-events: none;
  }

  .card-hoverable:hover::after {
    opacity: 1;
  }
</style>