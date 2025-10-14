#!/bin/bash

# Hızlı Çözüm Script
# API'yi sunucuda başlatmak için

echo "🚀 Hızlı Çözüm - API Başlatma"
echo "============================="

SERVER_IP="69.62.114.202"
SERVER_USER="root"  # Kullanıcı adınızı buraya yazın

# Renklendirme
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Sunucuya bağlanıp API'yi başlatıyoruz...${NC}"

# SSH ile sunucuda işlemleri yap
ssh $SERVER_USER@$SERVER_IP << 'EOF'
echo "🔧 Sunucuda işlemler başlatılıyor..."

# 1. Uygulama dizinine git
cd /opt/product-api || { echo "❌ Uygulama dizini bulunamadı"; exit 1; }

echo "✅ Uygulama dizinindeyiz: $(pwd)"

# 2. Eğer servis varsa durdur
echo "🛑 Mevcut servisi durduruyor..."
sudo systemctl stop product-api 2>/dev/null || echo "ℹ️  Servis zaten durdurulmuş"

# 3. Eğer process çalışıyorsa öldür
echo "🔪 Çalışan process'leri temizliyor..."
sudo pkill -f product-api 2>/dev/null || echo "ℹ️  Çalışan process bulunamadı"

# 4. Port 6000'i kullanan process'i öldür
echo "🔌 Port 6000'i temizliyor..."
sudo lsof -ti:6000 | xargs sudo kill -9 2>/dev/null || echo "ℹ️  Port 6000 zaten boş"

# 5. Environment dosyasını oluştur
echo "📝 Environment dosyası oluşturuluyor..."
cat > .env << 'ENVEOF'
PORT=6000
DATABASE_URL=postgres://productuser:productpass@localhost:5432/productdb?sslmode=disable
GIN_MODE=release
ENVEOF

# 6. PostgreSQL'in çalıştığından emin ol
echo "🐘 PostgreSQL kontrol ediliyor..."
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 7. Veritabanını oluştur (eğer yoksa)
echo "🗄️  Veritabanı kontrol ediliyor..."
sudo -u postgres psql -c "CREATE USER productuser WITH PASSWORD 'productpass';" 2>/dev/null || echo "ℹ️  User zaten mevcut"
sudo -u postgres psql -c "CREATE DATABASE productdb OWNER productuser;" 2>/dev/null || echo "ℹ️  Database zaten mevcut"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE productdb TO productuser;" 2>/dev/null

# 8. Go modüllerini güncelle
echo "📦 Go modülleri güncelleniyor..."
go mod tidy

# 9. Uygulamayı derle
echo "🔨 Uygulama derleniyor..."
go build -o product-api .

# 10. Uygulamayı arka planda başlat
echo "🚀 Uygulama başlatılıyor..."
nohup ./product-api > app.log 2>&1 &

# 11. Birkaç saniye bekle
sleep 3

# 12. Çalışıp çalışmadığını kontrol et
echo "🔍 Uygulama durumu kontrol ediliyor..."
if ps aux | grep -v grep | grep product-api > /dev/null; then
    echo "✅ Uygulama başarıyla başlatıldı!"
    
    # Port kontrolü
    if netstat -tlnp | grep :6000 > /dev/null; then
        echo "✅ Port 6000 dinleniyor!"
    else
        echo "⚠️  Port 6000 henüz dinlenmiyor, birkaç saniye daha bekleyin..."
    fi
    
    # Log'ları göster
    echo "📋 Son log kayıtları:"
    tail -10 app.log
else
    echo "❌ Uygulama başlatılamadı!"
    echo "📋 Hata logları:"
    tail -20 app.log
fi

# 13. Firewall kontrolü
echo "🔥 Firewall kontrol ediliyor..."
sudo ufw allow 6000 2>/dev/null || echo "ℹ️  Firewall kuralı zaten mevcut"

echo "🎉 İşlemler tamamlandı!"
echo "🌐 Test URL: http://69.62.114.202:6000/health"
EOF

echo -e "${GREEN}Hızlı çözüm scripti tamamlandı!${NC}"
echo -e "${YELLOW}Şimdi test edin: curl http://69.62.114.202:6000/health${NC}"