package s3

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Endpoint        string // For S3 compatible storage
	UseSSL          bool
	PathStyle       bool // Force path-style URLs
}

type S3Client struct {
	client *s3.Client
	config S3Config
}

// NewS3Client creates a new S3 client based on environment configuration
func NewS3Client() *S3Client {
	cfg := S3Config{
		AccessKeyID:     getEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
		Region:          getEnvOrDefault("AWS_REGION", "us-east-1"),
		Bucket:          getEnvOrDefault("AWS_S3_BUCKET", ""),
		Endpoint:        getEnvOrDefault("AWS_S3_ENDPOINT", ""), // For S3 compatible storage
		UseSSL:          getEnvOrDefault("AWS_S3_USE_SSL", "true") == "true",
		PathStyle:       getEnvOrDefault("AWS_S3_PATH_STYLE", "false") == "true",
	}

	var awsConfig aws.Config
	var err error

	// Check if using S3 compatible storage (has endpoint) or Amazon S3
	if cfg.Endpoint != "" {
		// S3 Compatible Storage (MinIO, DigitalOcean Spaces, etc.)
		log.Printf("Configuring S3 compatible storage with endpoint: %s", cfg.Endpoint)

		awsConfig, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			)),
		)
		if err != nil {
			log.Printf("Failed to load S3 compatible config: %v", err)
			return nil
		}

		// Create S3 client with custom endpoint
		client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = cfg.PathStyle

			// Handle SSL settings
			if !cfg.UseSSL {
				o.BaseEndpoint = aws.String(strings.Replace(cfg.Endpoint, "https://", "http://", 1))
			}
		})

		return &S3Client{
			client: client,
			config: cfg,
		}
	} else {
		// Amazon S3
		log.Println("Configuring Amazon S3")

		awsConfig, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.Region),
		)
		if err != nil {
			log.Printf("Failed to load AWS config: %v", err)
			return nil
		}

		// If credentials are provided via environment, use them
		if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
			awsConfig.Credentials = credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			)
		}

		client := s3.NewFromConfig(awsConfig)

		return &S3Client{
			client: client,
			config: cfg,
		}
	}
}

// UploadFile uploads a file to S3 with dynamic path and returns the file URL
func (c *S3Client) UploadFile(fileData []byte, modelType string, modelID uint, ext string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	if c.config.Bucket == "" {
		return "", fmt.Errorf("S3 bucket not configured")
	}

	// Get environment for path prefix
	appEnv := getEnvOrDefault("APP_ENV", "development")
	envPrefix := "dev"
	if appEnv == "production" {
		envPrefix = "prod"
	} else if appEnv == "staging" {
		envPrefix = "staging"
	}

	// Generate dynamic path: env/modelType/modelID.ext
	fileName := fmt.Sprintf("%s/%s/%d.%s", envPrefix, modelType, modelID, ext)

	// Upload to S3
	_, err := c.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileData),
		Metadata: map[string]string{
			"model-type":  modelType,
			"model-id":    fmt.Sprintf("%d", modelID),
			"environment": appEnv,
			"upload-time": time.Now().Format(time.RFC3339),
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Generate file URL
	fileURL := c.generateFileURL(fileName)

	log.Printf("File uploaded successfully: %s", fileURL)
	return fileURL, nil
}

// DeleteFile deletes a file from S3
func (c *S3Client) DeleteFile(fileURL string) error {
	if c.client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	if c.config.Bucket == "" {
		return fmt.Errorf("S3 bucket not configured")
	}

	// Extract key from URL
	key := c.extractKeyFromURL(fileURL)
	if key == "" {
		return fmt.Errorf("invalid file URL")
	}

	_, err := c.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	log.Printf("File deleted successfully: %s", fileURL)
	return nil
}

// GetFileURL generates a presigned URL for file access
func (c *S3Client) GetFileURL(key string, expiration time.Duration) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	presignClient := s3.NewPresignClient(c.client)

	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return req.URL, nil
}

// TestConnection tests the S3 connection
func (c *S3Client) TestConnection() error {
	if c.client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	if c.config.Bucket == "" {
		return fmt.Errorf("S3 bucket not configured")
	}

	// Try to head the bucket
	_, err := c.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(c.config.Bucket),
	})

	if err != nil {
		return fmt.Errorf("S3 connection test failed: %v", err)
	}

	log.Printf("S3 connection test successful for bucket: %s", c.config.Bucket)
	return nil
}

// Helper functions

func (c *S3Client) generateFileURL(key string) string {
	if c.config.Endpoint != "" {
		// S3 Compatible Storage
		if c.config.PathStyle {
			return fmt.Sprintf("%s/%s/%s", c.config.Endpoint, c.config.Bucket, key)
		}
		return fmt.Sprintf("https://%s.%s/%s", c.config.Bucket, strings.TrimPrefix(c.config.Endpoint, "https://"), key)
	} else {
		// Amazon S3
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.config.Bucket, c.config.Region, key)
	}
}

func (c *S3Client) extractKeyFromURL(fileURL string) string {
	// For S3 compatible storage
	if c.config.Endpoint != "" {
		if c.config.PathStyle {
			// Format: https://endpoint/bucket/key
			parts := strings.Split(fileURL, "/")
			if len(parts) >= 4 {
				return strings.Join(parts[4:], "/")
			}
		} else {
			// Format: https://bucket.endpoint/key
			parts := strings.Split(fileURL, "/")
			if len(parts) >= 4 {
				return strings.Join(parts[3:], "/")
			}
		}
	} else {
		// Amazon S3 format: https://bucket.s3.region.amazonaws.com/key
		parts := strings.Split(fileURL, "/")
		if len(parts) >= 4 {
			return strings.Join(parts[3:], "/")
		}
	}
	return ""
}

func cleanFileName(fileName string) string {
	// Remove extension and clean the name
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	// Replace spaces and special characters
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	// Keep only alphanumeric and underscore
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetBucketName returns the configured bucket name
func (s *S3Client) GetBucketName() string {
	return s.config.Bucket
}

// GetEndpoint returns the configured endpoint
func (s *S3Client) GetEndpoint() string {
	if s.config.Endpoint != "" {
		return s.config.Endpoint
	}
	return "Amazon S3"
}

// IsSSLEnabled returns whether SSL is enabled
func (s *S3Client) IsSSLEnabled() bool {
	return s.config.UseSSL
}

// IsPathStyleEnabled returns whether path-style URLs are enabled
func (s *S3Client) IsPathStyleEnabled() bool {
	return s.config.PathStyle
}
