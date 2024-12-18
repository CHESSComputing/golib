package s3

import (
	"bytes"
	"fmt"
	"io"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	aws3 "github.com/aws/aws-sdk-go/service/s3"
)

// AWSClient represents AWS S3 client
type AWSClient struct {
	S3Client *aws3.S3
}

// Initialize initializes the S3 client for Ceph
func (c *AWSClient) Initialize() error {
	endpoint := srvConfig.Config.DataManagement.S3.Endpoint
	accessKey := srvConfig.Config.DataManagement.S3.AccessKey
	secretKey := srvConfig.Config.DataManagement.S3.AccessSecret
	region := srvConfig.Config.DataManagement.S3.Region
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region), // Region is needed even for Ceph.
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true), // Needed for Ceph's S3 compatibility
	})
	if err != nil {
		return err
	}
	c.S3Client = aws3.New(sess)
	return nil
}

// cephListBuckets retrieves all available buckets
func (c *AWSClient) ListBuckets() ([]BucketInfo, error) {
	output, err := c.S3Client.ListBuckets(&aws3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to list buckets: %v", err)
	}

	var buckets []BucketInfo
	for _, b := range output.Buckets {
		buckets = append(buckets, BucketInfo{
			Name:         aws.StringValue(b.Name),
			CreationDate: aws.TimeValue(b.CreationDate),
		})
	}
	return buckets, nil
}

// cephCreateBucket creates a new bucket
func (c *AWSClient) CreateBucket(bucket string) error {
	_, err := c.S3Client.CreateBucket(&aws3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("unable to create bucket %s: %v", bucket, err)
	}
	return nil
}

// cephDeleteBucket deletes an existing bucket
func (c *AWSClient) DeleteBucket(bucket string) error {
	_, err := c.S3Client.DeleteBucket(&aws3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("unable to delete bucket %s: %v", bucket, err)
	}
	return nil
}

// cephListObjects lists all objects in a bucket
func (c *AWSClient) ListObjects(bucket string) ([]ObjectInfo, error) {
	output, err := c.S3Client.ListObjectsV2(&aws3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list objects in bucket %s: %v", bucket, err)
	}

	var objects []ObjectInfo
	for _, o := range output.Contents {
		objects = append(objects, ObjectInfo{
			Name:         aws.StringValue(o.Key),
			LastModified: aws.TimeValue(o.LastModified),
			Size:         aws.Int64Value(o.Size),
		})
	}
	return objects, nil
}

// cephBucketContent retrieves all objects in a bucket
func (c *AWSClient) BucketContent(bucket string) (BucketObject, error) {
	objects, err := c.ListObjects(bucket)
	if err != nil {
		return BucketObject{}, err
	}
	return BucketObject{
		Bucket:  bucket,
		Objects: objects,
	}, nil
}

// cephUploadObject uploads an object to a bucket
func (c *AWSClient) UploadObject(bucket, objectName, contentType string, reader io.Reader, size int64) error {
	// Wrap the reader using aws.ReadSeekCloser
	readSeeker := aws.ReadSeekCloser(reader)

	_, err := c.S3Client.PutObject(&aws3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectName),
		Body:          readSeeker,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return fmt.Errorf("unable to upload object %s to bucket %s: %v", objectName, bucket, err)
	}
	return nil
}

// cephGetObject retrieves an object from a bucket
func (c *AWSClient) GetObject(bucket, objectName string) ([]byte, error) {
	output, err := c.S3Client.GetObject(&aws3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get object %s from bucket %s: %v", objectName, bucket, err)
	}
	defer output.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, output.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read object %s: %v", objectName, err)
	}
	return buf.Bytes(), nil
}

// cephDeleteObject deletes an object from a bucket
func (c *AWSClient) DeleteObject(bucket, objectName, versionId string) error {
	_, err := c.S3Client.DeleteObject(&aws3.DeleteObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(objectName),
		VersionId: aws.String(versionId),
	})
	if err != nil {
		return fmt.Errorf("unable to delete object %s from bucket %s: %v", objectName, bucket, err)
	}
	return nil
}

// GetS3Link generates a URL for an object in the S3 bucket or a bucket itself if objectName is empty.
// If expiresIn is 0, it generates a permanent link (for public buckets or objects with appropriate ACL).
func (c *AWSClient) GetS3Link(bucket, objectName string, expiresIn time.Duration) (string, error) {
	endpoint := c.S3Client.Endpoint // Get the endpoint configured in the S3 client

	// Permanent URL
	if expiresIn == 0 {
		if objectName == "" {
			// Generate link to the bucket
			return fmt.Sprintf("%s/%s", endpoint, bucket), nil
		}
		// Generate link to the object
		return fmt.Sprintf("%s/%s/%s", endpoint, bucket, objectName), nil
	}

	// Pre-signed URL with expiration
	if objectName == "" {
		return "", fmt.Errorf("cannot generate a pre-signed URL for the bucket itself with an expiration time")
	}

	// Create a request to get the object
	req, _ := c.S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectName),
	})

	// Generate a pre-signed URL
	url, err := req.Presign(expiresIn)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL for object %s in bucket %s: %v", objectName, bucket, err)
	}

	return url, nil
}
