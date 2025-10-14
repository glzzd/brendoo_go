#!/bin/bash

# Dosyaları Ubuntu sunucusuna yükleme scripti
# IP: 69.62.114.202

echo "📤 Dosyalar sunucuya yükleniyor..."

# Değişkenler
SERVER_IP="69.62.114.202"
SERVER_USER="root"  
APP_DIR="/opt/product-api"

# Renklendirme
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}1. Sunucuda uygulama dizini oluşturuluyor...${NC}"
ssh $SERVER_USER@$SERVER_IP "mkdir -p $APP_DIR"

echo -e "${YELLOW}2. Uygulama dosyaları yükleniyor...${NC}"
scp -r ./* $SERVER_USER@$SERVER_IP:$APP_DIR/

echo -e "${YELLOW}3. Deployment scripti çalıştırılıyor...${NC}"
ssh $SERVER_USER@$SERVER_IP "cd $APP_DIR && chmod +x deploy.sh && ./deploy.sh"

echo -e "${GREEN}✅ Yükleme tamamlandı!${NC}"
echo -e "${GREEN}🌐 API Endpoint: http://$SERVER_IP:6000/api/stock/integration/store${NC}"