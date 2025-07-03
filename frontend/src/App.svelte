<script lang="ts">
  import { onMount } from 'svelte';
  import { router } from './lib/router';
  import { itemsStore, jobTemplatesStore, jobsStore } from './lib/stores';
  import { Briefcase, Package, FileText, ChevronLeft, Menu, Zap, Settings } from 'lucide-svelte';
  
  import JobsPage from './pages/JobsPage.svelte';
  import ItemsPage from './pages/ItemsPage.svelte';
  import JobTemplatesPage from './pages/JobTemplatesPage.svelte';
  import JobDetailPage from './pages/JobDetailPage.svelte';
  import SettingsPage from './pages/SettingsPage.svelte';
  
  let sidebarCollapsed = false;
  let mobileMenuOpen = false;
  
  // Load mock data on mount
  onMount(() => {
    // Mock items
    itemsStore.loadItems([
      {
        id: '1',
        name: 'Outlet',
        unit: 'each',
        unitPrice: 15.50,
        category: 'Electrical',
        createdAt: new Date(),
        updatedAt: new Date()
      },
      {
        id: '2',
        name: 'Switch',
        unit: 'each',
        unitPrice: 12.00,
        category: 'Electrical',
        createdAt: new Date(),
        updatedAt: new Date()
      },
      {
        id: '3',
        name: 'Wire 12 AWG',
        unit: 'ft',
        unitPrice: 0.85,
        category: 'Wire',
        createdAt: new Date(),
        updatedAt: new Date()
      },
      {
        id: '4',
        name: 'Circuit Breaker 20A',
        unit: 'each',
        unitPrice: 45.00,
        category: 'Panel',
        createdAt: new Date(),
        updatedAt: new Date()
      }
    ]);
    
    // Mock job templates
    jobTemplatesStore.loadTemplates([
      {
        id: '1',
        name: 'Kitchen Remodel',
        description: 'Standard kitchen electrical upgrade',
        items: [
          { itemId: '1', defaultQuantity: 6 },
          { itemId: '2', defaultQuantity: 3 },
          { itemId: '3', defaultQuantity: 200 }
        ],
        phases: [
          { id: '1', name: 'Rough-in', order: 1, isCompleted: false },
          { id: '2', name: 'Trim out', order: 2, isCompleted: false },
          { id: '3', name: 'Final inspection', order: 3, isCompleted: false }
        ],
        createdAt: new Date(),
        updatedAt: new Date()
      },
      {
        id: '2',
        name: 'Bathroom Remodel',
        description: 'Standard bathroom electrical work',
        items: [
          { itemId: '1', defaultQuantity: 3 },
          { itemId: '2', defaultQuantity: 2 },
          { itemId: '3', defaultQuantity: 100 }
        ],
        phases: [
          { id: '1', name: 'Rough-in', order: 1, isCompleted: false },
          { id: '2', name: 'Trim out', order: 2, isCompleted: false }
        ],
        createdAt: new Date(),
        updatedAt: new Date()
      }
    ]);
  });
  
  const navItems = [
    { route: 'jobs', label: 'Jobs', icon: Briefcase },
    { route: 'items', label: 'Items', icon: Package },
    { route: 'templates', label: 'Templates', icon: FileText },
    { route: 'settings', label: 'Settings', icon: Settings }
  ];
  
  function handleNavClick(route: any) {
    router.navigate(route);
    mobileMenuOpen = false;
  }
</script>

<div class="app">
  <!-- Mobile Header -->
  <header class="mobile-header">
    <button 
      class="mobile-menu-btn"
      on:click={() => mobileMenuOpen = !mobileMenuOpen}
    >
      <Menu size={24} />
    </button>
    <div class="logo-mobile">
      <div class="logo-icon">
        <Zap size={24} />
      </div>
      <span class="logo-text">Bremray</span>
    </div>
  </header>

  <!-- Sidebar -->
  <aside class="sidebar" class:collapsed={sidebarCollapsed} class:mobile-open={mobileMenuOpen}>
    <div class="sidebar-header">
      {#if !sidebarCollapsed}
        <div class="logo">
          <div class="logo-icon">
            <Zap size={24} />
          </div>
          <span class="logo-text">Bremray</span>
        </div>
      {:else}
        <div class="logo-collapsed">
          <div class="logo-icon">
            <Zap size={20} />
          </div>
        </div>
      {/if}
      <button 
        class="collapse-btn desktop-only"
        on:click={() => sidebarCollapsed = !sidebarCollapsed}
      >
        <ChevronLeft size={20} class={sidebarCollapsed ? 'rotate-180' : ''} />
      </button>
    </div>
    
    <nav class="nav">
      {#each navItems as item}
        <button
          class="nav-item"
          class:active={$router.route === item.route}
          on:click={() => handleNavClick(item.route)}
        >
          <svelte:component this={item.icon} size={20} />
          {#if !sidebarCollapsed}
            <span class="nav-label">{item.label}</span>
          {/if}
        </button>
      {/each}
    </nav>
    
    <div class="sidebar-footer">
      {#if !sidebarCollapsed}
        <div class="user-info">
          <div class="user-avatar">JD</div>
          <div class="user-details">
            <div class="user-name">John Doe</div>
            <div class="user-role">Master Electrician</div>
          </div>
        </div>
      {:else}
        <div class="user-avatar-small">JD</div>
      {/if}
    </div>
  </aside>
  
  <!-- Mobile Menu Overlay -->
  {#if mobileMenuOpen}
    <div class="mobile-overlay" on:click={() => mobileMenuOpen = false}></div>
  {/if}

  <!-- Main Content -->
  <main class="main-content" class:sidebar-collapsed={sidebarCollapsed}>
    {#if $router.route === 'jobs'}
      <JobsPage />
    {:else if $router.route === 'items'}
      <ItemsPage />
    {:else if $router.route === 'templates'}
      <JobTemplatesPage />
    {:else if $router.route === 'job-detail'}
      <JobDetailPage jobId={$router.params?.id} />
    {:else if $router.route === 'settings'}
      <SettingsPage />
    {/if}
  </main>
</div>

<style>
  .app {
    display: flex;
    min-height: 100vh;
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
    background: var(--card-bg);
    border-bottom: 1px solid var(--gray-200);
    padding: 0 1rem;
    align-items: center;
    z-index: 100;
  }
  
  .mobile-menu-btn {
    background: none;
    border: none;
    padding: 0.5rem;
    color: var(--gray-300);
    cursor: pointer;
    border-radius: 8px;
  }
  
  .mobile-menu-btn:hover {
    background: var(--bg-tertiary);
  }
  
  .logo-mobile {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-left: 1rem;
  }
  
  /* Sidebar */
  .sidebar {
    width: 280px;
    background: var(--card-bg);
    border-right: 1px solid var(--gray-200);
    display: flex;
    flex-direction: column;
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 200;
    box-shadow: 0 0 20px rgba(0, 0, 0, 0.05);
  }
  
  .sidebar.collapsed {
    width: 80px;
  }
  
  .sidebar-header {
    padding: 2rem 1.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .logo-collapsed {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
  }
  
  .logo-icon {
    width: 44px;
    height: 44px;
    background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
    color: white;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  }
  
  @keyframes shine {
    0% { transform: translateX(-100%) translateY(-100%) rotate(45deg); }
    100% { transform: translateX(100%) translateY(100%) rotate(45deg); }
  }
  
  .logo-text {
    font-size: 1.375rem;
    font-weight: 700;
    color: var(--gray-900);
    letter-spacing: -0.025em;
  }
  
  .collapse-btn {
    background: none;
    border: none;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--gray-400);
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.2s;
  }
  
  .collapse-btn:hover {
    background: var(--gray-100);
    color: var(--gray-700);
  }
  
  :global(.rotate-180) {
    transform: rotate(180deg);
  }
  
  /* Navigation */
  .nav {
    flex: 1;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.875rem;
    padding: 0.875rem 1rem;
    background: none;
    border: none;
    border-radius: 16px;
    color: var(--gray-600);
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
    width: 100%;
    position: relative;
    font-size: 0.9375rem;
  }
  
  .sidebar.collapsed .nav-item {
    justify-content: center;
    padding: 0.75rem;
  }
  
  .nav-item:hover {
    background: var(--gray-100);
    color: var(--gray-900);
  }
  
  .nav-item.active {
    background: var(--primary-100);
    color: var(--primary-700);
    font-weight: 600;
  }
  
  .nav-item.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 4px;
    background: var(--primary-500);
    border-radius: 0 4px 4px 0;
  }
  
  .nav-label {
    font-size: 0.9375rem;
  }
  
  /* User Info */
  .sidebar-footer {
    padding: 1.5rem;
    margin-top: auto;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .user-avatar,
  .user-avatar-small {
    width: 40px;
    height: 40px;
    background: var(--gray-200);
    color: var(--gray-700);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 0.875rem;
  }
  
  .user-avatar-small {
    margin: 0 auto;
  }
  
  .user-details {
    flex: 1;
  }
  
  .user-name {
    font-weight: 600;
    color: var(--gray-900);
    font-size: 0.9375rem;
  }
  
  .user-role {
    font-size: 0.8125rem;
    color: var(--gray-500);
  }
  
  /* Main Content */
  .main-content {
    flex: 1;
    margin-left: 280px;
    transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    min-height: 100vh;
  }
  
  .main-content.sidebar-collapsed {
    margin-left: 80px;
  }
  
  /* Mobile Styles */
  @media (max-width: 768px) {
    .mobile-header {
      display: flex;
    }
    
    .desktop-only {
      display: none;
    }
    
    .sidebar {
      transform: translateX(-100%);
      transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }
    
    .sidebar.mobile-open {
      transform: translateX(0);
    }
    
    .mobile-overlay {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.7);
      backdrop-filter: blur(4px);
      z-index: 150;
    }
    
    .main-content {
      margin-left: 0;
      padding-top: 60px;
    }
    
    .main-content.sidebar-collapsed {
      margin-left: 0;
    }
  }
</style>