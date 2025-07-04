<script lang="ts">
  import { onMount } from 'svelte';
  import { userStore, permissions, effectiveRole } from '../lib/stores';
  import { Shield, Users, Mail, ChevronDown } from 'lucide-svelte';
  
  // Mock user data - in production this would come from your backend
  let users = [
    { email: 'brenthall@gmail.com', role: 'admin', name: 'Brent Hall' },
    { email: 'tech1@example.com', role: 'tech', name: 'John Smith' },
    { email: 'tech2@example.com', role: 'tech', name: 'Jane Doe' },
  ];
  
  let newUserEmail = '';
  let newUserName = '';
  let newUserRole: 'admin' | 'tech' = 'tech';
  let showAddForm = false;
  
  function updateUserRole(email: string, role: 'admin' | 'tech') {
    users = users.map(u => 
      u.email === email ? { ...u, role } : u
    );
    // In production, save to backend
    userStore.setUserRole(email, role);
  }
  
  function addUser() {
    if (newUserEmail && newUserName) {
      users = [...users, {
        email: newUserEmail,
        role: newUserRole,
        name: newUserName
      }];
      // Reset form
      newUserEmail = '';
      newUserName = '';
      newUserRole = 'tech';
      showAddForm = false;
    }
  }
  
  function removeUser(email: string) {
    if (email === 'brenthall@gmail.com') {
      alert('Cannot remove master admin');
      return;
    }
    users = users.filter(u => u.email !== email);
  }
</script>

{#if permissions.canEditJobs($effectiveRole)}
  <div class="users-page">
    <div class="page-header">
      <h1>User Management</h1>
      <button class="add-user-btn" on:click={() => showAddForm = !showAddForm}>
        <Users size={18} />
        <span>Add User</span>
      </button>
    </div>
    
    {#if showAddForm}
      <div class="add-form">
        <h3>Add New User</h3>
        <div class="form-grid">
          <input
            type="text"
            placeholder="Name"
            bind:value={newUserName}
            class="form-input"
          />
          <input
            type="email"
            placeholder="Email"
            bind:value={newUserEmail}
            class="form-input"
          />
          <select bind:value={newUserRole} class="form-select">
            <option value="tech">Tech</option>
            <option value="admin">Admin</option>
          </select>
          <div class="form-actions">
            <button class="btn-secondary" on:click={() => showAddForm = false}>
              Cancel
            </button>
            <button class="btn-primary" on:click={addUser}>
              Add User
            </button>
          </div>
        </div>
      </div>
    {/if}
    
    <div class="users-list">
      {#each users as user}
        <div class="user-card">
          <div class="user-info">
            <div class="user-avatar">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div class="user-details">
              <h3>{user.name}</h3>
              <span class="user-email">
                <Mail size={14} />
                {user.email}
              </span>
            </div>
          </div>
          
          <div class="user-actions">
            <div class="role-selector">
              <select 
                value={user.role} 
                on:change={(e) => updateUserRole(user.email, e.target.value)}
                class="role-select"
                disabled={user.email === 'brenthall@gmail.com'}
              >
                <option value="admin">Admin</option>
                <option value="tech">Tech</option>
              </select>
            </div>
            
            {#if user.email !== 'brenthall@gmail.com'}
              <button 
                class="remove-btn"
                on:click={() => removeUser(user.email)}
              >
                Remove
              </button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>
{:else}
  <div class="access-denied">
    <Shield size={48} />
    <h2>Access Denied</h2>
    <p>You don't have permission to manage users.</p>
  </div>
{/if}

<style>
  .users-page {
    padding: 2rem 1.5rem;
    max-width: 1000px;
    margin: 0 auto;
  }
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }
  
  .page-header h1 {
    font-size: 2.25rem;
    font-weight: 800;
    color: #0a0a0a;
    margin: 0;
    letter-spacing: -0.03em;
  }
  
  .add-user-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    background: #0a0a0a;
    color: white;
    border: none;
    border-radius: 10px;
    font-size: 0.9375rem;
    font-weight: 500;
    letter-spacing: -0.01em;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .add-user-btn:hover {
    background: #171717;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
  
  .add-form {
    background: white;
    border: 1px solid #e5e5e5;
    border-radius: 12px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }
  
  .add-form h3 {
    margin: 0 0 1rem 0;
    font-size: 1.125rem;
    font-weight: 600;
    color: #0a0a0a;
  }
  
  .form-grid {
    display: grid;
    gap: 1rem;
  }
  
  .form-input,
  .form-select {
    padding: 0.75rem 1rem;
    border: 1px solid #e5e5e5;
    border-radius: 8px;
    font-size: 0.9375rem;
    background: #fafafa;
    transition: all 0.2s ease;
  }
  
  .form-input:focus,
  .form-select:focus {
    border-color: #d4d4d4;
    background: white;
    outline: none;
    box-shadow: 0 0 0 4px rgba(0, 0, 0, 0.02);
  }
  
  .form-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
  }
  
  .btn-primary,
  .btn-secondary {
    padding: 0.625rem 1.25rem;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    border: none;
  }
  
  .btn-primary {
    background: #0a0a0a;
    color: white;
  }
  
  .btn-primary:hover {
    background: #171717;
  }
  
  .btn-secondary {
    background: #f5f5f5;
    color: #737373;
  }
  
  .btn-secondary:hover {
    background: #e5e5e5;
    color: #404040;
  }
  
  .users-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .user-card {
    background: white;
    border: 1px solid #e5e5e5;
    border-radius: 12px;
    padding: 1.25rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: all 0.2s ease;
  }
  
  .user-card:hover {
    border-color: #d4d4d4;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .user-avatar {
    width: 3rem;
    height: 3rem;
    border-radius: 50%;
    background: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 1.125rem;
    color: #737373;
  }
  
  .user-details h3 {
    margin: 0 0 0.25rem 0;
    font-size: 1rem;
    font-weight: 600;
    color: #0a0a0a;
    letter-spacing: -0.01em;
  }
  
  .user-email {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    color: #737373;
    font-size: 0.875rem;
  }
  
  .user-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .role-select {
    padding: 0.5rem 2.5rem 0.5rem 1rem;
    border: 1px solid #e5e5e5;
    border-radius: 8px;
    background: #fafafa;
    font-size: 0.875rem;
    font-weight: 500;
    color: #0a0a0a;
    cursor: pointer;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg width='12' height='8' viewBox='0 0 12 8' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1.5L6 6.5L11 1.5' stroke='%23737373' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    transition: all 0.2s ease;
  }
  
  .role-select:hover:not(:disabled) {
    border-color: #d4d4d4;
    background-color: white;
  }
  
  .role-select:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .remove-btn {
    padding: 0.5rem 1rem;
    border: 1px solid #ef4444;
    background: white;
    color: #ef4444;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .remove-btn:hover {
    background: #ef4444;
    color: white;
  }
  
  .access-denied {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 60vh;
    text-align: center;
    color: #737373;
  }
  
  .access-denied h2 {
    margin: 1rem 0 0.5rem 0;
    color: #0a0a0a;
  }
  
  @media (max-width: 768px) {
    .users-page {
      padding: 1.25rem 1rem;
    }
    
    .page-header {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }
    
    .user-card {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }
    
    .user-actions {
      justify-content: space-between;
    }
  }
</style>
