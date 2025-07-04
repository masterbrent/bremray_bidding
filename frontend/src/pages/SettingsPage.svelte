<script lang="ts">
  import { onMount } from 'svelte';
  import { companySettingsStore } from '../lib/stores';
  import { Card, Button } from '../lib/components';
  import { Upload, X, Save } from 'lucide-svelte';
  import type { CompanySettings } from '../lib/stores/companySettings';
  
  let settings: CompanySettings | null = null;
  let logoFile: File | null = null;
  let logoPreview: string | null = null;
  let isSaving = false;
  let loading = true;
  let error: string | null = null;
  
  onMount(async () => {
    await companySettingsStore.load();
  });
  
  $: if ($companySettingsStore.settings) {
    settings = { ...$companySettingsStore.settings };
    logoPreview = settings.logo;
    loading = false;
  }
  
  $: loading = $companySettingsStore.loading;
  $: error = $companySettingsStore.error;
  
  function handleLogoUpload(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];
    
    if (file) {
      logoFile = file;
      const reader = new FileReader();
      reader.onload = (e) => {
        logoPreview = e.target?.result as string;
      };
      reader.readAsDataURL(file);
    }
  }
  
  function removeLogo() {
    if (!settings) return;
    logoFile = null;
    logoPreview = null;
    settings.logo = null;
  }
  
  async function saveSettings() {
    if (!settings) return;
    
    isSaving = true;
    
    try {
      // In production, upload logo to R2 if logoFile exists
      if (logoFile && logoPreview) {
        settings.logo = logoPreview; // TODO: use data URL for now, implement R2 upload
      }
      
      await companySettingsStore.updateSettings(settings);
      
      // Show success message
      alert('Settings saved successfully!');
    } catch (error) {
      console.error('Error saving settings:', error);
      alert('Error saving settings. Please try again.');
    } finally {
      isSaving = false;
    }
  }
  
</script>

<div class="settings-page">
  <div class="page-header">
    <h1>Company Settings</h1>
  </div>
  
  {#if loading}
    <div class="loading-state">
      <p>Loading settings...</p>
    </div>
  {:else if error}
    <div class="error-state">
      <p>{error}</p>
    </div>
  {:else if settings}
  <div class="settings-content">
    <Card>
      <h2>Company Information</h2>
      
      <div class="form-section">
        <div class="logo-section">
          <label>Company Logo</label>
          <div class="logo-upload">
            {#if logoPreview}
              <div class="logo-preview">
                <img src={logoPreview} alt="Company logo" />
                <button class="remove-logo" on:click={removeLogo}>
                  <X size={16} />
                </button>
              </div>
            {:else}
              <label for="logo-input" class="upload-area">
                <Upload size={24} />
                <span>Upload Logo</span>
                <input
                  id="logo-input"
                  type="file"
                  accept="image/*"
                  on:change={handleLogoUpload}
                  hidden
                />
              </label>
            {/if}
          </div>
        </div>
        
        <div class="form-group">
          <label for="company-name">Company Name</label>
          <input
            id="company-name"
            type="text"
            bind:value={settings.name}
            disabled={!settings}
            placeholder="Enter company name"
          />
        </div>
        
        <div class="form-group">
          <label for="license">License Number</label>
          <input
            id="license"
            type="text"
            bind:value={settings.license}
            disabled={!settings}
            placeholder="Enter license number"
          />
        </div>
      </div>
    </Card>
    
    <Card>
      <h2>Contact Information</h2>
      
      <div class="form-section">
        <div class="form-group">
          <label for="address">Street Address</label>
          <input
            id="address"
            type="text"
            bind:value={settings.address}
            disabled={!settings}
            placeholder="Enter street address"
          />
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label for="city">City</label>
            <input
              id="city"
              type="text"
              bind:value={settings.city}
            disabled={!settings}
              placeholder="Enter city"
            />
          </div>
          
          <div class="form-group">
            <label for="state">State</label>
            <input
              id="state"
              type="text"
              bind:value={settings.state}
            disabled={!settings}
              placeholder="State"
              maxlength="2"
            />
          </div>
          
          <div class="form-group">
            <label for="zip">ZIP Code</label>
            <input
              id="zip"
              type="text"
              bind:value={settings.zip}
            disabled={!settings}
              placeholder="ZIP"
              maxlength="10"
            />
          </div>
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label for="phone">Phone Number</label>
            <input
              id="phone"
              type="tel"
              bind:value={settings.phone}
            disabled={!settings}
              placeholder="(555) 123-4567"
            />
          </div>
          
          <div class="form-group">
            <label for="email">Email Address</label>
            <input
              id="email"
              type="email"
              bind:value={settings.email}
            disabled={!settings}
              placeholder="company@example.com"
            />
          </div>
        </div>
        
        <div class="form-group">
          <label for="website">Website</label>
          <input
            id="website"
            type="url"
            bind:value={settings.website}
            disabled={!settings}
            placeholder="https://www.example.com"
          />
        </div>
      </div>
    </Card>
    
    <div class="actions">
      <Button
        variant="primary"
        on:click={saveSettings}
        disabled={isSaving || !settings}
      >
        <Save size={16} />
        {isSaving ? 'Saving...' : 'Save Settings'}
      </Button>
    </div>
  </div>
  {/if}
</div>

<style>
  .settings-page {
    padding: 2rem;
    max-width: 800px;
    margin: 0 auto;
  }
  
  .page-header {
    margin-bottom: 2rem;
  }
  
  .page-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: var(--gray-900);
    margin: 0;
  }
  
  .settings-content {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
  
  h2 {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--gray-900);
    margin: 0 0 1.5rem 0;
  }
  
  .form-section {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .logo-section {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .logo-upload {
    width: 200px;
  }
  
  .upload-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    width: 200px;
    height: 200px;
    border: 2px dashed var(--gray-300);
    border-radius: 12px;
    background: var(--gray-50);
    color: var(--gray-600);
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .upload-area:hover {
    border-color: var(--primary-400);
    background: var(--primary-50);
    color: var(--primary-600);
  }
  
  .logo-preview {
    position: relative;
    width: 200px;
    height: 200px;
    border-radius: 12px;
    overflow: hidden;
    background: var(--gray-50);
  }
  
  .logo-preview img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
  
  .remove-logo {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    width: 32px;
    height: 32px;
    background: rgba(255, 255, 255, 0.9);
    border: none;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.2s;
  }
  
  .remove-logo:hover {
    background: white;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1;
  }
  
  .form-row {
    display: flex;
    gap: 1rem;
  }
  
  label {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--gray-700);
  }
  
  input {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid var(--gray-300);
    border-radius: 8px;
    font-size: 0.9375rem;
    transition: all 0.2s;
    background: white;
  }
  
  input:focus {
    outline: none;
    border-color: var(--primary-500);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
  
  input::placeholder {
    color: var(--gray-400);
  }
  
  .actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 1rem;
  }
  
  .loading-state,
  .error-state {
    text-align: center;
    padding: 4rem 2rem;
    color: #6b7280;
  }
  
  .error-state {
    color: #ef4444;
  }
  
  @media (max-width: 640px) {
    .settings-page {
      padding: 1rem;
    }
    
    .form-row {
      flex-direction: column;
    }
  }
</style>