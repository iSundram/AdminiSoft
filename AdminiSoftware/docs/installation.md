
# AdminiSoftware Installation Guide

This guide will walk you through the installation process of AdminiSoftware, a comprehensive control panel system for web hosting management.

## Prerequisites

Before installing AdminiSoftware, ensure your system meets the following requirements:

### System Requirements
- **Operating System**: Linux (Ubuntu 20.04+ recommended, CentOS 8+, Debian 10+)
- **CPU**: 2+ cores
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 20GB minimum, SSD recommended
- **Network**: Internet connection for package downloads

### Software Requirements
- **Go**: 1.21 or higher
- **Node.js**: 18.0 or higher
- **npm**: 8.0 or higher
- **PostgreSQL**: 12 or higher
- **Redis**: 6.0 or higher
- **Nginx**: 1.18 or higher (optional, for production)

## Installation Methods

### Method 1: Quick Installation Script

The easiest way to install AdminiSoftware is using our automated installation script.

```bash
# Download and run the installation script
curl -fsSL https://raw.githubusercontent.com/yourusername/AdminiSoftware/main/scripts/install.sh | sudo bash

# Or download first and review before running
wget https://raw.githubusercontent.com/yourusername/AdminiSoftware/main/scripts/install.sh
chmod +x install.sh
sudo ./install.sh
```

### Method 2: Manual Installation

#### Step 1: Install Dependencies

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y postgresql redis-server nginx supervisor curl wget git build-essential
```

**CentOS/RHEL:**
```bash
sudo yum update -y
sudo yum install -y postgresql-server redis nginx supervisor curl wget git gcc gcc-c++ make
```

#### Step 2: Install Go

```bash
cd /tmp
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
source /etc/profile
```

#### Step 3: Install Node.js

```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

#### Step 4: Setup Database

**PostgreSQL Setup:**
```bash
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres createuser adminisoftware
sudo -u postgres createdb adminisoftware_db
sudo -u postgres psql -c "ALTER USER adminisoftware WITH PASSWORD 'your_secure_password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE adminisoftware_db TO adminisoftware;"
```

**Redis Setup:**
```bash
sudo systemctl start redis
sudo systemctl enable redis
```

#### Step 5: Create Application User

```bash
sudo useradd -m -s /bin/bash adminisoftware
sudo usermod -aG sudo adminisoftware
```

#### Step 6: Download and Setup AdminiSoftware

```bash
sudo mkdir -p /opt/adminisoftware
cd /opt/adminisoftware

# Clone repository (replace with actual repository URL)
sudo git clone https://github.com/yourusername/AdminiSoftware.git .
sudo chown -R adminisoftware:adminisoftware /opt/adminisoftware

# Switch to adminisoftware user
sudo -u adminisoftware -s

# Build backend
cd /opt/adminisoftware/backend
go mod tidy
go build -o bin/adminisoftware cmd/main.go

# Build frontend
cd /opt/adminisoftware/frontend
npm install
npm run build
```

#### Step 7: Configuration

```bash
# Create configuration directory
sudo mkdir -p /etc/adminisoftware
sudo mkdir -p /var/log/adminisoftware
sudo chown adminisoftware:adminisoftware /var/log/adminisoftware

# Copy environment configuration
sudo cp /opt/adminisoftware/.env.example /etc/adminisoftware/config.env

# Edit configuration
sudo nano /etc/adminisoftware/config.env
```

**Edit the configuration file with your settings:**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=adminisoftware_db
DB_USER=adminisoftware
DB_PASSWORD=your_secure_password

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Application Configuration
APP_PORT=5000
APP_SECRET_KEY=your_secret_key_here
JWT_SECRET=your_jwt_secret_here

# Email Configuration (optional)
SMTP_HOST=localhost
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=noreply@yourdomain.com
```

#### Step 8: Setup Systemd Service

```bash
sudo tee /etc/systemd/system/adminisoftware.service > /dev/null <<EOF
[Unit]
Description=AdminiSoftware Control Panel
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=adminisoftware
Group=adminisoftware
WorkingDirectory=/opt/adminisoftware/backend
ExecStart=/opt/adminisoftware/backend/bin/adminisoftware
Restart=always
RestartSec=10
Environment=ENVIRONMENT=production
EnvironmentFile=/etc/adminisoftware/config.env

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable adminisoftware
```

#### Step 9: Setup Nginx (Production)

```bash
sudo tee /etc/nginx/sites-available/adminisoftware > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com;
    
    # Frontend
    location / {
        root /opt/adminisoftware/frontend/dist;
        try_files \$uri \$uri/ /index.html;
        
        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header Referrer-Policy "no-referrer-when-downgrade" always;
        add_header Content-Security-Policy "default-src 'self'" always;
    }
    
    # Backend API
    location /api/ {
        proxy_pass http://127.0.0.1:5000;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Increase timeout for long operations
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # WebSocket
    location /ws {
        proxy_pass http://127.0.0.1:5000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    # Static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
EOF

sudo ln -s /etc/nginx/sites-available/adminisoftware /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

#### Step 10: Setup SSL (Optional but Recommended)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal
sudo systemctl enable certbot.timer
```

#### Step 11: Start Services

```bash
sudo systemctl start adminisoftware
sudo systemctl start nginx
```

#### Step 12: Initial Setup

```bash
# Run database migrations
cd /opt/adminisoftware/backend
sudo -u adminisoftware ./bin/adminisoftware --migrate

# Create admin user (optional, can be done via web interface)
sudo -u adminisoftware ./bin/adminisoftware --create-admin
```

## Verification

### Check Service Status

```bash
# Check AdminiSoftware service
sudo systemctl status adminisoftware

# Check database connection
sudo systemctl status postgresql

# Check Redis
sudo systemctl status redis

# Check Nginx
sudo systemctl status nginx
```

### Test Installation

1. Open your web browser and navigate to your server's IP address or domain
2. You should see the AdminiSoftware login page
3. Use the default admin credentials:
   - Email: `admin@adminisoftware.com`
   - Password: `admin123`

### Change Default Credentials

**Important**: Immediately change the default admin password after first login.

## Post-Installation

### Security Hardening

1. **Change default passwords**
2. **Enable firewall**:
   ```bash
   sudo ufw allow ssh
   sudo ufw allow http
   sudo ufw allow https
   sudo ufw enable
   ```
3. **Setup fail2ban**:
   ```bash
   sudo apt install fail2ban
   sudo systemctl enable fail2ban
   sudo systemctl start fail2ban
   ```

### Backup Setup

```bash
# Create backup directory
sudo mkdir -p /var/backups/adminisoftware

# Setup automated backups (add to crontab)
sudo crontab -e
# Add: 0 2 * * * /opt/adminisoftware/scripts/backup.sh
```

### Log Rotation

```bash
sudo tee /etc/logrotate.d/adminisoftware > /dev/null <<EOF
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
```

## Troubleshooting

### Common Issues

1. **Service fails to start**:
   ```bash
   sudo journalctl -u adminisoftware -f
   ```

2. **Database connection issues**:
   ```bash
   sudo -u postgres psql -c "\l"
   ```

3. **Permission issues**:
   ```bash
   sudo chown -R adminisoftware:adminisoftware /opt/adminisoftware
   ```

4. **Port conflicts**:
   ```bash
   sudo netstat -tulpn | grep :5000
   ```

### Log Files

- Application logs: `/var/log/adminisoftware/app.log`
- Nginx logs: `/var/log/nginx/access.log`, `/var/log/nginx/error.log`
- PostgreSQL logs: `/var/log/postgresql/`
- System logs: `sudo journalctl -u adminisoftware`

## Getting Help

- Documentation: [docs/](../docs/)
- GitHub Issues: [GitHub Issues](https://github.com/yourusername/AdminiSoftware/issues)
- Support: support@adminisoftware.com

## Next Steps

After successful installation:

1. Review the [Configuration Guide](configuration.md)
2. Read the [User Guide](user-guide.md)
3. Check the [Admin Guide](admin-guide.md)
4. Setup monitoring and alerts
5. Configure automated backups
6. Review security settings

---

**Note**: Replace `your-domain.com`, `your_secure_password`, and other placeholders with your actual values during installation.
