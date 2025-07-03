<script lang="ts">
  import { jobsStore } from '../lib/stores';
  import { router } from '../lib/router';
  import { Card, Button, QuantityPicker, PhotoCaptureModal, PhotoGallery } from '../lib/components';
  import { PhotoService } from '../lib/services/photoService';
  import { Camera } from 'lucide-svelte';
  
  export let jobId: string = '';
  
  $: job = $jobsStore.find(j => j.id === jobId);
  
  let showPhotoCapture = false;
  let isUploadingPhotos = false;
  
  function handleQuantityChange(jobItemId: string, newQuantity: number) {
    jobsStore.updateItemQuantity(jobId, jobItemId, newQuantity);
  }
  
  function handlePhaseToggle(phaseId: string, isCompleted: boolean) {
    jobsStore.updatePhase(jobId, phaseId, isCompleted);
  }
  
  function openGPS(address: string) {
    const encodedAddress = encodeURIComponent(address);
    window.open(`https://maps.google.com/?q=${encodedAddress}`, '_blank');
  }
  
  function calculateProgress() {
    if (!job) return 0;
    const totalItems = job.items.reduce((sum, item) => sum + item.quantity, 0);
    const installedItems = job.items.reduce((sum, item) => sum + item.installedQuantity, 0);
    return totalItems > 0 ? Math.round((installedItems / totalItems) * 100) : 0;
  }
  
  async function handlePhotosCapture(event: CustomEvent<{ files: File[] }>) {
    const files = event.detail.files;
    if (files.length === 0) return;
    
    isUploadingPhotos = true;
    try {
      const uploadedPhotos = await PhotoService.uploadPhotos(jobId, files);
      const photoUrls = uploadedPhotos.map(p => p.url);
      jobsStore.addPhotos(jobId, photoUrls);
    } catch (error) {
      console.error('Error uploading photos:', error);
    } finally {
      isUploadingPhotos = false;
    }
  }
  
  async function handlePhotosDelete(event: CustomEvent<{ ids: string[] }>) {
    if (!job) return;
    
    const photosToDelete = job.photos.filter((_, index) => 
      event.detail.ids.includes(`photo_${index}`)
    );
    
    await PhotoService.deletePhotos(event.detail.ids);
    jobsStore.removePhotos(jobId, photosToDelete);
  }
  
  async function handlePhotosDownload(event: CustomEvent<{ ids: string[] }>) {
    if (!job) return;
    
    const photosToDownload = job.photos
      .map((url, index) => ({
        id: `photo_${index}`,
        url,
        name: `job_${jobId}_photo_${index + 1}.jpg`
      }))
      .filter(p => event.detail.ids.includes(p.id));
    
    await PhotoService.downloadPhotos(photosToDownload);
  }
  
  $: progress = calculateProgress();
  $: galleryPhotos = job?.photos?.map((url, index) => ({
    id: `photo_${index}`,
    url,
    thumbnail: url,
    name: `Photo ${index + 1}`
  })) || [];
</script>

{#if !job}
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
      <h1>{job.customer.name}</h1>
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
                <div class="item-name">{jobItem.item.name}</div>
                <div class="item-meta">
                  Total: {jobItem.quantity} {jobItem.item.unit}
                  {#if jobItem.installedQuantity > 0}
                    <span class="installed">
                      ({jobItem.installedQuantity} installed)
                    </span>
                  {/if}
                </div>
              </div>
              <QuantityPicker
                value={jobItem.installedQuantity}
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
    bind:open={showPhotoCapture}
    on:capture={handlePhotosCapture}
  />
{/if}

<style>
  .job-detail {
    padding: 1rem;
    max-width: 800px;
    margin: 0 auto;
  }
  
  .error-state {
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