<script lang="ts">
  import { jobsStore, jobTemplatesStore } from '../lib/stores';
  import { router } from '../lib/router';
  import { Card, Button, Modal, PhotoGalleryModal, InvoiceModal } from '../lib/components';
  import { Plus, MapPin, Calendar, Package, AlertCircle, CheckCircle, Clock, XCircle, 
           Shield, ShieldOff, Camera, FileText, Send, ExternalLink, DollarSign } from 'lucide-svelte';
  import type { Customer, Job } from '../lib/types/models';
  
  let showCreateModal = false;
  let showPhotoGallery = false;
  let showInvoiceModal = false;
  let selectedJob: Job | null = null;
  let editingDateJobId: string | null = null;
  let tempStartDate: string = '';
  let tempEndDate: string = '';
  
  let newJob = {
    customerName: '',
    address: '',
    requiresPermit: false,
    startDate: '',
    endDate: '',
    templateId: ''
  };
  
  function openGPS(address: string) {
    // Encode the address for URL
    const encodedAddress = encodeURIComponent(address);
    // This will open the device's default map app
    window.open(`https://maps.google.com/?q=${encodedAddress}`, '_blank');
  }
  
  function togglePermit(event: Event, jobId: string) {
    event.stopPropagation();
    jobsStore.togglePermit(jobId);
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
  
  function sendToWave(event: Event, job: Job) {
    event.stopPropagation();
    // Mock Wave integration - in real app, this would call Wave API
    const mockInvoiceId = `WAV-${Date.now()}`;
    const mockInvoiceUrl = `https://wave.app/invoice/${mockInvoiceId}`;
    jobsStore.setWaveInvoice(job.id, mockInvoiceId, mockInvoiceUrl);
  }
  
  function openWaveInvoice(event: Event, url: string) {
    event.stopPropagation();
    window.open(url, '_blank');
  }
  
  function formatCurrency(amount: number): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  }
  
  let formError = '';
  
  function handleCreateJob() {
    if (!newJob.customerName || !newJob.address || !newJob.templateId) {
      formError = 'Please fill in all required fields';
      return;
    }
    
    formError = '';
    
    const customer: Customer = {
      id: crypto.randomUUID(),
      name: newJob.customerName
    };
    
    const jobId = jobsStore.createFromTemplate(
      customer,
      newJob.address,
      newJob.templateId,
      newJob.requiresPermit,
      newJob.startDate ? new Date(newJob.startDate) : undefined,
      newJob.endDate ? new Date(newJob.endDate) : undefined
    );
    
    // Reset form
    newJob = {
      customerName: '',
      address: '',
      requiresPermit: false,
      startDate: '',
      endDate: '',
      templateId: ''
    };
    
    showCreateModal = false;
    formError = '';
    
    // Navigate to job detail
    router.navigate('job-detail', { id: jobId });
  }
  
  function getStatusIcon(status: string) {
    switch (status) {
      case 'pending': return Clock;
      case 'in-progress': return AlertCircle;
      case 'completed': return CheckCircle;
      case 'cancelled': return XCircle;
      default: return Clock;
    }
  }
  
  function getStatusColor(status: string) {
    switch (status) {
      case 'pending': return 'warning';
      case 'in-progress': return 'primary';
      case 'completed': return 'success';
      case 'cancelled': return 'danger';
      default: return 'gray';
    }
  }
  
  function formatDate(date?: Date) {
    if (!date) return 'Not set';
    return new Date(date).toLocaleDateString();
  }
  
  function formatDateForInput(date?: Date): string {
    if (!date) return '';
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  }
  
  function startEditingDate(event: Event, job: Job) {
    event.stopPropagation();
    editingDateJobId = job.id;
    tempStartDate = formatDateForInput(job.startDate);
    tempEndDate = formatDateForInput(job.endDate);
  }
  
  function saveDates(event: Event, jobId: string) {
    event.stopPropagation();
    const startDate = tempStartDate ? new Date(tempStartDate) : undefined;
    const endDate = tempEndDate ? new Date(tempEndDate) : undefined;
    jobsStore.updateDates(jobId, startDate, endDate);
    editingDateJobId = null;
  }
  
  function cancelEditDate(event: Event) {
    event.stopPropagation();
    editingDateJobId = null;
    tempStartDate = '';
    tempEndDate = '';
  }
  
  // Calculate stats
  $: stats = {
    total: $jobsStore.length,
    active: $jobsStore.filter(j => j.status === 'in-progress').length,
    pending: $jobsStore.filter(j => j.status === 'pending').length,
    completed: $jobsStore.filter(j => j.status === 'completed').length
  };
</script>

<div class="jobs-page">
  <div class="page-header">
    <div>
      <h1>Jobs Dashboard</h1>
      <p class="subtitle">Manage your electrical projects</p>
    </div>
    <Button on:click={() => showCreateModal = true}>
      <Plus size={20} />
      Create New Job
    </Button>
  </div>
  
  <!-- Stats Cards -->
  <div class="stats-grid">
    <Card variant="elevated" padding="small">
      <div class="stat-card">
        <div class="stat-icon primary">
          <Package size={24} />
        </div>
        <div class="stat-content">
          <div class="stat-value">{stats.total}</div>
          <div class="stat-label">Total Jobs</div>
        </div>
      </div>
    </Card>
    
    <Card variant="elevated" padding="small">
      <div class="stat-card">
        <div class="stat-icon primary">
          <AlertCircle size={24} />
        </div>
        <div class="stat-content">
          <div class="stat-value">{stats.active}</div>
          <div class="stat-label">Active</div>
        </div>
      </div>
    </Card>
    
    <Card variant="elevated" padding="small">
      <div class="stat-card">
        <div class="stat-icon warning">
          <Clock size={24} />
        </div>
        <div class="stat-content">
          <div class="stat-value">{stats.pending}</div>
          <div class="stat-label">Pending</div>
        </div>
      </div>
    </Card>
    
    <Card variant="elevated" padding="small">
      <div class="stat-card">
        <div class="stat-icon success">
          <CheckCircle size={24} />
        </div>
        <div class="stat-content">
          <div class="stat-value">{stats.completed}</div>
          <div class="stat-label">Completed</div>
        </div>
      </div>
    </Card>
  </div>
  
  <!-- Jobs Grid -->
  {#if $jobsStore.length === 0}
    <Card variant="elevated">
      <div class="empty-state">
        <Package size={48} strokeWidth={1.5} />
        <h3>No jobs yet</h3>
        <p>Create your first job to get started</p>
        <Button on:click={() => showCreateModal = true}>
          <Plus size={20} />
          Create First Job
        </Button>
      </div>
    </Card>
  {:else}
    <div class="jobs-grid">
      {#each $jobsStore as job}
        <Card 
          variant="elevated"
          clickable
          on:click={() => router.navigate('job-detail', { id: job.id })}
        >
          <div class="job-card">
            <div class="job-header">
              <div>
                <h3>{job.customer.name}</h3>
                <div class="badges">
                  <span class="badge template-badge">{job.template.name}</span>
                  <button 
                    class="badge permit-badge"
                    class:required={job.requiresPermit}
                    on:click={(e) => togglePermit(e, job.id)}
                    title="Click to toggle permit status"
                  >
                    {#if job.requiresPermit}
                      <Shield size={14} />
                      Permit
                    {:else}
                      <ShieldOff size={14} />
                      No Permit
                    {/if}
                  </button>
                </div>
              </div>
              <div class="job-total">
                <DollarSign size={20} />
                <span>{formatCurrency(jobsStore.calculateJobTotal(job))}</span>
              </div>
            </div>
            
            <button 
              class="address-link"
              on:click|stopPropagation={() => openGPS(job.address)}
              title="Open in maps"
            >
              <MapPin size={16} />
              {job.address}
            </button>
            
            <div class="job-meta-row">
              {#if editingDateJobId === job.id}
                <div class="date-edit-form">
                  <div class="date-input-group">
                    <label for="start-{job.id}">Start:</label>
                    <input 
                      id="start-{job.id}"
                      type="date" 
                      bind:value={tempStartDate}
                      on:click|stopPropagation
                    />
                  </div>
                  <div class="date-input-group">
                    <label for="end-{job.id}">End:</label>
                    <input 
                      id="end-{job.id}"
                      type="date" 
                      bind:value={tempEndDate}
                      on:click|stopPropagation
                    />
                  </div>
                  <div class="date-edit-actions">
                    <button 
                      class="date-btn save"
                      on:click={(e) => saveDates(e, job.id)}
                      title="Save dates"
                    >
                      <CheckCircle size={16} />
                    </button>
                    <button 
                      class="date-btn cancel"
                      on:click={cancelEditDate}
                      title="Cancel"
                    >
                      <XCircle size={16} />
                    </button>
                  </div>
                </div>
              {:else}
                <button 
                  class="job-meta clickable-date"
                  on:click={(e) => startEditingDate(e, job)}
                  title="Click to edit dates"
                >
                  <Calendar size={16} />
                  <span>{formatDate(job.startDate)}</span>
                  {#if job.endDate}
                    <span class="date-separator">→</span>
                    <span>{formatDate(job.endDate)}</span>
                  {/if}
                </button>
              {/if}
              <div class="job-meta">
                <div class="status-badge {getStatusColor(job.status)}">
                  <svelte:component this={getStatusIcon(job.status)} size={14} />
                  {job.status}
                </div>
              </div>
            </div>
            
            <div class="job-actions">
              <button 
                class="action-btn"
                on:click={(e) => openPhotoGallery(e, job)}
                title="View photos"
              >
                <Camera size={18} />
                <span class="action-count">{job.photos?.length || 0}</span>
              </button>
              
              <button 
                class="action-btn"
                on:click={(e) => openInvoice(e, job)}
                title="View invoice"
              >
                <FileText size={18} />
              </button>
              
              {#if job.waveInvoiceId}
                <button 
                  class="action-btn sent"
                  on:click={(e) => openWaveInvoice(e, job.waveInvoiceUrl || '')}
                  title="View in Wave"
                >
                  <ExternalLink size={18} />
                  <span class="invoice-id">{job.waveInvoiceId}</span>
                </button>
              {:else}
                <button 
                  class="action-btn primary"
                  on:click={(e) => sendToWave(e, job)}
                  title="Send to Wave"
                >
                  <Send size={18} />
                  Wave
                </button>
              {/if}
            </div>
          </div>
        </Card>
      {/each}
    </div>
  {/if}
</div>

<Modal bind:isOpen={showCreateModal} title="Create New Job" size="medium">
  {#if formError}
    <div class="error-message">
      <AlertCircle size={20} />
      {formError}
    </div>
  {/if}
  
  <form on:submit|preventDefault={handleCreateJob}>
    <div class="form-group">
      <label for="customer-name">Customer Name *</label>
      <input
        id="customer-name"
        type="text"
        bind:value={newJob.customerName}
        placeholder="John Doe"
        required
      />
    </div>
    
    <div class="form-group">
      <label for="address">Job Address *</label>
      <input
        id="address"
        type="text"
        bind:value={newJob.address}
        placeholder="123 Main St, City, State"
        required
      />
    </div>
    
    <div class="form-group">
      <label for="template">Job Template *</label>
      <select id="template" bind:value={newJob.templateId} required>
        <option value="">Select a template...</option>
        {#each $jobTemplatesStore as template}
          <option value={template.id}>{template.name}</option>
        {/each}
      </select>
    </div>
    
    <div class="form-row">
      <div class="form-group">
        <label for="start-date">Start Date</label>
        <input
          id="start-date"
          type="date"
          bind:value={newJob.startDate}
        />
      </div>
      
      <div class="form-group">
        <label for="end-date">End Date</label>
        <input
          id="end-date"
          type="date"
          bind:value={newJob.endDate}
        />
      </div>
    </div>
    
    <div class="form-group checkbox">
      <label>
        <input
          type="checkbox"
          bind:checked={newJob.requiresPermit}
        />
        <span>Requires Permit</span>
      </label>
    </div>
  </form>
  
  <div slot="footer">
    <Button variant="ghost" on:click={() => showCreateModal = false}>
      Cancel
    </Button>
    <Button on:click={handleCreateJob}>
      Create Job
    </Button>
  </div>
</Modal>

<PhotoGalleryModal 
  bind:isOpen={showPhotoGallery} 
  photos={selectedJob?.photos || []}
  title={selectedJob ? `${selectedJob.customer.name} - Photos` : 'Photos'}
/>

<InvoiceModal 
  bind:isOpen={showInvoiceModal}
  job={selectedJob}
/>

<style>
  .jobs-page {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
  }
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }
  
  h1 {
    margin: 0;
  }
  
  .subtitle {
    margin: 0.25rem 0 0 0;
    color: var(--gray-500);
    font-size: 1rem;
  }
  
  /* Stats Grid */
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
  }
  
  .stat-card {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .stat-icon.primary {
    background: var(--primary-100);
    color: var(--primary-600);
  }
  
  .stat-icon.warning {
    background: var(--warning-100);
    color: var(--warning-600);
  }
  
  .stat-icon.success {
    background: var(--success-100);
    color: var(--success-600);
  }
  
  .stat-content {
    flex: 1;
  }
  
  .stat-value {
    font-size: 2rem;
    font-weight: 700;
    color: var(--gray-900);
    line-height: 1;
  }
  
  .stat-label {
    font-size: 0.9375rem;
    color: var(--gray-500);
    margin-top: 0.25rem;
  }
  
  /* Empty State */
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: var(--gray-500);
  }
  
  .empty-state h3 {
    margin: 1rem 0 0.5rem 0;
    color: var(--gray-900);
  }
  
  .empty-state p {
    margin: 0 0 1.5rem 0;
  }
  
  /* Jobs Grid */
  .jobs-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
    gap: 1.5rem;
  }
  
  .job-card {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .job-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .job-header h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.125rem;
    color: var(--gray-900);
    font-weight: 600;
  }
  
  .badges {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
  
  .badge {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 500;
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
  }
  
  .template-badge {
    background: var(--gray-100);
    color: var(--gray-700);
  }
  
  .permit-badge {
    background: var(--gray-100);
    color: var(--gray-600);
    border: none;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .permit-badge:hover {
    background: var(--gray-200);
  }
  
  .permit-badge.required {
    background: var(--warning-100);
    color: var(--warning-700);
  }
  
  .permit-badge.required:hover {
    background: var(--warning-200);
  }
  
  .job-total {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--gray-900);
  }
  
  .status-badge {
    padding: 0.25rem 0.625rem;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 500;
    text-transform: capitalize;
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
  }
  
  .status-badge.primary {
    background: var(--primary-100);
    color: var(--primary-700);
  }
  
  .status-badge.warning {
    background: var(--warning-100);
    color: var(--warning-700);
  }
  
  .status-badge.success {
    background: var(--success-100);
    color: var(--success-600);
  }
  
  .status-badge.danger {
    background: var(--danger-100);
    color: var(--danger-700);
  }
  
  .status-badge.gray {
    background: var(--gray-100);
    color: var(--gray-700);
  }
  
  .address-link {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--gray-100);
    border: none;
    padding: 0.75rem 1rem;
    margin: 0 -0.25rem;
    color: var(--gray-700);
    cursor: pointer;
    font-size: 0.9375rem;
    text-align: left;
    transition: all 0.2s;
    border-radius: 12px;
    font-family: inherit;
  }
  
  .address-link:hover {
    background: var(--primary-100);
    color: var(--primary-700);
  }
  
  .job-meta-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 0;
    border-top: 1px solid var(--gray-200);
  }
  
  .job-meta {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.875rem;
    color: var(--gray-500);
  }
  
  .job-meta.clickable-date {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    transition: all 0.2s;
    font-family: inherit;
    border-radius: 8px;
    padding: 0.25rem 0.5rem;
    margin: -0.25rem -0.5rem;
  }
  
  .job-meta.clickable-date:hover {
    background: var(--gray-100);
    color: var(--gray-700);
  }
  
  .date-separator {
    color: var(--gray-400);
    margin: 0 0.25rem;
  }
  
  .date-edit-form {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    flex: 1;
  }
  
  .date-input-group {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }
  
  .date-input-group label {
    font-size: 0.75rem;
    color: var(--gray-600);
    font-weight: 500;
  }
  
  .date-input-group input[type="date"] {
    padding: 0.375rem 0.625rem;
    border: 1px solid var(--gray-300);
    border-radius: 8px;
    font-size: 0.8125rem;
    background: white;
    color: var(--gray-900);
    min-width: 120px;
  }
  
  .date-input-group input[type="date"]:focus {
    outline: none;
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
  
  .date-edit-actions {
    display: flex;
    gap: 0.375rem;
  }
  
  .date-btn {
    width: 28px;
    height: 28px;
    padding: 0;
    background: none;
    border: none;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .date-btn.save {
    color: var(--success-600);
  }
  
  .date-btn.save:hover {
    background: var(--success-100);
  }
  
  .date-btn.cancel {
    color: var(--danger-600);
  }
  
  .date-btn.cancel:hover {
    background: var(--danger-100);
  }
  
  .job-actions {
    display: flex;
    gap: 0.5rem;
    padding-top: 0.75rem;
    border-top: 1px solid var(--gray-200);
  }
  
  .action-btn {
    flex: 1;
    padding: 0.625rem 0.875rem;
    background: var(--gray-100);
    border: none;
    border-radius: 12px;
    color: var(--gray-700);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    font-family: inherit;
  }
  
  .action-btn:hover {
    background: var(--gray-200);
  }
  
  .action-btn.primary {
    background: var(--primary-100);
    color: var(--primary-700);
  }
  
  .action-btn.primary:hover {
    background: var(--primary-200);
  }
  
  .action-btn.sent {
    background: var(--success-100);
    color: var(--success-700);
  }
  
  .action-btn.sent:hover {
    background: var(--success-200);
  }
  
  .action-count {
    background: var(--gray-700);
    color: white;
    padding: 0.125rem 0.375rem;
    border-radius: 10px;
    font-size: 0.75rem;
    min-width: 1.25rem;
    text-align: center;
  }
  
  .invoice-id {
    font-size: 0.75rem;
    font-family: monospace;
  }
  
  /* Form styles */
  form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  label {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--gray-700);
  }
  
  input[type="text"],
  input[type="date"],
  select {
    padding: 0.875rem 1rem;
    border: 1px solid var(--gray-300);
    border-radius: 12px;
    font-size: 0.9375rem;
    transition: all 0.2s;
    background: white;
    color: var(--gray-900);
  }
  
  input[type="text"]:focus,
  input[type="date"]:focus,
  select:focus {
    outline: none;
    border-color: var(--primary-500);
    box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
  }
  
  .checkbox label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: normal;
    cursor: pointer;
  }
  
  input[type="checkbox"] {
    width: 1.25rem;
    height: 1.25rem;
    cursor: pointer;
    accent-color: var(--primary-500);
  }
  
  .error-message {
    background: var(--danger-50);
    border: 1px solid var(--danger-200);
    color: var(--danger-700);
    padding: 0.75rem 1rem;
    border-radius: 8px;
    font-size: 0.875rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  @media (max-width: 768px) {
    .jobs-page {
      padding: 1rem;
    }
    
    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
    }
    
    .stats-grid {
      grid-template-columns: repeat(2, 1fr);
    }
    
    .jobs-grid {
      grid-template-columns: 1fr;
    }
  }
</style>