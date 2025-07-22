
#!/bin/bash

# AdminiSoftware Database Migration Script
echo "=== AdminiSoftware Database Migration ==="

# Set environment
ENVIRONMENT=${ENVIRONMENT:-development}
CONFIG_FILE="/etc/adminisoftware/config.env"

if [ "$ENVIRONMENT" = "production" ] && [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
fi

# Default values for development
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-adminisoftware_db}
DB_USER=${DB_USER:-adminisoftware}
DB_PASSWORD=${DB_PASSWORD:-adminisoftware123}

echo "Running migrations for $ENVIRONMENT environment..."
echo "Database: $DB_HOST:$DB_PORT/$DB_NAME"

# Check if database is accessible
export PGPASSWORD="$DB_PASSWORD"
pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"
if [ $? -ne 0 ]; then
    echo "Error: Cannot connect to database"
    exit 1
fi

# Run migrations
echo "Creating database schema..."

psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << 'EOF'
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    package_id INTEGER,
    quota_limit BIGINT,
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    two_factor_secret VARCHAR(32),
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Packages table
CREATE TABLE IF NOT EXISTS packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    disk_quota BIGINT NOT NULL,
    bandwidth_quota BIGINT NOT NULL,
    email_accounts INTEGER NOT NULL,
    databases INTEGER NOT NULL,
    subdomains INTEGER NOT NULL,
    ftp_accounts INTEGER NOT NULL,
    price DECIMAL(10,2),
    billing_cycle VARCHAR(20),
    features JSONB,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Domains table
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    document_root VARCHAR(255),
    php_version VARCHAR(10),
    ssl_enabled BOOLEAN DEFAULT FALSE,
    ssl_certificate_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Email accounts table
CREATE TABLE IF NOT EXISTS email_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    domain_id INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    quota_mb INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Databases table
CREATE TABLE IF NOT EXISTS databases (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'mysql',
    host VARCHAR(255) DEFAULT 'localhost',
    port INTEGER,
    username VARCHAR(100),
    password_hash VARCHAR(255),
    size_mb BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SSL certificates table
CREATE TABLE IF NOT EXISTS ssl_certificates (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    provider VARCHAR(50),
    certificate TEXT,
    private_key TEXT,
    chain TEXT,
    expires_at TIMESTAMP,
    auto_renew BOOLEAN DEFAULT TRUE,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Backups table
CREATE TABLE IF NOT EXISTS backups (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    size_mb BIGINT,
    file_path VARCHAR(500),
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- DNS records table
CREATE TABLE IF NOT EXISTS dns_records (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL,
    value TEXT NOT NULL,
    ttl INTEGER DEFAULT 3600,
    priority INTEGER,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- System settings table
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_domains_user_id ON domains(user_id);
CREATE INDEX IF NOT EXISTS idx_domains_name ON domains(name);
CREATE INDEX IF NOT EXISTS idx_email_accounts_user_id ON email_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_databases_user_id ON databases(user_id);
CREATE INDEX IF NOT EXISTS idx_dns_records_domain_id ON dns_records(domain_id);

-- Insert default data
INSERT INTO packages (name, description, disk_quota, bandwidth_quota, email_accounts, databases, subdomains, ftp_accounts, price, billing_cycle)
VALUES 
    ('Basic', 'Basic hosting package', 1073741824, 10737418240, 5, 2, 5, 2, 9.99, 'monthly'),
    ('Standard', 'Standard hosting package', 5368709120, 53687091200, 20, 10, 20, 10, 19.99, 'monthly'),
    ('Premium', 'Premium hosting package', 21474836480, 214748364800, -1, -1, -1, -1, 39.99, 'monthly')
ON CONFLICT DO NOTHING;

-- Insert default admin user
INSERT INTO users (username, email, password_hash, role, status)
VALUES ('admin', 'admin@adminisoftware.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', 'active')
ON CONFLICT DO NOTHING;

-- Insert default system settings
INSERT INTO system_settings (key, value, description)
VALUES 
    ('site_name', 'AdminiSoftware', 'Site name'),
    ('default_theme', 'default', 'Default theme'),
    ('enable_registration', 'false', 'Enable user registration'),
    ('maintenance_mode', 'false', 'Maintenance mode'),
    ('backup_retention_days', '30', 'Backup retention period'),
    ('max_file_upload_size', '100', 'Maximum file upload size in MB')
ON CONFLICT (key) DO NOTHING;

EOF

echo "Migration completed successfully!"

# Run Go application migrations if binary exists
if [ -f "/opt/adminisoftware/bin/adminisoftware" ]; then
    echo "Running application-level migrations..."
    cd /opt/adminisoftware && ./bin/adminisoftware --migrate
fi

echo "=== Migration Complete ==="
