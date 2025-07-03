<script lang="ts">
  import { itemsStore } from '../lib/stores';
  import { Card, Button, Modal, ConfirmModal } from '../lib/components';
  import type { Item } from '../lib/types/models';
  import { onMount } from 'svelte';
  
  // Load items when component mounts
  onMount(() => {
    itemsStore.load();
  });
  
  let showAddModal = false;
  let editingItem: Item | null = null;
  let deleteConfirm = {
    show: false,
    item: null as Item | null
  };
  let formError = '';
  
  let formData = {
    name: '',
    description: '',
    unit: 'each' as Item['unit'],
    unitPrice: 0,
    category: ''
  };
  
  function openAddModal() {
    formData = {
      name: '',
      description: '',
      unit: 'each',
      unitPrice: 0,
      category: ''
    };
    editingItem = null;
    formError = '';
    showAddModal = true;
  }
  
  function openEditModal(item: Item) {
    formData = {
      name: item.name,
      description: item.description || '',
      unit: item.unit,
      unitPrice: item.unitPrice,
      category: item.category || ''
    };
    editingItem = item;
    formError = '';
    showAddModal = true;
  }
  
  async function handleSave() {
    if (!formData.name || formData.unitPrice <= 0) {
      formError = 'Please fill in all required fields';
      return;
    }
    
    formError = '';
    
    try {
      if (editingItem) {
        await itemsStore.update(editingItem.id, formData);
      } else {
        await itemsStore.add(formData);
      }
      showAddModal = false;
    } catch (error) {
      formError = error instanceof Error ? error.message : 'Failed to save item';
    }
  }
  
  function confirmDelete(item: Item) {
    deleteConfirm = {
      show: true,
      item
    };
  }
  
  async function handleDelete() {
    if (deleteConfirm.item) {
      try {
        await itemsStore.remove(deleteConfirm.item.id);
        deleteConfirm = { show: false, item: null };
      } catch (error) {
        console.error('Failed to delete item:', error);
        // You might want to show an error message to the user
        deleteConfirm = { show: false, item: null };
      }
    }
  }
  
  function formatCurrency(amount: number) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  }
  
  // Subscribe to the store to get reactive updates
  $: ({ items, loading, error } = $itemsStore);
  
  // Group items by category
  $: groupedItems = items.reduce((acc, item) => {
    const category = item.category || 'Uncategorized';
    if (!acc[category]) acc[category] = [];
    acc[category].push(item);
    return acc;
  }, {} as Record<string, Item[]>);
</script>

<div class="items-page">
  <div class="page-header">
    <h1>Items Management</h1>
    <Button on:click={openAddModal}>
      Add Item
    </Button>
  </div>
  
  {#if loading}
    <div class="loading-state">
      <p>Loading items...</p>
    </div>
  {:else if error}
    <div class="error-state">
      <p>Error: {error}</p>
      <Button on:click={() => itemsStore.load()}>Try Again</Button>
    </div>
  {:else if items.length === 0}
    <div class="empty-state">
      <p>No items yet. Add your first item to get started.</p>
    </div>
  {:else}
    <div class="categories">
      {#each Object.entries(groupedItems) as [category, items]}
        <Card>
          <h2>{category}</h2>
          <div class="items-table">
            <div class="table-header">
              <span>Item Name</span>
              <span>Unit</span>
              <span>Price</span>
              <span>Actions</span>
            </div>
            {#each items as item}
              <div class="table-row">
                <div class="item-name">
                  <span class="name">{item.name}</span>
                  {#if item.description}
                    <span class="description">{item.description}</span>
                  {/if}
                </div>
                <span class="unit">{item.unit}</span>
                <span class="price">{formatCurrency(item.unitPrice)}</span>
                <div class="actions">
                  <button
                    class="action-btn edit"
                    on:click={() => openEditModal(item)}
                    title="Edit item"
                  >
                    ✏️
                  </button>
                  <button
                    class="action-btn delete"
                    on:click={() => confirmDelete(item)}
                    title="Delete item"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            {/each}
          </div>
        </Card>
      {/each}
    </div>
  {/if}
</div>

<Modal 
  bind:isOpen={showAddModal} 
  title={editingItem ? 'Edit Item' : 'Add New Item'}
>
  {#if formError}
    <div class="error-message">
      {formError}
    </div>
  {/if}
  
  <form on:submit|preventDefault={handleSave}>
    <div class="form-group">
      <label for="item-name">Item Name *</label>
      <input
        id="item-name"
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
    
    <div class="form-row">
      <div class="form-group">
        <label for="unit">Unit *</label>
        <select id="unit" bind:value={formData.unit} required>
          <option value="each">Each</option>
          <option value="ft">Feet</option>
          <option value="hr">Hours</option>
          <option value="lot">Lot</option>
        </select>
      </div>
      
      <div class="form-group">
        <label for="price">Unit Price *</label>
        <input
          id="price"
          type="number"
          bind:value={formData.unitPrice}
          min="0"
          step="0.01"
          required
        />
      </div>
    </div>
    
    <div class="form-group">
      <label for="category">Category</label>
      <input
        id="category"
        type="text"
        bind:value={formData.category}
        placeholder="e.g., Electrical, Wire, Panel"
      />
    </div>
  </form>
  
  <div slot="footer">
    <Button variant="secondary" on:click={() => showAddModal = false}>
      Cancel
    </Button>
    <Button on:click={handleSave}>
      {editingItem ? 'Update' : 'Add'} Item
    </Button>
  </div>
</Modal>

<ConfirmModal
  bind:isOpen={deleteConfirm.show}
  title="Delete Item"
  message={`Are you sure you want to delete "${deleteConfirm.item?.name}"? This action cannot be undone.`}
  confirmText="Delete"
  onConfirm={handleDelete}
  onCancel={() => deleteConfirm = { show: false, item: null }}
/>

<style>
  .items-page {
    padding: 2rem;
    max-width: 1200px;
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
  
  .empty-state,
  .loading-state,
  .error-state {
    text-align: center;
    padding: 4rem 2rem;
    color: #6b7280;
  }
  
  .error-state {
    color: #dc2626;
  }
  
  .error-state p {
    margin-bottom: 1rem;
  }
  
  .categories {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .items-table {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .table-header {
    display: grid;
    grid-template-columns: 1fr auto auto auto;
    gap: 1rem;
    padding: 0.75rem 1rem;
    background-color: #f9fafb;
    border-radius: 6px;
    font-weight: 500;
    font-size: 0.875rem;
    color: #6b7280;
  }
  
  .table-row {
    display: grid;
    grid-template-columns: 1fr auto auto auto;
    gap: 1rem;
    padding: 0.75rem 1rem;
    align-items: center;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    transition: background-color 0.2s;
  }
  
  .table-row:hover {
    background-color: #f9fafb;
  }
  
  .item-name {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .item-name .name {
    font-weight: 500;
    color: #111827;
  }
  
  .item-name .description {
    font-size: 0.875rem;
    color: #6b7280;
  }
  
  .unit {
    color: #374151;
  }
  
  .price {
    font-weight: 500;
    color: #111827;
  }
  
  .actions {
    display: flex;
    gap: 0.5rem;
  }
  
  .action-btn {
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
  
  .action-btn:hover {
    background-color: #f3f4f6;
  }
  
  .action-btn.edit:hover {
    border-color: #3b82f6;
    background-color: #eff6ff;
  }
  
  .action-btn.delete:hover {
    border-color: #ef4444;
    background-color: #fef2f2;
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
  
  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  label {
    font-size: 0.875rem;
    font-weight: 500;
    color: #374151;
  }
  
  input[type="text"],
  input[type="number"],
  select,
  textarea {
    padding: 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 1rem;
    font-family: inherit;
  }
  
  input[type="text"]:focus,
  input[type="number"]:focus,
  select:focus,
  textarea:focus {
    outline: none;
    border-color: #3b82f6;
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
</style>