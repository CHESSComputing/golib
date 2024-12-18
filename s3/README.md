# S3 Client
This module provides two implementations based on AWS and Minio SDK for generic S3 client.

Both **AWS SDKs** and the **MinIO client** are compatible with Ceph Object
Storage, as Ceph implements the **Amazon S3-compatible API**. This means you
can interact with Ceph Object Storage using tools and SDKs designed for AWS S3,
including the AWS SDKs or the MinIO client.

However, there are some considerations to ensure compatibility:

---

### 1. **Compatibility**
   - **AWS SDK**: Fully compatible with Ceph as it uses the S3 API. You can use it for operations like creating buckets, uploading files, etc. Just make sure to configure the endpoint properly.
   - **MinIO Client**: Also fully compatible with Ceph and is often more lightweight for scenarios focused on S3-like operations.
   - Both libraries work with Ceph for typical S3 operations such as creating buckets, uploading/downloading objects, deleting objects, and generating pre-signed URLs.

### 2. **Configuration**
   - When using **AWS SDK**, you need to specify the Ceph endpoint (e.g., `http://my-ceph-server:9000`) and credentials (access key and secret key). 
   - For **MinIO**, you do the same by configuring the endpoint and credentials in the `minio.Options`.

---

### 3. **When to Choose AWS vs MinIO**
   - **AWS SDK**:
     - Best if you already use AWS in your application or if you rely on AWS-specific S3 features (e.g., IAM policies, specific S3 bucket policies).
     - Provides extensive AWS ecosystem integration, such as CloudFront, Lambda, etc.
   - **MinIO Client**:
     - Lightweight and focused entirely on object storage operations. It’s simpler for tasks like uploading files or managing buckets in Ceph.
     - Often used with self-hosted S3-compatible object storage like Ceph or MinIO itself.

---

### 4. **Key Differences**
   - **AWS SDK** supports advanced AWS-specific features like S3 Object Lock, Transfer Acceleration, etc., which may not be applicable to Ceph.
   - **MinIO Client** is easier to configure and explicitly designed for S3-compatible storage like Ceph or MinIO.

---

### 5. **Single Implementation**
If you're targeting only **S3-compatible APIs** (like Ceph, MinIO, or AWS S3) and don’t need AWS-specific features, you can implement a **generic interface** that abstracts bucket and object operations. Then, you can use either the MinIO or AWS SDK under the hood. Here's a rough structure:

```go
type S3Client interface {
    CreateBucket(bucket string) error
    UploadFile(bucket, objectName, contentType string, reader io.Reader, size int64) error
    DeleteObject(bucket, objectName string) error
    GetObject(bucket, objectName string) ([]byte, error)
}

type AWSClient struct { /* ... */ }
type MinioClient struct { /* ... */ }

func (c *AWSClient) CreateBucket(bucket string) error {
    // Implement using AWS SDK
}

func (c *MinioClient) CreateBucket(bucket string) error {
    // Implement using MinIO SDK
}

// Use S3Client interface in your application
```

---

### 6. **Ceph-Specific Features**
While both clients work for general S3 operations, if you need to work with **Ceph-specific features** (e.g., native filesystem integration with CephFS or advanced RADOS Gateway features), you may need to use Ceph-specific tools like:

- **radosgw-admin** for low-level management.
- Native libraries like `go-ceph` for CephFS or RADOS-specific functionality.

---

### Summary
- You can use either the AWS SDK or MinIO client to access Ceph Object Storage, and both are compatible with S3 APIs.
- MinIO is lightweight and works great for self-hosted S3-like storage.
- AWS SDK is more feature-rich but may include AWS-specific functionality that’s not applicable to Ceph.
- If you’re only targeting S3-compatible APIs, you can create a unified abstraction that supports both SDKs interchangeably.

### Usage
To use either client you may initialize it as following:
```
package main

import (
	"log"
    s3 "github.com/CHESSComputing/golib/s3"
)

func main() {
	// Assume clientType is provided by user or configuration file.
	clientType := "minio" // Could be "aws" or "minio".

	// Initialize the appropriate S3 client.
	s3Client, err := s3.InitializeS3Client(clientType)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Use the S3 client.
	bucketName := "example-bucket"
	err = s3Client.CreateBucket(bucketName)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	// Upload a file (for example, a simple string converted to io.Reader).
	fileContent := "Hello, S3!"
	err = s3Client.UploadFile(bucketName, "example.txt", "text/plain", bytes.NewReader([]byte(fileContent)), int64(len(fileContent)))
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	// Retrieve the file.
	data, err := s3Client.GetObject(bucketName, "example.txt")
	if err != nil {
		log.Fatalf("Failed to get object: %v", err)
	}
	log.Printf("Retrieved file content: %s", string(data))

	// Delete the file.
	err = s3Client.DeleteObject(bucketName, "example.txt")
	if err != nil {
		log.Fatalf("Failed to delete object: %v", err)
	}
}
```
