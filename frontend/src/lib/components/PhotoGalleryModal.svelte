<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import Modal from './Modal.svelte';
  import PhotoGallery from './PhotoGallery.svelte';
  import ConfirmModal from './ConfirmModal.svelte';
  import { PhotoService } from '../services/photoService';
  import { jobsStore, jobs } from '../stores';
  import type { Job } from '../types/models';
  
  export let isOpen = false;
  export let job: Job | null = null;
  
  const dispatch = createEventDispatcher();
  
  let isDeletingPhotos = false;
  let deleteConfirm = {
    show: false,
    photoIds: [] as string[],
    message: ''
  };
  
  function handleClose() {
    isOpen = false;
  }
  
  function handleDelete(event: CustomEvent<{ ids: string[] }>) {
    if (!job || isDeletingPhotos) return;
    
    const photoIds = event.detail.ids.map(id => id.replace('photo_', ''));
    
    if (photoIds.length === 0) return;
    
    deleteConfirm = {
      show: true,
      photoIds: photoIds,
      message: photoIds.length === 1 
        ? 'Are you sure you want to delete this photo?' 
        : `Are you sure you want to delete ${photoIds.length} photos?`
    };
  }
  
  async function confirmDelete() {
    if (!job || deleteConfirm.photoIds.length === 0) return;
    
    try {
      isDeletingPhotos = true;
      
      // Delete photos via API
      await PhotoService.deletePhotos(job.id, deleteConfirm.photoIds);
      
      // Reload job to get updated photos
      await jobsStore.loadById(job.id);
      
      // Update the job reference
      job = $jobs.find(j => j.id === job!.id) || null;
      
      // Dispatch event so parent can update if needed
      dispatch('photosDeleted', { jobId: job?.id, photoIds: deleteConfirm.photoIds });
      
      // Close the confirmation modal
      deleteConfirm = { show: false, photoIds: [], message: '' };
      
    } catch (error) {
      console.error('Error deleting photos:', error);
      // Close the modal and we could show an error modal here
      deleteConfirm = { show: false, photoIds: [], message: '' };
    } finally {
      isDeletingPhotos = false;
    }
  }
  
  function handleDownload(event: CustomEvent<{ ids: string[] }>) {
    if (!job) return;
    
    const photos = job.photos || [];
    const selectedPhotos = photos
      .filter(photo => event.detail.ids.includes(`photo_${photo.id}`))
      .map(photo => ({
        id: photo.id,
        url: photo.url,
        name: `job_${job.id}_photo_${photo.id}.jpg`
      }));
    
    // Download each photo
    selectedPhotos.forEach(photo => {
      const link = document.createElement('a');
      link.href = photo.url;
      link.download = photo.name;
      link.target = '_blank';
      link.click();
    });
  }
  
  $: photos = job?.photos?.map(photo => ({
    id: `photo_${photo.id}`,
    url: photo.url,
    thumbnail: photo.url,
    name: photo.caption || `Photo ${photo.id.slice(0, 8)}`
  })) || [];
</script>

<Modal bind:isOpen={isOpen} on:close={handleClose} title="{job?.customer?.name || 'Job'} Photos">
  <div class="gallery-modal-content">
    {#if job}
      <div class="job-info">
        <h3>{job.customer?.name || 'Unknown Customer'}</h3>
        <p class="job-address">{job.customer?.address || 'No address'}</p>
        <p class="photo-count">{photos.length} photo{photos.length !== 1 ? 's' : ''}</p>
      </div>
      
      <PhotoGallery 
        {photos}
        on:delete={handleDelete}
        on:download={handleDownload}
      />
    {/if}
  </div>
</Modal>

<ConfirmModal
  bind:isOpen={deleteConfirm.show}
  title="Delete Photos"
  message={deleteConfirm.message}
  confirmText="Delete"
  variant="danger"
  onConfirm={confirmDelete}
  onCancel={() => deleteConfirm = { show: false, photoIds: [], message: '' }}
/>

<style>
  .gallery-modal-content {
    padding: 1rem;
    max-height: 80vh;
    overflow-y: auto;
  }
  
  .job-info {
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #e5e7eb;
  }
  
  .job-info h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: #111827;
  }
  
  .job-address {
    margin: 0 0 0.25rem 0;
    color: #6b7280;
    font-size: 0.875rem;
  }
  
  .photo-count {
    margin: 0;
    color: #6b7280;
    font-size: 0.875rem;
    font-weight: 500;
  }
  
  @media (max-width: 640px) {
    .gallery-modal-content {
      padding: 0.5rem;
    }
  }
</style>
