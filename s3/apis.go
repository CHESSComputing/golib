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

// BucketContent provides content on given bucket
func BucketContent(bucket string) (BucketObject, error) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		var objects []ObjectInfo
		bobj, err := bucketContent(bucket)
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
	return BucketObject{}, nil
}

// ListBuckets provides list of buckets in S3 store
func ListBuckets() ([]BucketInfo, error) {
	var blist []BucketInfo
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		buckets, err := listBuckets()
		for _, bucket := range buckets {
			b := BucketInfo{
				Name:         bucket.Name,
				CreationDate: bucket.CreationDate,
			}
			blist = append(blist, b)
		}
		return blist, err
	}
	return blist, nil
}

// ListObjects provides list of buckets in S3 store
func ListObjects(bucket string) ([]ObjectInfo, error) {
	var olist []ObjectInfo
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		objects, err := listObjects(bucket)
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
	return olist, nil
}

// CreateBucket creates new bucket in S3 store
func CreateBucket(bucket string) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return createBucket(bucket)
	}
	return nil
}

// DeleteBucket deletes bucket in S3 store
func DeleteBucket(bucket string) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return deleteBucket(bucket)
	}
	return nil
}

// UploadObject uploads given object to S3 store
func UploadObject(bucket, objectName, contentType string, reader io.Reader, size int64) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		_, err := uploadObject(bucket, objectName, contentType, reader, size)
		return err
	}
	return nil
}

// DeleteObject deletes object from S3 storage
func DeleteObject(bucket, objectName, versionId string) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return deleteObject(bucket, objectName, versionId)
	}
	return nil
}

// GetObjects returns given object from S3 storage
func GetObject(bucket, objectName string) ([]byte, error) {
	var obj []byte
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name == "minio" {
		return getObject(bucket, objectName)
	}
	return obj, nil
}
