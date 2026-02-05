package storage

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "imgpdf"

func envBool(key string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes" || v == "y" || v == "on"
}

func New(ctx context.Context) (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_URL") // host:port
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useTLS := envBool("MINIO_USE_SSL", false)

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useTLS,
	}

	// Опционально: доверенный CA для self-signed сертификата
	// Если переменная не задана — используем системное доверие (по умолчанию).
	if useTLS {
		if caPath := strings.TrimSpace(os.Getenv("MINIO_CA_CERT_PATH")); caPath != "" {
			caPEM, err := os.ReadFile(caPath)
			if err != nil {
				return nil, err
			}
			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(caPEM); !ok {
				return nil, err
			}

			opts.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: pool,
				},
			}
		}
	}

	client, err := minio.New(endpoint, opts)
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		log.Printf("Bucket %s does not exist, creating...", bucketName)
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		log.Printf("Bucket %s created successfully", bucketName)
	}

	return client, nil
}
