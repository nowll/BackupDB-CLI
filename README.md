# 🗄️ Database Backup CLI

A production-ready, enterprise-grade command-line tool for backing up multiple database types with support for cloud storage, automated scheduling, and real-time notifications. Built with Go following Clean Architecture principles.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Usage Guide](#-usage-guide)
- [Configuration](#-configuration)
- [Storage Backends](#-storage-backends)
- [Scheduling](#-scheduling)
- [Security](#-security)
- [Examples](#-examples)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

### Database Support
- ✅ **MySQL** - Full support with mysqldump
- ✅ **PostgreSQL** - Custom format backups with pg_dump
- ✅ **MongoDB** - BSON archive format with mongodump
- ✅ **SQLite** - Direct file-based backups

### Backup Capabilities
- 🔄 **Multiple Backup Types**: Full, incremental, and differential
- 🗜️ **Automatic Compression**: Gzip compression to save storage space
- ✔️ **Integrity Verification**: SHA-256 checksum validation
- 📊 **Metadata Tracking**: Complete backup history and statistics
- 🎯 **Selective Backup**: Choose specific tables/collections
- ⚡ **Fast & Efficient**: Optimized for large databases

### Storage Options
- 💾 **Local Storage**: Filesystem-based storage
- ☁️ **AWS S3**: Amazon S3 bucket storage
- 🌩️ **Google Cloud Storage**: GCS bucket support
- 🔷 **Azure Blob Storage**: Microsoft Azure integration

### Advanced Features
- ⏰ **Automated Scheduling**: Cron-like scheduling for regular backups
- 📢 **Slack Notifications**: Real-time backup status updates
- 📝 **Comprehensive Logging**: Structured logging with Zap
- 🔐 **Secure Connections**: SSL/TLS support for all databases
- 🧪 **Connection Testing**: Validate credentials before backup
- 🔄 **Easy Restoration**: Simple one-command restore process

### Project Structure

```
db-backup-cli/
├── cmd/
│   └── backup/
│       └── main.go                 # CLI entry point
├── internal/
│   ├── domain/                     # Business entities
│   │   ├── backup.go               # Backup domain models
│   │   ├── database.go             # Database interfaces
│   │   └── storage.go              # Storage interfaces
│   ├── usecase/                    # Business logic
│   │   ├── backup_usecase.go       # Backup orchestration
│   │   └── restore_usecase.go      # Restore orchestration
│   ├── repository/                 # Database implementations
│   │   ├── mysql_repository.go
│   │   ├── postgresql_repository.go
│   │   ├── mongodb_repository.go
│   │   ├── sqlite_repository.go
│   │   ├── factory.go
│   │   └── metadata_repository.go
│   ├── storage/                    # Storage implementations
│   │   ├── local_storage.go
│   │   ├── s3_storage.go
│   │   ├── gcs_storage.go
│   │   ├── azure_storage.go
│   │   └── factory.go
│   └── infrastructure/             # External services
│       ├── logger/
│       │   └── logger.go
│       ├── notifier/
│       │   └── slack.go
│       └── scheduler/
│           └── scheduler.go
├── pkg/                            # Public packages
│   ├── config/
│   │   └── config.go
│   └── compression/
│       └── compression.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 📦 Installation

### Prerequisites

**Required:**
- Go 1.21 or higher
- Git

**Database Tools (based on what you'll backup):**
- MySQL: `mysqldump` and `mysql` client
- PostgreSQL: `pg_dump` and `pg_restore`
- MongoDB: `mongodump` and `mongorestore`
- SQLite: No additional tools needed

### Install from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/db-backup-cli.git
cd db-backup-cli

# Download dependencies
go mod download

# Build the binary
go build -o db-backup cmd/backup/main.go

# (Optional) Install globally
sudo mv db-backup /usr/local/bin/

# Verify installation
db-backup --help
```

### Install via Go

```bash
go install github.com/yourusername/db-backup-cli/cmd/backup@latest
```

### Build with Makefile

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install

# Run tests
make test
```

## 🚀 Quick Start

### 1. Test Database Connection

Before creating backups, verify your database connection:

```bash
db-backup test \
  --db-type mysql \
  --host localhost \
  --port 3306 \
  --user root \
  --password mypassword \
  --database mydb
```

**Expected output:**
```
Testing connection to mysql database at localhost:3306...
✓ Connection successful!
```

### 2. Create Your First Backup

```bash
db-backup backup \
  --db-type mysql \
  --host localhost \
  --port 3306 \
  --user root \
  --password mypassword \
  --database mydb \
  --type full \
  --compress true \
  --storage local \
  --path ./backups
```

**Expected output:**
```
Backup completed successfully!
Backup ID: mysql_mydb_20240119_143022
Size: 15728640 bytes (compressed: 4194304 bytes)
Duration: 3.45s
```

### 3. List Available Backups

```bash
db-backup list --path ./backups
```

**Expected output:**
```
Found 3 backup(s):

ID: mysql_mydb_20240119_143022
  Type: mysql (full)
  Database: mydb
  Date: 2024-01-19 14:30:22
  Size: 15728640 bytes (compressed: 4194304 bytes)
  Status: completed

ID: mysql_mydb_20240118_020000
  Type: mysql (full)
  Database: mydb
  Date: 2024-01-18 02:00:00
  Size: 15200000 bytes (compressed: 4000000 bytes)
  Status: completed
```

### 4. Restore from Backup

```bash
db-backup restore mysql_mydb_20240119_143022 \
  --db-type mysql \
  --host localhost \
  --port 3306 \
  --user root \
  --password mypassword \
  --database mydb_restored
```

## 📚 Usage Guide

### Command Reference

#### `backup` - Create a Database Backup

```bash
db-backup backup [flags]
```

**Required Flags:**
- `--db-type` - Database type (mysql, postgresql, mongodb, sqlite)
- `--user` - Database username
- `--password` - Database password
- `--database` - Database name

**Optional Flags:**
- `--host` - Database host (default: localhost)
- `--port` - Database port (default: 3306 for MySQL)
- `--type` - Backup type: full, incremental, differential (default: full)
- `--compress` - Enable compression (default: true)
- `--storage` - Storage type: local, s3, gcs, azure (default: local)
- `--path` - Storage path (default: ./backups)

**Examples:**

```bash
# MySQL full backup
db-backup backup \
  --db-type mysql \
  --host db.example.com \
  --port 3306 \
  --user admin \
  --password secret \
  --database production

# PostgreSQL backup to S3
db-backup backup \
  --db-type postgresql \
  --host 10.0.1.5 \
  --port 5432 \
  --user postgres \
  --password pgpass \
  --database appdb \
  --storage s3 \
  --path s3://my-backups/postgres

# MongoDB backup with specific collections
db-backup backup \
  --db-type mongodb \
  --host mongo.example.com \
  --port 27017 \
  --user mongoadmin \
  --password mongopass \
  --database analytics

# SQLite backup
db-backup backup \
  --db-type sqlite \
  --database /var/lib/myapp/app.db \
  --storage local \
  --path /backup/sqlite
```

#### `restore` - Restore from Backup

```bash
db-backup restore [backup-id] [flags]
```

**Required Arguments:**
- `backup-id` - The ID of the backup to restore

**Required Flags:**
- `--db-type` - Database type
- `--user` - Database username
- `--password` - Database password

**Optional Flags:**
- `--host` - Database host (default: localhost)
- `--port` - Database port
- `--database` - Target database name

**Examples:**

```bash
# Restore MySQL backup
db-backup restore mysql_mydb_20240119_143022 \
  --db-type mysql \
  --host localhost \
  --user root \
  --password secret \
  --database mydb_restored

# Restore PostgreSQL to different server
db-backup restore postgresql_appdb_20240119_020000 \
  --db-type postgresql \
  --host new-server.example.com \
  --port 5432 \
  --user postgres \
  --password pgpass \
  --database appdb_new
```

#### `list` - List All Backups

```bash
db-backup list [flags]
```

**Flags:**
- `--path` - Storage path to list backups from (default: ./backups)

**Example:**

```bash
db-backup list --path ./backups
db-backup list --path s3://my-backups
```

#### `schedule` - Schedule Automatic Backups

```bash
db-backup schedule [flags]
```

**Required Flags:**
- `--cron` - Cron expression for schedule
- `--db-type` - Database type
- `--user` - Database username
- `--password` - Database password
- `--database` - Database name

**Cron Expression Examples:**
- `"0 2 * * *"` - Daily at 2 AM
- `"0 */6 * * *"` - Every 6 hours
- `"0 0 * * 0"` - Weekly on Sunday at midnight
- `"*/30 * * * *"` - Every 30 minutes
- `"0 3 * * 1-5"` - Weekdays at 3 AM

**Example:**

```bash
# Daily backup at 2 AM
db-backup schedule \
  --cron "0 2 * * *" \
  --db-type mysql \
  --host localhost \
  --user root \
  --password secret \
  --database production \
  --compress true \
  --storage s3 \
  --path s3://daily-backups/mysql
```

#### `test` - Test Database Connection

```bash
db-backup test [flags]
```

**Required Flags:**
- `--db-type` - Database type
- `--user` - Database username
- `--password` - Database password

**Example:**

```bash
db-backup test \
  --db-type postgresql \
  --host db.example.com \
  --port 5432 \
  --user postgres \
  --password pgpass \
  --database testdb
```

## ⚙️ Configuration

### Using Configuration File

Create a `config.yaml` file:

```yaml
database:
  type: mysql
  host: localhost
  port: 3306
  username: root
  password: mypassword
  database: mydb
  sslmode: disable

storage:
  type: s3
  bucket: my-backup-bucket
  region: us-east-1
  access_key: AKIAIOSFODNN7EXAMPLE
  secret_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
  path: backups/mysql

logging:
  level: info
  file_path: logs/backup.log

slack:
  enabled: true
  token: xoxb-your-slack-token
  channel: #backups
```

**Use the config file:**

```bash
db-backup backup --config config.yaml
```

### Environment Variables

Set sensitive credentials via environment variables:

```bash
# AWS S3
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-east-1

# Google Cloud Storage
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json

# Azure Blob Storage
export AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=https;AccountName=..."

# Slack Notifications
export SLACK_TOKEN=xoxb-your-slack-bot-token
export SLACK_CHANNEL=#backups

# Database Credentials (not recommended for production)
export DB_USER=admin
export DB_PASSWORD=secret
```

## ☁️ Storage Backends

### Local Storage

Store backups on the local filesystem:

```bash
db-backup backup \
  --storage local \
  --path /mnt/backups/mysql
```

**Directory structure:**
```
/mnt/backups/mysql/
├── backups/
│   └── mysql/
│       └── mysql_mydb_20240119_143022.backup
└── metadata/
    └── mysql_mydb_20240119_143022.json
```

### AWS S3

Store backups in Amazon S3:

```bash
# Set credentials
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret

# Create backup
db-backup backup \
  --storage s3 \
  --path s3://my-backup-bucket/mysql
```

**S3 Bucket Policy Example:**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::my-backup-bucket/*",
        "arn:aws:s3:::my-backup-bucket"
      ]
    }
  ]
}
```

### Google Cloud Storage

Store backups in GCS:

```bash
# Set credentials
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json

# Create backup
db-backup backup \
  --storage gcs \
  --path gs://my-backup-bucket/mysql
```

**Service Account Permissions:**
- `storage.objects.create`
- `storage.objects.delete`
- `storage.objects.get`
- `storage.objects.list`

### Azure Blob Storage

Store backups in Azure:

```bash
# Set connection string
export AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=https;..."

# Create backup
db-backup backup \
  --storage azure \
  --path azure://my-container/mysql
```

## ⏰ Scheduling

### Systemd Service (Linux)

Create `/etc/systemd/system/db-backup.service`:

```ini
[Unit]
Description=Database Backup Service
After=network.target

[Service]
Type=simple
User=backup
Group=backup
WorkingDirectory=/opt/db-backup
ExecStart=/usr/local/bin/db-backup schedule \
  --cron "0 2 * * *" \
  --db-type mysql \
  --host localhost \
  --user root \
  --password secret \
  --database production
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Enable and start:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable db-backup
sudo systemctl start db-backup
sudo systemctl status db-backup
```

### Cron Job (Linux/macOS)

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /usr/local/bin/db-backup backup --config /etc/db-backup/config.yaml >> /var/log/db-backup.log 2>&1
```

### Task Scheduler (Windows)

```powershell
# Create scheduled task
$action = New-ScheduledTaskAction -Execute "C:\Program Files\db-backup\db-backup.exe" -Argument "backup --config C:\db-backup\config.yaml"
$trigger = New-ScheduledTaskTrigger -Daily -At 2am
Register-ScheduledTask -TaskName "DatabaseBackup" -Action $action -Trigger $trigger
```

### Docker Container

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o db-backup cmd/backup/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates mysql-client postgresql-client
COPY --from=builder /app/db-backup /usr/local/bin/
ENTRYPOINT ["db-backup"]
```

**Run with Docker:**

```bash
docker run -v /backups:/backups \
  -e DB_USER=root \
  -e DB_PASSWORD=secret \
  db-backup:latest backup \
  --db-type mysql \
  --host host.docker.internal \
  --database mydb
```

## 🔐 Security

### Best Practices

1. **Never Commit Credentials**
   ```bash
   # Use .gitignore
   echo "config.yaml" >> .gitignore
   echo "*.env" >> .gitignore
   ```

2. **Use Environment Variables**
   ```bash
   export DB_PASSWORD=$(cat /run/secrets/db_password)
   ```

3. **Encrypt Backups**
   ```bash
   # Encrypt before upload
   gpg --encrypt --recipient backup@example.com backup.sql
   ```

4. **Restrict File Permissions**
   ```bash
   chmod 600 config.yaml
   chmod 700 /backups
   ```

5. **Use IAM Roles** (AWS)
   - Attach IAM role to EC2 instance
   - No need for access keys

6. **Enable Audit Logging**
   ```yaml
   logging:
     level: debug
     file_path: /var/log/db-backup/audit.log
   ```

### SSL/TLS Connections

**MySQL:**
```bash
db-backup backup \
  --db-type mysql \
  --host secure-db.example.com \
  --sslmode require
```

**PostgreSQL:**
```bash
db-backup backup \
  --db-type postgresql \
  --host secure-db.example.com \
  --sslmode require
```

## 📖 Examples

### Example 1: Daily Production Backup to S3

```bash
#!/bin/bash
# daily-backup.sh

export AWS_ACCESS_KEY_ID=$(cat /run/secrets/aws_key)
export AWS_SECRET_ACCESS_KEY=$(cat /run/secrets/aws_secret)
export SLACK_TOKEN=$(cat /run/secrets/slack_token)

db-backup backup \
  --db-type postgresql \
  --host prod-db.example.com \
  --port 5432 \
  --user backup_user \
  --password "$(cat /run/secrets/db_password)" \
  --database production \
  --type full \
  --compress true \
  --storage s3 \
  --path s3://company-backups/postgresql/production

# Check exit code
if [ $? -eq 0 ]; then
    echo "Backup completed successfully"
else
    echo "Backup failed!"
    exit 1
fi
```

### Example 2: Multi-Database Backup Script

```bash
#!/bin/bash
# backup-all.sh

DATABASES=("userdb" "productdb" "analyticsdb")
BACKUP_PATH="./backups"

for db in "${DATABASES[@]}"; do
    echo "Backing up $db..."

    db-backup backup \
      --db-type mysql \
      --host localhost \
      --port 3306 \
      --user root \
      --password secret \
      --database "$db" \
      --compress true \
      --storage local \
      --path "$BACKUP_PATH"

    if [ $? -ne 0 ]; then
        echo "Failed to backup $db"
    fi
done
```

### Example 3: Backup with Retention Policy

```bash
#!/bin/bash
# backup-with-retention.sh

# Create backup
db-backup backup \
  --db-type mysql \
  --host localhost \
  --user root \
  --password secret \
  --database mydb \
  --path /backups

# Delete backups older than 30 days
find /backups/backups -name "*.backup" -mtime +30 -delete
find /backups/metadata -name "*.json" -mtime +30 -delete

echo "Cleanup completed"
```

### Example 4: Automated Restore Testing

```bash
#!/bin/bash
# test-restore.sh

LATEST_BACKUP=$(db-backup list --path ./backups | grep "ID:" | head -1 | awk '{print $2}')

echo "Testing restore of backup: $LATEST_BACKUP"

db-backup restore "$LATEST_BACKUP" \
  --db-type mysql \
  --host localhost \
  --user root \
  --password secret \
  --database test_restore

if [ $? -eq 0 ]; then
    echo "✓ Restore test passed"
    mysql -u root -psecret -e "DROP DATABASE test_restore"
else
    echo "✗ Restore test failed"
    exit 1
fi
```

## 🔧 Troubleshooting

### Common Issues

#### 1. Connection Refused

**Problem:**
```
Error: test connection: dial tcp 127.0.0.1:3306: connect: connection refused
```

**Solution:**
- Verify database is running: `systemctl status mysql`
- Check host and port: `netstat -tlnp | grep 3306`
- Verify firewall rules: `sudo ufw status`

#### 2. Authentication Failed

**Problem:**
```
Error: test connection: Access denied for user 'backup'@'localhost'
```

**Solution:**
```sql
-- Grant permissions
GRANT SELECT, LOCK TABLES, SHOW VIEW ON *.* TO 'backup'@'localhost' IDENTIFIED BY 'password';
FLUSH PRIVILEGES;
```

#### 3. Insufficient Permissions

**Problem:**
```
Error: backup: mysqldump: Got error: 1044: Access denied
```

**Solution:**
```sql
-- MySQL: Grant necessary privileges
GRANT SELECT, LOCK TABLES, SHOW VIEW, EVENT, TRIGGER ON database.* TO 'user'@'host';

-- PostgreSQL: Grant usage
GRANT USAGE ON SCHEMA public TO backup_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO backup_user;
```

#### 4. Out of Disk Space

**Problem:**
```
Error: upload: write /backups/...: no space left on device
```

**Solution:**
```bash
# Check disk usage
df -h

# Clean old backups
find /backups -name "*.backup" -mtime +7 -delete

# Enable compression
db-backup backup --compress true
```

#### 5. S3 Access Denied

**Problem:**
```
Error: upload to s3: AccessDenied: Access Denied
```

**Solution:**
- Verify AWS credentials: `aws s3 ls`
- Check bucket policy and IAM permissions
- Ensure bucket exists: `aws s3 mb s3://bucket-name`

### Debug Mode

Enable verbose logging:

```bash
# Set log level to debug
export LOG_LEVEL=debug

db-backup backup \
  --db-type mysql \
  --host localhost \
  --user root \
  --password secret \
  --database mydb
```

### Check Logs

```bash
# View logs
tail -f logs/backup.log

# Search for errors
grep ERROR logs/backup.log

# View specific backup logs
cat logs/backup_20240119.log
```

## 🧪 Testing

### Run Tests

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...

# Specific package
go test ./internal/usecase

# Verbose output
go test -v ./...
```

### Manual Testing

```bash
# 1. Test connection
db-backup test --db-type mysql --host localhost --user root --password secret

# 2. Create backup
db-backup backup --db-type mysql --host localhost --user root --password secret --database testdb

# 3. List backups
db-backup list

# 4. Restore backup
db-backup restore <backup-id> --db-type mysql --host localhost --user root --password secret --database testdb_restored

# 5. Verify restore
mysql -u root -psecret testdb_restored -e "SHOW TABLES;"
```

### Development Setup

```bash
# Fork and clone
git clone https://github.com/yourusername/db-backup-cli.git
cd db-backup-cli

# Install dependencies
go mod download

# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
go test ./...

# Commit changes
git commit -m "Add amazing feature"

# Push and create PR
git push origin feature/amazing-feature
```

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Write tests for new features
- Update documentation

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Zap](https://github.com/uber-go/zap) - Structured logging
- [AWS SDK](https://github.com/aws/aws-sdk-go) - S3 integration
- Clean Architecture by Robert C. Martin

## 📞 Support

- 📧 Email: support@example.com
- 💬 Slack: [Join our community](https://slack.example.com)
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/db-backup-cli/issues)
- 📖 Documentation: [Wiki](https://github.com/yourusername/db-backup-cli/wiki)

## 🗺️ Roadmap

### Version 2.0
- [ ] Support for Microsoft SQL Server
- [ ] Support for Oracle Database
- [ ] Backup encryption (AES-256)
- [ ] Email notifications
- [ ] Web UI for management

### Version 2.1
- [ ] Backup rotation policies
- [ ] Point-in-time recovery
- [ ] Backup deduplication
- [ ] Parallel backup support
- [ ] Backup verification tests

### Version 3.0
- [ ] Multi-region replication
- [ ] Backup analytics dashboard
- [ ] Kubernetes operator
- [ ] Disaster recovery automation
- [ ] Compliance reporting

---


*Star ⭐ this repository if you find it helpful!*
