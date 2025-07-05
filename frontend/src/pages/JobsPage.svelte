<script lang="ts">
  import { onMount } from 'svelte';
  import { jobsStore, jobTemplatesStore, customersStore, userStore, effectiveRole, permissions } from '../lib/stores';
  import { router } from '../lib/router';
  import { Plus, Search, Calendar, Clock, CheckCircle, AlertCircle,
           MapPin, DollarSign, Camera, FileText, Trash2, ChevronRight,
           Users, Package, Shield, ShieldOff, Activity, Send,
           Edit3, ExternalLink, LayoutGrid, List } from 'lucide-svelte';
  import type { Customer, Job, TemplatePhase } from '../lib/types/models';
  import { WaveService } from '../lib/services/waveService';
  
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
  let categoryFilter = 'all'; // New filter for skyview/contractors/rayno
  
  // Initialize viewMode from localStorage if available
  let viewMode: 'cards' | 'list' = 'cards';
  if (typeof window !== 'undefined' && window.localStorage) {
    const savedViewMode = localStorage.getItem('jobsViewMode');
    if (savedViewMode === 'cards' || savedViewMode === 'list') {
      viewMode = savedViewMode;
    }
  }
  
  let formError = '';
  let sendingToWave: { [jobId: string]: boolean } = {};
  let waveError: { [jobId: string]: string } = {};
  
  // Update localStorage when viewMode changes
  $: if (typeof window !== 'undefined' && window.localStorage && viewMode) {
    localStorage.setItem('jobsViewMode', viewMode);
  }
  
  // Load data
  onMount(async () => {
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
        job.customer?.name?.toLowerCase().startsWith(searchQuery.toLowerCase());
      const matchesStatus = statusFilter === 'all' || job.status === statusFilter;
      return matchesSearch && matchesStatus;
    });
  
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
  
  function getCurrentPhase(job: Job): TemplatePhase | null {
    if (!job.currentPhaseId || !job.template?.phases) return null;
    return job.template.phases.find(p => p.id === job.currentPhaseId) || null;
  }
  
  async function updateJobPhase(event: Event, job: Job, phaseId: string) {
    event.stopPropagation();
    
    try {
      await jobsStore.update(job.id, {
        currentPhaseId: phaseId || null
      });
    } catch (error) {
      console.error('Failed to update job phase:', error);
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
  
  async function sendToWave(event: Event, job: Job) {
    event.stopPropagation();
    
    if (sendingToWave[job.id]) return;
    
    // Clear any previous error
    delete waveError[job.id];
    
    try {
      sendingToWave[job.id] = true;
      
      const result = await WaveService.sendToWave(job.id);
      
      // Update the job in the store with Wave invoice info
      await jobsStore.updateWaveInfo(job.id, result.invoiceNumber, result.invoiceUrl);
      
      // Reload jobs to get updated data
      await jobsStore.load();
      
    } catch (error) {
      console.error('Failed to send to Wave:', error);
      waveError[job.id] = error instanceof Error ? error.message : 'Failed to create Wave invoice';
    } finally {
      sendingToWave[job.id] = false;
    }
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
  <!-- Admin view mode indicator -->
  {#if $userStore?.role === 'admin' && $userStore.isViewingAsTech}
    <div class="view-mode-banner">
      <Users size={16} />
      <span>Viewing as Tech</span>
    </div>
  {/if}

  <!-- Header -->
  <div class="dashboard-header">
    <h1>Jobs</h1>
    {#if permissions.canCreateJobs($effectiveRole)}
      <button class="new-job-btn" on:click={() => showCreateModal = true}>
        <Plus size={18} />
        <span>New Job</span>
      </button>
    {/if}
  </div>

  <!-- Tabs -->
  <div class="category-tabs">
    <button 
      class="category-tab"
      class:active={categoryFilter === 'all'}
      on:click={() => categoryFilter = 'all'}
    >
      All
    </button>
    <button 
      class="category-tab"
      class:active={categoryFilter === 'skyview'}
      on:click={() => categoryFilter = 'skyview'}
    >
      Skyview
    </button>
    <button 
      class="category-tab"
      class:active={categoryFilter === 'contractors'}
      on:click={() => categoryFilter = 'contractors'}
    >
      Contractors
    </button>
    <button 
      class="category-tab"
      class:active={categoryFilter === 'rayno'}
      on:click={() => categoryFilter = 'rayno'}
    >
      Rayno
    </button>
  </div>

  <!-- Search -->
  <div class="search-section">
    <div class="search-wrapper">
      <Search size={18} class="search-icon" />
      <input
        type="text"
        placeholder="Search by customer name..."
        bind:value={searchQuery}
        class="search-input"
      />
    </div>
    <div class="view-toggle">
      <button
        class="view-btn"
        class:active={viewMode === 'cards'}
        on:click={() => viewMode = 'cards'}
        title="Card view"
      >
        <LayoutGrid size={16} />
      </button>
      <button
        class="view-btn"
        class:active={viewMode === 'list'}
        on:click={() => viewMode = 'list'}
        title="List view"
      >
        <List size={16} />
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
            <!-- Phase Dropdown instead of Status Badge -->
            {#if job.template?.phases && job.template.phases.length > 0}
              <select 
                class="phase-dropdown-top"
                value={job.currentPhaseId || ''}
                on:change={(e) => updateJobPhase(e, job, e.target.value)}
                on:click|stopPropagation
              >
                <option value="">Select Phase</option>
                {#each job.template.phases.sort((a, b) => a.order - b.order) as phase}
                  <option value={phase.id}>{phase.name}</option>
                {/each}
              </select>
            {:else}
              <div class="no-phase-badge">No Phases</div>
            {/if}

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
            {#if permissions.canSeePrices($effectiveRole)}
              <div class="job-amount-section">
                <span class="amount">{formatCurrency(job.totalAmount || 0)}</span>
              </div>
            {/if}

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
            {#if permissions.canSeePrices($effectiveRole) && job.waveInvoiceId}
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
              {#if permissions.canSeePrices($effectiveRole)}
                {#if !job.waveInvoiceId && !sendingToWave[job.id]}
                  <button 
                    class="action-btn" 
                    on:click={(e) => sendToWave(e, job)}
                    title="Send to Wave"
                  >
                    <Send size={18} />
                    <span>Send to Wave</span>
                  </button>
                {:else if sendingToWave[job.id]}
                  <button class="action-btn" disabled>
                    <span>Sending...</span>
                  </button>
                {/if}
              {/if}
              <button class="action-btn" on:click={(e) => openPhotoGallery(e, job)}>
                <Camera size={18} />
                <span>Gallery {#if job.photos?.length > 0}({job.photos.length}){/if}</span>
              </button>
              {#if permissions.canSeePrices($effectiveRole)}
                <button class="action-btn" on:click={(e) => openInvoice(e, job)}>
                  <FileText size={18} />
                  <span>Invoice</span>
                </button>
              {/if}
            </div>
            
            {#if waveError[job.id]}
              <div class="wave-error">
                {waveError[job.id]}
              </div>
            {/if}

            <!-- Footer Actions -->
            <div class="job-footer">
              <Button fullWidth on:click={() => openJobDetail(job)}>
                View Job
              </Button>
              {#if permissions.canDeleteJobs($effectiveRole)}
                <button class="delete-btn" on:click={(e) => confirmDeleteJob(e, job)}>
                  <Trash2 size={18} />
                  Delete Job
                </button>
              {/if}
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
                  {#if job.template?.phases && job.template.phases.length > 0}
                    <select 
                      class="list-phase-dropdown"
                      value={job.currentPhaseId || ''}
                      on:change={(e) => updateJobPhase(e, job, e.target.value)}
                      on:click|stopPropagation
                    >
                      <option value="">Select Phase</option>
                      {#each job.template.phases.sort((a, b) => a.order - b.order) as phase}
                        <option value={phase.id}>{phase.name}</option>
                      {/each}
                    </select>
                  {:else}
                    <span class="list-no-phase">No Phases</span>
                  {/if}
                </div>
                <div class="list-item-details">
                  <div class="list-meta">
                    <MapPin size={14} />
                    <span>{job.address}</span>
                  </div>
                  <div class="list-meta">
                    <Calendar size={14} />
                    <span>
                      {#if job.scheduledDate}
                        {formatDate(job.scheduledDate)}
                      {:else}
                        Not scheduled
                      {/if}
                    </span>
                  </div>
                </div>
              </div>
              <div class="list-item-right">
                {#if permissions.canSeePrices($effectiveRole)}
                  <span class="list-amount">{formatCurrency(job.totalAmount || 0)}</span>
                {/if}
                <ChevronRight size={18} />
              </div>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
  
  <!-- Hidden admin toggle (only for admins) -->
  {#if $userStore?.role === 'admin'}
    <button 
      class="admin-toggle"
      on:click={() => userStore.toggleViewMode()}
      title={$userStore.isViewingAsTech ? 'Switch to Admin View' : 'Switch to Tech View'}
    >
      {#if $userStore.isViewingAsTech}
        <Shield size={16} />
      {:else}
        <Users size={16} />
      {/if}
    </button>
  {/if}
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
  <PhotoGalleryModal 
    bind:isOpen={showPhotoGallery} 
    job={selectedJob} 
    on:photosDeleted={() => jobsStore.load()}
  />
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
    padding: 2rem 1.5rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  /* View mode banner */
  .view-mode-banner {
    background: #fef3c7;
    color: #92400e;
    padding: 0.5rem 1rem;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
    letter-spacing: -0.01em;
  }

  .dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 3rem;
  }

  .dashboard-header h1 {
    font-size: 2.25rem;
    font-weight: 800;
    color: #0a0a0a;
    margin: 0;
    letter-spacing: -0.03em;
  }

  /* Category Tabs */
  .category-tabs {
    display: flex;
    gap: 2rem;
    margin-bottom: 2.5rem;
    border-bottom: 1px solid #e5e5e5;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
  }

  .category-tabs::-webkit-scrollbar {
    display: none;
  }

  .category-tab {
    padding: 0 0 1rem 0;
    border: none;
    background: none;
    font-size: 0.9375rem;
    font-weight: 500;
    color: #737373;
    cursor: pointer;
    transition: color 0.2s ease;
    border-bottom: 2px solid transparent;
    white-space: nowrap;
    letter-spacing: -0.01em;
    position: relative;
  }

  .category-tab:hover {
    color: #404040;
  }

  .category-tab.active {
    color: #0a0a0a;
    font-weight: 600;
  }

  .category-tab.active::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    right: 0;
    height: 2px;
    background: #0a0a0a;
    animation: tabSlide 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes tabSlide {
    from {
      transform: scaleX(0);
    }
    to {
      transform: scaleX(1);
    }
  }

  /* Search Section */
  .search-section {
    margin-bottom: 2rem;
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .search-wrapper {
    position: relative;
    flex: 1;
  }

  :global(.search-icon) {
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: #a3a3a3;
    pointer-events: none;
  }

  .search-input {
    width: 100%;
    padding: 0.875rem 1rem 0.875rem 2.875rem;
    border: 1px solid #e5e5e5;
    border-radius: 10px;
    font-size: 0.9375rem;
    font-weight: 400;
    background: #fafafa;
    transition: all 0.2s ease;
    letter-spacing: -0.01em;
  }

  .search-input::placeholder {
    color: #a3a3a3;
    font-weight: 400;
  }

  .search-input:focus {
    border-color: #d4d4d4;
    background: white;
    outline: none;
    box-shadow: 0 0 0 4px rgba(0, 0, 0, 0.02);
  }

  /* View Toggle */
  .view-toggle {
    display: flex;
    background: #f5f5f5;
    padding: 0.125rem;
    border-radius: 8px;
  }

  .view-btn {
    padding: 0.5rem 0.625rem;
    border: none;
    background: transparent;
    border-radius: 6px;
    cursor: pointer;
    color: #737373;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .view-btn:hover {
    color: #404040;
  }

  .view-btn.active {
    background: white;
    color: #0a0a0a;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }

  /* Jobs Container */
  .jobs-container {
    min-height: 400px;
  }

  .jobs-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
    gap: 1.25rem;
  }

  /* New Job Button */
  .new-job-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    background: #0a0a0a;
    color: white;
    border: none;
    border-radius: 10px;
    font-size: 0.9375rem;
    font-weight: 500;
    letter-spacing: -0.01em;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .new-job-btn:hover {
    background: #171717;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .new-job-btn span {
    line-height: 1;
  }

  /* List View */
  .jobs-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .job-list-item {
    background: white;
    border: 1px solid #e5e5e5;
    border-radius: 12px;
    padding: 1.25rem 1.5rem;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .job-list-item:hover {
    border-color: #d4d4d4;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    transform: translateY(-1px);
  }

  .list-item-left {
    flex: 1;
  }

  .list-item-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.625rem;
  }

  .list-item-header h4 {
    margin: 0;
    font-size: 1.0625rem;
    font-weight: 600;
    color: #0a0a0a;
    letter-spacing: -0.02em;
    line-height: 1.3;
  }

  .list-phase-dropdown {
    padding: 0.25rem 0.5rem;
    padding-right: 1.75rem;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 500;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    color: #475569;
    cursor: pointer;
    transition: all 0.2s ease;
    appearance: none;
    background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' width='10' height='10' viewBox='0 0 12 12'%3e%3cpath fill='%23475569' d='M6 9L1 4h10z'/%3e%3c/svg%3e");
    background-repeat: no-repeat;
    background-position: right 0.5rem center;
    min-width: 100px;
  }

  .list-phase-dropdown:hover {
    background-color: #e2e8f0;
    border-color: #cbd5e1;
  }

  .list-phase-dropdown:focus {
    outline: none;
    border-color: #5b5bd6;
    box-shadow: 0 0 0 2px rgba(91, 91, 214, 0.1);
  }

  .list-no-phase {
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 500;
    background: #f1f5f9;
    color: #64748b;
  }

  .list-item-details {
    display: flex;
    gap: 1.75rem;
    align-items: center;
  }

  .list-meta {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #737373;
    font-size: 0.875rem;
    font-weight: 400;
    letter-spacing: -0.01em;
  }

  .list-meta :global(svg) {
    color: #a3a3a3;
    flex-shrink: 0;
  }

  .list-item-right {
    display: flex;
    align-items: center;
    gap: 1.25rem;
  }

  .list-amount {
    font-size: 1.25rem;
    font-weight: 700;
    color: #0a0a0a;
    letter-spacing: -0.02em;
    font-variant-numeric: tabular-nums;
  }

  .list-item-right :global(svg) {
    color: #d4d4d4;
    transition: color 0.2s ease;
  }

  .job-list-item:hover .list-item-right :global(svg) {
    color: #737373;
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
    border: 1px solid #e5e5e5;
    border-radius: 14px;
    padding: 1.5rem;
    transition: all 0.2s ease;
    cursor: pointer;
    position: relative;
  }

  .job-card:hover {
    border-color: #d4d4d4;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
    transform: translateY(-2px);
  }

  /* Phase Dropdown at Top */
  .phase-dropdown-top {
    position: absolute;
    top: 1rem;
    right: 1rem;
    padding: 0.5rem 0.75rem;
    padding-right: 2rem;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    font-weight: 500;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    color: #475569;
    cursor: pointer;
    transition: all 0.2s ease;
    appearance: none;
    background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3e%3cpath fill='%23475569' d='M6 9L1 4h10z'/%3e%3c/svg%3e");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    min-width: 120px;
  }

  .phase-dropdown-top:hover {
    background-color: #e2e8f0;
    border-color: #cbd5e1;
  }

  .phase-dropdown-top:focus {
    outline: none;
    border-color: #5b5bd6;
    box-shadow: 0 0 0 3px rgba(91, 91, 214, 0.1);
  }

  .no-phase-badge {
    position: absolute;
    top: 1rem;
    right: 1rem;
    padding: 0.5rem 0.75rem;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    font-weight: 500;
    background: #f1f5f9;
    color: #64748b;
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
  
  .phase-dropdown {
    padding: 0.375rem 0.75rem;
    font-size: 0.875rem;
    border: 1px solid var(--border-color);
    border-radius: 0.375rem;
    background: var(--bg-tertiary);
    color: var(--text-primary);
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .phase-dropdown:hover {
    background: var(--bg-hover);
    border-color: var(--border-hover);
  }
  
  .phase-dropdown:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(14, 165, 233, 0.1);
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
  
  .wave-error {
    margin-top: 0.5rem;
    padding: 0.5rem;
    background-color: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 6px;
    color: #dc2626;
    font-size: 0.875rem;
    text-align: center;
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

  /* Mobile Responsive */
  @media (max-width: 768px) {
    .dashboard {
      padding: 1.25rem 1rem;
    }

    .dashboard-header {
      margin-bottom: 2rem;
    }

    .dashboard-header h1 {
      font-size: 1.875rem;
    }

    .new-job-btn {
      padding: 0.625rem 1rem;
      font-size: 0.875rem;
    }

    .new-job-btn :global(svg) {
      width: 16px;
      height: 16px;
    }

    .category-tabs {
      gap: 1.5rem;
      margin-bottom: 1.75rem;
    }

    .category-tab {
      font-size: 0.875rem;
    }

    .search-section {
      flex-direction: column;
      gap: 0.75rem;
      align-items: stretch;
      margin-bottom: 1.5rem;
    }

    .search-input {
      font-size: 16px; /* Prevents zoom on iOS */
    }

    .view-toggle {
      align-self: flex-start;
    }

    .jobs-grid {
      grid-template-columns: 1fr;
      gap: 1rem;
    }

    .job-card {
      padding: 1.25rem;
    }

    .job-list-item {
      padding: 1rem;
      flex-direction: column;
      align-items: stretch;
      gap: 0.75rem;
    }

    .list-item-header h4 {
      font-size: 1rem;
    }

    .list-item-right {
      flex-direction: row-reverse;
      justify-content: space-between;
      width: 100%;
      padding-top: 0.875rem;
      border-top: 1px solid #f5f5f5;
    }

    .list-amount {
      font-size: 1.125rem;
    }

    .list-item-details {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.625rem;
    }

    .list-meta {
      font-size: 0.8125rem;
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

  /* Admin Toggle - Hidden but accessible */
  .admin-toggle {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    background: #f5f5f5;
    border: 1px solid #e5e5e5;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    opacity: 0.3;
    transition: all 0.2s ease;
    color: #737373;
  }

  .admin-toggle:hover {
    opacity: 1;
    background: #0a0a0a;
    color: white;
    border-color: #0a0a0a;
  }
</style>