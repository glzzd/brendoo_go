package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ImagesDir = "uploads/images"
	MaxImageSize = 10 * 1024 * 1024 // 10MB
)

// ImageProcessor handles image processing and storage
type ImageProcessor struct {
	BaseDir string
}

// NewImageProcessor creates a new image processor
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		BaseDir: ImagesDir,
	}
}

// EnsureImageDir creates the images directory if it doesn't exist
func (ip *ImageProcessor) EnsureImageDir() error {
	if err := os.MkdirAll(ip.BaseDir, 0755); err != nil {
		return fmt.Errorf("failed to create images directory: %v", err)
	}
	return nil
}

// ProcessImages processes an array of images (base64 or URLs) and returns file paths
func (ip *ImageProcessor) ProcessImages(images []interface{}) ([]string, error) {
	if err := ip.EnsureImageDir(); err != nil {
		return nil, err
	}

	var processedImages []string
	
	for i, img := range images {
		if img == nil {
			continue
		}
		
		imgStr, ok := img.(string)
		if !ok {
			log.Printf("[WARN] ProcessImages: Image at index %d is not a string, skipping", i)
			continue
		}
		
		if imgStr == "" {
			continue
		}
		
		// Process the image based on its format
		filePath, err := ip.ProcessSingleImage(imgStr)
		if err != nil {
			log.Printf("[ERROR] ProcessImages: Failed to process image at index %d: %v", i, err)
			continue
		}
		
		if filePath != "" {
			processedImages = append(processedImages, filePath)
		}
	}
	
	return processedImages, nil
}

// ProcessSingleImage processes a single image (base64 or URL) and returns file path
func (ip *ImageProcessor) ProcessSingleImage(imageData string) (string, error) {
	// Check if it's a base64 image
	if strings.HasPrefix(imageData, "data:image/") {
		return ip.SaveBase64Image(imageData)
	}
	
	// Check if it's a URL
	if strings.HasPrefix(imageData, "http://") || strings.HasPrefix(imageData, "https://") {
		return ip.DownloadAndSaveImage(imageData)
	}
	
	// If it's already a file path, return as is
	if strings.HasPrefix(imageData, "/uploads/") || strings.HasPrefix(imageData, "uploads/") {
		return imageData, nil
	}
	
	log.Printf("[WARN] ProcessSingleImage: Unknown image format: %s", imageData[:min(50, len(imageData))])
	return imageData, nil // Return original if we can't process it
}

// SaveBase64Image saves a base64 encoded image to disk
func (ip *ImageProcessor) SaveBase64Image(base64Data string) (string, error) {
	// Parse the base64 data
	parts := strings.Split(base64Data, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 image format")
	}
	
	// Extract mime type
	mimeType := ""
	if strings.Contains(parts[0], "image/jpeg") || strings.Contains(parts[0], "image/jpg") {
		mimeType = "jpg"
	} else if strings.Contains(parts[0], "image/png") {
		mimeType = "png"
	} else if strings.Contains(parts[0], "image/webp") {
		mimeType = "webp"
	} else if strings.Contains(parts[0], "image/gif") {
		mimeType = "gif"
	} else {
		mimeType = "jpg" // Default to jpg
	}
	
	// Decode base64
	imageBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}
	
	// Check file size
	if len(imageBytes) > MaxImageSize {
		return "", fmt.Errorf("image size exceeds maximum allowed size of %d bytes", MaxImageSize)
	}
	
	// Generate filename
	hash := md5.Sum(imageBytes)
	filename := fmt.Sprintf("%x_%d.%s", hash, time.Now().Unix(), mimeType)
	filePath := filepath.Join(ip.BaseDir, filename)
	
	// Save to disk
	if err := os.WriteFile(filePath, imageBytes, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}
	
	log.Printf("[DEBUG] SaveBase64Image: Saved base64 image to %s (size: %d bytes)", filePath, len(imageBytes))
	return "/" + filePath, nil
}

// DownloadAndSaveImage downloads an image from URL and saves it to disk
func (ip *ImageProcessor) DownloadAndSaveImage(imageURL string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Download the image
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}
	
	// Check content type
	contentType := resp.Header.Get("Content-Type")
	var extension string
	switch contentType {
	case "image/jpeg", "image/jpg":
		extension = "jpg"
	case "image/png":
		extension = "png"
	case "image/webp":
		extension = "webp"
	case "image/gif":
		extension = "gif"
	default:
		// Try to guess from URL
		if strings.Contains(imageURL, ".jpg") || strings.Contains(imageURL, ".jpeg") {
			extension = "jpg"
		} else if strings.Contains(imageURL, ".png") {
			extension = "png"
		} else if strings.Contains(imageURL, ".webp") {
			extension = "webp"
		} else if strings.Contains(imageURL, ".gif") {
			extension = "gif"
		} else {
			extension = "jpg" // Default
		}
	}
	
	// Read image data with size limit
	limitedReader := io.LimitReader(resp.Body, MaxImageSize+1)
	imageBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}
	
	if len(imageBytes) > MaxImageSize {
		return "", fmt.Errorf("image size exceeds maximum allowed size of %d bytes", MaxImageSize)
	}
	
	// Generate filename
	hash := md5.Sum(imageBytes)
	filename := fmt.Sprintf("%x_%d.%s", hash, time.Now().Unix(), extension)
	filePath := filepath.Join(ip.BaseDir, filename)
	
	// Save to disk
	if err := os.WriteFile(filePath, imageBytes, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}
	
	log.Printf("[DEBUG] DownloadAndSaveImage: Downloaded and saved image from %s to %s (size: %d bytes)", imageURL, filePath, len(imageBytes))
	return "/" + filePath, nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}