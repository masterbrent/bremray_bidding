<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
  export let value: number = 0;
  export let min: number = 0;
  export let max: number = 999;
  export let step: number = 1;
  export let disabled: boolean = false;
  
  const dispatch = createEventDispatcher<{
    change: number;
  }>();
  
  function increment() {
    if (value < max) {
      value += step;
      dispatch('change', value);
    }
  }
  
  function decrement() {
    if (value > min) {
      value -= step;
      dispatch('change', value);
    }
  }
  
  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    const newValue = parseInt(target.value) || 0;
    if (newValue >= min && newValue <= max) {
      value = newValue;
      dispatch('change', value);
    }
  }
</script>

<div class="quantity-picker" class:disabled>
  <button 
    type="button"
    class="decrement"
    on:click={decrement}
    disabled={disabled || value <= min}
    aria-label="Decrease quantity"
  >
    −
  </button>
  
  <input
    type="number"
    {value}
    {min}
    {max}
    {step}
    {disabled}
    on:input={handleInput}
    class="value"
  />
  
  <button 
    type="button"
    class="increment"
    on:click={increment}
    disabled={disabled || value >= max}
    aria-label="Increase quantity"
  >
    +
  </button>
</div>

<style>
  .quantity-picker {
    display: inline-flex;
    align-items: center;
    gap: 0.75rem;
    background: var(--gray-100);
    border-radius: 16px;
    padding: 0.375rem;
  }
  
  .quantity-picker.disabled {
    opacity: 0.5;
  }
  
  button {
    width: 40px;
    height: 40px;
    border: none;
    background: white;
    color: var(--gray-700);
    border-radius: 12px;
    font-size: 1.25rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    line-height: 1;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }
  
  button:hover:not(:disabled) {
    background: var(--primary-500);
    color: white;
    transform: scale(1.05);
  }
  
  button:active:not(:disabled) {
    transform: scale(0.95);
  }
  
  button:disabled {
    cursor: not-allowed;
    color: var(--gray-400);
    opacity: 0.5;
  }
  
  .value {
    width: 60px;
    text-align: center;
    background: transparent;
    border: none;
    padding: 0.5rem 0;
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--gray-900);
  }
  
  .value:focus {
    outline: none;
  }
  
  input[type="number"]::-webkit-inner-spin-button,
  input[type="number"]::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }
  
  input[type="number"] {
    -moz-appearance: textfield;
  }
</style>