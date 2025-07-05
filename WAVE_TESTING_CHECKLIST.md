# Wave Integration Testing Checklist

**CRITICAL**: Test each item in order. If any test fails, DO NOT proceed until fixed.

## Pre-Test Setup

- [ ] Set environment variables in `.env`:
  ```
  WAVE_TOKEN=your_wave_api_token
  WAVE_BUSINESS_ID=your_wave_business_id
  SKYVIEW_CUSTOMER_ID=customer_id_from_wave (optional)
  ```

- [ ] Ensure "Skyview" customer exists in Wave (exact spelling)

- [ ] Ensure these EXACT products exist in Wave OR will be auto-created:
  - "Custom Work"
  - "Electrical Permit"
  - Any template item names from your items table

## Test Sequence

### 1. Test Wave Connection
- [ ] Call test endpoint to verify credentials work
- [ ] Verify it can fetch products
- [ ] Verify it finds "Skyview" customer

### 2. Test Product Finding/Creation
- [ ] Test finding existing product (case-insensitive)
- [ ] Test creating new product when not found
- [ ] Verify new product has income account assigned

### 3. Test Simple Invoice (Template Items Only)
Create test job with:
- [ ] One template item with completed quantity > 0
- [ ] Verify invoice creates successfully
- [ ] Verify invoice shows in Wave dashboard
- [ ] Verify PO number format: "CustomerName - City"

### 4. Test Custom Items
Create test job with:
- [ ] Custom item with name and description
- [ ] Verify product name is "Custom Work"
- [ ] Verify description appears under product name
- [ ] Verify price calculates correctly

### 5. Test Mixed Invoice
Create test job with:
- [ ] Template items (quantity > 0)
- [ ] Custom items
- [ ] Items with quantity = 0 (should be excluded)
- [ ] permitRequired = true
- [ ] Verify all items appear correctly
- [ ] Verify permit adds as separate line item

### 6. Test Edge Cases
- [ ] Job with only custom items
- [ ] Job with only permit (no other items)
- [ ] Address without commas (city = full address)
- [ ] Very long custom descriptions
- [ ] Special characters in customer names

### 7. Test Error Handling
- [ ] Job with no completed items (should fail gracefully)
- [ ] Invalid customer ID
- [ ] Network timeout handling
- [ ] Wave API errors

## Common Issues to Watch For

1. **Product Not Found Error**
   - Check exact spelling and case
   - Verify product exists in Wave
   - Check logs for product search attempts

2. **Customer Not Found**
   - Must be exactly "Skyview" (capital S)
   - Check if SKYVIEW_CUSTOMER_ID env var is set

3. **Invoice Creation Fails**
   - Check all prices are formatted as "X.XX" strings
   - Verify all required fields are present
   - Check debug logs for exact error

4. **Missing Line Items**
   - Ensure completedQuantity > 0
   - Check isCustomItem flag is set correctly
   - Verify item data structure matches expected format

## Debug Log Locations

- Development: `/tmp/wave-invoice-debug.log`
- Production: `/var/log/wave-invoice-debug.log`

## Verification in Wave Dashboard

After successful invoice creation:
1. Log into Wave
2. Navigate to Sales > Invoices
3. Find invoice by number
4. Verify:
   - Customer is "Skyview"
   - PO number shows end customer and city
   - All line items present with correct quantities
   - Descriptions appear under product names
   - Total amount is correct

## Production Deployment Checklist

- [ ] All tests pass in development
- [ ] Environment variables set in production
- [ ] Log file permissions configured
- [ ] Test with one real job first
- [ ] Monitor logs for first 10 invoices
- [ ] Verify invoices appear correctly in Wave