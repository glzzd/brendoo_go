#!/bin/bash

# Ubuntu Sunucusu Deployment Script
# IP: 69.62.114.202, Port: 6000

echo "🚀 Product API Deployment Script"
echo "================================"

# Renklendirme
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Hata durumunda scripti durdur
set -e

# Değişkenler
APP_NAME="product-api"
APP_DIR="/opt/product-api"
SERVICE_NAME="product-api"
PORT=6000

echo -e "${YELLOW}1. Sistem güncellemeleri kontrol ediliyor...${NC}"
sudo apt update

echo -e "${YELLOW}2. Gerekli paketler yükleniyor...${NC}"
sudo apt install -y postgresql postgresql-contrib golang-go git nginx

echo -e "${YELLOW}3. PostgreSQL ayarları yapılandırılıyor...${NC}"
sudo systemctl start postgresql
sudo systemctl enable postgresql

# PostgreSQL kullanıcısı ve veritabanı oluştur
sudo -u postgres psql -c "CREATE USER productuser WITH PASSWORD 'productpass';" || echo "User already exists"
sudo -u postgres psql -c "CREATE DATABASE productdb OWNER productuser;" || echo "Database already exists"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE productdb TO productuser;"

echo -e "${YELLOW}4. Uygulama dizini hazırlanıyor...${NC}"
sudo mkdir -p $APP_DIR
sudo chown $USER:$USER $APP_DIR

echo -e "${YELLOW}5. Uygulama dosyaları kopyalanıyor...${NC}"
# Bu adımda dosyaları sunucuya kopyalamanız gerekecek
# scp -r ./* user@69.62.114.202:/opt/product-api/

echo -e "${YELLOW}6. Go modülleri indiriliyor...${NC}"
cd $APP_DIR
go mod tidy

echo -e "${YELLOW}7. Uygulama derleniyor...${NC}"
go build -o $APP_NAME .

echo -e "${YELLOW}8. Environment dosyası ayarlanıyor...${NC}"
cat > .env << EOF
PORT=6000
DATABASE_URL=postgres://productuser:productpass@localhost:5432/productdb?sslmode=disable
GIN_MODE=release
EOF

echo -e "${YELLOW}9. Systemd service dosyası oluşturuluyor...${NC}"
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null << EOF
[Unit]
Description=Product API Service
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/$APP_NAME
Restart=always
RestartSec=5
Environment=PORT=6000
Environment=DATABASE_URL=postgres://productuser:productpass@localhost:5432/productdb?sslmode=disable
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

echo -e "${YELLOW}10. Service başlatılıyor...${NC}"
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl start $SERVICE_NAME

echo -e "${YELLOW}11. Nginx reverse proxy ayarlanıyor...${NC}"
sudo tee /etc/nginx/sites-available/$SERVICE_NAME > /dev/null << EOF
server {
    listen 80;
    server_name 69.62.114.202;

    location / {
        proxy_pass http://localhost:6000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF

sudo ln -sf /etc/nginx/sites-available/$SERVICE_NAME /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

echo -e "${YELLOW}12. Firewall ayarları...${NC}"
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 6000
sudo ufw --force enable

echo -e "${GREEN}✅ Deployment tamamlandı!${NC}"
echo -e "${GREEN}🌐 API Endpoint: http://69.62.114.202:6000/api/stock/integration/store${NC}"
echo -e "${GREEN}🌐 Nginx Proxy: http://69.62.114.202/api/stock/integration/store${NC}"

echo -e "${YELLOW}Service durumunu kontrol etmek için:${NC}"
echo "sudo systemctl status $SERVICE_NAME"

echo -e "${YELLOW}Logları görmek için:${NC}"
echo "sudo journalctl -u $SERVICE_NAME -f"