
# AdminiSoftware

A comprehensive control panel software similar to cPanel, built with Go backend and Vue.js frontend. AdminiSoftware provides three distinct panels: AdminiCore (system administration), AdminiReseller (reseller management), and AdminiPanel (end-user hosting management).

## 🚀 Features

### AdminiCore (System Administrator Panel)
- **Server Configuration**: Web servers, database servers, PHP selector, service management
- **Security Center**: 2FA, brute force protection, ModSecurity, CSF, Imunify360
- **Account Management**: Create, modify, suspend accounts, bandwidth monitoring
- **Package Management**: Hosting packages with resource limits and features
- **DNS Management**: Zone management, email routing, domain parking
- **Email Management**: Mail delivery, queue management, spam protection
- **SSL Management**: Certificate generation, installation, automatic SSL
- **Backup Management**: Scheduled backups, file restoration, remote backups
- **Clustering**: DNS clustering, remote access management
- **Monitoring**: System stats, load averages, service monitoring, alerts
- **Integrations**: CloudLinux, WordPress Toolkit, plugin management

### AdminiReseller (Reseller Panel)
- **Account Management**: Create and manage customer accounts
- **Package Management**: Create hosting packages for customers
- **Resource Management**: IP assignment, resource allocation monitoring
- **Communication**: Message users, announcements, support tickets
- **Branding**: Custom logos, themes, and links
- **Tools**: Terminal access, system information, API access

### AdminiPanel (End-User Panel)
- **Domain Management**: Subdomains, redirects, DNS zone editor, error pages
- **File Management**: Web-based file manager, FTP accounts, disk usage
- **Email Management**: Email accounts, forwarders, filters, webmail
- **Database Management**: MySQL, PostgreSQL, MongoDB databases
- **Application Management**: One-click installers, WordPress manager, Node.js apps
- **SSL Management**: Free SSL certificates, Let's Encrypt integration
- **Security Tools**: 2FA, IP blocking, password protection
- **Backup**: Backup wizard, restore functionality
- **Statistics**: AWStats, bandwidth usage, error logs
- **Advanced Tools**: Cron jobs, terminal access, CloudFlare integration

## 🏗️ Architecture

```
AdminiSoftware/
├── backend/          # Go backend API
├── frontend/         # Vue.js frontend
├── docs/            # Documentation
├── docker/          # Docker configuration
├── scripts/         # Installation and utility scripts
```

## 🛠️ Technology Stack

### Backend
- **Go 1.21+** - High-performance backend API
- **Gin** - Web framework
- **GORM** - ORM for database operations
- **JWT** - Authentication and authorization
- **Redis** - Session management and caching
- **PostgreSQL** - Primary database
- **WebSocket** - Real-time communications

### Frontend
- **Vue.js 3** - Progressive web framework
- **Vue Router** - Client-side routing
- **Pinia** - State management
- **Tailwind CSS** - Utility-first CSS framework
- **Headless UI** - Unstyled, accessible UI components
- **Chart.js** - Data visualization
- **Socket.io** - Real-time communication

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Node.js 18+ and npm
- PostgreSQL 12+
- Redis 6+

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/AdminiSoftware.git
   cd AdminiSoftware
   ```

2. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install backend dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

4. **Install frontend dependencies**
   ```bash
   cd frontend
   npm install
   ```

5. **Start the services**
   ```bash
   # Terminal 1: Start backend
   cd backend
   go run cmd/main.go
   
   # Terminal 2: Start frontend
   cd frontend
   npm run dev
   ```

6. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:5000
   - Default admin login: admin@adminisoftware.com / admin123

## 📖 Documentation

- [Installation Guide](docs/installation.md)
- [Configuration](docs/configuration.md)
- [API Documentation](docs/api-documentation.md)
- [User Guide](docs/user-guide.md)
- [Admin Guide](docs/admin-guide.md)
- [Troubleshooting](docs/troubleshooting.md)

## 🐳 Docker Deployment

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 🔒 Security Features

- **JWT Authentication** with refresh tokens
- **Two-Factor Authentication** (TOTP)
- **Brute Force Protection** with automatic IP blocking
- **Rate Limiting** on API endpoints
- **CORS Protection** with configurable origins
- **Input Validation** and sanitization
- **SQL Injection Protection** via GORM
- **XSS Protection** with Content Security Policy

## 🌟 Key Features

- **Multi-tenant Architecture** - Support for admin, reseller, and user roles
- **Real-time Updates** - WebSocket-based live updates
- **Responsive Design** - Mobile-friendly interface
- **Theme Support** - Multiple themes including cPanel and WHM styles
- **API-First Design** - RESTful API with comprehensive endpoints
- **Extensible Plugin System** - Easy to extend with custom features
- **Comprehensive Logging** - Detailed system and user activity logs
- **Backup & Restore** - Automated and manual backup solutions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💬 Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/yourusername/AdminiSoftware/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/AdminiSoftware/discussions)

## 🚧 Roadmap

- [ ] Email server management (Postfix, Dovecot)
- [ ] Web server management (Apache, Nginx)
- [ ] Database server management (MySQL, PostgreSQL)
- [ ] DNS server management (BIND)
- [ ] File server management (FTP, SFTP)
- [ ] Monitoring and alerting system
- [ ] Mobile application
- [ ] Advanced security features
- [ ] Multi-server management
- [ ] Billing integration

---

**AdminiSoftware** - Empowering web hosting management with modern technology.
# AdminiSoftware

A modern, open-source control panel alternative to cPanel/WHM built with Go and Vue.js.

## Quick Start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/adminisoftware/adminisoftware.git
   cd AdminiSoftware
   ```

2. **Run setup**:
   ```bash
   sudo ./scripts/setup.sh
   ```

3. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

4. **Start the application**:
   ```bash
   cd backend && go run cmd/main.go &
   cd ../frontend && npm run dev
   ```

## Features

- Multi-level user management (Admin/Reseller/User)
- Domain and DNS management
- Email management with spam protection
- SSL certificate management (Let's Encrypt integration)
- File manager and FTP accounts
- Database management (MySQL, PostgreSQL, MongoDB)
- Backup and restore functionality
- Security features (2FA, brute force protection)
- Real-time monitoring and statistics

## Documentation

See [docs/README.md](docs/README.md) for complete documentation.

## License

MIT License - see [LICENSE](LICENSE) file for details.
