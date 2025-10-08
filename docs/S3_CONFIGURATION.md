# S3 Storage Configuration Guide

File management system mendukung 2 jenis storage:

## 🚀 **1. Amazon S3 (AWS)**

Untuk menggunakan Amazon S3 storage, set environment variables berikut di `.env`:

```env
# Amazon S3 Configuration
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION=us-east-1
AWS_S3_BUCKET=my-wms-bucket
AWS_S3_ENDPOINT=
AWS_S3_USE_SSL=true
AWS_S3_PATH_STYLE=false
```

**Cara mendapatkan AWS credentials:**
1. Login ke AWS Console
2. Buka IAM service
3. Buat user baru atau gunakan existing user
4. Attach policy `AmazonS3FullAccess` atau custom policy
5. Generate Access Key dan Secret Key

## 🏠 **2. S3 Compatible Storage**

### **MinIO (Self-hosted)**

```env
# MinIO Configuration
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin123
AWS_REGION=us-east-1
AWS_S3_BUCKET=wms-files
AWS_S3_ENDPOINT=http://localhost:9000
AWS_S3_USE_SSL=false
AWS_S3_PATH_STYLE=true
```

**Setup MinIO:**
```bash
# Download MinIO
wget https://dl.min.io/server/minio/release/windows-amd64/minio.exe

# Start MinIO server
./minio.exe server C:\minio-data --console-address ":9001"

# Access web console: http://localhost:9001
# Default login: minioadmin / minioadmin
```

### **DigitalOcean Spaces**

```env
# DigitalOcean Spaces Configuration
AWS_ACCESS_KEY_ID=YOUR_SPACES_ACCESS_KEY
AWS_SECRET_ACCESS_KEY=YOUR_SPACES_SECRET_KEY
AWS_REGION=nyc3
AWS_S3_BUCKET=your-space-name
AWS_S3_ENDPOINT=https://nyc3.digitaloceanspaces.com
AWS_S3_USE_SSL=true
AWS_S3_PATH_STYLE=false
```

### **Wasabi**

```env
# Wasabi Configuration
AWS_ACCESS_KEY_ID=YOUR_WASABI_ACCESS_KEY
AWS_SECRET_ACCESS_KEY=YOUR_WASABI_SECRET_KEY
AWS_REGION=us-east-1
AWS_S3_BUCKET=your-wasabi-bucket
AWS_S3_ENDPOINT=https://s3.wasabisys.com
AWS_S3_USE_SSL=true
AWS_S3_PATH_STYLE=false
```

### **Backblaze B2**

```env
# Backblaze B2 Configuration
AWS_ACCESS_KEY_ID=YOUR_B2_KEY_ID
AWS_SECRET_ACCESS_KEY=YOUR_B2_APPLICATION_KEY
AWS_REGION=us-west-002
AWS_S3_BUCKET=your-b2-bucket
AWS_S3_ENDPOINT=https://s3.us-west-002.backblazeb2.com
AWS_S3_USE_SSL=true
AWS_S3_PATH_STYLE=false
```

## ⚙️ **Configuration Parameters**

| Parameter | Description | Example |
|-----------|-------------|---------|
| `AWS_ACCESS_KEY_ID` | Access key untuk authentication | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | Secret key untuk authentication | `wJalrXUtnFEMI/K7MDENG...` |
| `AWS_REGION` | Region storage | `us-east-1` |
| `AWS_S3_BUCKET` | Nama bucket | `my-wms-files` |
| `AWS_S3_ENDPOINT` | Custom endpoint (kosong untuk AWS S3) | `http://localhost:9000` |
| `AWS_S3_USE_SSL` | Gunakan HTTPS/SSL | `true` atau `false` |
| `AWS_S3_PATH_STYLE` | Force path-style URLs | `true` atau `false` |

## 🔄 **Auto-Detection Logic**

System akan otomatis detect storage type berdasarkan `AWS_S3_ENDPOINT`:

- **Jika `AWS_S3_ENDPOINT` kosong** → Gunakan Amazon S3
- **Jika `AWS_S3_ENDPOINT` ada** → Gunakan S3 Compatible Storage

## 📂 **File Structure di Storage**

Files akan disimpan dengan struktur:
```
bucket/
├── 2024/10/04/
│   ├── uuid1_product_image.jpg
│   ├── uuid2_document.pdf
│   └── uuid3_avatar.png
├── 2024/10/05/
│   └── ...
```

## 🧪 **Testing Configuration**

Setelah set environment variables, test connection dengan:

```bash
# Start aplikasi
go run main.go

# Check logs untuk konfirmasi
# Amazon S3: "Configuring Amazon S3"
# S3 Compatible: "Configuring S3 compatible storage with endpoint: ..."
```

## 🚨 **Troubleshooting**

### **Connection Failed**
- Pastikan credentials benar
- Check network connectivity
- Verify endpoint URL (untuk S3 compatible)

### **Permission Denied**
- Pastikan user/key punya permission pada bucket
- Check bucket policy

### **SSL Certificate Error**
- Set `AWS_S3_USE_SSL=false` untuk development
- Atau gunakan valid SSL certificate

## 🔐 **Security Best Practices**

1. **Jangan commit credentials** ke version control
2. **Gunakan IAM roles** untuk production (AWS)
3. **Limit bucket permissions** sesuai kebutuhan
4. **Enable versioning** pada bucket
5. **Set up CORS** jika diperlukan untuk web access

## 🎯 **Recommended Setup**

**Development:**
- MinIO local untuk development
- AWS S3 untuk staging/production

**Production:**
- Amazon S3 dengan IAM roles
- CloudFront untuk CDN (optional)
- Backup strategy