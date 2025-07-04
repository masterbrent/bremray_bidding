<script lang="ts">
  import { onMount } from 'svelte';
  
  let message = 'App is loading...';
  
  onMount(() => {
    message = 'App mounted successfully!';
    console.log('AppMinimal mounted');
  });
</script>

<div>
  <h1>Minimal App Test</h1>
  <p>{message}</p>
</div>

<style>
  div {
    padding: 20px;
    font-family: sans-serif;
  }
  
  h1 {
    color: #3b82f6;
  }
</style>