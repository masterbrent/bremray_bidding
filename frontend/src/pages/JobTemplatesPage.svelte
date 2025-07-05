<script lang="ts">
  import { onMount } from 'svelte';
  import { jobTemplatesStore, templates } from '../lib/stores/jobTemplates';
  import { itemsStore, items } from '../lib/stores/items';
  import type { JobTemplate, Item, TemplatePhase } from '../lib/types/models';
  import Button from '../lib/components/Button.svelte';
  import Modal from '../lib/components/Modal.svelte';
  import ConfirmModal from '../lib/components/ConfirmModal.svelte';
  
  let selectedTemplateId: string | null = null;
  let showCreateModal = false;
  let editingTemplate: JobTemplate | null = null;
  let deleteConfirm = {
    show: false,
    template: null as JobTemplate | null
  };
  let formError = '';
  let newPhaseForTemplate = ''; // For adding phases to existing templates
  let draggedPhaseIndex: number | null = null;
  let templateDraggedPhaseIndex: number | null = null; // For dragging phases in existing templates
  
  let formData = {
    name: '',
    description: '',
    items: [] as { itemId: string; defaultQuantity: number }[],
    phases: [] as { name: string; order: number; description?: string }[]
  };
  
  function handleDragStart(e: DragEvent, index: number) {
    draggedPhaseIndex = index;
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = 'move';
    }
  }
  
  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'move';
    }
  }
  
  function handleDrop(e: DragEvent, dropIndex: number) {
    e.preventDefault();
    if (draggedPhaseIndex === null || draggedPhaseIndex === dropIndex) return;
    
    const phases = [...formData.phases];
    const [draggedPhase] = phases.splice(draggedPhaseIndex, 1);
    phases.splice(dropIndex, 0, draggedPhase);
    
    // Update order numbers
    phases.forEach((phase, i) => {
      phase.order = i + 1;
    });
    
    formData.phases = phases;
    draggedPhaseIndex = null;
  }
  
  // Drag handlers for template detail view
  function handleTemplateDragStart(e: DragEvent, index: number) {
    templateDraggedPhaseIndex = index;
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = 'move';
    }
  }
  
  function handleTemplateDragOver(e: DragEvent) {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'move';
    }
  }
  
  async function handleTemplateDrop(e: DragEvent, dropIndex: number) {
    e.preventDefault();
    if (!selectedTemplate || templateDraggedPhaseIndex === null || templateDraggedPhaseIndex === dropIndex) return;
    
    const phases = [...(selectedTemplate.phases || [])];
    const [draggedPhase] = phases.splice(templateDraggedPhaseIndex, 1);
    phases.splice(dropIndex, 0, draggedPhase);
    
    // Update order numbers
    phases.forEach((phase, i) => {
      phase.order = i + 1;
    });
    
    // Update the template
    try {
      await jobTemplatesStore.update(selectedTemplate.id, {
        name: selectedTemplate.name,
        description: selectedTemplate.description,
        isActive: selectedTemplate.isActive,
        phases
      });
    } catch (error) {
      console.error('Failed to reorder phases:', error);
    }
    
    templateDraggedPhaseIndex = null;
  }
  
  async function deletePhaseFromTemplate(phaseIndex: number) {
    if (!selectedTemplate) return;
    
    // Remove the phase at the given index
    const phases = selectedTemplate.phases.filter((_, index) => index !== phaseIndex);
    
    // Reorder remaining phases
    phases.forEach((phase, i) => {
      phase.order = i + 1;
    });
    
    // Update the template
    try {
      await jobTemplatesStore.update(selectedTemplate.id, {
        name: selectedTemplate.name,
        description: selectedTemplate.description,
        isActive: selectedTemplate.isActive,
        phases
      });
    } catch (error) {
      console.error('Failed to delete phase:', error);
    }
  }
  
  onMount(() => {
    jobTemplatesStore.load();
    itemsStore.load();
  });
  
  function openCreateModal() {
    formData = {
      name: '',
      description: '',
      items: [],
      phases: []
    };
    editingTemplate = null;
    formError = '';
    showCreateModal = true;
  }
  
  async function handleSave() {
    if (!formData.name) {
      formError = 'Please enter a template name';
      return;
    }
    
    if (formData.items.length === 0) {
      formError = 'Please add at least one item to the template';
      return;
    }
    
    formError = '';
    
    try {
      await jobTemplatesStore.add(
        formData.name,
        formData.description,
        formData.items,
        formData.phases
      );
      showCreateModal = false;
    } catch (error) {
      formError = error instanceof Error ? error.message : 'Failed to create template';
    }
  }
  
  function addItemToTemplate(item: Item) {
    if (!formData.items.some(i => i.itemId === item.id)) {
      formData.items = [...formData.items, { itemId: item.id, defaultQuantity: 0 }];
    }
  }
  
  function removeItemFromTemplate(itemId: string) {
    formData.items = formData.items.filter(i => i.itemId !== itemId);
  }
  
  function addPhase() {
    const maxOrder = Math.max(0, ...formData.phases.map(p => p.order));
    formData.phases = [...formData.phases, {
      name: '',
      order: maxOrder + 1,
      description: ''
    }];
  }
  
  function removePhase(index: number) {
    formData.phases = formData.phases.filter((_, i) => i !== index);
    // Reorder remaining phases
    formData.phases = formData.phases.map((phase, i) => ({
      ...phase,
      order: i + 1
    }));
  }
  
  function movePhaseUp(index: number) {
    if (index === 0) return;
    const phases = [...formData.phases];
    [phases[index], phases[index - 1]] = [phases[index - 1], phases[index]];
    phases.forEach((phase, i) => phase.order = i + 1);
    formData.phases = phases;
  }
  
  function movePhaseDown(index: number) {
    if (index === formData.phases.length - 1) return;
    const phases = [...formData.phases];
    [phases[index], phases[index + 1]] = [phases[index + 1], phases[index]];
    phases.forEach((phase, i) => phase.order = i + 1);
    formData.phases = phases;
  }
  
  function confirmDelete(template: JobTemplate) {
    deleteConfirm = {
      show: true,
      template
    };
  }
  
  async function handleDelete() {
    if (deleteConfirm.template) {
      const id = deleteConfirm.template.id;
      try {
        await jobTemplatesStore.remove(id);
        if (selectedTemplateId === id) {
          selectedTemplateId = null;
        }
        deleteConfirm = { show: false, template: null };
      } catch (error) {
        console.error('Failed to delete template:', error);
      }
    }
  }
  
  async function toggleItemInTemplate(templateId: string, itemId: string) {
    const template = $templates.find(t => t.id === templateId);
    if (!template) return;
    
    const hasItem = template.items.some(i => i.itemId === itemId);
    try {
      if (hasItem) {
        await jobTemplatesStore.removeItemFromTemplate(templateId, itemId);
      } else {
        await jobTemplatesStore.addItemToTemplate(templateId, itemId, 0); // Always 0 quantity
      }
    } catch (error) {
      console.error('Failed to update template:', error);
    }
  }
  
  async function addPhaseToTemplate() {
    if (!selectedTemplate || !newPhaseForTemplate.trim()) return;
    
    console.log('Adding phase to template:', selectedTemplate.id);
    console.log('Current phases:', selectedTemplate.phases);
    
    const maxOrder = Math.max(0, ...(selectedTemplate.phases?.map(p => p.order) || []));
    const phases = [...(selectedTemplate.phases || []), {
      name: newPhaseForTemplate.trim(),
      order: maxOrder + 1,
      description: ''
    }];
    
    console.log('New phases array:', phases);
    
    // Update template with new phases
    try {
      // Send full update including existing data
      await jobTemplatesStore.update(selectedTemplate.id, {
        name: selectedTemplate.name,
        description: selectedTemplate.description,
        isActive: selectedTemplate.isActive,
        phases
      });
      newPhaseForTemplate = '';
      console.log('Phase added successfully');
    } catch (error) {
      console.error('Failed to add phase:', error);
    }
  }
  
  $: selectedTemplate = selectedTemplateId ? $templates.find(t => t.id === selectedTemplateId) : null;
  $: availableItems = $items.filter(item => !formData.items.some(ti => ti.itemId === item.id));
  $: selectedItems = formData.items.map(ti => {
    const item = $items.find(i => i.id === ti.itemId);
    return item ? { ...item, defaultQuantity: ti.defaultQuantity } : null;
  }).filter(Boolean);
</script>
<div class="templates-page">
  <div class="page-header">
    <h1>Job Templates</h1>
    <Button on:click={openCreateModal}>
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
      Create Template
    </Button>
  </div>
  
  <div class="layout">
    <div class="templates-list">
      {#if $jobTemplatesStore.loading}
        <p class="empty">Loading templates...</p>
      {:else if $jobTemplatesStore.error}
        <p class="error-message">{$jobTemplatesStore.error}</p>
      {:else if $templates.length === 0}
        <p class="empty">No templates yet</p>
      {:else}
        {#each $templates as template}
          <div 
            class="template-card"
            class:selected={selectedTemplateId === template.id}
            on:click={() => selectedTemplateId = template.id}
          >
            <div class="template-header">
              <h3>{template.name}</h3>
              <button 
                class="delete-btn"
                on:click|stopPropagation={() => confirmDelete(template)}
                title="Delete template"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2" />
                </svg>
              </button>
            </div>
            <p class="template-description">{template.description || 'No description'}</p>
            <div class="template-stats">
              <span>{template.items.length} items</span>
              <span>{template.phases?.length || 0} phases</span>
            </div>
          </div>
        {/each}
      {/if}
    </div>
    
    {#if selectedTemplate}
      <div class="template-details">
        <h2>{selectedTemplate.name}</h2>
        <p>{selectedTemplate.description || 'No description'}</p>
        
        <div class="template-sections">
          <div class="section">
            <h3>Items</h3>
            <div class="items-manager">
              <div class="available-items">
                <h4>Available Items</h4>
                <div class="items-list">
                  {#each $items as item}
                    {@const hasItem = selectedTemplate.items.some(i => i.itemId === item.id)}
                    {#if !hasItem}
                      <div class="item-row" on:click={() => toggleItemInTemplate(selectedTemplate.id, item.id)}>
                        <span class="item-name">{item.name}</span>
                        <span class="add-arrow">→</span>
                      </div>
                    {/if}
                  {/each}
                  {#if $items.every(item => selectedTemplate.items.some(i => i.itemId === item.id))}
                    <p class="empty-message">All items added</p>
                  {/if}
                </div>
              </div>
              
              <div class="template-items">
                <h4>Template Items</h4>
                <div class="items-list">
                  {#each selectedTemplate.items as templateItem}
                    {@const item = $items.find(i => i.id === templateItem.itemId)}
                    {#if item}
                      <div class="item-row" on:click={() => toggleItemInTemplate(selectedTemplate.id, item.id)}>
                        <span class="remove-arrow">←</span>
                        <span class="item-name">{item.name}</span>
                      </div>
                    {/if}
                  {/each}
                  {#if selectedTemplate.items.length === 0}
                    <p class="empty-message">No items selected</p>
                  {/if}
                </div>
              </div>
            </div>
          </div>
          
          {#if selectedTemplate.phases && selectedTemplate.phases.length > 0}
            <div class="section">
              <h3>Phases</h3>
              <div class="phases-list">
                {#each selectedTemplate.phases.sort((a, b) => a.order - b.order) as phase, index}
                  <div 
                    class="phase-item"
                    draggable="true"
                    on:dragstart={(e) => handleTemplateDragStart(e, index)}
                    on:dragover={handleTemplateDragOver}
                    on:drop={(e) => handleTemplateDrop(e, index)}
                  >
                    <div class="drag-handle">⋮⋮</div>
                    <span class="phase-order">{phase.order}</span>
                    <span class="phase-name">{phase.name}</span>
                    <button 
                      class="delete-phase-btn"
                      on:click={() => deletePhaseFromTemplate(index)}
                      title="Delete phase"
                    >
                      ×
                    </button>
                  </div>
                {/each}
              </div>
              
              <div class="add-phase-section">
                <input
                  type="text"
                  bind:value={newPhaseForTemplate}
                  placeholder="Add a new phase..."
                  class="phase-add-input"
                  on:keydown={(e) => e.key === 'Enter' && addPhaseToTemplate()}
                />
                <Button 
                  size="small" 
                  variant="secondary" 
                  on:click={addPhaseToTemplate}
                  disabled={!newPhaseForTemplate.trim()}
                >
                  Add Phase
                </Button>
              </div>
            </div>
          {:else}
            <div class="section">
              <h3>Phases</h3>
              <p class="empty-message">No phases defined</p>
              
              <div class="add-phase-section">
                <input
                  type="text"
                  bind:value={newPhaseForTemplate}
                  placeholder="Add the first phase..."
                  class="phase-add-input"
                  on:keydown={(e) => e.key === 'Enter' && addPhaseToTemplate()}
                />
                <Button 
                  size="small" 
                  variant="secondary" 
                  on:click={addPhaseToTemplate}
                  disabled={!newPhaseForTemplate.trim()}
                >
                  Add Phase
                </Button>
              </div>
            </div>
          {/if}
        </div>
      </div>
    {:else}
      <div class="empty-state">
        <p>Select a template to view details</p>
      </div>
    {/if}
  </div>
</div>
<Modal bind:isOpen={showCreateModal} title="Create Job Template">
  {#if formError}
    <div class="error-message">
      {formError}
    </div>
  {/if}
  
  <form on:submit|preventDefault={handleSave}>
    <div class="form-group">
      <label for="template-name">Template Name *</label>
      <input
        id="template-name"
        type="text"
        bind:value={formData.name}
        required
      />
    </div>
    
    <div class="form-group">
      <label for="description">Description</label>
      <textarea
        id="description"
        bind:value={formData.description}
        rows="2"
      ></textarea>
    </div>
    
    <div class="form-section">
      <h3>Job Phases</h3>
      <div class="phases-container">
        {#each formData.phases as phase, index}
          <div 
            class="phase-input-row" 
            draggable="true"
            on:dragstart={(e) => handleDragStart(e, index)}
            on:dragover={handleDragOver}
            on:drop={(e) => handleDrop(e, index)}
          >
            <div class="drag-handle">⋮⋮</div>
            <input
              type="text"
              bind:value={phase.name}
              placeholder="Phase name"
              class="phase-name-input"
            />
            <button type="button" on:click={() => removePhase(index)} class="remove-phase">×</button>
          </div>
        {/each}
        <button type="button" class="add-phase-btn" on:click={addPhase}>
          + Add Phase
        </button>
      </div>
    </div>
    
    <div class="form-section">
      <h3>Items *</h3>
      <div class="items-picker">
        <div class="picker-column">
          <h4>Available Items</h4>
          <div class="items-scroll">
            {#each availableItems as item}
              <div class="pickable-item" on:click={() => addItemToTemplate(item)}>
                <span class="item-name">{item.name}</span>
                <span class="add-icon">→</span>
              </div>
            {/each}
            {#if availableItems.length === 0}
              <p class="empty-message">All items added</p>
            {/if}
          </div>
        </div>
        
        <div class="picker-column">
          <h4>Selected Items ({selectedItems.length})</h4>
          <div class="items-scroll">
            {#each formData.items as templateItem, index}
              {@const item = $items.find(i => i.id === templateItem.itemId)}
              {#if item}
                <div class="selected-item">
                  <span class="remove-icon" on:click={() => removeItemFromTemplate(item.id)}>←</span>
                  <span class="item-name">{item.name}</span>
                </div>
              {/if}
            {/each}
            {#if selectedItems.length === 0}
              <p class="empty-message">Click items to add</p>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </form>
  
  <div slot="footer">
    <Button variant="secondary" on:click={() => showCreateModal = false}>
      Cancel
    </Button>
    <Button on:click={handleSave}>
      Create Template
    </Button>
  </div>
</Modal>
<ConfirmModal
  bind:isOpen={deleteConfirm.show}
  title="Delete Template"
  message={`Are you sure you want to delete "${deleteConfirm.template?.name}"? This action cannot be undone.`}
  confirmText="Delete"
  onConfirm={handleDelete}
  onCancel={() => deleteConfirm = { show: false, template: null }}
/>

<style>
  .templates-page {
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
  
  .page-header h1 {
    margin: 0;
    color: #333;
  }
  
  .icon {
    width: 1.25rem;
    height: 1.25rem;
    margin-right: 0.5rem;
  }
  
  .layout {
    display: grid;
    grid-template-columns: 350px 1fr;
    gap: 2rem;
    height: calc(100vh - 200px);
  }
  
  .templates-list {
    overflow-y: auto;
    padding-right: 1rem;
  }
  
  .template-card {
    background: white;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 1rem;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .template-card:hover {
    border-color: #3b82f6;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }
  
  .template-card.selected {
    border-color: #3b82f6;
    background: #f0f7ff;
  }
  
  .template-header {
    display: flex;
    justify-content: space-between;
    align-items: start;
  }
  
  .template-header h3 {
    margin: 0 0 0.5rem 0;
    color: #333;
    font-size: 1.1rem;
  }
  
  .delete-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.25rem;
    color: #666;
    transition: color 0.2s;
  }
  
  .delete-btn:hover {
    color: #dc2626;
  }
  
  .delete-btn svg {
    width: 1.25rem;
    height: 1.25rem;
  }
  
  .template-description {
    color: #666;
    font-size: 0.875rem;
    margin: 0 0 0.75rem 0;
  }
  
  .template-stats {
    display: flex;
    gap: 1rem;
    font-size: 0.875rem;
    color: #999;
  }
  
  .template-details,
  .empty-state {
    background: white;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    padding: 2rem;
    overflow-y: auto;
  }
  
  .template-details h2 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }
  
  .template-sections {
    margin-top: 2rem;
  }
  
  .section {
    margin-bottom: 2rem;
  }
  
  .section h3 {
    margin: 0 0 1rem 0;
    color: #333;
    font-size: 1.1rem;
  }
  
  .items-manager {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  .available-items,
  .template-items {
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    padding: 1rem;
    background: #fafafa;
  }
  
  .available-items h4,
  .template-items h4 {
    margin: 0 0 0.75rem 0;
    color: #666;
    font-size: 0.875rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  
  .items-list {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .item-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    margin: 0;
    background: white;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .item-row:hover {
    background: #f5f5f5;
    border-color: #3b82f6;
    transform: translateX(2px);
  }
  
  .item-name {
    font-weight: 500;
    color: #333;
  }
  
  .add-arrow,
  .remove-arrow {
    color: #3b82f6;
    font-weight: bold;
    font-size: 1.2rem;
  }
  
  .remove-arrow {
    color: #ef4444;
  }
  
  .empty-message {
    text-align: center;
    color: #999;
    font-style: italic;
    padding: 1rem;
  }
  
  .price {
    color: #666;
    font-size: 0.875rem;
  }
  
  .quantity {
    color: #3b82f6;
    font-size: 0.875rem;
  }
  
  .phases-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .phases-container {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }
  
  .phase-input-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 0.375rem;
    transition: all 0.2s;
  }
  
  .phase-input-row:hover {
    background: var(--bg-hover);
    border-color: var(--border-hover);
  }
  
  .drag-handle {
    cursor: grab;
    color: var(--text-secondary);
    user-select: none;
    padding: 0 0.25rem;
  }
  
  .drag-handle:active {
    cursor: grabbing;
  }
  
  .phase-name-input {
    flex: 1;
    padding: 0.375rem 0.75rem;
    font-size: 0.875rem;
    border: 1px solid var(--border-color);
    border-radius: 0.375rem;
    background: var(--bg-primary);
    color: var(--text-primary);
  }
  
  .phase-name-input:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(14, 165, 233, 0.1);
  }
  
  .remove-phase {
    padding: 0.25rem 0.5rem;
    background: none;
    border: none;
    color: var(--text-secondary);
    font-size: 1.25rem;
    cursor: pointer;
    transition: color 0.2s;
  }
  
  .remove-phase:hover {
    color: var(--danger);
  }
  
  .add-phase-btn {
    padding: 0.5rem 1rem;
    background: none;
    border: 2px dashed var(--border-color);
    border-radius: 0.375rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.875rem;
  }
  
  .add-phase-btn:hover {
    background: var(--bg-tertiary);
    border-color: var(--primary);
    color: var(--primary);
  }
  
  .phase-item,
  .phase-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem;
    background: #f9f9f9;
    border-radius: 4px;
    cursor: move;
    transition: all 0.2s;
  }
  
  .phase-item:hover {
    background: #f0f0f0;
  }
  
  .phase-item.dragging {
    opacity: 0.5;
  }
  
  .phase-order,
  .phase-number {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #3b82f6;
    color: white;
    border-radius: 50%;
    font-size: 0.75rem;
    font-weight: bold;
  }
  
  .phase-name {
    flex: 1;
  }
  
  .delete-phase-btn {
    padding: 0.25rem 0.5rem;
    background: none;
    border: none;
    color: #999;
    font-size: 1.5rem;
    cursor: pointer;
    transition: color 0.2s;
    margin-left: auto;
  }
  
  .delete-phase-btn:hover {
    color: #ef4444;
  }
  
  .add-phase-section {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
  }
  
  .phase-add-input {
    flex: 1;
    padding: 0.5rem 0.75rem;
    font-size: 0.875rem;
    border: 1px solid #e0e0e0;
    border-radius: 0.375rem;
    background: white;
    color: #333;
  }
  
  .phase-add-input:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
  
  .empty {
    text-align: center;
    color: #999;
    padding: 2rem;
  }
  
  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #999;
  }
  
  .error-message {
    background-color: #fef2f2;
    border: 1px solid #fecaca;
    color: #dc2626;
    padding: 0.75rem;
    border-radius: 6px;
    margin-bottom: 1rem;
    font-size: 0.875rem;
  }
  
  /* Modal Styles */
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #333;
  }
  
  .form-group input,
  .form-group textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    font-size: 1rem;
  }
  
  .form-section {
    margin-bottom: 2rem;
  }
  
  .form-section h3 {
    margin: 0 0 1rem 0;
    color: #333;
    font-size: 1.1rem;
  }
  
  .items-picker {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    height: 300px;
  }
  
  .picker-column {
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    background: #fafafa;
  }
  
  .picker-column h4 {
    margin: 0 0 0.75rem 0;
    color: #666;
    font-size: 0.875rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  
  .items-scroll {
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .pickable-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.75rem;
    margin: 0;
    cursor: pointer;
    transition: all 0.2s;
    border-radius: 6px;
    background: white;
    border: 1px solid #e0e0e0;
  }
  
  .pickable-item:hover {
    background: #f0f7ff;
    border-color: #3b82f6;
    transform: translateX(2px);
  }
  
  .pickable-item .item-name {
    flex: 1;
    font-weight: 500;
    color: #333;
  }
  
  .pickable-item .item-price {
    color: #666;
    font-size: 0.875rem;
    margin-right: 0.5rem;
  }
  
  .add-icon,
  .remove-icon {
    font-weight: bold;
    color: #3b82f6;
    font-size: 1.2rem;
  }
  
  .selected-item {
    display: flex;
    align-items: center;
    padding: 0.75rem;
    margin: 0;
    border-radius: 6px;
    background: white;
    border: 1px solid #e0e0e0;
    gap: 0.5rem;
  }
  
  .selected-item:hover {
    background: #f5f5f5;
  }
  
  .selected-item .remove-icon {
    cursor: pointer;
    margin-right: 0.75rem;
    color: #ef4444;
  }
  
  .selected-item .item-name {
    flex: 1;
    font-weight: 500;
    color: #333;
  }
  
  .qty-input {
    width: 60px;
    padding: 0.25rem;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    text-align: center;
  }
  
  .empty-message {
    text-align: center;
    color: #999;
    padding: 2rem;
    font-style: italic;
  }
  
  .phase-input {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  
  .phase-input input {
    flex: 1;
    padding: 0.75rem;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
  }
  
  .phase-actions {
    display: flex;
    gap: 0.25rem;
    margin-left: auto;
  }
  
  .phase-actions button {
    background: none;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    width: 28px;
    height: 28px;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .phase-actions button:hover:not(:disabled) {
    background: #f5f5f5;
    border-color: #3b82f6;
  }
  
  .phase-actions button:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }
  
  .phase-actions button.remove {
    color: #dc2626;
  }
  
  @media (max-width: 768px) {
    .layout {
      grid-template-columns: 1fr;
    }
    
    .template-list {
      max-height: 300px;
    }
    
    .items-manager {
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }
    
    .items-picker {
      grid-template-columns: 1fr;
      height: 400px;
    }
    
    .available-items,
    .template-items {
      max-height: 300px;
      overflow-y: auto;
    }
    
    .item-row {
      padding: 0.875rem;
      font-size: 0.9375rem;
    }
  }
</style>