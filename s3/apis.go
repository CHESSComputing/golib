package s3

// S3 APIs
//
// Copyright (c) 2024 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

// ObjectInfo provides information about S3 object
type ObjectInfo struct {
	Name         string    `json:"name"`          // Name of the object
	LastModified time.Time `json:"last_modified"` // Date and time the object was last modified.
	Size         int64     `json:"size"`          // Size in bytes of the object.
	ContentType  string    `json:"content_type"`  // A standard MIME type describing the format of the object data.
	Expires      time.Time `json:"expires"`       // The date and time at which the object is no longer able to be cached.
}

// BucketInfo provides information about S3 bucket
type BucketInfo struct {
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
}

// BucketObject represents S3 bucket object
type BucketObject struct {
	Bucket  string       `json:"bucket"`
	Objects []ObjectInfo `json:"objects"`
}

// Generic S3Client interface
type S3Client interface {
	Initialize() error
	BucketContent(bucket string) (BucketObject, error)
	ListBuckets() ([]BucketInfo, error)
	ListObjects(bucket string) ([]ObjectInfo, error)
	CreateBucket(bucket string) error
	DeleteBucket(bucket string) error
	UploadObject(bucket, objectName, contentType string, reader io.Reader, size int64) error
	DeleteObject(bucket, objectName, versionId string) error
	GetObject(bucket, objectName string) ([]byte, error)
	GetS3Link(bucket, objectName string, expiresIn time.Duration) (string, error)
}

// InitializeS3Client initializes either AWSClient or MinioClient based on the option.
func InitializeS3Client(clientType string) (S3Client, error) {
	var s3Client S3Client
	var err error
	switch clientType {
	case "aws":
		log.Println("Initializing AWS Client")
		// Initialize and return AWSClient
		s3Client = &AWSClient{}
	case "minio":
		log.Println("Initializing MinIO Client")
		// Initialize and return MinioClient
		s3Client = &MinioClient{}
	default:
		err = errors.New(fmt.Sprintf("Unsupported client type: %s", clientType))
	}
	if s3Client != nil {
		s3Client.Initialize()
	}
	return s3Client, err
}
