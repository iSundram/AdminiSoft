
#!/bin/bash

# AdminiSoftware Setup Script
echo "=== AdminiSoftware Setup ==="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

# Update system
echo "Updating system packages..."
apt-get update && apt-get upgrade -y

# Install required packages
echo "Installing required packages..."
apt-get install -y \
    postgresql \
    redis-server \
    nginx \
    supervisor \
    curl \
    wget \
    git \
    build-essential

# Install Go
echo "Installing Go..."
cd /tmp
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile

# Install Node.js
echo "Installing Node.js..."
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt-get install -y nodejs

# Setup PostgreSQL
echo "Setting up PostgreSQL..."
sudo -u postgres createuser adminisoftware
sudo -u postgres createdb adminisoftware_db
sudo -u postgres psql -c "ALTER USER adminisoftware WITH PASSWORD 'adminisoftware123';"

# Setup Redis
echo "Setting up Redis..."
systemctl enable redis-server
systemctl start redis-server

# Create application user
echo "Creating application user..."
useradd -m -s /bin/bash adminisoftware
usermod -aG sudo adminisoftware

# Create application directories
echo "Creating application directories..."
mkdir -p /opt/adminisoftware
mkdir -p /var/log/adminisoftware
mkdir -p /etc/adminisoftware

# Set permissions
chown -R adminisoftware:adminisoftware /opt/adminisoftware
chown -R adminisoftware:adminisoftware /var/log/adminisoftware

# Setup systemd service
echo "Setting up systemd service..."
cat > /etc/systemd/system/adminisoftware.service << 'EOF'
[Unit]
Description=AdminiSoftware Control Panel
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=adminisoftware
Group=adminisoftware
WorkingDirectory=/opt/adminisoftware
ExecStart=/opt/adminisoftware/bin/adminisoftware
Restart=always
RestartSec=10
Environment=ENVIRONMENT=production

[Install]
WantedBy=multi-user.target
EOF

# Setup Nginx
echo "Setting up Nginx..."
cat > /etc/nginx/sites-available/adminisoftware << 'EOF'
server {
    listen 80;
    server_name _;
    
    # Frontend
    location / {
        root /opt/adminisoftware/frontend/dist;
        try_files $uri $uri/ /index.html;
    }
    
    # Backend API
    location /api/ {
        proxy_pass http://127.0.0.1:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # WebSocket
    location /ws {
        proxy_pass http://127.0.0.1:5000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
EOF

ln -sf /etc/nginx/sites-available/adminisoftware /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Setup logrotate
echo "Setting up log rotation..."
cat > /etc/logrotate.d/adminisoftware << 'EOF'
/var/log/adminisoftware/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 0644 adminisoftware adminisoftware
    postrotate
        systemctl reload adminisoftware
    endscript
}
EOF

# Setup environment file
echo "Creating environment configuration..."
cat > /etc/adminisoftware/config.env << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=adminisoftware_db
DB_USER=adminisoftware
DB_PASSWORD=adminisoftware123

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Application Configuration
APP_PORT=5000
APP_SECRET_KEY=change_this_secret_key
JWT_SECRET=change_this_jwt_secret

# Email Configuration
SMTP_HOST=localhost
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=noreply@example.com

# Security Configuration
ENABLE_2FA=true
RATE_LIMIT_ENABLED=true
BRUTE_FORCE_PROTECTION=true
EOF

# Enable services
echo "Enabling services..."
systemctl enable postgresql
systemctl enable redis-server
systemctl enable nginx
systemctl enable adminisoftware

echo "=== Setup Complete ==="
echo "Please:"
echo "1. Update configuration in /etc/adminisoftware/config.env"
echo "2. Deploy the application to /opt/adminisoftware"
echo "3. Run: systemctl start adminisoftware"
echo "4. Run: systemctl start nginx"
