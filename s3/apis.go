package s3

// S3 APIs
//
// Copyright (c) 2024 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"io"
	"log"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
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

// Init function initializes S3 backend storage
func Init() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	var err error
	if srvConfig.Config.S3.Name == "minio" {
		err = minioInitialize()
	} else {
		err = cephInitialize(
			srvConfig.Config.DataManagement.S3.Endpoint,
			srvConfig.Config.DataManagement.S3.AccessKey,
			srvConfig.Config.DataManagement.S3.AccessSecret,
			srvConfig.Config.DataManagement.S3.Region,
		)
	}
	if err != nil {
		log.Fatal("ERROR", err)
	}
}

// BucketContent provides content on given bucket
func BucketContent(bucket string) (BucketObject, error) {
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		var objects []ObjectInfo
		bobj, err := minioBucketContent(bucket)
		for _, obj := range bobj.Objects {
			o := ObjectInfo{
				Name:         obj.Key,
				LastModified: obj.LastModified,
				Size:         obj.Size,
				ContentType:  obj.ContentType,
				Expires:      obj.Expires,
			}
			objects = append(objects, o)
		}
		b := BucketObject{
			Bucket:  bobj.Bucket,
			Objects: objects,
		}
		return b, err
	}
	return cephBucketContent(bucket)
}

// ListBuckets provides list of buckets in S3 store
func ListBuckets() ([]BucketInfo, error) {
	var blist []BucketInfo
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		buckets, err := minioListBuckets()
		for _, bucket := range buckets {
			b := BucketInfo{
				Name:         bucket.Name,
				CreationDate: bucket.CreationDate,
			}
			blist = append(blist, b)
		}
		return blist, err
	}
	return cephListBuckets()
}

// ListObjects provides list of buckets in S3 store
func ListObjects(bucket string) ([]ObjectInfo, error) {
	var olist []ObjectInfo
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		objects, err := minioListObjects(bucket)
		for _, obj := range objects {
			o := ObjectInfo{
				Name:         obj.Key,
				LastModified: obj.LastModified,
				Size:         obj.Size,
				ContentType:  obj.ContentType,
				Expires:      obj.Expires,
			}
			olist = append(olist, o)
		}
		return olist, err
	}
	return cephListObjects(bucket)
}

// CreateBucket creates new bucket in S3 store
func CreateBucket(bucket string) error {
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return minioCreateBucket(bucket)
	}
	return cephCreateBucket(bucket)
}

// DeleteBucket deletes bucket in S3 store
func DeleteBucket(bucket string) error {
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return minioDeleteBucket(bucket)
	}
	return cephDeleteBucket(bucket)
}

// UploadObject uploads given object to S3 store
func UploadObject(bucket, objectName, contentType string, reader io.Reader, size int64) error {
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		_, err := minioUploadObject(bucket, objectName, contentType, reader, size)
		return err
	}
	return cephUploadObject(bucket, objectName, contentType, reader, size)
}

// DeleteObject deletes object from S3 storage
func DeleteObject(bucket, objectName, versionId string) error {
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return minioDeleteObject(bucket, objectName, versionId)
	}
	return cephDeleteObject(bucket, objectName, versionId)
}

// GetObjects returns given object from S3 storage
func GetObject(bucket, objectName string) ([]byte, error) {
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return minioGetObject(bucket, objectName)
	}
	return cephGetObject(bucket, objectName)
}
