<script lang="ts">
  import { onMount } from 'svelte';
  import { jobsStore, jobTemplatesStore, customersStore } from '../lib/stores';
  import { router } from '../lib/router';
  import { Plus, Search, Calendar, Clock, CheckCircle, AlertCircle,
           MapPin, DollarSign, Camera, FileText, Trash2, ChevronRight,
           TrendingUp, Users, Package, Shield, ShieldOff, Activity,
           Edit3, ExternalLink, LayoutGrid, List } from 'lucide-svelte';
  import type { Customer, Job } from '../lib/types/models';
  
  // Component imports
  import Button from '../lib/components/Button.svelte';
  import Card from '../lib/components/Card.svelte';
  import Modal from '../lib/components/Modal.svelte';
  import PhotoGalleryModal from '../lib/components/PhotoGalleryModal.svelte';
  import InvoiceModal from '../lib/components/InvoiceModal.svelte';
  import ConfirmModal from '../lib/components/ConfirmModal.svelte';
  
  let showCreateModal = false;
  let showPhotoGallery = false;
  let showInvoiceModal = false;
  let showDateModal = false;
  let selectedJob: Job | null = null;
  let editingJob: Job | null = null;
  let deleteConfirm = {
    show: false,
    job: null as Job | null
  };
  
  let dateForm = {
    scheduledDate: '',
    startDate: '',
    endDate: ''
  };
  
  let newJob = {
    customerName: '',
    address: '',
    templateId: '',
    startDate: '',
    endDate: '',
    permitRequired: false
  };
  
  let searchQuery = '';
  let statusFilter = 'all';
  let formError = '';
  let viewMode: 'cards' | 'list' = 'cards';
  
  // Load view preference from localStorage
  function loadViewPreference() {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('jobsViewMode');
      if (saved === 'list' || saved === 'cards') {
        viewMode = saved;
      }
    }
  }
  
  // Save view preference to localStorage
  function saveViewPreference(mode: 'cards' | 'list') {
    viewMode = mode;
    if (typeof window !== 'undefined') {
      localStorage.setItem('jobsViewMode', mode);
    }
  }
  
  // Load data
  onMount(async () => {
    // Load view preference
    loadViewPreference();
    
    // Load data
    await Promise.all([
      jobsStore.load(),
      jobTemplatesStore.loadActive(),
      customersStore.load()
    ]);
  });
  
  // Subscriptions
  $: ({ jobs, loading, error } = $jobsStore);
  $: ({ templates } = $jobTemplatesStore);
  $: ({ customers } = $customersStore);
  
  // Filtered jobs with enriched data
  $: filteredJobs = jobs
    .map(job => ({
      ...job,
      customer: customers.find(c => c.id === job.customerId),
      template: templates.find(t => t.id === job.templateId)
    }))
    .filter(job => {
      const matchesSearch = !searchQuery || 
        job.customer?.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        job.address.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesStatus = statusFilter === 'all' || job.status === statusFilter;
      return matchesSearch && matchesStatus;
    });
  
  // Calculate stats
  $: totalRevenue = jobs.reduce((sum, job) => sum + (job.totalAmount || 0), 0);
  $: activeJobs = jobs.filter(j => j.status === 'in_progress').length;
  $: completedThisMonth = jobs.filter(j => {
    if (j.status !== 'completed' || !j.endDate) return false;
    const endDate = new Date(j.endDate);
    const now = new Date();
    return endDate.getMonth() === now.getMonth() && endDate.getFullYear() === now.getFullYear();
  }).length;
  
  function formatCurrency(amount: number): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  }
  
  function formatDate(date: string | Date | undefined): string {
    if (!date) return 'Not scheduled';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  }
  
  function formatDateRange(start?: string | Date, end?: string | Date): string {
    if (!start && !end) return '';
    if (start && !end) return `Started ${formatDate(start)}`;
    if (!start && end) return `Ended ${formatDate(end)}`;
    
    const startDate = new Date(start!);
    const endDate = new Date(end!);
    
    // Same day
    if (startDate.toDateString() === endDate.toDateString()) {
      return formatDate(start);
    }
    
    // Different days
    return `${formatDate(start)} - ${formatDate(end)}`;
  }
  
  function getStatusBadgeClass(status: string): string {
    switch (status) {
      case 'scheduled': return 'badge-scheduled';
      case 'in_progress': return 'badge-progress';
      case 'completed': return 'badge-completed';
      default: return 'badge-secondary';
    }
  }
  
  function getStatusLabel(status: string): string {
    switch (status) {
      case 'scheduled': return 'Scheduled';
      case 'in_progress': return 'In Progress';
      case 'completed': return 'Completed';
      default: return status;
    }
  }
  
  function openJobDetail(job: Job) {
    router.navigate('job-detail', { id: job.id });
  }
  
  function openPhotoGallery(event: Event, job: Job) {
    event.stopPropagation();
    selectedJob = job;
    showPhotoGallery = true;
  }
  
  function openInvoice(event: Event, job: Job) {
    event.stopPropagation();
    selectedJob = job;
    showInvoiceModal = true;
  }
  
  function openDateEdit(event: Event, job: Job) {
    event.stopPropagation();
    editingJob = job;
    dateForm = {
      scheduledDate: job.scheduledDate ? new Date(job.scheduledDate).toISOString().split('T')[0] : '',
      startDate: job.startDate ? new Date(job.startDate).toISOString().split('T')[0] : '',
      endDate: job.endDate ? new Date(job.endDate).toISOString().split('T')[0] : ''
    };
    showDateModal = true;
  }
  
  async function togglePermitRequired(event: Event, job: Job) {
    event.stopPropagation();
    try {
      await jobsStore.update(job.id, {
        permitRequired: !job.permitRequired,
        permitNumber: !job.permitRequired ? '' : job.permitNumber
      });
    } catch (error) {
      console.error('Failed to update permit status:', error);
    }
  }
  
  async function handleDateUpdate() {
    if (!editingJob) return;
    
    try {
      const updates: any = {};
      
      if (dateForm.scheduledDate) {
        updates.scheduledDate = new Date(dateForm.scheduledDate);
      }
      if (dateForm.startDate) {
        updates.startDate = new Date(dateForm.startDate);
        updates.status = 'in_progress';
      }
      if (dateForm.endDate) {
        updates.endDate = new Date(dateForm.endDate);
        updates.status = 'completed';
      }
      
      await jobsStore.update(editingJob.id, updates);
      showDateModal = false;
      editingJob = null;
    } catch (error) {
      console.error('Failed to update dates:', error);
    }
  }
  
  function confirmDeleteJob(event: Event, job: Job) {
    event.stopPropagation();
    deleteConfirm = { show: true, job };
  }
  
  async function handleDeleteJob() {
    if (deleteConfirm.job) {
      try {
        await jobsStore.remove(deleteConfirm.job.id);
        deleteConfirm = { show: false, job: null };
      } catch (error) {
        console.error('Failed to delete job:', error);
      }
    }
  }
  
  async function handleCreateJob() {
    formError = '';
    
    if (!newJob.customerName || !newJob.address || !newJob.templateId) {
      formError = 'Please fill in all required fields';
      return;
    }
    
    try {
      // Create or find customer
      let customer = customers.find(c => 
        c.name.toLowerCase() === newJob.customerName.toLowerCase()
      );
      
      if (!customer) {
        // Create new customer with minimal info
        customer = await customersStore.create(
          newJob.customerName,
          'no-email@placeholder.com', // Required by API but we don't have it
          ''  // No phone
        );
      }
      
      // Prepare dates
      const scheduledDate = newJob.startDate ? new Date(newJob.startDate) : undefined;
      
      // Create the job
      const jobId = await jobsStore.createFromTemplate(
        customer.id,
        newJob.address,
        newJob.templateId,
        scheduledDate
      );
      
      // If permit is required or end date is set, update the job
      if (newJob.permitRequired || newJob.endDate) {
        const updates: any = {};
        if (newJob.permitRequired) {
          updates.permitRequired = true;
        }
        if (newJob.endDate) {
          updates.endDate = new Date(newJob.endDate);
        }
        await jobsStore.update(jobId, updates);
      }
      
      // Reset form
      newJob = {
        customerName: '',
        address: '',
        templateId: '',
        startDate: '',
        endDate: '',
        permitRequired: false
      };
      showCreateModal = false;
      
      // Navigate to the new job
      router.navigate('job-detail', { id: jobId });
    } catch (error) {
      console.error('Failed to create job:', error);
      formError = error instanceof Error ? error.message : 'Failed to create job';
    }
  }
  
  // Helper to count completed items
  function getCompletedItemsCount(job: Job): number {
    return job.items?.filter(item => item.quantity > 0).length || 0;
  }
  
  // Helper to get progress percentage
  function getProgressPercentage(job: Job): number {
    if (!job.items || job.items.length === 0) return 0;
    const completed = getCompletedItemsCount(job);
    return Math.round((completed / job.items.length) * 100);
  }
  
  // Open GPS/Maps for address
  function openGPS(address: string) {
    const encodedAddress = encodeURIComponent(address);
    window.open(`https://maps.google.com/?q=${encodedAddress}`, '_blank');
  }
</script>

<div class="dashboard">
  <!-- Header -->
  <div class="dashboard-header">
    <div>
      <h1>Jobs</h1>
      <p>Manage your electrical projects</p>
    </div>
    <Button on:click={() => showCreateModal = true}>
      <Plus size={20} />
      Create Job
    </Button>
  </div>

  <!-- Stats Grid -->
  <div class="stats-grid">
    <div class="stat-card">
      <div class="stat-icon stat-icon-primary">
        <Package size={24} />
      </div>
      <div class="stat-content">
        <p class="stat-label">Total Jobs</p>
        <p class="stat-value">{jobs.length}</p>
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-icon stat-icon-success">
        <TrendingUp size={24} />
      </div>
      <div class="stat-content">
        <p class="stat-label">Total Revenue</p>
        <p class="stat-value">{formatCurrency(totalRevenue)}</p>
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

  <!-- Search and Filters -->
  <div class="controls">
    <div class="search-wrapper">
      <Search size={20} class="search-icon" />
      <input
        type="text"
        placeholder="Search by customer or address..."
        bind:value={searchQuery}
        class="search-input"
      />
    </div>

    <div class="filter-tabs">
      <button 
        class="filter-tab"
        class:active={statusFilter === 'all'}
        on:click={() => statusFilter = 'all'}
      >
        All ({jobs.length})
      </button>
      <button 
        class="filter-tab"
        class:active={statusFilter === 'scheduled'}
        on:click={() => statusFilter = 'scheduled'}
      >
        Scheduled ({jobs.filter(j => j.status === 'scheduled').length})
      </button>
      <button 
        class="filter-tab"
        class:active={statusFilter === 'in_progress'}
        on:click={() => statusFilter = 'in_progress'}
      >
        In Progress ({jobs.filter(j => j.status === 'in_progress').length})
      </button>
      <button 
        class="filter-tab"
        class:active={statusFilter === 'completed'}
        on:click={() => statusFilter = 'completed'}
      >
        Completed ({jobs.filter(j => j.status === 'completed').length})
      </button>
    </div>
    
    <div class="view-toggle">
      <button
        class="view-btn"
        class:active={viewMode === 'cards'}
        on:click={() => saveViewPreference('cards')}
        title="Card view"
      >
        <LayoutGrid size={18} />
      </button>
      <button
        class="view-btn"
        class:active={viewMode === 'list'}
        on:click={() => saveViewPreference('list')}
        title="List view"
      >
        <List size={18} />
      </button>
    </div>
  </div>

  <!-- Jobs List -->
  <div class="jobs-container">
    {#if loading}
      <div class="empty-state">
        <div class="spinner"></div>
        <p>Loading jobs...</p>
      </div>
    {:else if error}
      <div class="empty-state">
        <AlertCircle size={48} class="empty-icon" />
        <h3>Error loading jobs</h3>
        <p>{error}</p>
        <Button on:click={() => jobsStore.load()} variant="secondary">
          Try Again
        </Button>
      </div>
    {:else if filteredJobs.length === 0}
      <div class="empty-state">
        <Package size={48} class="empty-icon" />
        <h3>No jobs found</h3>
        <p>{searchQuery ? 'Try adjusting your search' : 'Create your first job to get started'}</p>
        {#if !searchQuery}
          <Button on:click={() => showCreateModal = true}>
            <Plus size={20} />
            Create Job
          </Button>
        {/if}
      </div>
    {:else}
      {#if viewMode === 'cards'}
        <div class="jobs-grid">
          {#each filteredJobs as job}
          <div class="job-card">
            <!-- Status Badge -->
            <div class="job-status-badge {getStatusBadgeClass(job.status)}">
              {getStatusLabel(job.status)}
            </div>

            <!-- Customer & Template -->
            <div class="job-header">
              <h3>{job.customer?.name || 'Unknown Customer'} - {job.address.split(',')[1]?.trim() || 'Location'}</h3>
              <div class="job-template-row">
                <span class="template-badge">{job.template?.name || 'Custom Job'}</span>
                <button 
                  class="permit-toggle"
                  class:permit-required={job.permitRequired}
                  on:click={(e) => togglePermitRequired(e, job)}
                >
                  {#if job.permitRequired}
                    <Shield size={16} />
                    <span>Permit Required</span>
                  {:else}
                    <ShieldOff size={16} />
                    <span>No Permit</span>
                  {/if}
                </button>
              </div>
            </div>

            <!-- Amount -->
            <div class="job-amount-section">
              <span class="amount">{formatCurrency(job.totalAmount || 0)}</span>
            </div>

            <!-- Address -->
            <div class="job-address" on:click={(e) => { e.stopPropagation(); openGPS(job.address); }}>
              <MapPin size={16} />
              <span>{job.address}</span>
            </div>

            <!-- Date -->
            <div class="job-date" on:click={(e) => openDateEdit(e, job)}>
              <Calendar size={16} />
              <span>
                {#if job.startDate && job.endDate}
                  {formatDateRange(job.startDate, job.endDate)}
                {:else if job.scheduledDate}
                  {formatDate(job.scheduledDate)}
                {:else}
                  Not scheduled
                {/if}
              </span>
              <Edit3 size={14} class="date-edit-icon" />
            </div>

            <!-- Progress -->
            <div class="job-progress">
              <div class="progress-info">
                <span>{getCompletedItemsCount(job)} items completed</span>
                <span class="progress-percent">{getProgressPercentage(job)}%</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" style="width: {getProgressPercentage(job)}%"></div>
              </div>
            </div>

            <!-- Wave Invoice -->
            {#if job.waveInvoiceId}
              <div class="wave-invoice">
                <span>Wave Invoice</span>
                <a href={job.waveInvoiceUrl} target="_blank" on:click|stopPropagation class="wave-link">
                  {job.waveInvoiceId}
                  <ExternalLink size={14} />
                </a>
              </div>
            {/if}

            <!-- Actions -->
            <div class="job-actions">
              <button class="action-btn" on:click={(e) => openPhotoGallery(e, job)}>
                <Camera size={18} />
                <span>Gallery</span>
              </button>
              <button class="action-btn" on:click={(e) => openInvoice(e, job)}>
                <FileText size={18} />
                <span>Invoice</span>
              </button>
            </div>

            <!-- Footer Actions -->
            <div class="job-footer">
              <Button fullWidth on:click={() => openJobDetail(job)}>
                View Job
              </Button>
              <button class="delete-btn" on:click={(e) => confirmDeleteJob(e, job)}>
                <Trash2 size={18} />
                Delete Job
              </button>
            </div>
          </div>
        {/each}
      </div>
      {:else}
        <div class="jobs-list">
          {#each filteredJobs as job}
            <div class="job-list-item" on:click={() => openJobDetail(job)}>
              <div class="list-item-left">
                <div class="list-item-header">
                  <h4>{job.customer?.name || 'Unknown Customer'}</h4>
                  <span class="list-badge {getStatusBadgeClass(job.status)}">{getStatusLabel(job.status)}</span>
                </div>
                <div class="list-item-details">
                  <span class="list-template">{job.template?.name || 'Custom Job'}</span>
                  <span class="list-address">{job.address}</span>
                </div>
                <div class="list-item-meta">
                  <span>{formatCurrency(job.totalAmount || 0)}</span>
                  <span>•</span>
                  <span>{getCompletedItemsCount(job)}/{job.items?.length || 0} items</span>
                  <span>•</span>
                  <span>
                    {#if job.startDate && job.endDate}
                      {formatDateRange(job.startDate, job.endDate)}
                    {:else if job.scheduledDate}
                      {formatDate(job.scheduledDate)}
                    {:else}
                      Not scheduled
                    {/if}
                  </span>
                </div>
              </div>
              <div class="list-item-right">
                <ChevronRight size={20} />
              </div>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</div>
<!-- Date Edit Modal -->
<Modal bind:isOpen={showDateModal} title="Edit Job Dates" size="sm">
  <form on:submit|preventDefault={handleDateUpdate}>
    <div class="form-group">
      <label>Scheduled Date</label>
      <input
        type="date"
        bind:value={dateForm.scheduledDate}
      />
    </div>
    
    <div class="form-group">
      <label>Start Date</label>
      <input
        type="date"
        bind:value={dateForm.startDate}
      />
    </div>
    
    <div class="form-group">
      <label>End Date</label>
      <input
        type="date"
        bind:value={dateForm.endDate}
      />
    </div>
  </form>
  
  <div slot="footer" class="modal-footer">
    <Button variant="ghost" on:click={() => showDateModal = false}>
      Cancel
    </Button>
    <Button on:click={handleDateUpdate}>
      Update Dates
    </Button>
  </div>
</Modal>

<!-- Create Job Modal -->
<Modal bind:isOpen={showCreateModal} title="Add New Job">
  {#if formError}
    <div class="error-message">
      {formError}
    </div>
  {/if}
  
  <form on:submit|preventDefault={handleCreateJob}>
    <div class="form-group">
      <label>Customer Name *</label>
      <input
        type="text"
        bind:value={newJob.customerName}
        placeholder="Enter customer name"
        required
      />
    </div>
    
    <div class="form-group">
      <label>Address *</label>
      <textarea
        bind:value={newJob.address}
        placeholder="Enter full address"
        rows="3"
        required
      ></textarea>
    </div>
    
    <div class="form-group">
      <label>Project Template *</label>
      <select bind:value={newJob.templateId} required>
        <option value="">Select project template</option>
        {#each templates as template}
          <option value={template.id}>{template.name}</option>
        {/each}
      </select>
    </div>
    
    <div class="form-row">
      <div class="form-group">
        <label>Start Date</label>
        <input
          type="date"
          bind:value={newJob.startDate}
        />
      </div>
      
      <div class="form-group">
        <label>End Date</label>
        <input
          type="date"
          bind:value={newJob.endDate}
        />
      </div>
    </div>
    
    <div class="form-group">
      <label class="checkbox-label">
        <input
          type="checkbox"
          bind:checked={newJob.permitRequired}
        />
        Permit Required
      </label>
    </div>
  </form>
  
  <div slot="footer" class="modal-footer">
    <Button variant="ghost" on:click={() => showCreateModal = false}>
      Cancel
    </Button>
    <Button on:click={handleCreateJob}>
      Create Job
    </Button>
  </div>
</Modal>

{#if selectedJob}
  <PhotoGalleryModal bind:isOpen={showPhotoGallery} job={selectedJob} />
  <InvoiceModal bind:isOpen={showInvoiceModal} job={selectedJob} />
{/if}

<ConfirmModal
  bind:isOpen={deleteConfirm.show}
  title="Delete Job"
  message={`Are you sure you want to delete this job for ${deleteConfirm.job?.customer?.name || 'this customer'}?`}
  confirmText="Delete"
  onConfirm={handleDeleteJob}
  onCancel={() => deleteConfirm = { show: false, job: null }}
/>

<style>
  .dashboard {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
  }

  .dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .dashboard-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0;
  }

  .dashboard-header p {
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

  /* Search and Controls */
  .controls {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
    align-items: center;
  }

  .search-wrapper {
    position: relative;
    flex: 1;
    max-width: 500px;
  }

  :global(.search-icon) {
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-tertiary);
    pointer-events: none;
  }

  .search-input {
    width: 100%;
    padding: 0.875rem 1rem 0.875rem 3rem;
    border: 2px solid var(--gray-200);
    border-radius: var(--radius-lg);
    font-size: 0.9375rem;
    background: white;
    transition: all var(--transition-base);
  }

  .search-input::placeholder {
    color: var(--text-tertiary);
  }

  .search-input:focus {
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(91, 91, 214, 0.1);
    outline: none;
  }

  .filter-tabs {
    display: flex;
    gap: 0.5rem;
    background: white;
    padding: 0.25rem;
    border-radius: var(--radius-lg);
    border: 1px solid var(--gray-200);
  }

  .filter-tab {
    padding: 0.5rem 1rem;
    border: none;
    background: transparent;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all var(--transition-base);
    white-space: nowrap;
  }

  .filter-tab:hover {
    color: var(--text-primary);
  }

  .filter-tab.active {
    background: var(--primary-500);
    color: white;
  }

  /* View Toggle */
  .view-toggle {
    display: flex;
    gap: 0.25rem;
    background: white;
    padding: 0.25rem;
    border-radius: var(--radius-lg);
    border: 1px solid var(--gray-200);
  }

  .view-btn {
    padding: 0.5rem;
    border: none;
    background: transparent;
    border-radius: var(--radius-md);
    cursor: pointer;
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all var(--transition-base);
  }

  .view-btn:hover {
    color: var(--text-primary);
    background: var(--gray-50);
  }

  .view-btn.active {
    background: var(--primary-500);
    color: white;
  }

  /* Jobs Container */
  .jobs-container {
    min-height: 400px;
  }

  .jobs-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
    gap: 1.5rem;
  }

  /* List View */
  .jobs-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .job-list-item {
    background: white;
    border: 1px solid var(--gray-100);
    border-radius: var(--radius-lg);
    padding: 1.25rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    transition: all var(--transition-base);
    box-shadow: var(--shadow-sm);
  }

  .job-list-item:hover {
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
    border-color: var(--gray-200);
  }

  .list-item-left {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .list-item-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .list-item-header h4 {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-primary);
  }

  .list-badge {
    padding: 0.25rem 0.625rem;
    border-radius: var(--radius-md);
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
  }

  .list-item-details {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .list-template {
    background: var(--gray-100);
    color: var(--text-secondary);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-md);
    font-size: 0.8125rem;
    font-weight: 500;
  }

  .list-address {
    color: var(--text-secondary);
    font-size: 0.875rem;
  }

  .list-item-meta {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-tertiary);
    font-size: 0.8125rem;
  }

  .list-item-right {
    color: var(--text-tertiary);
  }

  .job-card {
    background: white;
    border: 1px solid var(--gray-100);
    border-radius: var(--radius-xl);
    padding: 1.5rem;
    position: relative;
    transition: all var(--transition-base);
    box-shadow: var(--shadow-sm);
  }

  .job-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
    border-color: var(--gray-200);
  }

  /* Status Badge */
  .job-status-badge {
    position: absolute;
    top: 1rem;
    right: 1rem;
    padding: 0.375rem 0.75rem;
    border-radius: var(--radius-md);
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .badge-scheduled {
    background: rgba(245, 158, 11, 0.1);
    color: var(--warning-500);
  }

  .badge-progress {
    background: rgba(91, 91, 214, 0.1);
    color: var(--primary-500);
  }

  .badge-completed {
    background: rgba(16, 185, 129, 0.1);
    color: var(--success-500);
  }

  /* Job Header */
  .job-header {
    margin-bottom: 1rem;
  }

  .job-header h3 {
    font-size: 1.125rem;
    font-weight: 600;
    margin: 0 0 0.5rem 0;
    color: var(--text-primary);
    padding-right: 100px;
  }

  .job-template-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .template-badge {
    background: var(--gray-100);
    color: var(--text-secondary);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-md);
    font-size: 0.8125rem;
    font-weight: 500;
  }

  .permit-toggle {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.25rem 0.75rem;
    background: var(--gray-100);
    border: none;
    border-radius: var(--radius-md);
    font-size: 0.8125rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all var(--transition-base);
  }

  .permit-toggle:hover {
    background: var(--gray-200);
  }

  .permit-toggle.permit-required {
    background: rgba(245, 158, 11, 0.1);
    color: var(--warning-500);
  }

  /* Amount */
  .job-amount-section {
    margin-bottom: 1rem;
  }

  .amount {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--text-primary);
  }

  /* Address */
  .job-address {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-secondary);
    font-size: 0.875rem;
    margin-bottom: 0.75rem;
    cursor: pointer;
    padding: 0.375rem;
    margin: 0 -0.375rem 0.75rem;
    border-radius: var(--radius-md);
    transition: all var(--transition-base);
  }

  .job-address:hover {
    background: var(--gray-50);
    color: var(--primary-500);
  }

  .job-address :global(svg) {
    color: var(--text-tertiary);
    flex-shrink: 0;
  }

  .job-address:hover :global(svg) {
    color: var(--primary-500);
  }

  /* Date */
  .job-date {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-secondary);
    font-size: 0.875rem;
    margin-bottom: 1rem;
    cursor: pointer;
    padding: 0.375rem;
    margin: 0 -0.375rem 1rem;
    border-radius: var(--radius-md);
    transition: all var(--transition-base);
  }

  .job-date:hover {
    background: var(--gray-50);
  }

  .job-date :global(svg) {
    color: var(--text-tertiary);
    flex-shrink: 0;
  }

  :global(.date-edit-icon) {
    margin-left: auto;
    opacity: 0;
    transition: opacity var(--transition-base);
  }

  .job-date:hover :global(.date-edit-icon) {
    opacity: 1;
  }

  /* Progress */
  .job-progress {
    margin-bottom: 1rem;
  }

  .progress-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .progress-percent {
    font-weight: 600;
    color: var(--text-primary);
  }

  .progress-bar {
    height: 8px;
    background: var(--gray-200);
    border-radius: var(--radius-full);
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: var(--primary-500);
    border-radius: var(--radius-full);
    transition: width var(--transition-base);
  }

  /* Wave Invoice */
  .wave-invoice {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: var(--gray-50);
    border-radius: var(--radius-md);
    margin-bottom: 1rem;
    font-size: 0.875rem;
  }

  .wave-invoice span {
    color: var(--text-secondary);
  }

  .wave-link {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    color: var(--primary-500);
    text-decoration: none;
    font-weight: 500;
    transition: color var(--transition-base);
  }

  .wave-link:hover {
    color: var(--primary-600);
  }

  /* Actions */
  .job-actions {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .action-btn {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.625rem;
    background: var(--gray-100);
    border: none;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all var(--transition-base);
  }

  .action-btn:hover {
    background: var(--gray-200);
    color: var(--text-primary);
  }

  /* Footer */
  .job-footer {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .delete-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    width: 100%;
    padding: 0.75rem;
    background: transparent;
    border: 1px solid var(--danger-500);
    border-radius: var(--radius-md);
    color: var(--danger-500);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-base);
  }

  .delete-btn:hover {
    background: var(--danger-500);
    color: white;
  }

  /* Empty State */
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
  }

  :global(.empty-icon) {
    color: var(--text-tertiary);
    margin-bottom: 1rem;
  }

  .empty-state h3 {
    font-size: 1.25rem;
    margin: 0 0 0.5rem 0;
  }

  .empty-state p {
    margin: 0 0 1.5rem 0;
  }

  /* Spinner */
  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid var(--gray-200);
    border-top-color: var(--primary-500);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Forms */
  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
  }

  .form-group input[type="text"],
  .form-group input[type="date"],
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 0.625rem 0.875rem;
    border: 2px solid var(--gray-200);
    border-radius: var(--radius-md);
    font-size: 0.9375rem;
    background: white;
    transition: all var(--transition-base);
  }

  .form-group input:focus,
  .form-group select:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(91, 91, 214, 0.1);
  }

  .form-group textarea {
    resize: vertical;
    min-height: 80px;
  }

  .form-group small {
    display: block;
    font-size: 0.8125rem;
    color: var(--text-tertiary);
    margin-top: 0.25rem;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    font-size: 0.9375rem;
    color: var(--text-primary);
  }

  .checkbox-label input[type="checkbox"] {
    width: auto;
    margin: 0;
    cursor: pointer;
  }

  .form-section {
    margin-bottom: 1.5rem;
  }

  .form-section h4 {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 1rem 0;
  }

  .form-tabs {
    display: flex;
    gap: 0;
    background: var(--gray-100);
    padding: 0.25rem;
    border-radius: var(--radius-md);
    margin-bottom: 1rem;
  }

  .form-tab {
    flex: 1;
    padding: 0.5rem;
    border: none;
    background: transparent;
    border-radius: var(--radius-sm);
    font-size: 0.875rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all var(--transition-base);
  }

  .form-tab.active {
    background: white;
    color: var(--text-primary);
    box-shadow: var(--shadow-sm);
  }

  .form-group {
    margin-bottom: 1rem;
  }

  .form-group label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
  }

  .form-group input,
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--gray-200);
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    background: white;
    transition: all var(--transition-base);
  }

  .form-group input:focus,
  .form-group select:focus,
  .form-group textarea:focus {
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(91, 91, 214, 0.1);
  }

  .form-group textarea {
    resize: vertical;
    font-family: inherit;
  }

  .form-group small {
    display: block;
    font-size: 0.75rem;
    color: var(--text-tertiary);
    margin-top: 0.25rem;
  }

  .error-message {
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.2);
    color: var(--danger-500);
    padding: 0.75rem;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }

  .modal-footer {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
  }

  /* Responsive */
  @media (max-width: 1200px) {
    .jobs-grid {
      grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    }
  }

  @media (max-width: 768px) {
    .dashboard {
      padding: 1rem;
    }

    .dashboard-header {
      flex-direction: column;
      align-items: start;
      gap: 1rem;
    }

    .stats-grid {
      grid-template-columns: 1fr 1fr;
    }

    .controls {
      flex-direction: column;
    }

    .search-wrapper {
      max-width: 100%;
    }

    .filter-tabs {
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
    }

    .jobs-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 480px) {
    .stats-grid {
      grid-template-columns: 1fr;
    }
    
    .job-header h3 {
      font-size: 1rem;
    }
    
    .amount {
      font-size: 1.5rem;
    }
  }
</style>