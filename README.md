# Product API - High Performance Go Backend

Bu proje, milyonlarca ürün verisini işleyebilen yüksek performanslı bir Go backend API'sidir. PostgreSQL veritabanı kullanır ve 3 ana API endpoint'i sunar.

## 🚀 Özellikler

- **Yüksek Performans**: Milyonlarca kayıt için optimize edilmiş bulk insert işlemleri
- **Hızlı Sorgular**: İndekslenmiş veritabanı sorguları ile saniyeler içinde sonuç
- **Pagination**: Büyük veri setleri için sayfalama desteği
- **Store Filtreleme**: Mağaza bazında hızlı filtreleme
- **Docker Desteği**: Kolay deployment için Docker ve Docker Compose

## 📋 API Endpoints

### 1. Bulk Product Insert
```
POST /api/v1/products
```
Milyonlarca ürün verisini toplu olarak ekler. Batch işleme ile optimize edilmiştir.

**Request Body**: Array of Product objects

### 2. Get All Products
```
GET /api/v1/products?page=1&limit=100
```
Tüm aktif ürünleri sayfalama ile getirir.

**Query Parameters**:
- `page`: Sayfa numarası (default: 1)
- `limit`: Sayfa başına kayıt sayısı (default: 100, max: 1000)

### 3. Get Products by Store
```
GET /api/v1/products/store/{store}?page=1&limit=100
```
Belirtilen mağazaya ait ürünleri getirir.

**Path Parameters**:
- `store`: Mağaza adı (örn: "zara", "h&m")

## 🛠️ Kurulum

### Gereksinimler
- Go 1.25+
- PostgreSQL 15+
- Docker & Docker Compose (opsiyonel)

### 1. Proje Klonlama
```bash
git clone <repository-url>
cd product-api
```

### 2. Dependencies Yükleme
```bash
go mod download
```

### 3. Environment Variables
```bash
cp .env.example .env
# .env dosyasını düzenleyin
```

### 4. Docker ile Çalıştırma (Önerilen)
```bash
docker-compose up -d
```

### 5. Manuel Çalıştırma
```bash
# PostgreSQL'i başlatın
# Environment variables'ları ayarlayın
go run main.go
```

## 🏗️ Proje Yapısı

```
├── main.go                 # Ana uygulama dosyası
├── config/
│   └── config.go          # Konfigürasyon yönetimi
├── database/
│   └── database.go        # Veritabanı bağlantısı ve migration
├── models/
│   └── product.go         # Product model tanımları
├── handlers/
│   └── product_handler.go # API handler'ları
├── routes/
│   └── routes.go          # Route tanımları
├── docker-compose.yml     # Docker Compose konfigürasyonu
├── Dockerfile            # Docker image tanımı
└── init.sql              # PostgreSQL başlangıç scripti
```

## 📊 Performans Optimizasyonları

### Veritabanı İndeksleri
- `store` ve `is_active` için composite index
- `brand`, `category`, `stock_status` için ayrı indexler
- `created_at` için time-based index
- Aktif ürünler için partial index

### Bulk Insert Optimizasyonu
- 1000'lik batch'ler halinde işleme
- GORM'un `CreateInBatches` metodunu kullanma
- Memory-efficient processing

### Connection Pool
- Maksimum 100 açık bağlantı
- 10 idle bağlantı
- Prepared statements kullanımı

## 🧪 Test Etme

### Health Check
```bash
curl http://localhost:8080/health
```

### Bulk Insert Test
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '[{"_id":"test1","name":"Test Product","brand":"test","price":100,"currency":"USD","store":"test","isActive":true}]'
```

### Get Products Test
```bash
curl http://localhost:8080/api/v1/products?page=1&limit=10
```

### Get by Store Test
```bash
curl http://localhost:8080/api/v1/products/store/zara?page=1&limit=10
```

## 🔧 Konfigürasyon

### Environment Variables
- `DATABASE_URL`: PostgreSQL bağlantı string'i
- `PORT`: Sunucu port'u (default: 8080)

### PostgreSQL Ayarları
- `max_connections`: 200
- `shared_buffers`: 256MB
- `effective_cache_size`: 1GB

## 📈 Monitoring

API performansını izlemek için:
- PostgreSQL `pg_stat_statements` extension'ı aktif
- Gin framework'ün built-in logging'i
- Health check endpoint'i

## 🚀 Production Deployment

1. Environment variables'ları production değerleri ile ayarlayın
2. PostgreSQL için SSL bağlantısını aktifleştirin
3. Load balancer arkasında multiple instance çalıştırın
4. Database connection pool'unu ihtiyaca göre ayarlayın

## 📝 Lisans

Bu proje MIT lisansı altında lisanslanmıştır.