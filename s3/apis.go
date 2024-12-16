package s3

import (
	"io"
	"log"

	srvConfig "github.com/CHESSComputing/golib/config"
	minio "github.com/minio/minio-go/v7"
)

// BucketContent provides content on given bucket
func BucketContent(bucket string) (BucketObject, error) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name != "minio" {
		return bucketContent(bucket)
	}
	return BucketObject{}, nil
}

// ListBuckets provides list of buckets in S3 store
func ListBuckets() ([]minio.BucketInfo, error) {
	var blist []minio.BucketInfo
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name != "minio" {
		return blist, nil
	}
	return blist, nil
}

// ListObjects provides list of buckets in S3 store
func ListObjects(bucket string) ([]minio.ObjectInfo, error) {
	var olist []minio.ObjectInfo
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name != "minio" {
		return olist, nil
	}
	return olist, nil
}

// CreateBucket creates new bucket in S3 store
func CreateBucket(bucket string) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name != "minio" {
		return nil
	}
	return nil
}

// DeleteBucket deletes bucket in S3 store
func DeleteBucket(bucket string) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name != "minio" {
		return nil
	}
	return nil
}

// UploadObject uploads given object to S3 store
func UploadObject(bucket, objectName, contentType string, reader io.Reader, size int64) (minio.UploadInfo, error) {
	var info minio.UploadInfo
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name != "minio" {
		return info, nil
	}
	return info, nil
}

// DeleteObject deletes object from S3 storage
func DeleteObject(bucket, objectName, versionId string) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	log.Printf("Use %s S3 storage", srvConfig.Config.S3.Name)
	if srvConfig.Config.S3.Name != "minio" {
		return nil
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
	if srvConfig.Config.S3.Name != "minio" {
		return obj, nil
	}
	return obj, nil
}
