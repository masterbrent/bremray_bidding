<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Eye, Download, Trash2, Check, X, ChevronLeft, ChevronRight } from 'lucide-svelte';
  
  export let photos: { id: string; url: string; thumbnail?: string; name?: string }[] = [];
  
  const dispatch = createEventDispatcher<{
    delete: { ids: string[] };
    download: { ids: string[] };
  }>();
  
  let selectedIds = new Set<string>();
  let isSelecting = false;
  let fullViewIndex: number | null = null;
  
  function toggleSelection(id: string, event?: Event) {
    if (event) {
      event.stopPropagation();
    }
    
    if (selectedIds.has(id)) {
      selectedIds.delete(id);
    } else {
      selectedIds.add(id);
    }
    selectedIds = selectedIds; // Trigger reactivity
    
    if (!isSelecting && selectedIds.size > 0) {
      isSelecting = true;
    } else if (selectedIds.size === 0) {
      isSelecting = false;
    }
  }
  
  function selectAll() {
    selectedIds = new Set(photos.map(p => p.id));
    isSelecting = true;
  }
  
  function clearSelection() {
    selectedIds.clear();
    selectedIds = selectedIds;
    isSelecting = false;
  }
  
  function openFullView(index: number) {
    if (!isSelecting) {
      fullViewIndex = index;
    }
  }
  
  function closeFullView() {
    fullViewIndex = null;
  }
  
  function navigateFullView(direction: 'prev' | 'next') {
    if (fullViewIndex === null) return;
    
    if (direction === 'prev' && fullViewIndex > 0) {
      fullViewIndex--;
    } else if (direction === 'next' && fullViewIndex < photos.length - 1) {
      fullViewIndex++;
    }
  }
  
  function handleKeydown(event: KeyboardEvent) {
    if (fullViewIndex !== null) {
      if (event.key === 'ArrowLeft') {
        navigateFullView('prev');
      } else if (event.key === 'ArrowRight') {
        navigateFullView('next');
      } else if (event.key === 'Escape') {
        closeFullView();
      }
    }
  }
  
  function deleteSelected() {
    if (selectedIds.size > 0) {
      dispatch('delete', { ids: Array.from(selectedIds) });
      clearSelection();
    }
  }
  
  function downloadSelected() {
    if (selectedIds.size > 0) {
      dispatch('download', { ids: Array.from(selectedIds) });
    }
  }
  
  function downloadSingle(id: string, event: Event) {
    event.stopPropagation();
    dispatch('download', { ids: [id] });
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="photo-gallery">
  {#if photos.length === 0}
    <div class="empty-state">
      <p>No photos yet. Click "Take Photos" to add some.</p>
    </div>
  {:else}
    <!-- Selection toolbar -->
    {#if isSelecting}
      <div class="selection-toolbar">
        <div class="selection-info">
          <button class="select-toggle" on:click={clearSelection}>
            <X size={20} />
          </button>
          <span>{selectedIds.size} selected</span>
          <button class="select-all-btn" on:click={selectAll}>
            Select all
          </button>
        </div>
        
        <div class="selection-actions">
          <button class="action-btn download" on:click={downloadSelected}>
            <Download size={18} />
            Download
          </button>
          <button class="action-btn delete" on:click={deleteSelected}>
            <Trash2 size={18} />
            Delete
          </button>
        </div>
      </div>
    {/if}
    
    <!-- Photo grid -->
    <div class="photo-grid">
      {#each photos as photo, index}
        <div 
          class="photo-item"
          class:selected={selectedIds.has(photo.id)}
          on:click={() => openFullView(index)}
        >
          <img 
            src={photo.thumbnail || photo.url} 
            alt={photo.name || `Photo ${index + 1}`}
            loading="lazy"
          />
          
          <div class="photo-overlay">
            <button 
              class="select-checkbox"
              class:visible={isSelecting || selectedIds.has(photo.id)}
              on:click={(e) => toggleSelection(photo.id, e)}
            >
              {#if selectedIds.has(photo.id)}
                <Check size={16} />
              {/if}
            </button>
            
            {#if !isSelecting}
              <div class="photo-actions">
                <button 
                  class="photo-action"
                  on:click={(e) => downloadSingle(photo.id, e)}
                  title="Download"
                >
                  <Download size={16} />
                </button>
              </div>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<!-- Full view modal -->
{#if fullViewIndex !== null}
  <div class="full-view-modal" on:click={closeFullView}>
    <div class="full-view-content" on:click|stopPropagation>
      <button class="close-btn" on:click={closeFullView}>
        <X size={24} />
      </button>
      
      <img 
        src={photos[fullViewIndex].url} 
        alt={photos[fullViewIndex].name || `Photo ${fullViewIndex + 1}`}
      />
      
      {#if fullViewIndex > 0}
        <button 
          class="nav-btn prev"
          on:click={() => navigateFullView('prev')}
        >
          <ChevronLeft size={24} />
        </button>
      {/if}
      
      {#if fullViewIndex < photos.length - 1}
        <button 
          class="nav-btn next"
          on:click={() => navigateFullView('next')}
        >
          <ChevronRight size={24} />
        </button>
      {/if}
      
      <div class="full-view-info">
        {fullViewIndex + 1} / {photos.length}
      </div>
    </div>
  </div>
{/if}

<style>
  .photo-gallery {
    margin-top: 2rem;
  }
  
  .empty-state {
    text-align: center;
    padding: 3rem;
    background: var(--gray-50);
    border-radius: 16px;
    color: var(--gray-500);
  }
  
  .selection-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: var(--gray-100);
    border-radius: 12px;
    margin-bottom: 1rem;
  }
  
  .selection-info {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .select-toggle {
    width: 32px;
    height: 32px;
    border: none;
    background: var(--gray-300);
    color: var(--gray-700);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .select-toggle:hover {
    background: var(--gray-400);
  }
  
  .select-all-btn {
    background: none;
    border: none;
    color: var(--primary-600);
    font-weight: 500;
    cursor: pointer;
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
    transition: all 0.2s;
  }
  
  .select-all-btn:hover {
    background: var(--primary-100);
  }
  
  .selection-actions {
    display: flex;
    gap: 0.5rem;
  }
  
  .action-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 8px;
    font-weight: 500;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: all 0.2s;
  }
  
  .action-btn.download {
    background: var(--primary-100);
    color: var(--primary-700);
  }
  
  .action-btn.download:hover {
    background: var(--primary-200);
  }
  
  .action-btn.delete {
    background: var(--danger-100);
    color: var(--danger-700);
  }
  
  .action-btn.delete:hover {
    background: var(--danger-200);
  }
  
  .photo-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .photo-item {
    position: relative;
    aspect-ratio: 1;
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    background: var(--gray-100);
    transition: all 0.2s;
  }
  
  .photo-item:hover {
    transform: scale(1.02);
  }
  
  .photo-item.selected {
    outline: 3px solid var(--primary-500);
    outline-offset: -3px;
  }
  
  .photo-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .photo-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(to bottom, rgba(0,0,0,0.3) 0%, transparent 50%, rgba(0,0,0,0.3) 100%);
    opacity: 0;
    transition: opacity 0.2s;
  }
  
  .photo-item:hover .photo-overlay,
  .photo-item.selected .photo-overlay {
    opacity: 1;
  }
  
  .select-checkbox {
    position: absolute;
    top: 0.5rem;
    left: 0.5rem;
    width: 24px;
    height: 24px;
    background: white;
    border: 2px solid var(--gray-300);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    opacity: 0;
    transition: all 0.2s;
  }
  
  .select-checkbox.visible,
  .photo-item:hover .select-checkbox {
    opacity: 1;
  }
  
  .photo-item.selected .select-checkbox {
    background: var(--primary-500);
    border-color: var(--primary-500);
    color: white;
  }
  
  .photo-actions {
    position: absolute;
    bottom: 0.5rem;
    right: 0.5rem;
    display: flex;
    gap: 0.5rem;
  }
  
  .photo-action {
    width: 32px;
    height: 32px;
    background: rgba(255, 255, 255, 0.9);
    border: none;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .photo-action:hover {
    background: white;
    transform: scale(1.1);
  }
  
  /* Full view modal */
  .full-view-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.9);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .full-view-content {
    position: relative;
    max-width: 90vw;
    max-height: 90vh;
  }
  
  .full-view-content img {
    max-width: 100%;
    max-height: 90vh;
    object-fit: contain;
  }
  
  .close-btn {
    position: absolute;
    top: -3rem;
    right: 0;
    background: rgba(255, 255, 255, 0.1);
    color: white;
    border: none;
    width: 40px;
    height: 40px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .close-btn:hover {
    background: rgba(255, 255, 255, 0.2);
  }
  
  .nav-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background: rgba(255, 255, 255, 0.1);
    color: white;
    border: none;
    width: 48px;
    height: 48px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .nav-btn:hover {
    background: rgba(255, 255, 255, 0.2);
  }
  
  .nav-btn.prev {
    left: -4rem;
  }
  
  .nav-btn.next {
    right: -4rem;
  }
  
  .full-view-info {
    position: absolute;
    bottom: -2rem;
    left: 50%;
    transform: translateX(-50%);
    color: white;
    font-size: 0.875rem;
  }
  
  @media (max-width: 768px) {
    .photo-grid {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    }
    
    .selection-toolbar {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }
    
    .nav-btn.prev {
      left: 1rem;
    }
    
    .nav-btn.next {
      right: 1rem;
    }
  }
</style>