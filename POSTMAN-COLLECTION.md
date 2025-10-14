# 📋 Postman Test Collection - Product API

## 🌐 Base URL
```
http://69.62.114.202:6000
```

## 📁 Collection Structure

### 1. **Health Check**
- **Method:** `GET`
- **URL:** `{{base_url}}/health`
- **Expected Response:**
```json
{
  "status": "ok"
}
```

### 2. **Get All Products**
- **Method:** `GET`
- **URL:** `{{base_url}}/api/stock/integration/store`
- **Query Parameters (Optional):**
  - `page`: Page number (default: 1)
  - `limit`: Items per page (default: 100)
- **Expected Response:**
```json
{
  "count": 1,
  "products": [
    {
      "_id": "test-product-001",
      "name": "Test Ürün",
      "brand": "Test Brand",
      "price": 99.99,
      "currency": "TRY",
      "priceInRubles": 299.97,
      "discountedPrice": null,
      "description": "Bu bir test ürünüdür",
      "images": ["https://example.com/image1.jpg"],
      "sizes": [{"onStock": true, "sizeName": "M"}],
      "colors": [{"hex": "#FF0000", "name": "Kırmızı"}],
      "productUrl": "https://example.com/product/test-001",
      "store": "test-store",
      "category": "test-category",
      "processedAt": "2025-10-03T19:30:00Z",
      "isActive": true,
      "stockStatus": "in_stock",
      "stock": {"quantity": 10, "isInStock": true},
      "createdAt": "2025-10-03T19:32:27.385437Z",
      "updatedAt": "2025-10-03T19:32:27.385437Z"
    }
  ]
}
```

### 3. **Get Products by Store**
- **Method:** `GET`
- **URL:** `{{base_url}}/api/stock/integration/store?store={{store_name}}`
- **Available Stores:**
  - `zara`
  - `nike`
  - `adidas`
  - `gosport`
  - `test-store`

**Example URLs:**
```
{{base_url}}/api/stock/integration/store?store=zara
{{base_url}}/api/stock/integration/store?store=nike
{{base_url}}/api/stock/integration/store?store=adidas
```

### 4. **Add Products (Bulk Insert)**
- **Method:** `POST`
- **URL:** `{{base_url}}/api/stock/add`
- **Headers:**
  - `Content-Type: application/json`
- **Body (JSON Array):**

#### ✅ **Valid Request Example:**
```json
[
  {
    "_id": "unique-product-id-001",
    "name": "Örnek Ürün",
    "brand": "Örnek Marka",
    "price": 199.99,
    "currency": "TRY",
    "priceInRubles": 599.97,
    "discountedPrice": 149.99,
    "description": "Bu örnek bir ürün açıklamasıdır",
    "images": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.jpg"
    ],
    "sizes": [
      {
        "onStock": true,
        "sizeName": "S"
      },
      {
        "onStock": true,
        "sizeName": "M"
      },
      {
        "onStock": false,
        "sizeName": "L"
      }
    ],
    "colors": [
      {
        "hex": "#FF0000",
        "name": "Kırmızı"
      },
      {
        "hex": "#0000FF",
        "name": "Mavi"
      }
    ],
    "productUrl": "https://example.com/product/unique-product-id-001",
    "store": "zara",
    "category": "man-shirts",
    "processedAt": "2025-10-03T19:30:00Z",
    "isActive": true,
    "stockStatus": "in_stock",
    "stock": {
      "quantity": 25,
      "isInStock": true
    }
  }
]
```

#### ✅ **Success Response:**
```json
{
  "message": "Products created successfully",
  "count": 1
}
```

#### ❌ **Error Responses:**

**Empty Array:**
```json
{
  "error": "No products provided"
}
```

**Invalid JSON:**
```json
{
  "error": "json: cannot unmarshal object into Go value of type []models.Product"
}
```

**Database Error:**
```json
{
  "error": "Failed to insert products batch",
  "details": "UNIQUE constraint failed: products._id",
  "batch_start": 0,
  "batch_end": 1000
}
```

## 🔧 **Postman Environment Variables**

Create a new environment with these variables:

| Variable | Value |
|----------|-------|
| `base_url` | `http://69.62.114.202:6000` |
| `store_name` | `zara` (or any store name) |

## 📊 **Test Scripts for Postman**

### **Health Check Test:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has status ok", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql("ok");
});
```

### **Get Products Test:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has count and products", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property("count");
    pm.expect(jsonData).to.have.property("products");
    pm.expect(jsonData.products).to.be.an("array");
});
```

### **Add Products Test:**
```javascript
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("Products created successfully", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.message).to.eql("Products created successfully");
    pm.expect(jsonData.count).to.be.above(0);
});
```

## 🚀 **Quick Start Guide**

1. **Import Collection:**
   - Create new collection in Postman
   - Add the requests above

2. **Set Environment:**
   - Create environment with `base_url` variable
   - Set value to `http://69.62.114.202:6000`

3. **Test Sequence:**
   1. Health Check
   2. Get All Products (should be empty initially)
   3. Add Sample Product
   4. Get All Products (should show added product)
   5. Get Products by Store

## 📈 **Performance Notes**

- **Bulk Insert:** Supports up to 1000 products per batch
- **Pagination:** Default limit is 100 products per page
- **Response Time:** Typically < 100ms for GET requests
- **Database:** PostgreSQL with optimized indexes

## 🔍 **Troubleshooting**

### **Common 400 Errors:**
1. **Empty array:** Send at least one product
2. **Invalid JSON:** Check JSON syntax
3. **Missing required fields:** Ensure all required fields are present
4. **Duplicate IDs:** Use unique `_id` values

### **Log Monitoring:**
```bash
# On server
sudo journalctl -u product-api -f
```

## 🎯 **API Status**
- ✅ **Health Check:** Working
- ✅ **GET Endpoints:** Working  
- ✅ **POST Endpoint:** Working
- ✅ **CORS:** Configured
- ✅ **Database:** Connected
- ✅ **Pagination:** Implemented