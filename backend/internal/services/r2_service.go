package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type R2Service struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewR2Service() (*R2Service, error) {
	// Get credentials from environment
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	endpoint := os.Getenv("R2_ENDPOINT")
	publicURL := os.Getenv("R2_PUBLIC_URL")

	if accountID == "" || accessKeyID == "" || secretAccessKey == "" || bucketName == "" || endpoint == "" {
		return nil, fmt.Errorf("missing required R2 environment variables")
	}

	// Create custom resolver for R2
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL: endpoint,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})

	// Configure AWS SDK for R2
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load R2 config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &R2Service{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

// UploadPhoto uploads a photo to R2 and returns the public URL
func (s *R2Service) UploadPhoto(file multipart.File, header *multipart.FileHeader, jobID string) (string, error) {
	// Read file content
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg" // Default to jpg if no extension
	}
	
	// Create path: jobs/{jobID}/{timestamp}_{uuid}{ext}
	filename := fmt.Sprintf("jobs/%s/%d_%s%s", 
		jobID, 
		time.Now().Unix(), 
		uuid.New().String(),
		ext,
	)

	// Determine content type
	contentType := "image/jpeg"
	switch strings.ToLower(ext) {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	// Upload to R2
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
		// Note: R2 doesn't support ACL, files are public based on bucket policy
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	// Generate a presigned URL for accessing the photo
	// This works even if the bucket isn't public
	presignClient := s3.NewPresignClient(s.client)
	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(7 * 24 * time.Hour) // 7 days
	})
	if err != nil {
		// Fallback to public URL if presigning fails
		publicURL := fmt.Sprintf("%s/%s", strings.TrimRight(s.publicURL, "/"), filename)
		return publicURL, nil
	}
	
	return presignResult.URL, nil
}

// DeletePhoto deletes a photo from R2
func (s *R2Service) DeletePhoto(photoURL string) error {
	// Extract the key from the URL
	key := strings.TrimPrefix(photoURL, s.publicURL+"/")
	
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	
	return err
}

// DeleteJobPhotos deletes all photos for a specific job
func (s *R2Service) DeleteJobPhotos(jobID string) error {
	// List all objects with the job prefix
	prefix := fmt.Sprintf("jobs/%s/", jobID)
	
	resp, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	// Delete each object
	for _, obj := range resp.Contents {
		_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    obj.Key,
		})
		if err != nil {
			return fmt.Errorf("failed to delete object %s: %w", *obj.Key, err)
		}
	}

	return nil
}

// TestConnection tests the R2 connection by checking bucket access
func (s *R2Service) TestConnection(ctx context.Context) error {
	// Try to head the bucket to test connectivity
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	})
	return err
}
