package routes

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
)

type S3Client struct {
    client *s3.Client
    bucket string
    region string
}

func NewS3Client(accessKey, secretKey, bucket, region string) (*S3Client, error) {
    // Create credentials
    creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
    
    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(region),
        config.WithCredentialsProvider(creds),
    )
    if err != nil {
        return nil, fmt.Errorf("unable to load SDK config: %v", err)
    }

    // Create S3 client
    client := s3.NewFromConfig(cfg)

    return &S3Client{
        client: client,
        bucket: bucket,
        region: region,
    }, nil
}

func (s *S3Client) GeneratePresignedURL(key string, contentType string, expirationTime time.Duration) (string, error) {
    presignClient := s3.NewPresignClient(s.client)

    input := &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        ContentType: aws.String(contentType),
    }

    // Generate presigned URL
    resp, err := presignClient.PresignPutObject(context.TODO(), input, func(opts *s3.PresignOptions) {
        opts.Expires = expirationTime
    })
    if err != nil {
        return "", fmt.Errorf("couldn't generate presigned URL: %v", err)
    }

    return resp.URL, nil
}

// Example usage in your main.go:
func initializeS3Client() (*S3Client, error) {
    S3AccessKey := os.Getenv("S3_ACCESS_KEY")
    S3SecretKey := os.Getenv("S3_ACCESS_SECRET")
    S3Bucket := os.Getenv("S3_UPLOADS_BUCKET")
    S3Region := os.Getenv("S3_UPLOADS_REGION")

    s3Client, err := NewS3Client(
        S3AccessKey,
        S3SecretKey,
        S3Bucket,
        S3Region,
    )
    if err != nil {
        return nil, err
    }
    return s3Client, nil
}

// Example route handler
func HandlePresignedURL(c *gin.Context) {
	filename := c.Query("filename")
    contentType := c.Query("contentType")

    if filename == "" || contentType == "" {
        c.JSON(400, gin.H{"error": "Missing filename or contentType"})
        return
    }
    s3Client, err := initializeS3Client()
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to initialize S3 client"})
        return
    }

    // Generate a unique key for the file
    key := fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01/02"), filename)
    
    // Generate presigned URL with 15-minute expiration
    url, err := s3Client.GeneratePresignedURL(key, contentType, 15*time.Minute)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate presigned URL"})
        return
    } 

    c.JSON(200, gin.H{"url": url})
}