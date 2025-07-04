<script lang="ts">
  import Modal from './Modal.svelte';
  import Button from './Button.svelte';
  import type { Job } from '../types/models';
  import { Printer, Download, Send } from 'lucide-svelte';
  
  export let isOpen: boolean = false;
  export let job: Job | null = null;
  
  // Use the totalAmount from the job or calculate from items
  $: subtotal = job ? (job.totalAmount / 1.08) : 0; // Assuming 8% tax is included
  $: tax = job ? (job.totalAmount - subtotal) : 0;
  $: total = job?.totalAmount || 0;
  
  function formatCurrency(amount: number): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  }
  
  function formatDate(date?: Date): string {
    if (!date) return 'N/A';
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    }).format(new Date(date));
  }
</script>

<Modal bind:isOpen title="Invoice Preview" size="large">
  {#if job}
    <div class="invoice">
      <!-- Invoice Header -->
      <div class="invoice-header">
        <div class="company-info">
          <h2>Bremray Electrical</h2>
          <p>123 Electric Avenue</p>
          <p>Anytown, ST 12345</p>
          <p>Phone: (555) 123-4567</p>
          <p>License #: ELC123456</p>
        </div>
        
        <div class="invoice-info">
          <h1>INVOICE</h1>
          <div class="invoice-meta">
            <div>
              <strong>Invoice #:</strong> 
              <span>{job.waveInvoiceId || `INV-${job.id.slice(-6).toUpperCase()}`}</span>
            </div>
            <div>
              <strong>Date:</strong> 
              <span>{formatDate(new Date())}</span>
            </div>
            <div>
              <strong>Due Date:</strong> 
              <span>{formatDate(new Date(Date.now() + 30 * 24 * 60 * 60 * 1000))}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Bill To -->
      <div class="bill-to">
        <h3>Bill To:</h3>
        <p><strong>Customer #{job.customerId}</strong></p>
        <p>{job.address}</p>
      </div>
      
      <!-- Job Details -->
      <div class="job-details">
        <h3>Job Details</h3>
        <div class="detail-row">
          <span>Job ID:</span>
          <span>{job.id.slice(-8).toUpperCase()}</span>
        </div>
        <div class="detail-row">
          <span>Scheduled Date:</span>
          <span>{formatDate(job.scheduledDate)}</span>
        </div>
        <div class="detail-row">
          <span>Start Date:</span>
          <span>{formatDate(job.startDate)}</span>
        </div>
        <div class="detail-row">
          <span>Completion Date:</span>
          <span>{formatDate(job.endDate)}</span>
        </div>
        {#if job.permitRequired}
          <div class="detail-row">
            <span>Permit:</span>
            <span>{job.permitNumber || 'Required'}</span>
          </div>
        {/if}
      </div>
      
      <!-- Line Items -->
      <table class="line-items">
        <thead>
          <tr>
            <th>Description</th>
            <th>Quantity</th>
            <th>Rate</th>
            <th>Amount</th>
          </tr>
        </thead>
        <tbody>
          {#each job.items as item}
            <tr>
              <td>{item.name}</td>
              <td>{item.quantity}</td>
              <td>{formatCurrency(item.price)}</td>
              <td>{formatCurrency(item.total)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
      
      <!-- Totals -->
      <div class="totals">
        <div class="total-row">
          <span>Subtotal:</span>
          <span>{formatCurrency(subtotal)}</span>
        </div>
        <div class="total-row">
          <span>Tax (8%):</span>
          <span>{formatCurrency(tax)}</span>
        </div>
        <div class="total-row total">
          <span>Total:</span>
          <span>{formatCurrency(total)}</span>
        </div>
      </div>
      
      <!-- Terms -->
      <div class="terms">
        <h4>Terms & Conditions</h4>
        <p>Payment is due within 30 days. Please include invoice number with payment.</p>
        <p>Thank you for your business!</p>
      </div>
    </div>
  {/if}
  
  <div slot="footer">
    <Button variant="ghost">
      <Printer size={20} />
      Print
    </Button>
    <Button variant="ghost">
      <Download size={20} />
      Download PDF
    </Button>
    <Button variant="primary" on:click={() => isOpen = false}>
      Close
    </Button>
  </div>
</Modal>

<style>
  .invoice {
    font-size: 0.9375rem;
    line-height: 1.6;
    color: var(--gray-900);
  }
  
  .invoice-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 2px solid var(--gray-200);
  }
  
  .company-info h2 {
    font-size: 1.5rem;
    margin: 0 0 0.5rem 0;
    color: var(--primary-700);
  }
  
  .company-info p {
    margin: 0.25rem 0;
    color: var(--gray-600);
  }
  
  .invoice-info h1 {
    font-size: 2.5rem;
    margin: 0;
    text-align: right;
    color: var(--gray-900);
  }
  
  .invoice-meta {
    margin-top: 1rem;
    text-align: right;
  }
  
  .invoice-meta div {
    margin: 0.25rem 0;
  }
  
  .invoice-meta strong {
    color: var(--gray-700);
  }
  
  .bill-to {
    margin-bottom: 2rem;
  }
  
  .bill-to h3 {
    margin: 0 0 0.75rem 0;
    color: var(--gray-700);
    font-size: 1.125rem;
  }
  
  .bill-to p {
    margin: 0.25rem 0;
  }
  
  .job-details {
    background: var(--gray-50);
    padding: 1.5rem;
    border-radius: 12px;
    margin-bottom: 2rem;
  }
  
  .job-details h3 {
    margin: 0 0 1rem 0;
    color: var(--gray-700);
    font-size: 1.125rem;
  }
  
  .detail-row {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--gray-200);
  }
  
  .detail-row:last-child {
    border-bottom: none;
  }
  
  .line-items {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 2rem;
  }
  
  .line-items th {
    background: var(--gray-100);
    padding: 0.75rem;
    text-align: left;
    font-weight: 600;
    color: var(--gray-700);
    border-bottom: 2px solid var(--gray-300);
  }
  
  .line-items th:last-child,
  .line-items td:last-child {
    text-align: right;
  }
  
  .line-items td {
    padding: 0.75rem;
    border-bottom: 1px solid var(--gray-200);
  }
  
  .totals {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    margin-bottom: 2rem;
  }
  
  .total-row {
    display: flex;
    gap: 3rem;
    padding: 0.5rem 0;
    min-width: 250px;
    justify-content: space-between;
  }
  
  .total-row.total {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--gray-900);
    border-top: 2px solid var(--gray-300);
    padding-top: 1rem;
    margin-top: 0.5rem;
  }
  
  .terms {
    background: var(--gray-50);
    padding: 1.5rem;
    border-radius: 12px;
  }
  
  .terms h4 {
    margin: 0 0 0.75rem 0;
    color: var(--gray-700);
  }
  
  .terms p {
    margin: 0.5rem 0;
    color: var(--gray-600);
  }
</style>