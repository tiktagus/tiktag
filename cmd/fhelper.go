package cmd

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func getFileContentType(ouput *os.File) (string, error) {
	// to sniff the content type only the first
	// 512 bytes are used.
	buf := make([]byte, 512)
	_, err := ouput.Read(buf)
	if err != nil {
		return "", err
	}
	// the function that actually does the trick
	contentType := http.DetectContentType(buf)
	return contentType, nil
}

func getFileHash(fn string) (string, string) {
	// sha256
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	contentType, err := getFileContentType(f)
	if err != nil {
		log.Fatal(err)
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fHash := fmt.Sprintf("%x", h.Sum(nil))
	return fHash, contentType
}

const (
	Endpoint        string = "minio.endpoint"
	AccessKeyID            = "minio.accessKey"
	SecretAccessKey        = "minio.secretKey"
	UseSSL                 = "minio.useSSL"
	BucketName             = "minio.bucketName"
)

func publishFile(id uint64, fn string, contentType string) (string, string) {
	// Minio
	ctx := context.Background()
	endpoint := viper.GetString(Endpoint)
	accessKeyID := viper.GetString(AccessKeyID)
	secretAccessKey := viper.GetString(SecretAccessKey)
	useSSL := viper.GetBool(UseSSL)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := viper.GetString(BucketName)

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			// log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the file
	// _, file := filepath.Split(fn)
	ext := filepath.Ext(fn)
	objectName := fmt.Sprintf("%d%s", id, ext)
	filePath := fn

	// FIXME: check object name exists
	_, err = minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	// log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	url := fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, objectName)
	return url, ext
}
