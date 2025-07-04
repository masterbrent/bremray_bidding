<script lang="ts">
  import { onMount } from 'svelte';
  import { jobsStore, customersStore } from '../lib/stores';
  import { TrendingUp, Users, Package, CheckCircle, Clock, DollarSign, Calendar, Activity } from 'lucide-svelte';
  
  // Component imports
  import Card from '../lib/components/Card.svelte';
  
  // Subscriptions
  $: ({ jobs } = $jobsStore);
  $: ({ customers } = $customersStore);
  
  // Calculate stats
  $: totalRevenue = jobs.reduce((sum, job) => sum + (job.totalAmount || 0), 0);
  $: totalJobs = jobs.length;
  $: activeJobs = jobs.filter(j => j.status === 'in_progress').length;
  $: scheduledJobs = jobs.filter(j => j.status === 'scheduled').length;
  $: completedJobs = jobs.filter(j => j.status === 'completed').length;
  $: completedThisMonth = jobs.filter(j => {
    if (j.status !== 'completed' || !j.endDate) return false;
    const endDate = new Date(j.endDate);
    const now = new Date();
    return endDate.getMonth() === now.getMonth() && endDate.getFullYear() === now.getFullYear();
  }).length;
  
  // Revenue by month
  $: revenueByMonth = calculateRevenueByMonth(jobs);
  
  // Jobs by status
  $: jobsByStatus = {
    scheduled: scheduledJobs,
    in_progress: activeJobs,
    completed: completedJobs
  };
  
  // Top customers by revenue
  $: topCustomers = calculateTopCustomers(jobs, customers);
  
  function formatCurrency(amount: number): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  }
  
  function calculateRevenueByMonth(jobs: any[]) {
    const monthlyRevenue: Record<string, number> = {};
    const now = new Date();
    
    // Initialize last 6 months
    for (let i = 5; i >= 0; i--) {
      const date = new Date(now.getFullYear(), now.getMonth() - i, 1);
      const key = date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
      monthlyRevenue[key] = 0;
    }
    
    // Calculate revenue
    jobs.forEach(job => {
      if (job.status === 'completed' && job.endDate) {
        const date = new Date(job.endDate);
        const key = date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
        if (monthlyRevenue.hasOwnProperty(key)) {
          monthlyRevenue[key] += job.totalAmount || 0;
        }
      }
    });
    
    return Object.entries(monthlyRevenue).map(([month, revenue]) => ({
      month,
      revenue
    }));
  }
  
  function calculateTopCustomers(jobs: any[], customers: any[]) {
    const customerRevenue: Record<string, number> = {};
    
    jobs.forEach(job => {
      if (job.customerId) {
        customerRevenue[job.customerId] = (customerRevenue[job.customerId] || 0) + (job.totalAmount || 0);
      }
    });
    
    return Object.entries(customerRevenue)
      .map(([customerId, revenue]) => ({
        customer: customers.find(c => c.id === customerId),
        revenue
      }))
      .filter(item => item.customer)
      .sort((a, b) => b.revenue - a.revenue)
      .slice(0, 5);
  }
  
  onMount(async () => {
    await Promise.all([
      jobsStore.load(),
      customersStore.load()
    ]);
  });
</script>

<div class="stats-page">
  <div class="page-header">
    <h1>Statistics</h1>
    <p>Overview of your business performance</p>
  </div>

  <!-- Key Metrics -->
  <div class="stats-grid">
    <div class="stat-card">
      <div class="stat-icon stat-icon-success">
        <DollarSign size={24} />
      </div>
      <div class="stat-content">
        <p class="stat-label">Total Revenue</p>
        <p class="stat-value">{formatCurrency(totalRevenue)}</p>
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-icon stat-icon-primary">
        <Package size={24} />
      </div>
      <div class="stat-content">
        <p class="stat-label">Total Jobs</p>
        <p class="stat-value">{totalJobs}</p>
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-icon stat-icon-warning">
        <Clock size={24} />
      </div>
      <div class="stat-content">
        <p class="stat-label">Active Jobs</p>
        <p class="stat-value">{activeJobs}</p>
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-icon stat-icon-info">
        <CheckCircle size={24} />
      </div>
      <div class="stat-content">
        <p class="stat-label">Completed This Month</p>
        <p class="stat-value">{completedThisMonth}</p>
      </div>
    </div>
  </div>

  <!-- Detailed Stats -->
  <div class="detailed-stats">
    <!-- Jobs by Status -->
    <Card>
      <h2>Jobs by Status</h2>
      <div class="status-breakdown">
        <div class="status-item">
          <div class="status-bar scheduled" style="width: {(scheduledJobs / totalJobs) * 100}%"></div>
          <div class="status-info">
            <span>Scheduled</span>
            <span>{scheduledJobs}</span>
          </div>
        </div>
        <div class="status-item">
          <div class="status-bar in-progress" style="width: {(activeJobs / totalJobs) * 100}%"></div>
          <div class="status-info">
            <span>In Progress</span>
            <span>{activeJobs}</span>
          </div>
        </div>
        <div class="status-item">
          <div class="status-bar completed" style="width: {(completedJobs / totalJobs) * 100}%"></div>
          <div class="status-info">
            <span>Completed</span>
            <span>{completedJobs}</span>
          </div>
        </div>
      </div>
    </Card>

    <!-- Revenue Trend -->
    <Card>
      <h2>Revenue Trend (Last 6 Months)</h2>
      <div class="revenue-chart">
        {#each revenueByMonth as { month, revenue }}
          <div class="revenue-bar">
            <div class="bar" style="height: {(revenue / Math.max(...revenueByMonth.map(r => r.revenue))) * 200}px">
              <span class="bar-value">{formatCurrency(revenue)}</span>
            </div>
            <span class="bar-label">{month}</span>
          </div>
        {/each}
      </div>
    </Card>

    <!-- Top Customers -->
    <Card>
      <h2>Top Customers by Revenue</h2>
      <div class="customer-list">
        {#each topCustomers as { customer, revenue }, index}
          <div class="customer-item">
            <div class="customer-rank">{index + 1}</div>
            <div class="customer-info">
              <h4>{customer.name}</h4>
              <p>{formatCurrency(revenue)}</p>
            </div>
          </div>
        {/each}
        {#if topCustomers.length === 0}
          <p class="empty-message">No customer data available</p>
        {/if}
      </div>
    </Card>
  </div>
</div>

<style>
  .stats-page {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
  }

  .page-header {
    margin-bottom: 2rem;
  }

  .page-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0;
  }

  .page-header p {
    color: var(--text-secondary);
    margin: 0.25rem 0 0 0;
  }

  /* Stats Grid */
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
  }

  .stat-card {
    background: white;
    border-radius: var(--radius-lg);
    padding: 1.5rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    box-shadow: var(--shadow-sm);
    border: 1px solid var(--gray-100);
  }

  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: var(--radius-md);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .stat-icon-primary {
    background: rgba(91, 91, 214, 0.1);
    color: var(--primary-500);
  }

  .stat-icon-success {
    background: rgba(16, 185, 129, 0.1);
    color: var(--success-500);
  }

  .stat-icon-warning {
    background: rgba(245, 158, 11, 0.1);
    color: var(--warning-500);
  }

  .stat-icon-info {
    background: rgba(59, 130, 246, 0.1);
    color: #3B82F6;
  }

  .stat-content {
    flex: 1;
  }

  .stat-label {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin: 0;
  }

  .stat-value {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0;
  }

  /* Detailed Stats */
  .detailed-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 1.5rem;
  }

  .detailed-stats h2 {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 1rem 0;
  }

  /* Status Breakdown */
  .status-breakdown {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .status-item {
    position: relative;
  }

  .status-bar {
    height: 32px;
    border-radius: var(--radius-md);
    transition: width 0.3s ease;
  }

  .status-bar.scheduled {
    background: rgba(245, 158, 11, 0.2);
  }

  .status-bar.in-progress {
    background: rgba(91, 91, 214, 0.2);
  }

  .status-bar.completed {
    background: rgba(16, 185, 129, 0.2);
  }

  .status-info {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    left: 1rem;
    right: 1rem;
    display: flex;
    justify-content: space-between;
    font-size: 0.875rem;
    font-weight: 500;
  }

  /* Revenue Chart */
  .revenue-chart {
    display: flex;
    align-items: flex-end;
    justify-content: space-around;
    height: 250px;
    padding: 1rem 0;
    gap: 0.5rem;
  }

  .revenue-bar {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }

  .bar {
    width: 100%;
    max-width: 60px;
    background: var(--primary-500);
    border-radius: var(--radius-md) var(--radius-md) 0 0;
    position: relative;
    min-height: 20px;
    transition: height 0.3s ease;
  }

  .bar-value {
    position: absolute;
    top: -1.5rem;
    left: 50%;
    transform: translateX(-50%);
    font-size: 0.75rem;
    font-weight: 600;
    white-space: nowrap;
  }

  .bar-label {
    font-size: 0.75rem;
    color: var(--text-secondary);
  }

  /* Customer List */
  .customer-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .customer-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem;
    background: var(--gray-50);
    border-radius: var(--radius-md);
  }

  .customer-rank {
    width: 32px;
    height: 32px;
    background: var(--primary-500);
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 0.875rem;
  }

  .customer-info h4 {
    margin: 0;
    font-size: 0.9375rem;
    font-weight: 500;
    color: var(--text-primary);
  }

  .customer-info p {
    margin: 0.25rem 0 0 0;
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .empty-message {
    text-align: center;
    color: var(--text-tertiary);
    font-style: italic;
    padding: 2rem;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .stats-page {
      padding: 1rem;
    }

    .detailed-stats {
      grid-template-columns: 1fr;
    }
  }
</style>
