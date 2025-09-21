# 🚀 Deployment Guide

Complete deployment instructions for GO-WMS (Warehouse Management System) in production environments.

## 📋 Prerequisites


### System Requirements
- **OS**: Linux (Ubuntu 20.04+ recommended), Windows Server, or macOS
- **Go**: Version 1.23.12 or higher
- **PostgreSQL**: Version 16 or higher
- **Memory**: Minimum 2GB RAM (4GB+ recommended)
- **Storage**: 10GB+ available space
- **Network**: Stable internet connection

### Required Software
```bash
# Install Go
wget https://go.dev/dl/go1.23.12.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.12.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib
```

## 🏗️ Production Setup

### 1. Clone and Build Application
```bash
# Clone repository
git clone https://github.com/noname-03/GO-WMS.git
cd GO-WMS

# Install dependencies
go mod tidy

# Build for production
go build -o go-wms main.go
```

### 2. Environment Configuration
```bash
# Create production environment file
cp .env.example .env.production

# Edit production environment
nano .env.production
```

**Production Environment Variables:**
```env
# Database Configuration
DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=go_wms_user
DB_PASSWORD=your-secure-password
DB_NAME=go_wms_production

# Application Configuration  
APP_PORT=8080
APP_ENV=production
JWT_SECRET=your-super-secure-jwt-secret-key-minimum-32-chars

# Optional Configuration
LOG_LEVEL=info
MAX_CONNECTIONS=100
IDLE_CONNECTIONS=10
```

### 3. Database Setup
```bash
# Create production database
sudo -u postgres psql
CREATE DATABASE go_wms_production;
CREATE USER go_wms_user WITH ENCRYPTED PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE go_wms_production TO go_wms_user;
\q

# Run migrations (application will auto-migrate on startup)
```

### 4. Security Configuration

#### Firewall Setup
```bash
# Allow HTTP traffic
sudo ufw allow 8080/tcp

# Allow SSH (if needed)
sudo ufw allow 22/tcp

# Enable firewall
sudo ufw enable
```

#### SSL/TLS Configuration
```bash
# Install Nginx for reverse proxy
sudo apt install nginx

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/go-wms
```

**Nginx Configuration:**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/go-wms /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Install SSL certificate (Let's Encrypt)
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## 🔄 Process Management

### Manual Process Management

#### Check Running Processes
```bash
# Find main process
ps aux | grep main
```

#### Kill Process
```bash
# Kill specific process by PID
kill -9 547218

# Or kill by name
pkill -f "go-wms"
```

#### Start Application in Background
```bash
# Start with nohup for persistent running
nohup ./go-wms &

# Or with environment file
nohup ./go-wms --env=.env.production &
```

#### Monitor Application Logs
```bash
# Monitor real-time logs
tail -f nohup.out

# View last 100 lines
tail -n 100 nohup.out

# Follow logs with grep filtering
tail -f nohup.out | grep "ERROR"
```

#### Complete Process Management Workflow
```bash
# 1. Check if application is running
ps aux | grep main

# 2. Stop existing process (if running)
kill -9 547218

# 3. Start application in background
nohup ./go-wms &

# 4. Monitor application startup
tail -f nohup.out
```

### Systemd Service (Recommended)

#### Create Service File
```bash
sudo nano /etc/systemd/system/go-wms.service
```

**Service Configuration:**
```ini
[Unit]
Description=GO-WMS Application
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/go-wms
ExecStart=/opt/go-wms/go-wms
Restart=always
RestartSec=5
Environment=ENV=production
EnvironmentFile=/opt/go-wms/.env.production

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/go-wms

[Install]
WantedBy=multi-user.target
```

#### Service Management
```bash
# Reload systemd configuration
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable go-wms

# Start service
sudo systemctl start go-wms

# Check service status
sudo systemctl status go-wms

# Stop service
sudo systemctl stop go-wms

# Restart service
sudo systemctl restart go-wms

# View service logs
sudo journalctl -u go-wms -f
```

## 🐳 Docker Deployment

### Docker Compose Production
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=your-secure-password
      - DB_NAME=go_wms
      - JWT_SECRET=your-super-secure-jwt-secret
      - APP_ENV=production
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:16
    environment:
      - POSTGRES_DB=go_wms
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=your-secure-password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
```

### Docker Commands
```bash
# Build and start production containers
docker-compose -f docker-compose.prod.yml up -d

# View running containers
docker ps

# View application logs
docker-compose -f docker-compose.prod.yml logs -f app

# Stop containers
docker-compose -f docker-compose.prod.yml down

# Update application
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d
```

## 📊 Monitoring and Logging

### Application Monitoring
```bash
# Check application health
curl http://localhost:8080/health

# Monitor resource usage
htop
df -h
free -m
```

### Database Monitoring
```bash
# PostgreSQL status
sudo systemctl status postgresql

# Database connections
sudo -u postgres psql -c "SELECT count(*) FROM pg_stat_activity;"

# Database size
sudo -u postgres psql -c "SELECT pg_size_pretty(pg_database_size('go_wms_production'));"
```

### Log Management
```bash
# Rotate logs to prevent disk space issues
sudo nano /etc/logrotate.d/go-wms
```

**Log Rotation Configuration:**
```
/opt/go-wms/nohup.out {
    daily
    rotate 30
    compress
    missingok
    notifempty
    create 644 www-data www-data
    postrotate
        systemctl reload go-wms
    endscript
}
```

## 🔒 Security Best Practices

### Application Security
- **Use HTTPS** in production
- **Strong JWT secrets** (minimum 32 characters)
- **Environment-specific configurations**
- **Regular security updates**

### Database Security
- **Strong database passwords**
- **Limited database user permissions**
- **Regular database backups**
- **Connection encryption**

### Server Security
- **Firewall configuration**
- **Regular OS updates**
- **SSH key authentication**
- **Fail2ban for intrusion prevention**

## 💾 Backup and Recovery

### Database Backup
```bash
# Create backup script
nano /opt/scripts/backup-db.sh
```

**Backup Script:**
```bash
#!/bin/bash
BACKUP_DIR="/opt/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
DB_NAME="go_wms_production"

# Create backup
pg_dump -h localhost -U go_wms_user $DB_NAME > $BACKUP_DIR/backup_$TIMESTAMP.sql

# Compress backup
gzip $BACKUP_DIR/backup_$TIMESTAMP.sql

# Remove backups older than 30 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
```

```bash
# Make script executable
chmod +x /opt/scripts/backup-db.sh

# Add to crontab for daily backups
crontab -e
# Add: 0 2 * * * /opt/scripts/backup-db.sh
```

### Application Backup
```bash
# Backup application files
tar -czf /opt/backups/app_$(date +%Y%m%d).tar.gz /opt/go-wms
```

### Recovery Process
```bash
# Restore database
gunzip backup_20240315_020000.sql.gz
psql -h localhost -U go_wms_user go_wms_production < backup_20240315_020000.sql

# Restore application
tar -xzf app_20240315.tar.gz -C /opt/
```

## 🔧 Performance Tuning

### Application Optimization
```env
# Optimize connection pooling
MAX_CONNECTIONS=100
IDLE_CONNECTIONS=10
MAX_LIFETIME=3600

# Enable production optimizations
GOMEMLIMIT=2GiB
GOGC=100
```

### Database Optimization
```sql
-- PostgreSQL configuration tuning
-- Edit postgresql.conf

shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB
max_connections = 100
```

### Nginx Optimization
```nginx
# nginx.conf optimizations
worker_processes auto;
worker_connections 1024;

gzip on;
gzip_types text/plain application/json application/javascript text/css;

client_max_body_size 10M;
```

## 📈 Scaling Considerations

### Horizontal Scaling
- **Load balancer setup** (HAProxy/Nginx)
- **Multiple application instances**
- **Database read replicas**
- **Shared session storage** (Redis)

### Vertical Scaling
- **Increase server resources** (CPU/RAM)
- **Database connection tuning**
- **Application optimization**

## 🚨 Troubleshooting

### Common Issues

#### Application Won't Start
```bash
# Check port availability
netstat -tlnp | grep :8080

# Check database connectivity
pg_isready -h localhost -p 5432

# Verify environment variables
printenv | grep DB_
```

#### High Memory Usage
```bash
# Monitor memory usage
ps aux --sort=-%mem | head

# Check for memory leaks
go tool pprof http://localhost:8080/debug/pprof/heap
```

#### Database Connection Issues
```bash
# Test database connection
psql -h localhost -U go_wms_user -d go_wms_production

# Check active connections
sudo -u postgres psql -c "SELECT count(*) FROM pg_stat_activity;"
```

### Health Checks
```bash
# Application health
curl -f http://localhost:8080/health || exit 1

# Database health
pg_isready -h localhost -p 5432 || exit 1

# Disk space check
df -h | grep -E '^/dev/' | awk '{print $5}' | sed 's/%//' | grep -q '^[0-9][0-9]$' || exit 1
```

## 📞 Support and Maintenance

### Regular Maintenance Tasks
- **Weekly security updates**
- **Monthly dependency updates**
- **Quarterly performance reviews**
- **Annual security audits**

### Monitoring Alerts
- **Application downtime**
- **High memory/CPU usage**
- **Database connection failures**
- **Disk space warnings**