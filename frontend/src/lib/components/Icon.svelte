<script lang="ts">
  import * as icons from 'lucide-svelte';
  
  export let name: keyof typeof icons;
  export let size: number = 20;
  export let strokeWidth: number = 2;
  export let color: string = 'currentColor';
  export let class: string = '';
  
  $: Icon = icons[name];
</script>

{#if Icon}
  <svelte:component 
    this={Icon} 
    {size} 
    {strokeWidth} 
    {color} 
    class={`icon ${$$props.class || ''}`}
  />
{/if}

<style>
  .icon {
    display: inline-block;
    vertical-align: middle;
  }
</style>