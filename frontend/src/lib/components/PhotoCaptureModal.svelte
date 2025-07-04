<script lang="ts">
  import Modal from './Modal.svelte';
  import Button from './Button.svelte';
  import { Camera, Upload, X, Check } from 'lucide-svelte';
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  
  export let isOpen: boolean = false;
  
  const dispatch = createEventDispatcher<{
    capture: { files: File[] };
  }>();
  
  let capturedPhotos: { file: File; url: string }[] = [];
  let cameraStream: MediaStream | null = null;
  let videoElement: HTMLVideoElement | null = null;
  let canvasElement: HTMLCanvasElement | null = null;
  let fileInput: HTMLInputElement | null = null;
  let isUsingCamera = false;
  
  async function startCamera() {
    try {
      // Request camera permission
      cameraStream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: 'environment' }, // Use back camera on mobile
        audio: false
      });
      
      if (videoElement) {
        videoElement.srcObject = cameraStream;
        isUsingCamera = true;
      }
    } catch (error) {
      console.error('Error accessing camera:', error);
      alert('Unable to access camera. Please check permissions.');
    }
  }
  
  function stopCamera() {
    if (cameraStream) {
      cameraStream.getTracks().forEach(track => track.stop());
      cameraStream = null;
    }
    isUsingCamera = false;
  }
  
  function capturePhoto() {
    if (!videoElement || !canvasElement) return;
    
    const context = canvasElement.getContext('2d');
    if (!context) return;
    
    // Set canvas size to match video
    canvasElement.width = videoElement.videoWidth;
    canvasElement.height = videoElement.videoHeight;
    
    // Draw video frame to canvas
    context.drawImage(videoElement, 0, 0);
    
    // Convert canvas to blob
    canvasElement.toBlob((blob) => {
      if (blob) {
        const file = new File([blob], `photo_${Date.now()}.jpg`, { type: 'image/jpeg' });
        const url = URL.createObjectURL(blob);
        capturedPhotos = [...capturedPhotos, { file, url }];
      }
    }, 'image/jpeg', 0.9);
  }
  
  function handleFileSelect(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files) return;
    
    const newPhotos = Array.from(input.files).map(file => ({
      file,
      url: URL.createObjectURL(file)
    }));
    
    capturedPhotos = [...capturedPhotos, ...newPhotos];
  }
  
  function removePhoto(index: number) {
    URL.revokeObjectURL(capturedPhotos[index].url);
    capturedPhotos = capturedPhotos.filter((_, i) => i !== index);
  }
  
  function savePhotos() {
    const files = capturedPhotos.map(p => p.file);
    dispatch('capture', { files });
    closeModal();
  }
  
  function closeModal() {
    // Clean up
    capturedPhotos.forEach(photo => URL.revokeObjectURL(photo.url));
    capturedPhotos = [];
    stopCamera();
    isOpen = false;
  }
  
  onDestroy(() => {
    stopCamera();
    capturedPhotos.forEach(photo => URL.revokeObjectURL(photo.url));
  });
</script>

<Modal bind:isOpen title="Add Photos" size="large" on:close={closeModal}>
  <div class="photo-capture">
    {#if !isUsingCamera && capturedPhotos.length === 0}
      <!-- Initial options -->
      <div class="capture-options">
        <button class="option-card" on:click={startCamera}>
          <Camera size={48} />
          <h3>Use Camera</h3>
          <p>Take photos with your device camera</p>
        </button>
        
        <button class="option-card" on:click={() => fileInput?.click()}>
          <Upload size={48} />
          <h3>Upload Photos</h3>
          <p>Select photos from your device</p>
        </button>
      </div>
      
      <input
        bind:this={fileInput}
        type="file"
        accept="image/*"
        multiple
        on:change={handleFileSelect}
        style="display: none;"
      />
    {:else if isUsingCamera}
      <!-- Camera view -->
      <div class="camera-view">
        <video 
          bind:this={videoElement}
          autoplay
          playsinline
          class="camera-feed"
        ></video>
        
        <canvas 
          bind:this={canvasElement}
          style="display: none;"
        ></canvas>
        
        <div class="camera-controls">
          <Button variant="ghost" on:click={stopCamera}>
            <X size={20} />
            Close Camera
          </Button>
          
          <button class="capture-button" on:click={capturePhoto}>
            <div class="capture-button-inner"></div>
          </button>
          
          <Button variant="ghost" disabled>
            {capturedPhotos.length} photos
          </Button>
        </div>
      </div>
    {/if}
    
    {#if capturedPhotos.length > 0}
      <!-- Photo preview -->
      <div class="photo-preview">
        <h3>Captured Photos ({capturedPhotos.length})</h3>
        
        <div class="preview-grid">
          {#each capturedPhotos as photo, index}
            <div class="preview-item">
              <img src={photo.url} alt="Captured photo {index + 1}" />
              <button 
                class="remove-button"
                on:click={() => removePhoto(index)}
                title="Remove photo"
              >
                <X size={16} />
              </button>
            </div>
          {/each}
        </div>
        
        {#if isUsingCamera}
          <div class="continue-actions">
            <Button variant="secondary" on:click={capturePhoto}>
              Take Another Photo
            </Button>
          </div>
        {:else}
          <div class="add-more-actions">
            <Button variant="ghost" on:click={startCamera}>
              <Camera size={20} />
              Use Camera
            </Button>
            <Button variant="ghost" on:click={() => fileInput?.click()}>
              <Upload size={20} />
              Add More
            </Button>
          </div>
        {/if}
      </div>
    {/if}
  </div>
  
  <div slot="footer">
    <Button variant="ghost" on:click={closeModal}>
      Cancel
    </Button>
    <Button 
      variant="primary" 
      on:click={savePhotos}
      disabled={capturedPhotos.length === 0}
    >
      <Check size={20} />
      Add {capturedPhotos.length} Photo{capturedPhotos.length !== 1 ? 's' : ''}
    </Button>
  </div>
</Modal>

<style>
  .photo-capture {
    min-height: 400px;
  }
  
  .capture-options {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    padding: 3rem 1rem;
  }
  
  .option-card {
    background: var(--gray-50);
    border: 2px solid var(--gray-200);
    border-radius: 16px;
    padding: 3rem 2rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;
    font-family: inherit;
  }
  
  .option-card:hover {
    border-color: var(--primary-500);
    background: var(--primary-50);
    transform: translateY(-2px);
  }
  
  .option-card h3 {
    margin: 1rem 0 0.5rem 0;
    color: var(--gray-900);
  }
  
  .option-card p {
    margin: 0;
    color: var(--gray-600);
  }
  
  .camera-view {
    position: relative;
    background: black;
    border-radius: 16px;
    overflow: hidden;
  }
  
  .camera-feed {
    width: 100%;
    height: auto;
    display: block;
  }
  
  .camera-controls {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 2rem;
    background: linear-gradient(to top, rgba(0,0,0,0.8), transparent);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .capture-button {
    width: 70px;
    height: 70px;
    border-radius: 50%;
    border: 4px solid white;
    background: transparent;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.1s;
  }
  
  .capture-button:active {
    transform: scale(0.9);
  }
  
  .capture-button-inner {
    width: 56px;
    height: 56px;
    border-radius: 50%;
    background: white;
  }
  
  .photo-preview {
    margin-top: 2rem;
  }
  
  .photo-preview h3 {
    margin: 0 0 1rem 0;
    color: var(--gray-900);
  }
  
  .preview-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }
  
  .preview-item {
    position: relative;
    aspect-ratio: 1;
    border-radius: 12px;
    overflow: hidden;
    background: var(--gray-100);
  }
  
  .preview-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .remove-button {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    width: 28px;
    height: 28px;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    border: none;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.2s;
  }
  
  .preview-item:hover .remove-button {
    opacity: 1;
  }
  
  .continue-actions,
  .add-more-actions {
    display: flex;
    gap: 1rem;
    justify-content: center;
  }
  
  @media (max-width: 768px) {
    .capture-options {
      grid-template-columns: 1fr;
      gap: 1rem;
      padding: 2rem 1rem;
    }
    
    .option-card {
      padding: 2rem 1.5rem;
    }
    
    .preview-grid {
      grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    }
  }
</style>