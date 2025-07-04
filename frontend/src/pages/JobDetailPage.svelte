<script lang="ts">
  import { onMount } from 'svelte';
  import { jobsStore, jobs, customersStore, templates, userStore, effectiveRole, permissions } from '../lib/stores';
  import { router } from '../lib/router';
  import { Card, Button, QuantityPicker, PhotoCaptureModal, PhotoGallery } from '../lib/components';
  import { PhotoService } from '../lib/services/photoService';
  import { Camera } from 'lucide-svelte';
  import type { Job, Customer, JobTemplate } from '../lib/types/models';
  
  export let jobId: string = '';
  
  let job: Job | undefined;
  let customer: Customer | undefined;
  let template: JobTemplate | undefined;
  let loading = true;
  let error: string | null = null;
  
  onMount(async () => {
    try {
      // Load job details
      await jobsStore.loadById(jobId);
      job = $jobs.find(j => j.id === jobId);
      
      if (job) {
        // Load customer details
        await customersStore.loadById(job.customerId);
        customer = $customersStore.customers.find(c => c.id === job!.customerId);
        
        // Load template details
        template = $templates.find(t => t.id === job!.templateId);
      }
      
      loading = false;
    } catch (e) {
      error = e instanceof Error ? e.message : 'Failed to load job details';
      loading = false;
    }
  });
  
  let showPhotoCapture = false;
  let isUploadingPhotos = false;
  
  async function handleQuantityChange(jobItemId: string, newQuantity: number) {
    try {
      await jobsStore.updateJobItem(jobId, jobItemId, newQuantity);
      // Reload job to get updated data
      await jobsStore.loadById(jobId);
      job = $jobs.find(j => j.id === jobId);
    } catch (error) {
      console.error('Failed to update item quantity:', error);
    }
  }
  
  function openGPS(address: string) {
    const encodedAddress = encodeURIComponent(address);
    window.open(`https://maps.google.com/?q=${encodedAddress}`, '_blank');
  }
  
  function calculateProgress() {
    if (!job) return 0;
    // For now, base progress on job status
    switch (job.status) {
      case 'scheduled': return 0;
      case 'in_progress': return 50;
      case 'completed': return 100;
      case 'cancelled': return 0;
      default: return 0;
    }
  }
  
  async function handlePhotosCapture(event: CustomEvent<{ files: File[] }>) {
    const files = event.detail.files;
    if (files.length === 0 || !job) return;
    
    isUploadingPhotos = true;
    try {
      // Upload photos to backend (which handles R2 upload)
      const uploadedPhotos = await PhotoService.uploadPhotos(jobId, files);
      
      // Reload job to get updated photos from backend
      await jobsStore.loadById(jobId);
      job = $jobs.find(j => j.id === jobId);
      
      // Update gallery photos
      galleryPhotos = job?.photos?.map(photo => ({
        id: photo.id,
        url: photo.url,
        thumbnail: photo.url,
        name: photo.caption || `Photo ${photo.id.slice(0, 8)}`
      })) || [];
    } catch (error) {
      console.error('Error uploading photos:', error);
      alert('Failed to upload photos. Please try again.');
    } finally {
      isUploadingPhotos = false;
    }
  }
  
  async function handlePhotosDelete(event: CustomEvent<{ ids: string[] }>) {
    if (!job) return;
    
    try {
      // Delete photos via API (which also removes from R2)
      await PhotoService.deletePhotos(jobId, event.detail.ids);
      
      // Reload job to get updated photos
      await jobsStore.loadById(jobId);
      job = $jobs.find(j => j.id === jobId);
      
      // Update gallery photos
      galleryPhotos = job?.photos?.map(photo => ({
        id: photo.id,
        url: photo.url,
        thumbnail: photo.url,
        name: photo.caption || `Photo ${photo.id.slice(0, 8)}`
      })) || [];
    } catch (error) {
      console.error('Failed to delete photos:', error);
      alert('Failed to delete photos. Please try again.');
    }
  }
  
  async function handlePhotosDownload(event: CustomEvent<{ ids: string[] }>) {
    if (!job) return;
    
    const photosToDownload = job.photos
      .filter(photo => event.detail.ids.includes(`photo_${photo.id}`))
      .map(photo => ({
        id: photo.id,
        url: photo.url,
        name: `job_${jobId}_photo_${photo.id}.jpg`
      }));
    
    await PhotoService.downloadPhotos(photosToDownload);
  }
  
  $: progress = calculateProgress();
  $: galleryPhotos = job?.photos?.map(photo => ({
    id: `photo_${photo.id}`,
    url: photo.url,
    thumbnail: photo.url,
    name: photo.caption || `Photo ${photo.id.slice(0, 8)}`
  })) || [];
</script>

{#if loading}
  <div class="loading-state">
    <p>Loading job details...</p>
  </div>
{:else if error}
  <div class="error-state">
    <p>{error}</p>
    <Button on:click={() => router.navigate('jobs')}>Back to Jobs</Button>
  </div>
{:else if !job}
  <div class="error-state">
    <p>Job not found</p>
    <Button on:click={() => router.navigate('jobs')}>Back to Jobs</Button>
  </div>
{:else}
  <div class="job-detail">
    <div class="detail-header">
      <Button variant="secondary" on:click={() => router.navigate('jobs')}>
        ← Back
      </Button>
      <h1>{customer?.name || 'Unknown Customer'}</h1>
    </div>
    
    <button 
      class="address-section"
      on:click={() => openGPS(job.address)}
    >
      <span class="icon">📍</span>
      <div>
        <div class="label">Job Address</div>
        <div class="address">{job.address}</div>
      </div>
      <span class="arrow">→</span>
    </button>
    
    <div class="progress-section">
      <div class="progress-header">
        <span>Overall Progress</span>
        <span>{progress}%</span>
      </div>
      <div class="progress-bar">
        <div class="progress-fill" style="width: {progress}%"></div>
      </div>
    </div>
    
    <div class="sections">
      <Card>
        <h2>What we did</h2>
        <div class="items-list">
          {#each job.items as jobItem}
            <div class="item-row">
              <div class="item-info">
                <div class="item-name">{jobItem.nickname || jobItem.name}</div>
                <div class="item-meta">
                  {#if permissions.canSeePrices($effectiveRole)}
                    Price: ${jobItem.price.toFixed(2)} × {jobItem.quantity} = ${jobItem.total.toFixed(2)}
                  {:else}
                    Quantity: {jobItem.quantity}
                  {/if}
                </div>
              </div>
              <QuantityPicker
                value={jobItem.quantity}
                min={0}
                max={jobItem.quantity * 2}
                on:change={(e) => handleQuantityChange(jobItem.id, e.detail)}
              />
            </div>
          {/each}
        </div>
      </Card>
      
      <Card>
        <div class="photos-header">
          <h2>Job Photos</h2>
          <Button 
            variant="primary" 
            on:click={() => showPhotoCapture = true}
            disabled={isUploadingPhotos}
          >
            <Camera size={16} />
            Take Photos
          </Button>
        </div>
        
        <PhotoGallery 
          photos={galleryPhotos}
          on:delete={handlePhotosDelete}
          on:download={handlePhotosDownload}
        />
      </Card>
    </div>
  </div>
  
  <PhotoCaptureModal 
    bind:isOpen={showPhotoCapture}
    on:capture={handlePhotosCapture}
  />
{/if}

<style>
  .job-detail {
    padding: 1rem;
    max-width: 800px;
    margin: 0 auto;
  }
  
  .error-state,
  .loading-state {
    text-align: center;
    padding: 4rem 2rem;
  }
  
  .detail-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }
  
  .detail-header h1 {
    margin: 0;
    font-size: 1.75rem;
    color: #111827;
  }
  
  .address-section {
    display: flex;
    align-items: center;
    gap: 1rem;
    width: 100%;
    padding: 1rem;
    margin-bottom: 1.5rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
    font-family: inherit;
  }
  
  .address-section:hover {
    background-color: #f9fafb;
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }
  
  .address-section .icon {
    font-size: 1.5rem;
  }
  
  .address-section .label {
    font-size: 0.875rem;
    color: #6b7280;
    margin-bottom: 0.25rem;
  }
  
  .address-section .address {
    font-size: 1rem;
    color: #111827;
    font-weight: 500;
  }
  
  .address-section .arrow {
    margin-left: auto;
    color: #6b7280;
  }
  
  .progress-section {
    margin-bottom: 1.5rem;
  }
  
  .progress-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
    font-size: 0.875rem;
    color: #6b7280;
  }
  
  .progress-bar {
    height: 8px;
    background-color: #e5e7eb;
    border-radius: 4px;
    overflow: hidden;
  }
  
  .progress-fill {
    height: 100%;
    background-color: #10b981;
    transition: width 0.3s ease;
  }
  
  .sections {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  h2 {
    margin: 0 0 1rem 0;
    font-size: 1.25rem;
    color: #111827;
  }
  
  .items-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .item-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background-color: #f9fafb;
    border-radius: 6px;
  }
  
  .item-info {
    flex: 1;
  }
  
  .item-name {
    font-weight: 500;
    color: #111827;
    margin-bottom: 0.25rem;
  }
  
  .item-meta {
    font-size: 0.875rem;
    color: #6b7280;
  }
  
  .installed {
    color: #10b981;
    font-weight: 500;
  }
  
  .photos-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  @media (max-width: 640px) {
    .job-detail {
      padding: 1rem;
    }
    
    .detail-header {
      flex-direction: column;
      align-items: flex-start;
    }
    
    .item-row {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
    }
    
    .photos-header {
      flex-direction: column;
      align-items: stretch;
      gap: 1rem;
    }
  }
</style>