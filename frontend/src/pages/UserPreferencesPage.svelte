<script lang="ts">
  import { userStore, effectiveRole } from '../lib/stores';
  import { Card, Button } from '../lib/components';
  import { User, Save, Key, Moon, Sun, Bell } from 'lucide-svelte';
  
  let name = '';
  let email = '';
  let currentPassword = '';
  let newPassword = '';
  let confirmPassword = '';
  let theme = 'light';
  let notifications = true;
  let isSaving = false;
  let error: string | null = null;
  let successMessage: string | null = null;
  
  $: if ($userStore.user) {
    name = $userStore.user.name || '';
    email = $userStore.user.email || '';
  }
  
  async function saveProfile() {
    error = null;
    successMessage = null;
    isSaving = true;
    
    try {
      // TODO: Implement API call to update user profile
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulated delay
      successMessage = 'Profile updated successfully!';
    } catch (err) {
      error = 'Failed to update profile';
    } finally {
      isSaving = false;
    }
  }
  
  async function changePassword() {
    error = null;
    successMessage = null;
    
    if (!currentPassword || !newPassword) {
      error = 'Please fill in all password fields';
      return;
    }
    
    if (newPassword !== confirmPassword) {
      error = 'New passwords do not match';
      return;
    }
    
    if (newPassword.length < 8) {
      error = 'Password must be at least 8 characters';
      return;
    }
    
    isSaving = true;
    
    try {
      // TODO: Implement API call to change password
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulated delay
      successMessage = 'Password changed successfully!';
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
    } catch (err) {
      error = 'Failed to change password';
    } finally {
      isSaving = false;
    }
  }
</script>

<div class="preferences-page">
  <h1>User Preferences</h1>
  
  <div class="preferences-content">
    <!-- Profile Section -->
    <Card>
      <div class="section">
        <div class="section-header">
          <User size={20} />
          <h2>Profile Information</h2>
        </div>
        
        <div class="form-group">
          <label for="name">Name</label>
          <input
            id="name"
            type="text"
            bind:value={name}
            placeholder="Your name"
          />
        </div>
        
        <div class="form-group">
          <label for="email">Email</label>
          <input
            id="email"
            type="email"
            bind:value={email}
            placeholder="your@email.com"
            disabled
          />
          <p class="help-text">Contact your administrator to change your email</p>
        </div>
        
        <div class="form-group">
          <label>Role</label>
          <div class="role-badge">{$effectiveRole}</div>
        </div>
        
        <Button on:click={saveProfile} disabled={isSaving}>
          <Save size={16} />
          Save Profile
        </Button>
      </div>
    </Card>
    
    <!-- Password Section -->
    <Card>
      <div class="section">
        <div class="section-header">
          <Key size={20} />
          <h2>Change Password</h2>
        </div>
        
        <div class="form-group">
          <label for="current-password">Current Password</label>
          <input
            id="current-password"
            type="password"
            bind:value={currentPassword}
            placeholder="Enter current password"
          />
        </div>
        
        <div class="form-group">
          <label for="new-password">New Password</label>
          <input
            id="new-password"
            type="password"
            bind:value={newPassword}
            placeholder="Enter new password"
          />
        </div>
        
        <div class="form-group">
          <label for="confirm-password">Confirm New Password</label>
          <input
            id="confirm-password"
            type="password"
            bind:value={confirmPassword}
            placeholder="Confirm new password"
          />
        </div>
        
        <Button on:click={changePassword} disabled={isSaving}>
          <Key size={16} />
          Change Password
        </Button>
      </div>
    </Card>
    
    <!-- Preferences Section -->
    <Card>
      <div class="section">
        <div class="section-header">
          <Bell size={20} />
          <h2>Preferences</h2>
        </div>
        
        <div class="preference-item">
          <div class="preference-info">
            <h3>Theme</h3>
            <p>Choose your preferred color theme</p>
          </div>
          <div class="theme-toggle">
            <button 
              class="theme-option"
              class:active={theme === 'light'}
              on:click={() => theme = 'light'}
            >
              <Sun size={16} />
              Light
            </button>
            <button 
              class="theme-option"
              class:active={theme === 'dark'}
              on:click={() => theme = 'dark'}
            >
              <Moon size={16} />
              Dark
            </button>
          </div>
        </div>
        
        <div class="preference-item">
          <div class="preference-info">
            <h3>Notifications</h3>
            <p>Receive updates about job changes</p>
          </div>
          <label class="toggle">
            <input type="checkbox" bind:checked={notifications} />
            <span class="toggle-slider"></span>
          </label>
        </div>
      </div>
    </Card>
    
    {#if error}
      <div class="message error">{error}</div>
    {/if}
    
    {#if successMessage}
      <div class="message success">{successMessage}</div>
    {/if}
  </div>
</div>

<style>
  .preferences-page {
    padding: 2rem;
    max-width: 800px;
    margin: 0 auto;
  }
  
  h1 {
    font-size: 1.75rem;
    font-weight: 700;
    margin-bottom: 2rem;
    color: var(--text-primary);
  }
  
  .preferences-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .section {
    padding: 1.5rem;
  }
  
  .section-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }
  
  .section-header h2 {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0;
    color: var(--text-primary);
  }
  
  .form-group {
    margin-bottom: 1.25rem;
  }
  
  label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-secondary);
    margin-bottom: 0.5rem;
  }
  
  input[type="text"],
  input[type="email"],
  input[type="password"] {
    width: 100%;
    padding: 0.625rem 0.875rem;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    transition: all 0.2s ease;
  }
  
  input[type="text"]:focus,
  input[type="email"]:focus,
  input[type="password"]:focus {
    outline: none;
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(91, 91, 214, 0.1);
  }
  
  input:disabled {
    background: var(--gray-50);
    cursor: not-allowed;
  }
  
  .help-text {
    font-size: 0.75rem;
    color: var(--text-tertiary);
    margin-top: 0.25rem;
  }
  
  .role-badge {
    display: inline-block;
    padding: 0.375rem 0.75rem;
    background: var(--primary-50);
    color: var(--primary-600);
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    font-weight: 500;
    text-transform: capitalize;
  }
  
  .preference-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 0;
    border-bottom: 1px solid var(--border-light);
  }
  
  .preference-item:last-child {
    border-bottom: none;
  }
  
  .preference-info h3 {
    font-size: 0.9375rem;
    font-weight: 500;
    margin: 0 0 0.25rem 0;
    color: var(--text-primary);
  }
  
  .preference-info p {
    font-size: 0.8125rem;
    color: var(--text-tertiary);
    margin: 0;
  }
  
  .theme-toggle {
    display: flex;
    gap: 0.5rem;
  }
  
  .theme-option {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.5rem 0.875rem;
    border: 1px solid var(--border-color);
    background: white;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .theme-option:hover {
    background: var(--gray-50);
  }
  
  .theme-option.active {
    background: var(--primary-500);
    color: white;
    border-color: var(--primary-500);
  }
  
  /* Toggle Switch */
  .toggle {
    position: relative;
    display: inline-block;
    width: 48px;
    height: 24px;
  }
  
  .toggle input {
    opacity: 0;
    width: 0;
    height: 0;
  }
  
  .toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #ccc;
    transition: 0.3s;
    border-radius: 24px;
  }
  
  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 16px;
    width: 16px;
    left: 4px;
    bottom: 4px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
  }
  
  input:checked + .toggle-slider {
    background-color: var(--primary-500);
  }
  
  input:checked + .toggle-slider:before {
    transform: translateX(24px);
  }
  
  .message {
    padding: 0.875rem;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    margin-top: 1rem;
  }
  
  .message.error {
    background: rgba(239, 68, 68, 0.1);
    color: var(--error-600);
  }
  
  .message.success {
    background: rgba(16, 185, 129, 0.1);
    color: var(--success-600);
  }
  
  @media (max-width: 768px) {
    .preferences-page {
      padding: 1rem;
    }
    
    h1 {
      font-size: 1.5rem;
    }
    
    .section {
      padding: 1rem;
    }
    
    .preference-item {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.75rem;
    }
  }
</style>
