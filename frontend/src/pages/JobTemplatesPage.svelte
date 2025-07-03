<script lang="ts">
  import { jobTemplatesStore, itemsStore } from '../lib/stores';
  import { Card, Button, Modal, QuantityPicker, ConfirmModal } from '../lib/components';
  import type { JobTemplate, JobPhase } from '../lib/types/models';
  
  let showCreateModal = false;
  let editingTemplate: JobTemplate | null = null;
  let selectedTemplateId: string | null = null;
  let deleteConfirm = {
    show: false,
    template: null as JobTemplate | null
  };
  let formError = '';
  
  let formData = {
    name: '',
    description: '',
    phases: [] as { name: string; order: number }[]
  };
  
  let newPhaseName = '';
  
  function openCreateModal() {
    formData = {
      name: '',
      description: '',
      phases: []
    };
    editingTemplate = null;
    formError = '';
    showCreateModal = true;
  }
  
  function addPhase() {
    if (newPhaseName.trim()) {
      formData.phases = [...formData.phases, {
        name: newPhaseName.trim(),
        order: formData.phases.length + 1
      }];
      newPhaseName = '';
    }
  }
  
  function removePhase(index: number) {
    formData.phases = formData.phases.filter((_, i) => i !== index);
    // Reorder phases
    formData.phases = formData.phases.map((phase, i) => ({
      ...phase,
      order: i + 1
    }));
  }
  
  function handleSave() {
    if (!formData.name) {
      formError = 'Please enter a template name';
      return;
    }
    
    formError = '';
    
    const phases: JobPhase[] = formData.phases.map(phase => ({
      id: crypto.randomUUID(),
      name: phase.name,
      order: phase.order,
      isCompleted: false
    }));
    
    jobTemplatesStore.add({
      name: formData.name,
      description: formData.description,
      items: [],
      phases
    });
    
    showCreateModal = false;
  }
  
  function confirmDelete(template: JobTemplate) {
    deleteConfirm = {
      show: true,
      template
    };
  }
  
  function handleDelete() {
    if (deleteConfirm.template) {
      const id = deleteConfirm.template.id;
      jobTemplatesStore.remove(id);
      if (selectedTemplateId === id) {
        selectedTemplateId = null;
      }
      deleteConfirm = { show: false, template: null };
    }
  }
  
  function toggleItemInTemplate(templateId: string, itemId: string) {
    const template = $jobTemplatesStore.find(t => t.id === templateId);
    if (!template) return;
    
    const hasItem = template.items.some(i => i.itemId === itemId);
    if (hasItem) {
      jobTemplatesStore.removeItemFromTemplate(templateId, itemId);
    } else {
      jobTemplatesStore.addItemToTemplate(templateId, itemId);
    }
  }
  
  function updateItemQuantity(templateId: string, itemId: string, quantity: number) {
    jobTemplatesStore.update(templateId, {
      items: $jobTemplatesStore
        .find(t => t.id === templateId)!
        .items.map(item => 
          item.itemId === itemId 
            ? { ...item, defaultQuantity: quantity }
            : item
        )
    });
  }
  
  $: selectedTemplate = selectedTemplateId 
    ? $jobTemplatesStore.find(t => t.id === selectedTemplateId)
    : null;
  
  $: templateItems = selectedTemplate
    ? $itemsStore.map(item => ({
        item,
        isInTemplate: selectedTemplate.items.some(ti => ti.itemId === item.id),
        quantity: selectedTemplate.items.find(ti => ti.itemId === item.id)?.defaultQuantity || 1
      }))
    : [];
</script>

<div class="templates-page">
  <div class="page-header">
    <h1>Job Templates</h1>
    <Button on:click={openCreateModal}>
      Create Template
    </Button>
  </div>
  
  <div class="layout">
    <div class="templates-section">
      <Card>
        <h2>Templates</h2>
        {#if $jobTemplatesStore.length === 0}
          <p class="empty">No templates yet</p>
        {:else}
          <div class="templates-list">
            {#each $jobTemplatesStore as template}
              <div
                class="template-item"
                class:selected={selectedTemplateId === template.id}
                on:click={() => selectedTemplateId = template.id}
                on:keydown={(e) => e.key === 'Enter' && (selectedTemplateId = template.id)}
                role="button"
                tabindex="0"
              >
                <div class="template-info">
                  <div class="template-name">{template.name}</div>
                  {#if template.description}
                    <div class="template-desc">{template.description}</div>
                  {/if}
                  <div class="template-meta">
                    {template.items.length} items • {template.phases.length} phases
                  </div>
                </div>
                <button
                  class="delete-btn"
                  on:click|stopPropagation={() => confirmDelete(template)}
                  title="Delete template"
                >
                  🗑️
                </button>
              </div>
            {/each}
          </div>
        {/if}
      </Card>
    </div>
    
    <div class="details-section">
      {#if selectedTemplate}
        <Card>
          <h2>{selectedTemplate.name} - Items</h2>
          <div class="items-selector">
            {#each templateItems as { item, isInTemplate, quantity }}
              <div class="item-selector-row" class:selected={isInTemplate}>
                <label class="item-checkbox">
                  <input
                    type="checkbox"
                    checked={isInTemplate}
                    on:change={() => toggleItemInTemplate(selectedTemplate.id, item.id)}
                  />
                  <span>{item.name}</span>
                  <span class="item-price">${item.unitPrice}/{item.unit}</span>
                </label>
                {#if isInTemplate}
                  <QuantityPicker
                    value={quantity}
                    min={1}
                    on:change={(e) => updateItemQuantity(selectedTemplate.id, item.id, e.detail)}
                  />
                {/if}
              </div>
            {/each}
          </div>
        </Card>
        
        <Card>
          <h2>{selectedTemplate.name} - Phases</h2>
          {#if selectedTemplate.phases.length === 0}
            <p class="empty">No phases defined</p>
          {:else}
            <div class="phases-list">
              {#each selectedTemplate.phases.sort((a, b) => a.order - b.order) as phase}
                <div class="phase-item">
                  <span class="phase-order">{phase.order}</span>
                  <span class="phase-name">{phase.name}</span>
                </div>
              {/each}
            </div>
          {/if}
        </Card>
      {:else}
        <Card>
          <p class="empty">Select a template to view details</p>
        </Card>
      {/if}
    </div>
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
    
    <div class="form-group">
      <label>Job Phases</label>
      <div class="phase-input">
        <input
          type="text"
          bind:value={newPhaseName}
          placeholder="Phase name"
          on:keydown={(e) => e.key === 'Enter' && (e.preventDefault(), addPhase())}
        />
        <Button size="small" on:click={addPhase}>Add</Button>
      </div>
      
      {#if formData.phases.length > 0}
        <div class="phases-preview">
          {#each formData.phases as phase, index}
            <div class="phase-preview-item">
              <span>{phase.order}. {phase.name}</span>
              <button
                type="button"
                on:click={() => removePhase(index)}
                class="remove-phase"
              >
                ×
              </button>
            </div>
          {/each}
        </div>
      {/if}
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
  
  h1 {
    margin: 0;
    font-size: 2rem;
    color: #111827;
  }
  
  h2 {
    margin: 0 0 1rem 0;
    font-size: 1.25rem;
    color: #374151;
  }
  
  .layout {
    display: grid;
    grid-template-columns: 400px 1fr;
    gap: 1.5rem;
  }
  
  .empty {
    color: #6b7280;
    text-align: center;
    padding: 2rem;
  }
  
  .templates-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .template-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    cursor: pointer;
    text-align: left;
    transition: all 0.2s;
    width: 100%;
    font-family: inherit;
  }
  
  .template-item:hover {
    background-color: #f9fafb;
  }
  
  .template-item.selected {
    border-color: #3b82f6;
    background-color: #eff6ff;
  }
  
  .template-info {
    flex: 1;
  }
  
  .template-name {
    font-weight: 500;
    color: #111827;
    margin-bottom: 0.25rem;
  }
  
  .template-desc {
    font-size: 0.875rem;
    color: #6b7280;
    margin-bottom: 0.25rem;
  }
  
  .template-meta {
    font-size: 0.75rem;
    color: #9ca3af;
  }
  
  .delete-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    border: 1px solid #e5e7eb;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .delete-btn:hover {
    border-color: #ef4444;
    background-color: #fef2f2;
  }
  
  .details-section {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .items-selector {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .item-selector-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    transition: all 0.2s;
  }
  
  .item-selector-row.selected {
    background-color: #f0f9ff;
    border-color: #3b82f6;
  }
  
  .item-checkbox {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
    cursor: pointer;
  }
  
  .item-checkbox input {
    width: 1.25rem;
    height: 1.25rem;
  }
  
  .item-price {
    margin-left: auto;
    color: #6b7280;
    font-size: 0.875rem;
  }
  
  .phases-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .phase-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem;
    background-color: #f9fafb;
    border-radius: 6px;
  }
  
  .phase-order {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #3b82f6;
    color: white;
    border-radius: 50%;
    font-weight: 500;
    font-size: 0.875rem;
  }
  
  .phase-name {
    color: #111827;
    font-weight: 500;
  }
  
  /* Form styles */
  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  label {
    font-size: 0.875rem;
    font-weight: 500;
    color: #374151;
  }
  
  input[type="text"],
  textarea {
    padding: 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 1rem;
    font-family: inherit;
  }
  
  input[type="text"]:focus,
  textarea:focus {
    outline: none;
    border-color: #3b82f6;
  }
  
  .phase-input {
    display: flex;
    gap: 0.5rem;
  }
  
  .phase-input input {
    flex: 1;
  }
  
  .phases-preview {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }
  
  .phase-preview-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0.75rem;
    background-color: #f9fafb;
    border-radius: 4px;
  }
  
  .remove-phase {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    border: none;
    color: #ef4444;
    cursor: pointer;
    font-size: 1.25rem;
    line-height: 1;
  }
  
  .remove-phase:hover {
    background-color: #fef2f2;
    border-radius: 4px;
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
  
  @media (max-width: 768px) {
    .layout {
      grid-template-columns: 1fr;
    }
  }
</style>