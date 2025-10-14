# 🚀 Ubuntu Sunucusu Deployment Rehberi

Bu rehber, Product API'yi Ubuntu sunucusunda (IP: 69.62.114.202) port 6000'de yayımlamak için hazırlanmıştır.

## 📋 Ön Gereksinimler

- Ubuntu sunucusu (IP: 69.62.114.202)
- SSH erişimi
- Root veya sudo yetkisi

## 🔧 Deployment Adımları

### 1. Dosyaları Sunucuya Yükleme

```bash
# Upload scriptini çalıştır (kullanıcı adını düzenleyin)
./upload-to-server.sh
```

**Önemli:** `upload-to-server.sh` dosyasındaki `SERVER_USER` değişkenini kendi kullanıcı adınızla değiştirin.

### 2. Manuel Deployment

Eğer script çalışmazsa, manuel olarak:

```bash
# Dosyaları kopyala
scp -r ./* user@69.62.114.202:/opt/product-api/

# Sunucuya bağlan
ssh user@69.62.114.202

# Deployment scriptini çalıştır
cd /opt/product-api
chmod +x deploy.sh
./deploy.sh
```

## 🌐 API Endpoint'leri

Deployment sonrası aşağıdaki endpoint'ler kullanılabilir olacak:

### Tüm Ürünleri Getir
```
GET http://69.62.114.202:6000/api/stock/integration/store
```

### Mağazaya Göre Ürünleri Getir
```
GET http://69.62.114.202:6000/api/stock/integration/store?store=zara
GET http://69.62.114.202:6000/api/stock/integration/store?store=nike
GET http://69.62.114.202:6000/api/stock/integration/store?store=adidas
```

### Yeni Ürün Ekleme
```
POST http://69.62.114.202:6000/api/stock/add
Content-Type: application/json

[
  {
    "id": "product-001",
    "name": "Ürün Adı",
    "brand": "Marka",
    "price": 99.99,
    "currency": "TRY",
    "store": "zara"
    // ... diğer alanlar
  }
]
```

### Health Check
```
GET http://69.62.114.202:6000/health
```

## 🔍 Servis Yönetimi

### Servis Durumunu Kontrol Etme
```bash
sudo systemctl status product-api
```

### Servis Loglarını Görme
```bash
sudo journalctl -u product-api -f
```

### Servisi Yeniden Başlatma
```bash
sudo systemctl restart product-api
```

### Servisi Durdurma
```bash
sudo systemctl stop product-api
```

## 🔧 Yapılandırma Dosyaları

### Environment Variables
- Port: 6000
- Database: PostgreSQL (productdb)
- User: productuser
- Password: productpass

### Nginx Reverse Proxy
Nginx ile port 80'den port 6000'e yönlendirme yapılmıştır:
```
http://69.62.114.202/api/stock/integration/store
```

## 🛡️ Güvenlik

- Firewall: Port 22, 80, 6000 açık
- CORS: Belirtilen origin'lere izin verilmiş
- PostgreSQL: Yerel bağlantı

## 🐛 Sorun Giderme

### API Çalışmıyor
```bash
# Servis durumunu kontrol et
sudo systemctl status product-api

# Logları kontrol et
sudo journalctl -u product-api -n 50

# Portu kontrol et
sudo netstat -tlnp | grep 6000
```

### Veritabanı Bağlantı Sorunu
```bash
# PostgreSQL durumunu kontrol et
sudo systemctl status postgresql

# Veritabanına bağlan
sudo -u postgres psql -d productdb
```

### Nginx Sorunu
```bash
# Nginx durumunu kontrol et
sudo systemctl status nginx

# Nginx yapılandırmasını test et
sudo nginx -t
```

## 📊 Test Etme

### Postman ile Test
1. Base URL: `http://69.62.114.202:6000`
2. Endpoint'leri yukarıdaki örneklere göre test edin

### cURL ile Test
```bash
# Tüm ürünleri getir
curl -X GET "http://69.62.114.202:6000/api/stock/integration/store"

# Zara ürünlerini getir
curl -X GET "http://69.62.114.202:6000/api/stock/integration/store?store=zara"

# Health check
curl -X GET "http://69.62.114.202:6000/health"
```

## 🔄 Güncelleme

Kod değişikliklerinden sonra:

```bash
# Yeni kodu yükle
scp -r ./* user@69.62.114.202:/opt/product-api/

# Sunucuda yeniden derle ve başlat
ssh user@69.62.114.202
cd /opt/product-api
go build -o product-api .
sudo systemctl restart product-api
```

---

**🎉 Deployment tamamlandı! API artık http://69.62.114.202:6000 adresinde çalışıyor.**