package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "imgpdf"

func New(ctx context.Context) (*minio.Client, error) {
	// client, err := minio.New("localhost:9000", &minio.Options{
	// 	Creds:  credentials.NewStaticV4("imgpdf", "imgpdf12345", ""),
	// 	Secure: false,
	// })
	client, err := minio.New(os.Getenv("MINIO_URL"), &minio.Options{
		Creds: credentials.NewStaticV4(
			os.Getenv("MINIO_ROOT_USER"),
			os.Getenv("MINIO_ROOT_PASSWORD"),
			"",
		),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	// Проверяем существование bucket, если нет - создаем
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		log.Printf("Bucket %s does not exist, creating...", bucketName)
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Bucket %s created successfully", bucketName)
	}

	return client, nil
}
