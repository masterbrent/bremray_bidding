package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const WAVE_API_ENDPOINT = "https://gql.waveapps.com/graphql/public"

// WaveCredentials holds the Wave API credentials
type WaveCredentials struct {
	APIKey     string
	BusinessID string
}

// LineItem represents an invoice line item
type LineItem struct {
	ProductName string  `json:"productName"`
	Description string  `json:"description,omitempty"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

// WaveProduct represents a product in Wave
type WaveProduct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WaveCustomer represents a customer in Wave
type WaveCustomer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WaveInvoice represents a created invoice
type WaveInvoice struct {
	ID            string `json:"id"`
	InvoiceNumber string `json:"invoiceNumber"`
	ViewURL       string `json:"viewUrl"`
}

// WaveAPIService handles all Wave API interactions
type WaveAPIService struct {
	credentials WaveCredentials
	httpClient  *http.Client
	logFile     *os.File
}

// NewWaveAPIService creates a new Wave API service instance
func NewWaveAPIService(credentials WaveCredentials) (*WaveAPIService, error) {
	logPath := "/tmp/wave-invoice-debug.log"
	if os.Getenv("NODE_ENV") == "production" {
		logPath = "/var/log/wave-invoice-debug.log"
	}
	
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	
	return &WaveAPIService{
		credentials: credentials,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logFile: logFile,
	}, nil
}

// Close closes the log file
func (w *WaveAPIService) Close() error {
	if w.logFile != nil {
		return w.logFile.Close()
	}
	return nil
}

// log writes to the debug log file
func (w *WaveAPIService) log(message string, data interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] [WaveAPI] %s", timestamp, message)
	if data != nil {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		logEntry += "\n" + string(jsonData)
	}
	logEntry += "\n\n"
	
	w.logFile.WriteString(logEntry)
	log.Printf("%s", message)
}

// makeGraphQLRequest makes a GraphQL request to Wave API
func (w *WaveAPIService) makeGraphQLRequest(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	w.log("Making GraphQL request:", map[string]interface{}{
		"query":     query[:min(200, len(query))] + "...",
		"variables": variables,
	})
	
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", WAVE_API_ENDPOINT, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+w.credentials.APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		w.log("Wave API Error Response:", map[string]interface{}{
			"status":       resp.StatusCode,
			"statusText":   resp.Status,
			"responseText": string(body),
		})
		return nil, fmt.Errorf("Wave API request failed: %s", resp.Status)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		w.log("Failed to parse Wave response:", string(body))
		return nil, fmt.Errorf("invalid JSON response from Wave API: %w", err)
	}
	
	if errors, ok := result["errors"]; ok {
		w.log("Wave GraphQL errors:", errors)
		return nil, fmt.Errorf("Wave GraphQL errors: %v", errors)
	}
	
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Wave API")
	}
	
	w.log("Wave API Response Success:", map[string]interface{}{
		"data": string(body)[:min(500, len(body))] + "...",
	})
	
	return data, nil
}

// FindCustomerByName finds a customer by name (case-insensitive)
func (w *WaveAPIService) FindCustomerByName(customerName string) (*WaveCustomer, error) {
	query := fmt.Sprintf(`
		query {
			business(id: "%s") {
				customers(page: 1, pageSize: 200) {
					edges {
						node {
							id
							name
						}
					}
				}
			}
		}
	`, w.credentials.BusinessID)
	
	data, err := w.makeGraphQLRequest(query, nil)
	if err != nil {
		return nil, err
	}
	
	business, ok := data["business"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("business not found in response")
	}
	
	customers, ok := business["customers"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("customers not found in response")
	}
	
	edges, ok := customers["edges"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("edges not found in response")
	}
	
	customerNameLower := strings.ToLower(customerName)
	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		
		name, ok := node["name"].(string)
		if !ok {
			continue
		}
		
		if strings.ToLower(name) == customerNameLower {
			id, _ := node["id"].(string)
			return &WaveCustomer{
				ID:   id,
				Name: name,
			}, nil
		}
	}
	
	return nil, nil
}

// GetProducts fetches all products from Wave
func (w *WaveAPIService) GetProducts() ([]WaveProduct, error) {
	query := fmt.Sprintf(`
		query {
			business(id: "%s") {
				products(page: 1, pageSize: 200) {
					edges {
						node {
							id
							name
						}
					}
				}
			}
		}
	`, w.credentials.BusinessID)
	
	data, err := w.makeGraphQLRequest(query, nil)
	if err != nil {
		return nil, err
	}
	
	business, ok := data["business"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("business not found in response")
	}
	
	products, ok := business["products"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("products not found in response")
	}
	
	edges, ok := products["edges"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("edges not found in response")
	}
	
	var result []WaveProduct
	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		
		id, _ := node["id"].(string)
		name, _ := node["name"].(string)
		
		result = append(result, WaveProduct{
			ID:   id,
			Name: name,
		})
	}
	
	return result, nil
}

// GetDefaultIncomeAccount finds the default income account
func (w *WaveAPIService) GetDefaultIncomeAccount() (string, error) {
	query := fmt.Sprintf(`
		query {
			business(id: "%s") {
				accounts(types: [INCOME], page: 1, pageSize: 10) {
					edges {
						node {
							id
							name
						}
					}
				}
			}
		}
	`, w.credentials.BusinessID)
	
	data, err := w.makeGraphQLRequest(query, nil)
	if err != nil {
		return "", err
	}
	
	business, ok := data["business"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("business not found in response")
	}
	
	accounts, ok := business["accounts"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("accounts not found in response")
	}
	
	edges, ok := accounts["edges"].([]interface{})
	if !ok || len(edges) == 0 {
		return "", fmt.Errorf("no income accounts found")
	}
	
	// Look for sales or income account
	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		
		name, _ := node["name"].(string)
		nameLower := strings.ToLower(name)
		
		if strings.Contains(nameLower, "sales") || strings.Contains(nameLower, "income") {
			id, _ := node["id"].(string)
			return id, nil
		}
	}
	
	// Return first account if no sales/income account found
	firstEdge := edges[0].(map[string]interface{})
	firstNode := firstEdge["node"].(map[string]interface{})
	id, _ := firstNode["id"].(string)
	return id, nil
}

// FindOrCreateProduct finds an existing product or creates a new one
func (w *WaveAPIService) FindOrCreateProduct(productName string) (string, error) {
	w.log("Looking for product:", productName)
	
	// First, try to find existing product
	existingProducts, err := w.GetProducts()
	if err != nil {
		return "", fmt.Errorf("failed to get products: %w", err)
	}
	
	w.log(fmt.Sprintf("Found %d products in Wave", len(existingProducts)), nil)
	
	// Case-insensitive search
	productNameLower := strings.ToLower(productName)
	for _, product := range existingProducts {
		if strings.ToLower(product.Name) == productNameLower {
			w.log("Found existing product:", map[string]interface{}{
				"productName": productName,
				"productId":   product.ID,
				"exactName":   product.Name,
			})
			return product.ID, nil
		}
	}
	
	// Product doesn't exist, create it
	w.log(fmt.Sprintf("Product \"%s\" not found in Wave. Creating new product...", productName), nil)
	
	// Get default income account
	incomeAccountID, err := w.GetDefaultIncomeAccount()
	if err != nil {
		return "", fmt.Errorf("could not find default income account: %w", err)
	}
	
	// Create the product
	mutation := `
		mutation CreateProduct($input: ProductCreateInput!) {
			productCreate(input: $input) {
				product {
					id
					name
				}
				didSucceed
				inputErrors {
					path
					message
					code
				}
			}
		}
	`
	
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"businessId":      w.credentials.BusinessID,
			"name":            productName,
			"description":     "", // MUST be empty string, not null
			"unitPrice":       "0.00", // MUST be string with 2 decimals
			"incomeAccountId": incomeAccountID,
		},
	}
	
	w.log("Creating product with input:", variables)
	
	data, err := w.makeGraphQLRequest(mutation, variables)
	if err != nil {
		return "", fmt.Errorf("failed to create product: %w", err)
	}
	
	productCreate, ok := data["productCreate"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("productCreate not found in response")
	}
	
	didSucceed, _ := productCreate["didSucceed"].(bool)
	if !didSucceed {
		inputErrors, _ := productCreate["inputErrors"].([]interface{})
		var errorMessages []string
		for _, err := range inputErrors {
			errMap, _ := err.(map[string]interface{})
			path, _ := errMap["path"].(string)
			message, _ := errMap["message"].(string)
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", path, message))
		}
		w.log("Product creation failed with errors:", inputErrors)
		return "", fmt.Errorf("product creation failed: %s", strings.Join(errorMessages, ", "))
	}
	
	product, ok := productCreate["product"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("product not found in response")
	}
	
	id, _ := product["id"].(string)
	name, _ := product["name"].(string)
	
	w.log("Product created successfully:", map[string]interface{}{
		"id":   id,
		"name": name,
	})
	
	return id, nil
}

// CreateInvoice creates an invoice in Wave
func (w *WaveAPIService) CreateInvoice(customerID string, lineItems []LineItem, poNumber string, invoiceDate string) (*WaveInvoice, error) {
	// Prepare invoice items
	var invoiceItems []map[string]interface{}
	
	for _, item := range lineItems {
		productID, err := w.FindOrCreateProduct(item.ProductName)
		if err != nil {
			return nil, fmt.Errorf("failed to find/create product %s: %w", item.ProductName, err)
		}
		
		invoiceItem := map[string]interface{}{
			"productId": productID,
			"quantity":  item.Quantity,
			"unitPrice": fmt.Sprintf("%.2f", item.Price),
		}
		
		// Add description if provided
		if item.Description != "" {
			invoiceItem["description"] = item.Description
		}
		
		invoiceItems = append(invoiceItems, invoiceItem)
	}
	
	// Use the documented GraphQL mutation structure
	mutation := `
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
	`
	
	if invoiceDate == "" {
		invoiceDate = time.Now().Format("2006-01-02")
	}
	
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"businessId":  w.credentials.BusinessID,
			"customerId":  customerID,
			"poNumber":    poNumber,
			"items":       invoiceItems,
			"invoiceDate": invoiceDate,
		},
	}
	
	w.log("Creating Wave invoice with variables:", variables)
	
	data, err := w.makeGraphQLRequest(mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}
	
	invoiceCreate, ok := data["invoiceCreate"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invoiceCreate not found in response")
	}
	
	didSucceed, _ := invoiceCreate["didSucceed"].(bool)
	if !didSucceed {
		inputErrors, _ := invoiceCreate["inputErrors"].([]interface{})
		var errorMessages []string
		for _, err := range inputErrors {
			errMap, _ := err.(map[string]interface{})
			path, _ := errMap["path"].(string)
			message, _ := errMap["message"].(string)
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", path, message))
		}
		return nil, fmt.Errorf("failed to create invoice: %s", strings.Join(errorMessages, ", "))
	}
	
	invoice, ok := invoiceCreate["invoice"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invoice not found in response")
	}
	
	id, _ := invoice["id"].(string)
	invoiceNumber, _ := invoice["invoiceNumber"].(string)
	viewURL, _ := invoice["viewUrl"].(string)
	
	return &WaveInvoice{
		ID:            id,
		InvoiceNumber: invoiceNumber,
		ViewURL:       viewURL,
	}, nil
}

// GetWaveCredentials gets Wave credentials from environment
func GetWaveCredentials() (*WaveCredentials, error) {
	apiKey := os.Getenv("WAVE_TOKEN")
	if apiKey == "" {
		apiKey = os.Getenv("WAVE_ACCESS_TOKEN")
	}
	if apiKey == "" {
		apiKey = os.Getenv("WAVEAPPS_FULL_ACCESS_ID")
	}
	
	businessID := os.Getenv("WAVE_BUSINESS_ID")
	if businessID == "" {
		businessID = os.Getenv("WAVEAPPS_BUSINESS_ID")
	}
	
	if apiKey == "" || businessID == "" {
		return nil, fmt.Errorf("Wave API credentials not configured. Please set WAVE_TOKEN and WAVE_BUSINESS_ID environment variables")
	}
	
	return &WaveCredentials{
		APIKey:     apiKey,
		BusinessID: businessID,
	}, nil
}

// IsWaveConfigured checks if Wave is properly configured
func IsWaveConfigured() bool {
	_, err := GetWaveCredentials()
	return err == nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PrepareJobForWaveInvoice prepares a job for Wave invoice creation
// This encapsulates the exact logic from production
func PrepareJobForWaveInvoice(job map[string]interface{}) (string, []LineItem, error) {
	// Extract customer ID - first check environment, then search by name
	customerID := os.Getenv("SKYVIEW_CUSTOMER_ID")
	
	credentials, err := GetWaveCredentials()
	if err != nil {
		return "", nil, err
	}
	
	waveService, err := NewWaveAPIService(*credentials)
	if err != nil {
		return "", nil, err
	}
	defer waveService.Close()
	
	if customerID == "" {
		// Search for Skyview customer
		customer, err := waveService.FindCustomerByName("Skyview")
		if err != nil {
			return "", nil, fmt.Errorf("failed to search for customer: %w", err)
		}
		if customer == nil {
			return "", nil, fmt.Errorf("Skyview customer not found in Wave. Please create a customer named 'Skyview' in your Wave account")
		}
		customerID = customer.ID
	}
	
	// Prepare line items
	var lineItems []LineItem
	
	// Process job items
	items, ok := job["items"].([]interface{})
	if ok {
		for _, itemInterface := range items {
			jobItem, ok := itemInterface.(map[string]interface{})
			if !ok {
				continue
			}
			
			completedQty, _ := jobItem["completedQuantity"].(float64)
			if completedQty <= 0 {
				continue
			}
			
			isCustom, _ := jobItem["isCustomItem"].(bool)
			
			if isCustom {
				// Custom items - use "Custom Work" as product name
				customPrice, _ := jobItem["customPrice"].(string)
				price, _ := strconv.ParseFloat(customPrice, 64)
				
				description := ""
				if desc, ok := jobItem["customDescription"].(string); ok && desc != "" {
					description = desc
				} else if name, ok := jobItem["customName"].(string); ok && name != "" {
					description = name
				}
				
				lineItems = append(lineItems, LineItem{
					ProductName: "Custom Work",
					Description: description,
					Quantity:    int(completedQty),
					Price:       price,
					Total:       completedQty * price,
				})
			} else {
				// Template items
				item, ok := jobItem["item"].(map[string]interface{})
				if !ok {
					continue
				}
				
				name, _ := item["name"].(string)
				priceStr, _ := item["price"].(string)
				price, _ := strconv.ParseFloat(priceStr, 64)
				description, _ := item["description"].(string)
				
				lineItems = append(lineItems, LineItem{
					ProductName: name,
					Description: description,
					Quantity:    int(completedQty),
					Price:       price,
					Total:       completedQty * price,
				})
			}
		}
	}
	
	// Add Electrical Permit if required
	if permitRequired, ok := job["permitRequired"].(bool); ok && permitRequired {
		// Default permit price is $250 unless specified differently
		permitPrice := 250.0
		
		lineItems = append(lineItems, LineItem{
			ProductName: "Electrical Permit",
			Description: "Required for this project",
			Quantity:    1,
			Price:       permitPrice,
			Total:       permitPrice,
		})
	}
	
	if len(lineItems) == 0 {
		return "", nil, fmt.Errorf("no billable items found in job")
	}
	
	return customerID, lineItems, nil
}

// FormatPONumber formats the PO number according to production logic
func FormatPONumber(customerName, address string) string {
	// Extract city from address
	addressParts := strings.Split(address, ",")
	city := address // Default to full address
	
	if len(addressParts) > 1 {
		// Use second-to-last part as city
		city = strings.TrimSpace(addressParts[len(addressParts)-2])
	}
	
	return fmt.Sprintf("%s - %s", customerName, city)
}