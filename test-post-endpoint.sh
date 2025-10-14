#!/bin/bash

# POST Endpoint Test Script
# API'ye örnek ürün gönderme testi

echo "🧪 POST Endpoint Test"
echo "===================="

SERVER_URL="http://69.62.114.202:6000"

# Test 1: Tek ürün gönderme
echo "📦 Test 1: Tek ürün gönderme..."

curl -X POST "$SERVER_URL/api/stock/add" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "_id": "test-product-001",
      "name": "Test Ürün",
      "brand": "Test Brand",
      "price": 99.99,
      "currency": "TRY",
      "priceInRubles": 299.97,
      "description": "Bu bir test ürünüdür",
      "images": ["https://example.com/image1.jpg"],
      "sizes": [
        {
          "onStock": true,
          "sizeName": "M"
        }
      ],
      "colors": [
        {
          "hex": "#FF0000",
          "name": "Kırmızı"
        }
      ],
      "productUrl": "https://example.com/product/test-001",
      "store": "test-store",
      "category": "test-category",
      "processedAt": "2025-10-03T19:30:00Z",
      "isActive": true,
      "stockStatus": "in_stock",
      "stock": {
        "quantity": 10,
        "isInStock": true
      }
    }
  ]'

echo -e "\n\n"

# Test 2: Boş array gönderme (hata testi)
echo "❌ Test 2: Boş array gönderme (hata bekleniyor)..."

curl -X POST "$SERVER_URL/api/stock/add" \
  -H "Content-Type: application/json" \
  -d '[]'

echo -e "\n\n"

# Test 3: Geçersiz JSON gönderme (hata testi)
echo "❌ Test 3: Geçersiz JSON gönderme (hata bekleniyor)..."

curl -X POST "$SERVER_URL/api/stock/add" \
  -H "Content-Type: application/json" \
  -d '{"invalid": "json"}'

echo -e "\n\n"

# Test 4: Eksik alanlar ile gönderme (hata testi)
echo "❌ Test 4: Eksik alanlar ile gönderme (hata bekleniyor)..."

curl -X POST "$SERVER_URL/api/stock/add" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "name": "Eksik Alan Ürün"
    }
  ]'

echo -e "\n\n"

# Test 5: Ürünleri kontrol etme
echo "📋 Test 5: Eklenen ürünleri kontrol etme..."

curl -X GET "$SERVER_URL/api/stock/integration/store?store=test-store"

echo -e "\n\n🎉 Test tamamlandı!"