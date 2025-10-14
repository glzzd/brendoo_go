#!/bin/bash

# Düzeltme ve Deploy Script
# Tüm sorunları çözerek uygulamayı başlatır

echo "🔧 Sorun Çözme ve Deploy İşlemi"
echo "==============================="

SERVER_IP="69.62.114.202"
SERVER_USER="root"  # Kullanıcı adınızı buraya yazın

# Renklendirme
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}1. Yerel olarak temiz derleme yapılıyor...${NC}"

# Yerel temizlik ve derleme
echo "🧹 Yerel temizlik..."
rm -f product-api
go clean -cache
go clean -modcache
go mod tidy

echo "🔨 Yerel derleme..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o product-api .

if [ ! -f "product-api" ]; then
    echo -e "${RED}❌ Yerel derleme başarısız!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Yerel derleme başarılı!${NC}"

echo -e "${YELLOW}2. Dosyalar sunucuya yükleniyor...${NC}"

# Sunucuya dosyaları yükle
scp product-api $SERVER_USER@$SERVER_IP:/tmp/
scp .env.production $SERVER_USER@$SERVER_IP:/tmp/

echo -e "${YELLOW}3. Sunucuda kurulum yapılıyor...${NC}"

# SSH ile sunucuda işlemleri yap
ssh $SERVER_USER@$SERVER_IP << 'EOF'
echo "🔧 Sunucuda kurulum başlatılıyor..."

# Mevcut process'leri temizle
echo "🛑 Mevcut process'leri temizliyor..."
sudo pkill -f product-api 2>/dev/null || echo "ℹ️  Çalışan process bulunamadı"
sudo lsof -ti:6000 | xargs sudo kill -9 2>/dev/null || echo "ℹ️  Port 6000 zaten boş"

# Uygulama dizinini oluştur
echo "📁 Uygulama dizini hazırlanıyor..."
sudo mkdir -p /opt/product-api
sudo chown $USER:$USER /opt/product-api

# Dosyaları taşı
echo "📦 Dosyalar taşınıyor..."
sudo mv /tmp/product-api /opt/product-api/
sudo mv /tmp/.env.production /opt/product-api/.env
sudo chmod +x /opt/product-api/product-api

# Dizine geç
cd /opt/product-api

# PostgreSQL'i başlat ve yapılandır
echo "🐘 PostgreSQL yapılandırılıyor..."
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Veritabanını oluştur
echo "🗄️  Veritabanı oluşturuluyor..."
sudo -u postgres psql -c "DROP DATABASE IF EXISTS productdb;" 2>/dev/null
sudo -u postgres psql -c "DROP USER IF EXISTS productuser;" 2>/dev/null
sudo -u postgres psql -c "CREATE USER productuser WITH PASSWORD 'productpass';"
sudo -u postgres psql -c "CREATE DATABASE productdb OWNER productuser;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE productdb TO productuser;"

# Firewall ayarları
echo "🔥 Firewall yapılandırılıyor..."
sudo ufw allow 6000/tcp
sudo ufw --force enable

# Systemd servis dosyası oluştur
echo "⚙️  Systemd servisi oluşturuluyor..."
sudo tee /etc/systemd/system/product-api.service > /dev/null << 'SERVICEEOF'
[Unit]
Description=Product API Service
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/product-api
ExecStart=/opt/product-api/product-api
Restart=always
RestartSec=5
Environment=PORT=6000
Environment=DATABASE_URL=postgres://productuser:productpass@localhost:5432/productdb?sslmode=disable
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
SERVICEEOF

# Systemd'yi yenile ve servisi başlat
echo "🚀 Servis başlatılıyor..."
sudo systemctl daemon-reload
sudo systemctl enable product-api
sudo systemctl start product-api

# Birkaç saniye bekle
sleep 5

# Durum kontrolü
echo "🔍 Servis durumu kontrol ediliyor..."
sudo systemctl status product-api --no-pager

# Port kontrolü
echo "🔌 Port kontrolü..."
if sudo netstat -tlnp | grep :6000; then
    echo "✅ Port 6000 başarıyla dinleniyor!"
else
    echo "❌ Port 6000 dinlenmiyor!"
    echo "📋 Servis logları:"
    sudo journalctl -u product-api -n 20 --no-pager
fi

# Test
echo "🧪 API testi..."
sleep 2
curl -s http://localhost:6000/health || echo "❌ Health check başarısız"

echo "🎉 Kurulum tamamlandı!"
echo "🌐 Test URL: http://69.62.114.202:6000/health"
EOF

echo -e "${GREEN}Deploy işlemi tamamlandı!${NC}"
echo -e "${YELLOW}Test için: curl http://69.62.114.202:6000/health${NC}"