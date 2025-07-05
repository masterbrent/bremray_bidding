# CRITICAL Wave Integration Documentation - DO NOT DEVIATE

**WARNING**: Wave's API is extremely sensitive. Even the smallest deviation from this exact implementation will cause invoice creation to fail. This document captures the EXACT working implementation from production.

## Environment Variables Required

```bash
# Wave API Credentials (MUST have both)
WAVE_TOKEN=your_wave_api_token_here
WAVE_BUSINESS_ID=your_wave_business_id_here

# Optional - If not set, system will search for "Skyview" customer by name
SKYVIEW_CUSTOMER_ID=customer_id_from_wave
```

## Critical Implementation Details

### 1. Wave API Endpoint
- **MUST** use: `https://gql.waveapps.com/graphql/public`
- **MUST** use POST method
- **MUST** include headers:
  ```
  Content-Type: application/json
  Authorization: Bearer ${WAVE_TOKEN}
  ```

### 2. Customer Handling
- **ALWAYS** use "Skyview" as the customer name
- System first checks for `SKYVIEW_CUSTOMER_ID` environment variable
- If not found, searches for customer by name using **case-insensitive** search
- Customer search query MUST use pagination: `customers(page: 1, pageSize: 200)`

### 3. Product/Service Handling

#### Finding Products
- **MUST** search products with **case-insensitive** comparison
- Use exact query structure:
  ```graphql
  products(page: 1, pageSize: 200) {
    edges {
      node {
        id
        name
      }
    }
  }
  ```

#### Creating Products When Not Found
- **MUST** get default income account first
- Income account search: Look for accounts with "sales" or "income" in name
- Product creation **MUST** use these exact fields:
  ```javascript
  {
    businessId: credentials.businessId,
    name: productName,
    description: "",  // MUST be empty string, not null
    unitPrice: "0.00",  // MUST be string with 2 decimals
    incomeAccountId: incomeAccountId  // REQUIRED
  }
  ```

### 4. Invoice Line Items Structure

#### For Template Items (from items table):
```javascript
{
  productName: item.name,  // Product name from items table
  description: item.description || undefined,  // Optional description
  quantity: completedQuantity,  // Integer
  price: parseFloat(item.price).toFixed(2),  // String with 2 decimals
  total: (completedQuantity * price)  // Calculated
}
```

#### For Custom Items:
```javascript
{
  productName: "Custom Work",  // ALWAYS use this exact name
  description: customDescription || customName || undefined,  // Put actual work description here
  quantity: completedQuantity,
  price: parseFloat(customPrice).toFixed(2),
  total: (completedQuantity * price)
}
```

#### For Electrical Permit:
```javascript
{
  productName: "Electrical Permit",  // EXACT name
  description: "Required for this project",  // EXACT description
  quantity: 1,
  price: permitPrice,  // From items table or default $250
  total: permitPrice
}
```

### 5. Invoice Creation Mutation

**EXACT GraphQL Mutation Structure:**
```graphql
mutation CreateInvoice($input: InvoiceCreateInput!) {
  invoiceCreate(input: $input) {
    invoice { 
      id 
      invoiceNumber
      viewUrl 
      status 
    }
    didSucceed
    inputErrors {
      path
      message
    }
  }
}
```

**EXACT Variables Structure:**
```javascript
{
  input: {
    businessId: credentials.businessId,
    customerId: customerId,
    poNumber: poNumber,  // Format: "{CustomerName} - {City}"
    items: invoiceItems,  // Array of items
    invoiceDate: invoiceDate || new Date().toISOString().split('T')[0]
  }
}
```

### 6. Invoice Item Structure for Mutation
```javascript
{
  productId: productId,  // From findOrCreateProduct
  quantity: quantity,    // Integer
  unitPrice: price.toFixed(2),  // String with 2 decimals
  description: description  // Optional - appears under product name
}
```

### 7. PO Number Format
- **MUST** extract city from address
- If address has commas: use second-to-last part as city
- If no commas: use entire address as city
- Format: `${job.customerName} - ${city}`

### 8. Error Handling
- Check for `didSucceed` in mutation response
- If false, extract errors from `inputErrors` array
- Format errors as: `${error.path}: ${error.message}`
- Log ALL requests and responses to debug log file

### 9. Post-Invoice Creation
- Update job status to 'invoiced'
- Store Wave invoice details:
  - `waveInvoiceId`
  - `waveInvoiceNumber`
  - `waveInvoiceUrl`

## Common Pitfalls That WILL Break Integration

1. **Using null instead of undefined for optional fields**
2. **Not formatting prices as strings with exactly 2 decimals**
3. **Case-sensitive product name matching**
4. **Missing incomeAccountId when creating products**
5. **Using wrong field names in GraphQL mutations**
6. **Not checking didSucceed before accessing response data**
7. **Forgetting to find/create products before creating invoice**

## Testing Checklist

Before deploying, verify:
- [ ] Can find Skyview customer
- [ ] Can find existing products (case-insensitive)
- [ ] Can create new products with income account
- [ ] Can create invoice with template items
- [ ] Can create invoice with custom items
- [ ] Can create invoice with permit
- [ ] PO number shows customer name and city
- [ ] Invoice URL is returned and stored

## Debug Logging

Production log location: `/var/log/wave-invoice-debug.log`
Development log location: `/tmp/wave-invoice-debug.log`

Log MUST include:
- All GraphQL queries and variables
- All API responses
- Product search/creation attempts
- Complete error messages and stack traces