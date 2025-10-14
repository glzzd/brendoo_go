#!/bin/bash

# Test script for Product API
echo "🚀 Testing Product API..."

# Base URL
BASE_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    
    echo -e "${YELLOW}Testing: $method $endpoint${NC}"
    
    if [ "$method" = "POST" ]; then
        response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
    fi
    
    # Extract status code (last line)
    status_code=$(echo "$response" | tail -n1)
    # Extract response body (all but last line)
    body=$(echo "$response" | head -n -1)
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ Success: $status_code${NC}"
        echo "Response: $body"
    else
        echo -e "${RED}❌ Failed: Expected $expected_status, got $status_code${NC}"
        echo "Response: $body"
    fi
    echo "---"
}

# Wait for server to start
echo "⏳ Waiting for server to start..."
sleep 3

# Test 1: Health Check
test_endpoint "GET" "/health" "" "200"

# Test 2: Bulk Insert Products
echo "📦 Testing bulk product insert..."
test_data='[
    {
        "_id": "test-product-1",
        "name": "Test Ürün 1",
        "brand": "test-brand",
        "price": 99.99,
        "currency": "TRY",
        "description": "Test açıklaması",
        "images": ["https://example.com/image1.jpg"],
        "sizes": [{"sizeName": "M", "onStock": true}],
        "colors": [{"name": "Mavi", "hex": "#0000FF"}],
        "productUrl": "https://example.com/product1",
        "store": "test-store",
        "category": "test-category",
        "processedAt": "12:00:00",
        "isActive": true,
        "stockStatus": "in_stock",
        "stock": {"quantity": 10, "isInStock": true}
    },
    {
        "_id": "test-product-2",
        "name": "Test Ürün 2",
        "brand": "test-brand",
        "price": 149.99,
        "currency": "TRY",
        "description": "Test açıklaması 2",
        "images": ["https://example.com/image2.jpg"],
        "sizes": [{"sizeName": "L", "onStock": true}],
        "colors": [{"name": "Kırmızı", "hex": "#FF0000"}],
        "productUrl": "https://example.com/product2",
        "store": "zara",
        "category": "test-category",
        "processedAt": "12:00:00",
        "isActive": true,
        "stockStatus": "in_stock",
        "stock": {"quantity": 5, "isInStock": true}
    }
]'

test_endpoint "POST" "/api/v1/products" "$test_data" "201"

# Test 3: Get All Products
test_endpoint "GET" "/api/v1/products?page=1&limit=10" "" "200"

# Test 4: Get Products by Store
test_endpoint "GET" "/api/v1/products/store/zara?page=1&limit=10" "" "200"

# Test 5: Get Products by Store (test-store)
test_endpoint "GET" "/api/v1/products/store/test-store?page=1&limit=10" "" "200"

echo -e "${GREEN}🎉 API testing completed!${NC}"