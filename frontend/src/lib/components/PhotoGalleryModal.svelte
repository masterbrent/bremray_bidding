<script lang="ts">
  import Modal from './Modal.svelte';
  import { ChevronLeft, ChevronRight, X, Image as ImageIcon } from 'lucide-svelte';
  
  export let isOpen: boolean = false;
  export let photos: string[] = [];
  export let title: string = 'Photos';
  
  let currentIndex = 0;
  
  function nextPhoto() {
    if (currentIndex < photos.length - 1) {
      currentIndex++;
    }
  }
  
  function previousPhoto() {
    if (currentIndex > 0) {
      currentIndex--;
    }
  }
  
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'ArrowLeft') {
      previousPhoto();
    } else if (event.key === 'ArrowRight') {
      nextPhoto();
    }
  }
  
  $: if (!isOpen) {
    currentIndex = 0;
  }
</script>

<Modal bind:isOpen {title} size="large">
  {#if photos.length === 0}
    <div class="empty-state">
      <ImageIcon size={48} />
      <p>No photos available</p>
    </div>
  {:else}
    <div class="gallery" on:keydown={handleKeydown} tabindex="0">
      <div class="photo-container">
        <img src={photos[currentIndex]} alt="Job photo {currentIndex + 1}" />
        
        {#if photos.length > 1}
          <button 
            class="nav-button prev" 
            on:click={previousPhoto}
            disabled={currentIndex === 0}
            aria-label="Previous photo"
          >
            <ChevronLeft size={24} />
          </button>
          
          <button 
            class="nav-button next" 
            on:click={nextPhoto}
            disabled={currentIndex === photos.length - 1}
            aria-label="Next photo"
          >
            <ChevronRight size={24} />
          </button>
        {/if}
      </div>
      
      {#if photos.length > 1}
        <div class="photo-counter">
          {currentIndex + 1} / {photos.length}
        </div>
        
        <div class="thumbnails">
          {#each photos as photo, index}
            <button 
              class="thumbnail"
              class:active={index === currentIndex}
              on:click={() => currentIndex = index}
              aria-label="View photo {index + 1}"
            >
              <img src={photo} alt="Thumbnail {index + 1}" />
            </button>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</Modal>

<style>
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    color: var(--gray-400);
    gap: 1rem;
  }
  
  .gallery {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    outline: none;
  }
  
  .photo-container {
    position: relative;
    background: var(--gray-100);
    border-radius: 16px;
    overflow: hidden;
    min-height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .photo-container img {
    width: 100%;
    height: auto;
    max-height: 600px;
    object-fit: contain;
  }
  
  .nav-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background: rgba(255, 255, 255, 0.9);
    border: none;
    width: 48px;
    height: 48px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .nav-button:hover:not(:disabled) {
    background: white;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  }
  
  .nav-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .nav-button.prev {
    left: 1rem;
  }
  
  .nav-button.next {
    right: 1rem;
  }
  
  .photo-counter {
    text-align: center;
    font-size: 0.875rem;
    color: var(--gray-600);
    font-weight: 500;
  }
  
  .thumbnails {
    display: flex;
    gap: 0.75rem;
    overflow-x: auto;
    padding: 0.25rem;
  }
  
  .thumbnail {
    flex-shrink: 0;
    width: 80px;
    height: 80px;
    border: 2px solid transparent;
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.2s;
    background: var(--gray-100);
  }
  
  .thumbnail img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .thumbnail:hover {
    transform: scale(1.05);
  }
  
  .thumbnail.active {
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
  }
</style>