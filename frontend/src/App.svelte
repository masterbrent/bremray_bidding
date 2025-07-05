<script lang="ts">
  import { onMount } from 'svelte';
  import { router } from './lib/router';
  import { itemsStore, jobsStore, jobTemplatesStore, companySettingsStore, userStore, effectiveRole, permissions } from './lib/stores';
  import { Menu, X, Briefcase, Package, FileText, Settings, Zap, BarChart3, Users, Shield } from 'lucide-svelte';
  import StatusIndicators from './lib/components/StatusIndicators.svelte';
  
  // Import pages
  import JobsPage from './pages/JobsPage.svelte';
  import ItemsPage from './pages/ItemsPage.svelte';
  import JobTemplatesPage from './pages/JobTemplatesPage.svelte';
  import SettingsPage from './pages/SettingsPage.svelte';
  import UserPreferencesPage from './pages/UserPreferencesPage.svelte';
  import JobDetailPage from './pages/JobDetailPage.svelte';
  import StatsPage from './pages/StatsPage.svelte';
  import UsersPage from './pages/UsersPage.svelte';
  
  let currentRoute = 'jobs';
  let routeParams: any = {};
  let error: any = null;
  let mobileMenuOpen = false;
  
  // Workspace state
  let currentWorkspace = 'skyview';
  
  // For techs, always force skyview
  $: if (!permissions.canSeePrices($effectiveRole)) {
    currentWorkspace = 'skyview';
  }
  
  $: navItems = [
    { route: 'jobs', label: 'Jobs', icon: Briefcase },
    ...(permissions.canSeePrices($effectiveRole) ? [
      { route: 'stats', label: 'Stats', icon: BarChart3 },
    ] : []),
    ...(permissions.canEditJobs($effectiveRole) ? [
      { route: 'items', label: 'Items', icon: Package },
      { route: 'templates', label: 'Templates', icon: FileText },
      { route: 'users', label: 'Users', icon: Users },
      { route: 'settings', label: 'Settings', icon: Settings },
    ] : [
      { route: 'settings', label: 'Preferences', icon: Settings },
    ]),
  ];
  
  // Subscribe to router
  router.subscribe(state => {
    currentRoute = state.route;
    routeParams = state.params || {};
  });
  
  onMount(() => {
    // Initialize user - in production this would come from auth
    // For testing, using a hardcoded email
    userStore.init('brenthall@gmail.com');
    
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
    
    {#if permissions.canSeePrices($effectiveRole)}
      <div class="workspace-selector">
        <label>Workspace</label>
        <select bind:value={currentWorkspace} class="workspace-dropdown">
          <option value="skyview">Skyview</option>
          <option value="contractors">Contractors</option>
          <option value="rayno">Rayno</option>
        </select>
      </div>
    {/if}
    
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
      <StatusIndicators />
      
      {#if $userStore?.role === 'admin'}
        <button 
          class="admin-toggle"
          on:click={() => userStore.toggleViewMode()}
          title={$userStore.isViewingAsTech ? 'Switch to Admin View' : 'Switch to Tech View'}
        >
          {#if $userStore.isViewingAsTech}
            <Shield size={16} />
            <span>Switch to Admin View</span>
          {:else}
            <Users size={16} />
            <span>Switch to Tech View</span>
          {/if}
        </button>
      {/if}
      
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
          <JobsPage workspace={currentWorkspace} />
        {:else if currentRoute === 'stats'}
          {#if permissions.canSeePrices($effectiveRole)}
            <StatsPage />
          {:else}
            <div class="access-denied">
              <h2>Access Denied</h2>
              <p>You don't have permission to view this page.</p>
            </div>
          {/if}
        {:else if currentRoute === 'items'}
          {#if permissions.canEditJobs($effectiveRole)}
            <ItemsPage />
          {:else}
            <div class="access-denied">
              <h2>Access Denied</h2>
              <p>You don't have permission to view this page.</p>
            </div>
          {/if}
        {:else if currentRoute === 'templates'}
          {#if permissions.canEditJobs($effectiveRole)}
            <JobTemplatesPage />
          {:else}
            <div class="access-denied">
              <h2>Access Denied</h2>
              <p>You don't have permission to view this page.</p>
            </div>
          {/if}
        {:else if currentRoute === 'users'}
          {#if permissions.canEditJobs($effectiveRole)}
            <UsersPage />
          {:else}
            <div class="access-denied">
              <h2>Access Denied</h2>
              <p>You don't have permission to view this page.</p>
            </div>
          {/if}
        {:else if currentRoute === 'settings'}
          {#if permissions.canEditJobs($effectiveRole)}
            <SettingsPage />
          {:else}
            <UserPreferencesPage />
          {/if}
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
  
  /* Workspace Selector */
  .workspace-selector {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--gray-100);
  }
  
  .workspace-selector label {
    display: block;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.5rem;
  }
  
  .workspace-dropdown {
    width: 100%;
    padding: 0.5rem 2.25rem 0.5rem 0.75rem;
    background: var(--gray-50);
    border: 1px solid var(--gray-200);
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-primary);
    appearance: none;
    background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3e%3cpath fill='%236b7280' d='M6 9L1 4h10z'/%3e%3c/svg%3e");
    background-repeat: no-repeat;
    background-position: right 0.5rem center;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .workspace-dropdown:hover {
    background-color: white;
    border-color: var(--gray-300);
  }
  
  .workspace-dropdown:focus {
    outline: none;
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(91, 91, 214, 0.1);
  }

  /* Sidebar Footer */
  .sidebar-footer {
    padding: 1.5rem;
    border-top: 1px solid var(--gray-100);
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  /* Admin Toggle */
  .admin-toggle {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    background: var(--gray-50);
    border: 1px solid var(--gray-200);
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .admin-toggle:hover {
    background: white;
    border-color: var(--gray-300);
    color: var(--text-primary);
  }
  
  .admin-toggle :global(svg) {
    width: 16px;
    height: 16px;
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
  
  /* Access Denied styles */
  .access-denied {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 50vh;
    text-align: center;
    padding: 2rem;
  }
  
  .access-denied h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
  }
  
  .access-denied p {
    color: var(--text-secondary);
    font-size: 0.9375rem;
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