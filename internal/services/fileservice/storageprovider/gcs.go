package storageprovider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSConfig struct {
	GoogleCredential string
	GoogleAccessID   string
	ProjectID        string
	BucketName       string
	UploadPath       string
	PrivateKey       []byte
}

type GCSProvider struct {
	Client *storage.Client
	GCSConfig
}

func NewGCS(config GCSConfig) StorageProvider {
	client, _ := storage.NewClient(
		context.Background(),
		option.WithCredentialsFile(config.GoogleCredential),
	)

	return NewGCSWithClient(client, config)
}

func NewGCSWithClient(client *storage.Client, config GCSConfig) StorageProvider {
	return &GCSProvider{
		Client:    client,
		GCSConfig: config,
	}
}

func (p *GCSProvider) Upload(file []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	wc := p.Client.
		Bucket(p.BucketName).
		Object(p.UploadPath + "/" + filename).
		NewWriter(ctx)
	defer wc.Close()

	if _, err := io.Copy(wc, bytes.NewReader(file)); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return filename, nil
}

func (p *GCSProvider) GetSigedURL(filename string, expired time.Duration) (string, error) {
	return storage.SignedURL(
		p.BucketName,
		p.UploadPath+"/"+filename,
		&storage.SignedURLOptions{
			GoogleAccessID: p.GoogleAccessID,
			PrivateKey:     p.PrivateKey,
			Method:         http.MethodGet,
			Expires:        time.Now().Add(expired),
		},
	)
}
