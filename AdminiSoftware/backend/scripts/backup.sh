
#!/bin/bash

# AdminiSoftware Backup Script
echo "=== AdminiSoftware Backup Script ==="

# Configuration
BACKUP_DIR="/var/backups/adminisoftware"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="adminisoftware_backup_$DATE"
RETENTION_DAYS=30

# Load environment
CONFIG_FILE="/etc/adminisoftware/config.env"
if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
fi

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-adminisoftware_db}
DB_USER=${DB_USER:-adminisoftware}
DB_PASSWORD=${DB_PASSWORD:-adminisoftware123}

# Create backup directory
mkdir -p "$BACKUP_DIR"

echo "Starting backup: $BACKUP_NAME"

# Database backup
echo "Backing up database..."
export PGPASSWORD="$DB_PASSWORD"
pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_DIR/${BACKUP_NAME}_database.sql"

if [ $? -eq 0 ]; then
    echo "Database backup completed successfully"
else
    echo "Database backup failed!"
    exit 1
fi

# Application files backup
echo "Backing up application files..."
if [ -d "/opt/adminisoftware" ]; then
    tar -czf "$BACKUP_DIR/${BACKUP_NAME}_application.tar.gz" -C /opt adminisoftware
    echo "Application files backup completed"
fi

# Configuration backup
echo "Backing up configuration..."
if [ -d "/etc/adminisoftware" ]; then
    tar -czf "$BACKUP_DIR/${BACKUP_NAME}_config.tar.gz" -C /etc adminisoftware
    echo "Configuration backup completed"
fi

# User data backup (if applicable)
echo "Backing up user data..."
if [ -d "/home" ]; then
    # Only backup adminisoftware user data
    if [ -d "/home/adminisoftware" ]; then
        tar -czf "$BACKUP_DIR/${BACKUP_NAME}_userdata.tar.gz" -C /home adminisoftware
        echo "User data backup completed"
    fi
fi

# Create manifest file
echo "Creating backup manifest..."
cat > "$BACKUP_DIR/${BACKUP_NAME}_manifest.txt" << EOF
Backup Created: $(date)
Backup Name: $BACKUP_NAME
Database: ${BACKUP_NAME}_database.sql
Application: ${BACKUP_NAME}_application.tar.gz
Configuration: ${BACKUP_NAME}_config.tar.gz
User Data: ${BACKUP_NAME}_userdata.tar.gz
EOF

# Compress all backup files
echo "Compressing backup..."
cd "$BACKUP_DIR"
tar -czf "${BACKUP_NAME}.tar.gz" ${BACKUP_NAME}_*

# Clean up individual files
rm -f ${BACKUP_NAME}_database.sql
rm -f ${BACKUP_NAME}_application.tar.gz
rm -f ${BACKUP_NAME}_config.tar.gz
rm -f ${BACKUP_NAME}_userdata.tar.gz
rm -f ${BACKUP_NAME}_manifest.txt

# Calculate backup size
BACKUP_SIZE=$(du -h "${BACKUP_NAME}.tar.gz" | cut -f1)
echo "Backup completed: ${BACKUP_NAME}.tar.gz ($BACKUP_SIZE)"

# Clean old backups
echo "Cleaning old backups (older than $RETENTION_DAYS days)..."
find "$BACKUP_DIR" -name "adminisoftware_backup_*.tar.gz" -mtime +$RETENTION_DAYS -delete

# Log backup
echo "$(date): Backup $BACKUP_NAME completed successfully ($BACKUP_SIZE)" >> "$BACKUP_DIR/backup.log"

echo "=== Backup Complete ==="

# Optional: Upload to remote storage
if [ ! -z "$REMOTE_BACKUP_SCRIPT" ] && [ -f "$REMOTE_BACKUP_SCRIPT" ]; then
    echo "Uploading to remote storage..."
    "$REMOTE_BACKUP_SCRIPT" "$BACKUP_DIR/${BACKUP_NAME}.tar.gz"
fi
#!/bin/bash

# AdminiSoftware Backup Script

# Configuration
BACKUP_DIR="/var/backups/adminisoftware"
DB_NAME="adminisoftware_db"
DB_USER="adminisoftware"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Create backup directory
mkdir -p $BACKUP_DIR

echo "Starting AdminiSoftware backup at $(date)"

# Database backup
echo "Backing up database..."
pg_dump -U $DB_USER -h localhost $DB_NAME > $BACKUP_DIR/database_$TIMESTAMP.sql
gzip $BACKUP_DIR/database_$TIMESTAMP.sql

# Application files backup
echo "Backing up application files..."
tar -czf $BACKUP_DIR/app_files_$TIMESTAMP.tar.gz /opt/adminisoftware

# Configuration backup
echo "Backing up configuration..."
tar -czf $BACKUP_DIR/config_$TIMESTAMP.tar.gz /etc/adminisoftware

# Log files backup
echo "Backing up logs..."
tar -czf $BACKUP_DIR/logs_$TIMESTAMP.tar.gz /var/log/adminisoftware

# Clean up old backups
echo "Cleaning up old backups..."
find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed at $(date)"

# Create backup report
cat > $BACKUP_DIR/backup_report_$TIMESTAMP.txt << EOF
AdminiSoftware Backup Report
Date: $(date)
Backup Location: $BACKUP_DIR

Files Created:
- database_$TIMESTAMP.sql.gz
- app_files_$TIMESTAMP.tar.gz
- config_$TIMESTAMP.tar.gz
- logs_$TIMESTAMP.tar.gz

Status: Success
EOF

echo "Backup report saved to $BACKUP_DIR/backup_report_$TIMESTAMP.txt"
