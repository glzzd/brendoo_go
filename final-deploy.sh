#!/bin/bash

# Final Deploy Script - Son Hali Deploy Etme
# Tüm güncellemeleri sunucuya yükler

echo "🚀 Final Deploy - Son Hali Yükleniyor"
echo "===================================="

SERVER_IP="69.62.114.202"
SERVER_USER="root"  # Kullanıcı adınızı buraya yazın

# Renklendirme
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}📋 Deploy edilecek dosyalar:${NC}"
echo "  ✓ Güncel kaynak kodlar"
echo "  ✓ Yapılandırma dosyaları"
echo "  ✓ Environment ayarları"
echo "  ✓ Postman collection dokümantasyonu"
echo ""

# 1. Yerel temizlik ve derleme
echo -e "${YELLOW}🧹 1. Yerel temizlik ve derleme...${NC}"
rm -f product-api
go clean -cache
go mod tidy

echo -e "${YELLOW}🔨 2. Linux için cross-compile...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o product-api .

if [ ! -f "product-api" ]; then
    echo -e "${RED}❌ Derleme başarısız!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Derleme başarılı!${NC}"

# 2. Dosyaları sunucuya yükle
echo -e "${YELLOW}📤 3. Dosyalar sunucuya yükleniyor...${NC}"

# Ana uygulama dosyaları
scp product-api $SERVER_USER@$SERVER_IP:/tmp/
scp .env.production $SERVER_USER@$SERVER_IP:/tmp/

# Dokümantasyon dosyaları
scp POSTMAN-COLLECTION.md $SERVER_USER@$SERVER_IP:/tmp/
scp README-DEPLOYMENT.md $SERVER_USER@$SERVER_IP:/tmp/

echo -e "${GREEN}✅ Dosyalar yüklendi!${NC}"

# 3. Sunucuda güncelleme
echo -e "${YELLOW}🔄 4. Sunucuda güncelleme yapılıyor...${NC}"

ssh $SERVER_USER@$SERVER_IP << 'EOF'
echo "🔧 Sunucuda güncelleme başlatılıyor..."

# Mevcut servisi durdur
echo "🛑 Mevcut servisi durduruyor..."
sudo systemctl stop product-api

# Backup al (isteğe bağlı)
echo "💾 Backup alınıyor..."
sudo cp /opt/product-api/product-api /opt/product-api/product-api.backup.$(date +%Y%m%d_%H%M%S) 2>/dev/null || echo "ℹ️  Backup alınamadı (ilk kurulum olabilir)"

# Yeni dosyaları yerleştir
echo "📦 Yeni dosyalar yerleştiriliyor..."
sudo mv /tmp/product-api /opt/product-api/
sudo mv /tmp/.env.production /opt/product-api/.env
sudo chmod +x /opt/product-api/product-api

# Dokümantasyon dosyalarını yerleştir
sudo mv /tmp/POSTMAN-COLLECTION.md /opt/product-api/ 2>/dev/null || echo "ℹ️  Postman collection dosyası bulunamadı"
sudo mv /tmp/README-DEPLOYMENT.md /opt/product-api/ 2>/dev/null || echo "ℹ️  README dosyası bulunamadı"

# Dosya sahipliklerini ayarla
sudo chown -R root:root /opt/product-api/

# Servisi başlat
echo "🚀 Servis başlatılıyor..."
sudo systemctl start product-api

# Birkaç saniye bekle
sleep 5

# Durum kontrolü
echo "🔍 Servis durumu kontrol ediliyor..."
if sudo systemctl is-active --quiet product-api; then
    echo "✅ Servis başarıyla başlatıldı!"
    
    # Port kontrolü
    if sudo netstat -tlnp | grep :8080 > /dev/null; then
        echo "✅ Port 8080 dinleniyor!"
    else
        echo "⚠️  Port 8080 henüz dinlenmiyor, birkaç saniye bekleyin..."
    fi
    
    # API testi
    echo "🧪 API testi yapılıyor..."
    sleep 2
    if curl -s http://localhost:8080/health | grep -q "ok"; then
        echo "✅ API başarıyla çalışıyor!"
    else
        echo "❌ API testi başarısız!"
    fi
    
else
    echo "❌ Servis başlatılamadı!"
    echo "📋 Servis logları:"
    sudo journalctl -u product-api -n 10 --no-pager
fi

echo ""
echo "📊 Servis bilgileri:"
echo "==================="
sudo systemctl status product-api --no-pager -l
EOF

# 4. Final test
echo -e "${YELLOW}🧪 5. Final test yapılıyor...${NC}"
sleep 3

echo -e "${BLUE}Health Check:${NC}"
curl -s http://69.62.114.202:8080/health || echo -e "${RED}❌ Health check başarısız${NC}"

echo -e "\n${BLUE}API Test:${NC}"
curl -s "http://69.62.114.202:8080/api/stock/integration/store" | head -c 100 || echo -e "${RED}❌ API testi başarısız${NC}"

echo -e "\n\n${GREEN}🎉 Deploy tamamlandı!${NC}"
echo -e "${BLUE}📋 Erişim bilgileri:${NC}"
echo "  🌐 API Base URL: http://69.62.114.202:8080"
echo "  ❤️  Health Check: http://69.62.114.202:8080/health"
echo "  📦 Products API: http://69.62.114.202:8080/api/stock/integration/store"
echo "  📋 Postman Collection: /opt/product-api/POSTMAN-COLLECTION.md"
echo ""
echo -e "${YELLOW}📖 Dokümantasyon:${NC}"
echo "  • POSTMAN-COLLECTION.md - Postman test collection"
echo "  • README-DEPLOYMENT.md - Deployment rehberi"
echo ""
echo -e "${BLUE}🔧 Servis yönetimi:${NC}"
echo "  • sudo systemctl status product-api"
echo "  • sudo systemctl restart product-api"
echo "  • sudo journalctl -u product-api -f"