package s3

// monio S3 module
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

// S3 represent S3 storage record
type S3 struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	UseSSL       bool
}

// MinioClient represents Minio S3 client
type MinioClient struct {
	S3Client *minio.Client
}

// helper function to get s3 minio client
func (c *MinioClient) Initialize() error {
	var err error

	// get s3 object without any buckets info
	s3 := S3{
		Endpoint:     srvConfig.Config.DataManagement.S3.Endpoint,
		AccessKey:    srvConfig.Config.DataManagement.S3.AccessKey,
		AccessSecret: srvConfig.Config.DataManagement.S3.AccessSecret,
		UseSSL:       srvConfig.Config.DataManagement.S3.UseSSL,
	}
	if srvConfig.Config.DataManagement.WebServer.Verbose > 1 {
		log.Printf("INFO: s3 object %+v", s3)
	}

	// Initialize minio client object.
	c.S3Client, err = minio.New(s3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3.AccessKey, s3.AccessSecret, ""),
		Secure: s3.UseSSL,
	})
	return err
}

// MinioBucketObject represents s3 object
type MinioBucketObject struct {
	Bucket  string             `json:"bucket"`
	Objects []minio.ObjectInfo `json:"objects"`
}

// BucketContent provides content on given bucket
func (c *MinioClient) BucketContent(bucket string) (BucketObject, error) {
	if srvConfig.Config.DataManagement.WebServer.Verbose > 0 {
		log.Printf("looking for bucket:'%s'", bucket)
	}
	objects, err := c.ListObjects(bucket)
	if err != nil {
		log.Printf("ERROR: unabel to list bucket '%s', error %v", bucket, err)
	}
	bobj := BucketObject{
		Bucket:  bucket,
		Objects: objects,
	}
	return bobj, err
}

// helper function to provide list of buckets in S3 store
func (c *MinioClient) ListBuckets() ([]BucketInfo, error) {
	ctx := context.Background()
	buckets, err := c.S3Client.ListBuckets(ctx)

	// convert minio buckets into generic list of BucketInfo objects
	var blist []BucketInfo
	for _, bucket := range buckets {
		b := BucketInfo{
			Name:         bucket.Name,
			CreationDate: bucket.CreationDate,
		}
		blist = append(blist, b)
	}
	return blist, err
}

// helper function to provide list of buckets in S3 store
func (c *MinioClient) ListObjects(bucket string) ([]ObjectInfo, error) {
	var olist []ObjectInfo
	ctx := context.Background()
	// list individual buckets
	objectCh := c.S3Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Recursive: true,
	})
	for obj := range objectCh {
		if obj.Err != nil {
			log.Printf("ERROR: unable to list objects in a bucket, error %v", obj.Err)
			return olist, obj.Err
		}
		//         obj := fmt.Sprintf("%v %s %10d %s\n", object.LastModified, object.ETag, object.Size, object.Key)
		// convert minio obj into generic ObjectInfo
		o := ObjectInfo{
			Name:         obj.Key,
			LastModified: obj.LastModified,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
			Expires:      obj.Expires,
		}
		olist = append(olist, o)
	}
	return olist, nil
}

// helper function to create new bucket in S3 store
func (c *MinioClient) CreateBucket(bucket string) error {
	// get s3 object without any buckets info
	ctx := context.Background()

	// create new bucket on s3 storage
	//     err = c.S3Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: location})
	err := c.S3Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := c.S3Client.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			if srvConfig.Config.DataManagement.WebServer.Verbose > 0 {
				log.Printf("WARNING: we already own %s\n", bucket)
			}
			return nil
		} else {
			log.Printf("ERROR: unable to create bucket, error %v", err)
		}
	} else {
		if srvConfig.Config.DataManagement.WebServer.Verbose > 0 {
			log.Printf("Successfully created %s\n", bucket)
		}
	}
	return err
}

// helper function to delete bucket in S3 store
func (c *MinioClient) DeleteBucket(bucket string) error {
	ctx := context.Background()
	err := c.S3Client.RemoveBucket(ctx, bucket)
	if err != nil {
		log.Printf("ERROR: unable to remove bucket %s, error, %v", bucket, err)
	}
	return err
}

// helper function to upload given object to S3 store
func (c *MinioClient) UploadObject(bucket, objectName, contentType string, reader io.Reader, size int64) error {
	ctx := context.Background()

	// Upload the zip file with PutObject
	options := minio.PutObjectOptions{}
	if contentType != "" {
		options = minio.PutObjectOptions{ContentType: contentType}
	}
	info, err := c.S3Client.PutObject(
		ctx,
		bucket,
		objectName,
		reader,
		size,
		options)
	if err != nil {
		log.Printf("ERROR: fail to upload file object, error %v", err)
	} else {
		if srvConfig.Config.DataManagement.WebServer.Verbose > 0 {
			log.Println("INFO: upload file", info)
		}
	}
	return err
}

// helper function to delete object from S3 storage
func (c *MinioClient) DeleteObject(bucket, objectName, versionId string) error {
	ctx := context.Background()

	// remove given object from our s3 store
	options := minio.RemoveObjectOptions{
		// Set the bypass governance header to delete an object locked with GOVERNANCE mode
		GovernanceBypass: true,
	}
	if versionId != "" {
		options.VersionID = versionId
	}
	err := c.S3Client.RemoveObject(
		ctx,
		bucket,
		objectName,
		options)
	if err != nil {
		log.Printf("ERROR: fail to delete file object, error %v", err)
	}
	return err
}

// helper function to get given object from S3 storage
func (c *MinioClient) GetObject(bucket, objectName string) ([]byte, error) {
	ctx := context.Background()

	// Upload the zip file with PutObject
	options := minio.GetObjectOptions{}
	object, err := c.S3Client.GetObject(
		ctx,
		bucket,
		objectName,
		options)
	if err != nil {
		log.Printf("ERROR: fail to download file object, error %v", err)
	}
	data, err := io.ReadAll(object)
	return data, err
}

// minioGetS3Link generates a URL for an object in the bucket or a bucket itself if objectName is empty.
// If expiresIn is 0, it generates a permanent link (for public buckets or objects with appropriate ACL).
func (c *MinioClient) GetS3Link(bucket, objectName string, expiresIn time.Duration) (string, error) {
	// Permanent URL
	if expiresIn == 0 {
		if objectName == "" {
			// Generate link to the bucket
			return fmt.Sprintf("%s/%s", c.S3Client.EndpointURL().String(), bucket), nil
		}
		// Generate link to the object
		return fmt.Sprintf("%s/%s/%s", c.S3Client.EndpointURL().String(), bucket, objectName), nil
	}

	// Pre-signed URL with expiration
	if objectName == "" {
		return "", fmt.Errorf("cannot generate a pre-signed URL for the bucket itself with an expiration time")
	}

	// Generate a pre-signed URL for the object
	ctx := context.Background()
	url, err := c.S3Client.PresignedGetObject(ctx, bucket, objectName, expiresIn, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL for object %s in bucket %s: %v", objectName, bucket, err)
	}

	return url.String(), nil
}
