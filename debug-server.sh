#!/bin/bash

# Sunucu Debug Script
# IP: 69.62.114.202

echo "🔍 Sunucu Debug Kontrolü"
echo "========================"

SERVER_IP="69.62.114.202"
SERVER_USER="root"  # Kullanıcı adınızı buraya yazın

# Renklendirme
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}1. Sunucu erişilebilirlik kontrolü...${NC}"
if ping -c 1 $SERVER_IP > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Sunucu erişilebilir${NC}"
else
    echo -e "${RED}❌ Sunucu erişilemez${NC}"
    exit 1
fi

echo -e "${YELLOW}2. SSH bağlantısı test ediliyor...${NC}"
echo "SSH ile sunucuya bağlanıp kontroller yapılacak..."

# SSH ile sunucuda kontroller yap
ssh $SERVER_USER@$SERVER_IP << 'EOF'
echo "🔍 Sunucu üzerinde kontroller:"
echo "=============================="

echo "1. Product-API servisi durumu:"
sudo systemctl status product-api --no-pager || echo "❌ Servis bulunamadı"

echo -e "\n2. Port 6000 dinleniyor mu:"
sudo netstat -tlnp | grep :6000 || echo "❌ Port 6000 dinlenmiyor"

echo -e "\n3. Uygulama dosyaları mevcut mu:"
ls -la /opt/product-api/ || echo "❌ Uygulama dizini bulunamadı"

echo -e "\n4. PostgreSQL durumu:"
sudo systemctl status postgresql --no-pager || echo "❌ PostgreSQL çalışmıyor"

echo -e "\n5. Nginx durumu:"
sudo systemctl status nginx --no-pager || echo "❌ Nginx çalışmıyor"

echo -e "\n6. Firewall durumu:"
sudo ufw status || echo "❌ UFW durumu kontrol edilemedi"

echo -e "\n7. Son 10 log kaydı:"
sudo journalctl -u product-api -n 10 --no-pager || echo "❌ Log kayıtları okunamadı"

echo -e "\n8. Process listesi (product-api):"
ps aux | grep product-api || echo "❌ Product-API process'i bulunamadı"

echo -e "\n9. Disk kullanımı:"
df -h /opt/product-api/ || echo "❌ Disk bilgisi alınamadı"

echo -e "\n10. Memory kullanımı:"
free -h || echo "❌ Memory bilgisi alınamadı"
EOF

echo -e "${YELLOW}Debug kontrolü tamamlandı.${NC}"