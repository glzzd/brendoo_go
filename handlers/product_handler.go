package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"product-api/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

// NewProductHandler creates a new product handler
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		DB: db,
	}
}

// checkIDExistsInDB checks if an ID already exists in the database
func (h *ProductHandler) checkIDExistsInDB(id string) bool {
	var count int64
	h.DB.Model(&models.Product{}).Where("id = ?", id).Count(&count)
	return count > 0
}

// BulkCreateProducts handles bulk insertion of products (optimized for millions of records)
// Accepts both single product object and array of products
func (h *ProductHandler) BulkCreateProducts(c *gin.Context) {
	// Log incoming request
	log.Printf("[DEBUG] BulkCreateProducts: Received POST request from %s", c.ClientIP())

	// Read raw body for logging
	rawBody, err := c.GetRawData()
	if err != nil {
		log.Printf("[ERROR] BulkCreateProducts: Failed to read raw body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Log raw request body (first 500 chars to avoid huge logs)
	bodyStr := string(rawBody)
	if len(bodyStr) > 500 {
		log.Printf("[DEBUG] BulkCreateProducts: Raw body (truncated): %s...", bodyStr[:500])
	} else {
		log.Printf("[DEBUG] BulkCreateProducts: Raw body: %s", bodyStr)
	}

	var products []models.Product

	// Try to parse as array first
	if err := json.Unmarshal(rawBody, &products); err != nil {
		log.Printf("[DEBUG] BulkCreateProducts: Array parsing failed, trying single object: %v", err)

		// If array parsing fails, try parsing as single object
		var singleProduct models.Product
		if err := json.Unmarshal(rawBody, &singleProduct); err != nil {
			log.Printf("[ERROR] BulkCreateProducts: Both array and single object parsing failed: %v", err)
			log.Printf("[ERROR] BulkCreateProducts: Raw body that failed: %s", bodyStr)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         "Invalid JSON format - must be either a product object or array of products",
				"details":       err.Error(),
				"received_body": bodyStr,
			})
			return
		}

		// Successfully parsed as single object, convert to array
		products = []models.Product{singleProduct}
		log.Printf("[DEBUG] BulkCreateProducts: Successfully parsed single product, converted to array")
	} else {
		log.Printf("[DEBUG] BulkCreateProducts: Successfully parsed as array with %d products", len(products))
	}

	if len(products) == 0 {
		log.Printf("[WARN] BulkCreateProducts: No products provided in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No products provided"})
		return
	}

	// Log first product for debugging
	if len(products) > 0 {
		firstProduct, _ := json.Marshal(products[0])
		log.Printf("[DEBUG] BulkCreateProducts: First product: %s", string(firstProduct))
	}

	// Implement upsert logic
	log.Printf("[DEBUG] BulkCreateProducts: Starting upsert process for %d products", len(products))

	// Generate unique IDs for products without ID and validate duplicates within the batch
	seenIDs := make(map[string]bool)
	seenNames := make(map[string]bool)
	seenNameUrlCombos := make(map[string]bool) // Track name+productUrl combinations
	var uniqueProducts []models.Product
	duplicateCount := 0

	for i := range products {
		log.Printf("[DEBUG] BulkCreateProducts: Product %d - ID: %s, Name: %s, ProductURL: %s",
			i, products[i].ID, products[i].Name, products[i].ProductURL)

		// Check if ProductURL already exists in database
		if products[i].ProductURL != "" {
			var existingProduct models.Product
			result := h.DB.Where("product_url = ?", products[i].ProductURL).First(&existingProduct)
			if result.Error == nil {
				log.Printf("[WARN] BulkCreateProducts: ProductURL already exists in database: %s, skipping product: %s", 
					products[i].ProductURL, products[i].Name)
				duplicateCount++
				continue // Skip this product as ProductURL already exists
			}
		}

		// Check for duplicate name+productUrl combination within the batch
		nameUrlKey := fmt.Sprintf("%s|%s", products[i].Name, products[i].ProductURL)
		if products[i].Name != "" && products[i].ProductURL != "" {
			if seenNameUrlCombos[nameUrlKey] {
				log.Printf("[WARN] BulkCreateProducts: Duplicate name+productUrl found in batch: %s + %s, skipping", 
					products[i].Name, products[i].ProductURL)
				duplicateCount++
				continue // Skip this duplicate product
			}
			seenNameUrlCombos[nameUrlKey] = true
		}

		// Generate unique ID if empty
		if products[i].ID == "" {
			if products[i].Name != "" {
				// Clean Turkish characters and create base ID
				baseID := strings.ToLower(products[i].Name)
				// Replace Turkish characters
				baseID = strings.ReplaceAll(baseID, "ç", "c")
				baseID = strings.ReplaceAll(baseID, "ğ", "g")
				baseID = strings.ReplaceAll(baseID, "ı", "i")
				baseID = strings.ReplaceAll(baseID, "ö", "o")
				baseID = strings.ReplaceAll(baseID, "ş", "s")
				baseID = strings.ReplaceAll(baseID, "ü", "u")
				// Replace spaces and special characters
				baseID = strings.ReplaceAll(baseID, " ", "-")
				baseID = strings.ReplaceAll(baseID, ".", "")
				baseID = strings.ReplaceAll(baseID, ",", "")
				baseID = strings.ReplaceAll(baseID, "(", "")
				baseID = strings.ReplaceAll(baseID, ")", "")
				baseID = strings.ReplaceAll(baseID, "/", "-")
				baseID = strings.ReplaceAll(baseID, "\\", "-")
				
				counter := 1
				newID := baseID
				
				// Check both in-memory and database for uniqueness
				for seenIDs[newID] || h.checkIDExistsInDB(newID) {
					newID = fmt.Sprintf("%s-%d", baseID, counter)
					counter++
				}
				products[i].ID = newID
				seenIDs[newID] = true
			} else {
				// Generate timestamp-based ID and ensure uniqueness
				baseID := fmt.Sprintf("product-%d-%d", time.Now().Unix(), i)
				counter := 1
				newID := baseID
				
				for seenIDs[newID] || h.checkIDExistsInDB(newID) {
					newID = fmt.Sprintf("%s-%d", baseID, counter)
					counter++
				}
				products[i].ID = newID
				seenIDs[products[i].ID] = true
			}
		} else {
			// Check if provided ID already exists in batch or database
			if seenIDs[products[i].ID] || h.checkIDExistsInDB(products[i].ID) {
				log.Printf("[WARN] BulkCreateProducts: Duplicate ID found (batch or DB): %s, generating new ID", products[i].ID)
				counter := 1
				baseID := products[i].ID
				newID := fmt.Sprintf("%s-%d", baseID, counter)
				for seenIDs[newID] || h.checkIDExistsInDB(newID) {
					counter++
					newID = fmt.Sprintf("%s-%d", baseID, counter)
				}
				products[i].ID = newID
			}
			seenIDs[products[i].ID] = true
		}

		// Track names for duplicate detection
		if products[i].Name != "" {
			if seenNames[products[i].Name] {
				log.Printf("[WARN] BulkCreateProducts: Duplicate name found in batch: %s", products[i].Name)
			}
			seenNames[products[i].Name] = true
		}

		// Calculate PriceInRubles based on Price
		if products[i].Price > 0 {
			price := products[i].Price
			var multiplier float64

			if price >= 1 && price <= 100 {
				multiplier = 120
			} else if price >= 101 && price <= 150 {
				multiplier = 100
			} else if price >= 151 && price <= 200 {
				multiplier = 90
			} else if price >= 201 && price <= 350 {
				multiplier = 85
			} else if price > 350 {
				multiplier = 80
			} else {
				multiplier = 120
			}

			priceInRubles := price * multiplier
			products[i].PriceInRubles = &priceInRubles

			log.Printf("[DEBUG] BulkCreateProducts: Product %s - Price: %.2f, Multiplier: %.0f, PriceInRubles: %.2f",
				products[i].ID, price, multiplier, priceInRubles)
		}

		// Add to unique products list
		uniqueProducts = append(uniqueProducts, products[i])
	}

	log.Printf("[INFO] BulkCreateProducts: Filtered %d duplicates from batch, processing %d unique products", 
		duplicateCount, len(uniqueProducts))

	// Update products slice to only contain unique products
	products = uniqueProducts

	// Delete existing products where (name + productUrl) matches
	for _, product := range products {
		if product.Name != "" && product.ProductURL != "" {
			deleteResult := h.DB.Where("name = ? AND product_url = ?", product.Name, product.ProductURL).
				Delete(&models.Product{})
			if deleteResult.Error != nil {
				log.Printf("[ERROR] BulkCreateProducts: Failed to delete existing product (name=%s, url=%s): %v",
					product.Name, product.ProductURL, deleteResult.Error)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to delete existing product",
					"details": deleteResult.Error.Error(),
				})
				return
			}
			if deleteResult.RowsAffected > 0 {
				log.Printf("[DEBUG] BulkCreateProducts: Deleted %d existing product(s) with (name=%s, url=%s)",
					deleteResult.RowsAffected, product.Name, product.ProductURL)
			}
		}
	}

	// Use batch size for optimal performance
	batchSize := 1000
	log.Printf("[DEBUG] BulkCreateProducts: Starting batch processing with batch size %d", batchSize)

	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		batch := products[i:end]
		log.Printf("[DEBUG] BulkCreateProducts: Processing batch %d-%d (%d products)", i, end-1, len(batch))

		if err := h.DB.CreateInBatches(batch, batchSize).Error; err != nil {
			log.Printf("[ERROR] BulkCreateProducts: Database error in batch %d-%d: %v", i, end-1, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       "Failed to insert products batch",
				"details":     err.Error(),
				"batch_start": i,
				"batch_end":   end,
			})
			return
		}
		log.Printf("[DEBUG] BulkCreateProducts: Successfully inserted batch %d-%d", i, end-1)
	}

	log.Printf("[SUCCESS] BulkCreateProducts: Successfully upserted %d unique products (filtered %d duplicates from original %d products)", 
		len(products), duplicateCount, len(products)+duplicateCount)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Products upserted successfully",
		"count":   len(products),
		"original_count": len(products) + duplicateCount,
		"duplicates_filtered": duplicateCount,
	})
}

// GetAllProducts retrieves all products with pagination
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 100
	}
	
	offset := (page - 1) * limit
	
	var products []models.Product
	var total int64
	
	// Select only necessary fields, exclude heavy images field for better performance
	selectFields := "id, name, brand, price, currency, price_in_rubles, discounted_price, description, sizes, colors, product_url, store, category, processed_at, is_active, stock_status, stock, created_at, updated_at"
	
	// Get total count
	if err := h.DB.Model(&models.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		log.Printf("[ERROR] GetAllProducts: Failed to count products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}
	
	// Get products with pagination
	if err := h.DB.Select(selectFields).Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		log.Printf("[ERROR] GetAllProducts: Failed to fetch products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	
	log.Printf("[DEBUG] GetAllProducts: Successfully fetched %d products (page %d, limit %d)", len(products), page, limit)
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetProductsByStore retrieves products filtered by store with pagination
func (h *ProductHandler) GetProductsByStore(c *gin.Context) {
	store := c.Param("store")
	if store == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Store parameter is required"})
		return
	}
	
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 100
	}
	
	offset := (page - 1) * limit
	
	var products []models.Product
	var total int64
	
	// Select only necessary fields, exclude heavy images field for better performance
	selectFields := "id, name, brand, price, currency, price_in_rubles, discounted_price, description, sizes, colors, product_url, store, category, processed_at, is_active, stock_status, stock, created_at, updated_at"
	
	// Get total count for the store
	if err := h.DB.Model(&models.Product{}).
		Where("store = ? AND is_active = ?", store, true).
		Count(&total).Error; err != nil {
		log.Printf("[ERROR] GetProductsByStore: Failed to count products for store %s: %v", store, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}
	
	// Get products filtered by store with pagination
	if err := h.DB.Select(selectFields).Where("store = ? AND is_active = ?", store, true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		log.Printf("[ERROR] GetProductsByStore: Failed to fetch products for store %s: %v", store, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	
	log.Printf("[DEBUG] GetProductsByStore: Successfully fetched %d products for store %s (page %d, limit %d)", len(products), store, page, limit)
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"store":    store,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetProductsByStoreOrAll retrieves all products or filters by store without pagination
func (h *ProductHandler) GetProductsByStoreOrAll(c *gin.Context) {
	store := c.Query("store")
	
	var products []models.Product
	var result *gorm.DB
	
	if store != "" {
		// Filter by store
		result = h.DB.Where("store = ? AND is_active = ?", store, true).Order("created_at DESC").Find(&products)
	} else {
		// Get all products
		result = h.DB.Where("is_active = ?", true).Order("created_at DESC").Find(&products)
	}
	
	if result.Error != nil {
		log.Printf("[ERROR] GetProductsByStoreOrAll: Failed to fetch products: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	
	log.Printf("[DEBUG] GetProductsByStoreOrAll: Successfully fetched %d products", len(products))
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}

// GetProductImages retrieves images for a specific product by ID
func (h *ProductHandler) GetProductImages(c *gin.Context) {
	productID := c.Param("id")
	
	if productID == "" {
		log.Printf("[ERROR] GetProductImages: Product ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}
	
	var product models.Product
	
	// Select only ID and images field
	if err := h.DB.Select("id, images").Where("id = ? AND is_active = ?", productID, true).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[WARN] GetProductImages: Product not found with ID: %s", productID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		log.Printf("[ERROR] GetProductImages: Failed to fetch product images for ID %s: %v", productID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product images"})
		return
	}
	
	log.Printf("[DEBUG] GetProductImages: Successfully fetched images for product ID: %s", productID)
	
	c.JSON(http.StatusOK, gin.H{
		"id":     product.ID,
		"images": product.Images,
	})
}

// GetMultipleProductImages retrieves images for multiple products by IDs
func (h *ProductHandler) GetMultipleProductImages(c *gin.Context) {
	var requestBody struct {
		ProductIDs []string `json:"product_ids" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("[ERROR] GetMultipleProductImages: Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	if len(requestBody.ProductIDs) == 0 {
		log.Printf("[ERROR] GetMultipleProductImages: No product IDs provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No product IDs provided"})
		return
	}
	
	// Limit to prevent abuse
	if len(requestBody.ProductIDs) > 50 {
		log.Printf("[WARN] GetMultipleProductImages: Too many product IDs requested: %d", len(requestBody.ProductIDs))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 50 product IDs allowed per request"})
		return
	}
	
	var products []models.Product
	
	// Select only ID and images fields
	if err := h.DB.Select("id, images").Where("id IN ? AND is_active = ?", requestBody.ProductIDs, true).Find(&products).Error; err != nil {
		log.Printf("[ERROR] GetMultipleProductImages: Failed to fetch product images: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product images"})
		return
	}
	
	// Create response map
	imageMap := make(map[string]interface{})
	for _, product := range products {
		imageMap[product.ID] = product.Images
	}
	
	log.Printf("[DEBUG] GetMultipleProductImages: Successfully fetched images for %d products", len(products))
	
	c.JSON(http.StatusOK, gin.H{
		"images": imageMap,
		"count":  len(products),
	})
}

// GetProductsWithImages returns products with images included - optimized for when images are needed
func (h *ProductHandler) GetProductsWithImages(c *gin.Context) {
	store := c.Query("store")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	
	// Default pagination
	limit := 100 // Default limit to prevent huge responses
	offset := 0
	
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 500 { // Max limit to prevent abuse
				limit = 500
			}
		}
	}
	
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	var products []models.Product
	var result *gorm.DB
	
	// Include ALL fields including images
	if store != "" {
		// Filter by store with pagination
		result = h.DB.Where("store = ? AND is_active = ?", store, true).
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&products)
	} else {
		// Get all products with pagination
		result = h.DB.Where("is_active = ?", true).
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&products)
	}
	
	if result.Error != nil {
		log.Printf("[ERROR] GetProductsWithImages: Failed to fetch products: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	
	log.Printf("[DEBUG] GetProductsWithImages: Successfully fetched %d products with images (limit: %d, offset: %d)", len(products), limit, offset)
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
		"limit":    limit,
		"offset":   offset,
	})
}