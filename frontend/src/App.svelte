<script lang="ts">
  import { onMount } from 'svelte';
  import { router } from './lib/router';
  import { itemsStore, jobsStore, jobTemplatesStore, companySettingsStore } from './lib/stores';
  import { Menu, X, Briefcase, Package, FileText, Settings, Zap } from 'lucide-svelte';
  
  // Import pages
  import JobsPage from './pages/JobsPage.svelte';
  import ItemsPage from './pages/ItemsPage.svelte';
  import JobTemplatesPage from './pages/JobTemplatesPage.svelte';
  import SettingsPage from './pages/SettingsPage.svelte';
  import JobDetailPage from './pages/JobDetailPage.svelte';
  
  let currentRoute = 'jobs';
  let routeParams: any = {};
  let error: any = null;
  let mobileMenuOpen = false;
  
  const navItems = [
    { route: 'jobs', label: 'Jobs', icon: Briefcase },
    { route: 'items', label: 'Items', icon: Package },
    { route: 'templates', label: 'Templates', icon: FileText },
    { route: 'settings', label: 'Settings', icon: Settings },
  ];
  
  // Subscribe to router
  router.subscribe(state => {
    currentRoute = state.route;
    routeParams = state.params || {};
  });
  
  onMount(() => {
    // Load data from backend
    itemsStore.load().catch(e => console.error('Failed to load items:', e));
    jobsStore.load().catch(e => console.error('Failed to load jobs:', e));
    jobTemplatesStore.load().catch(e => console.error('Failed to load templates:', e));
    companySettingsStore.load().catch(e => console.error('Failed to load company settings:', e));
  });
  
  function navigate(route: string) {
    router.navigate(route as any);
    mobileMenuOpen = false;
  }
  
  function toggleMobileMenu() {
    mobileMenuOpen = !mobileMenuOpen;
  }
</script>

<div class="app">
  <!-- Mobile Header -->
  <div class="mobile-header">
    <div class="mobile-header-content">
      <button class="logo-mobile" on:click={() => navigate('jobs')}>
        <Zap size={20} />
        <span>Bremray</span>
      </button>
      <button class="menu-toggle" on:click={toggleMobileMenu}>
        {#if mobileMenuOpen}
          <X size={24} />
        {:else}
          <Menu size={24} />
        {/if}
      </button>
    </div>
  </div>
  
  <!-- Sidebar -->
  <aside class="sidebar" class:open={mobileMenuOpen}>
    <div class="sidebar-header">
      <button class="logo" on:click={() => navigate('jobs')}>
        <Zap size={24} class="logo-icon" />
        <span>Bremray</span>
      </button>
    </div>
    
    <nav class="nav">
      {#each navItems as item}
        <button 
          class="nav-item"
          class:active={currentRoute === item.route}
          on:click={() => navigate(item.route)}
        >
          <svelte:component this={item.icon} size={20} class="nav-icon" />
          <span>{item.label}</span>
        </button>
      {/each}
    </nav>
    
    <div class="sidebar-footer">
      <div class="user-section">
        <div class="user-avatar">B</div>
        <div>
          <p class="user-name">Brent Hall</p>
          <p class="user-role">Admin</p>
        </div>
      </div>
    </div>
  </aside>
  
  <!-- Overlay for mobile -->
  {#if mobileMenuOpen}
    <div class="overlay" on:click={toggleMobileMenu}></div>
  {/if}
  
  <!-- Main Content -->
  <main class="main-content">
    {#if error}
      <div class="error">
        <h2>Error:</h2>
        <pre>{error}</pre>
      </div>
    {:else}
      <div class="page-wrapper">
        {#if currentRoute === 'jobs'}
          <JobsPage />
        {:else if currentRoute === 'items'}
          <ItemsPage />
        {:else if currentRoute === 'templates'}
          <JobTemplatesPage />
        {:else if currentRoute === 'settings'}
          <SettingsPage />
        {:else if currentRoute === 'job-detail'}
          <JobDetailPage jobId={routeParams.id} />
        {:else}
          <div class="page">
            <h1>Page not found</h1>
            <p>Route: {currentRoute}</p>
          </div>
        {/if}
      </div>
    {/if}
  </main>
</div>

<style>
  .app {
    display: flex;
    height: 100vh;
    background: var(--bg-secondary);
  }
  
  /* Mobile Header */
  .mobile-header {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 60px;
    background: white;
    border-bottom: 1px solid var(--gray-200);
    z-index: 100;
  }
  
  .mobile-header-content {
    height: 100%;
    padding: 0 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo-mobile {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 700;
    font-size: 1.125rem;
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    color: var(--text-primary);
    transition: opacity var(--transition-base);
  }
  
  .logo-mobile:hover {
    opacity: 0.8;
  }
  
  .logo-mobile:active {
    transform: scale(0.98);
  }
  
  .logo-mobile :global(svg) {
    color: var(--primary-500);
  }
  
  .menu-toggle {
    background: none;
    border: none;
    padding: 0.5rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-md);
    transition: all var(--transition-base);
  }
  
  .menu-toggle:hover {
    background: var(--gray-100);
  }
  
  /* Sidebar */
  .sidebar {
    width: 240px;
    background: white;
    border-right: 1px solid var(--gray-200);
    display: flex;
    flex-direction: column;
    transition: transform var(--transition-slow);
  }
  
  .sidebar-header {
    padding: 2rem 1.5rem;
    border-bottom: 1px solid var(--gray-100);
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-weight: 700;
    font-size: 1.25rem;
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    color: var(--text-primary);
    transition: opacity var(--transition-base);
  }
  
  .logo:hover {
    opacity: 0.8;
  }
  
  .logo:active {
    transform: scale(0.98);
  }
  
  :global(.logo-icon) {
    color: var(--primary-500);
  }
  
  /* Navigation */
  .nav {
    flex: 1;
    padding: 1rem;
  }
  
  .nav-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    margin-bottom: 0.25rem;
    background: none;
    border: none;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all var(--transition-base);
    text-align: left;
  }
  
  :global(.nav-icon) {
    color: var(--text-tertiary);
  }
  
  .nav-item:hover {
    background: var(--gray-50);
    color: var(--text-primary);
  }
  
  .nav-item:hover :global(.nav-icon) {
    color: var(--text-secondary);
  }
  
  .nav-item.active {
    background: var(--primary-500);
    color: white;
  }
  
  .nav-item.active :global(.nav-icon) {
    color: white;
  }
  
  /* Sidebar Footer */
  .sidebar-footer {
    padding: 1.5rem;
    border-top: 1px solid var(--gray-100);
  }
  
  .user-section {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .user-avatar {
    width: 36px;
    height: 36px;
    background: var(--primary-500);
    color: white;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 0.875rem;
  }
  
  .user-name {
    margin: 0;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--text-primary);
  }
  
  .user-role {
    margin: 0;
    font-size: 0.75rem;
    color: var(--text-tertiary);
  }
  
  /* Main Content */
  .main-content {
    flex: 1;
    overflow-y: auto;
  }
  
  .page-wrapper {
    min-height: 100%;
  }
  
  .overlay {
    display: none;
  }
  
  /* Error styles */
  .error {
    padding: 2rem;
    color: var(--danger-500);
  }
  
  .error pre {
    background: var(--gray-100);
    padding: 1rem;
    border-radius: var(--radius-md);
    margin-top: 1rem;
  }
  
  /* Responsive */
  @media (max-width: 768px) {
    .mobile-header {
      display: block;
    }
    
    .sidebar {
      position: fixed;
      left: 0;
      top: 0;
      bottom: 0;
      transform: translateX(-100%);
      z-index: 200;
      box-shadow: 4px 0 20px rgba(0, 0, 0, 0.1);
    }
    
    .sidebar.open {
      transform: translateX(0);
    }
    
    .main-content {
      padding-top: 60px;
    }
    
    .overlay {
      display: block;
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.3);
      z-index: 150;
    }
  }
</style>