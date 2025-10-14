package routes

import (
	"product-api/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Production IP'si ve localhost'a izin ver
		allowedOrigins := []string{
			"http://69.62.114.202",
			"http://69.62.114.202:6000",
			"https://69.62.114.202",
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:6000",
		}
		
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*") // Geliştirme için
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize handlers
	productHandler := handlers.NewProductHandler(db)

	// API routes
	api := r.Group("/api")
	{
		stock := api.Group("/stock")
		{
			// Bulk product insert
			stock.POST("/add", productHandler.BulkCreateProducts)
			
			// Integration endpoints
			integration := stock.Group("/integration")
			{
				// Get all products or filter by store
				integration.GET("/store", productHandler.GetProductsByStoreOrAll)
				integration.GET("/all", productHandler.GetAllProducts)
			}
			
			// Store specific endpoints
			stock.GET("/store/:store", productHandler.GetProductsByStore)
			
			// NEW: Products WITH images - use when images are needed
			stock.GET("/with-images", productHandler.GetProductsWithImages)
			
			// Image endpoints
			images := stock.Group("/images")
			{
				// Get images for a single product
				images.GET("/:id", productHandler.GetProductImages)
				// Get images for multiple products
				images.POST("/batch", productHandler.GetMultipleProductImages)
			}
		}
	}

	// Static file serving for uploaded images
	r.Static("/uploads", "./uploads")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}