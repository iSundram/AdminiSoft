
# AdminiSoftware

AdminiSoftware is a modern, open-source control panel alternative to cPanel/WHM, built with Go backend and Vue.js frontend. It provides comprehensive web hosting management capabilities including account management, DNS control, SSL certificates, email management, and more.

## Features

### Core Features
- **Multi-level Access**: Admin, Reseller, and User roles
- **Account Management**: Create, modify, suspend, and delete hosting accounts
- **Domain Management**: Add domains, subdomains, domain pointers, and DNS management
- **Email Management**: Email accounts, forwarders, autoresponders, and spam protection
- **Database Management**: MySQL, PostgreSQL, and MongoDB support
- **File Management**: Web-based file manager with FTP account management
- **SSL Management**: Let's Encrypt integration and custom SSL certificates
- **Backup System**: Automated and manual backup/restore functionality

### Security Features
- **Two-Factor Authentication**: TOTP support
- **Brute Force Protection**: IP-based login attempt monitoring
- **JWT Authentication**: Secure API access
- **Rate Limiting**: API and login rate limiting
- **SSL/TLS**: Full SSL support with automatic certificate generation

### Monitoring & Statistics
- **Resource Monitoring**: CPU, memory, disk usage tracking
- **Website Statistics**: Visitor analytics and error logs
- **System Health**: Service monitoring and alerts
- **Usage Reports**: Bandwidth and storage usage reports

### Integration Support
- **CloudLinux**: Full CloudLinux OS integration
- **Let's Encrypt**: Automatic SSL certificate generation
- **WordPress Toolkit**: WordPress management and security
- **Imunify360**: Advanced security integration

## Architecture

```
AdminiSoftware/
├── backend/          # Go backend API
│   ├── cmd/         # Application entry points
│   ├── internal/    # Internal packages
│   │   ├── api/     # HTTP handlers and routes
│   │   ├── auth/    # Authentication logic
│   │   ├── config/  # Configuration management
│   │   ├── models/  # Data models
│   │   ├── services/# Business logic
│   │   └── utils/   # Utility functions
│   └── pkg/         # External integrations
├── frontend/        # Vue.js frontend
│   ├── src/
│   │   ├── components/ # Reusable components
│   │   ├── views/     # Page components
│   │   ├── services/  # API services
│   │   └── store/     # State management
└── docs/           # Documentation
```

## Requirements

### System Requirements
- **OS**: Linux (Ubuntu 20.04+ / CentOS 8+ recommended)
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 20GB minimum
- **Network**: Public IP address

### Software Dependencies
- **Go**: 1.21+
- **Node.js**: 18+
- **PostgreSQL**: 13+
- **Redis**: 6+
- **Nginx**: 1.18+

## Quick Installation

1. **Download and extract**:
   ```bash
   wget https://github.com/adminisoftware/adminisoftware/archive/main.zip
   unzip main.zip
   cd AdminiSoftware
   ```

2. **Run setup script**:
   ```bash
   sudo chmod +x scripts/setup.sh
   sudo ./scripts/setup.sh
   ```

3. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Build and run**:
   ```bash
   cd backend
   go mod download
   go build -o adminisoftware cmd/main.go
   
   cd ../frontend
   npm install
   npm run build
   ```

## Configuration

### Environment Variables

```bash
# Database
DATABASE_URL=postgres://user:password@localhost/adminisoftware_db
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Server
PORT=5000
ENVIRONMENT=production

# External Services
LETSENCRYPT_EMAIL=admin@yourdomain.com
CLOUDLINUX_LICENSE_KEY=your-license-key
```

### Database Setup

```sql
-- Create database
CREATE DATABASE adminisoftware_db;
CREATE USER adminisoftware WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE adminisoftware_db TO adminisoftware;
```

## API Documentation

### Authentication

All API requests require authentication via JWT token:

```bash
curl -H "Authorization: Bearer <token>" https://api.adminisoftware.com/api/v1/accounts
```

### Core Endpoints

#### Account Management
- `GET /api/v1/accounts` - List all accounts
- `POST /api/v1/accounts` - Create new account
- `PUT /api/v1/accounts/:id` - Update account
- `DELETE /api/v1/accounts/:id` - Delete account

#### Domain Management
- `GET /api/v1/domains` - List domains
- `POST /api/v1/domains` - Add domain
- `PUT /api/v1/domains/:id` - Update domain
- `DELETE /api/v1/domains/:id` - Remove domain

#### DNS Management
- `GET /api/v1/dns/:domain` - Get DNS records
- `POST /api/v1/dns/:domain` - Add DNS record
- `PUT /api/v1/dns/:domain/:id` - Update DNS record
- `DELETE /api/v1/dns/:domain/:id` - Delete DNS record

## User Interface

### Admin Panel (AdminiCore)
- Dashboard with system overview
- Account management and creation
- Package management
- Server configuration
- Security center
- Monitoring and statistics

### Reseller Panel (AdminiReseller)
- Reseller dashboard
- Client account management
- Resource allocation
- Custom branding options

### User Panel (AdminiPanel)
- User dashboard
- Domain and subdomain management
- Email account management
- File manager
- Database management
- SSL certificate management

## Development

### Backend Development

```bash
cd backend
go mod download
go run cmd/main.go
```

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

### Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm run test
```

## Security

### Best Practices
1. **Regular Updates**: Keep system and dependencies updated
2. **Strong Passwords**: Use complex passwords and 2FA
3. **Firewall**: Configure iptables/ufw properly
4. **SSL**: Use HTTPS everywhere
5. **Backups**: Regular automated backups
6. **Monitoring**: Set up system monitoring

### Security Features
- JWT-based authentication
- Rate limiting on all endpoints
- Brute force protection
- Input validation and sanitization
- SQL injection prevention
- XSS protection

## Monitoring

### System Monitoring
- CPU and memory usage
- Disk space monitoring
- Service health checks
- Network statistics

### Application Monitoring
- API response times
- Error rates
- User activity logs
- Performance metrics

## Backup & Recovery

### Automated Backups
```bash
# Setup cron job for daily backups
0 2 * * * /opt/adminisoftware/scripts/backup.sh
```

### Manual Backup
```bash
sudo /opt/adminisoftware/scripts/backup.sh
```

### Restore Process
```bash
sudo /opt/adminisoftware/scripts/restore.sh /path/to/backup
```

## Troubleshooting

### Common Issues

1. **Service won't start**
   ```bash
   sudo journalctl -u adminisoftware -f
   ```

2. **Database connection issues**
   ```bash
   sudo -u postgres psql -c "\l"
   ```

3. **Permission issues**
   ```bash
   sudo chown -R adminisoftware:adminisoftware /opt/adminisoftware
   ```

### Log Files
- Application: `/var/log/adminisoftware/app.log`
- Nginx: `/var/log/nginx/adminisoftware.log`
- Database: `/var/log/postgresql/`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

AdminiSoftware is released under the MIT License. See [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: https://docs.adminisoftware.com
- **Community**: https://community.adminisoftware.com
- **Issues**: https://github.com/adminisoftware/adminisoftware/issues
- **Email**: support@adminisoftware.com

## Changelog

See [CHANGELOG.md](docs/changelog.md) for version history and updates.
