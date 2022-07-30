package storageprovider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pakkaparn/no-idea-api/internal/config"
	"google.golang.org/api/option"
)

type GCSProvider struct {
	Client         *storage.Client
	ProjectID      string
	BucketName     string
	UploadPath     string
	googleAccessID string
	privateKey     []byte
}

func NewGCS(options config.Options) StorageProvider {
	client, _ := storage.NewClient(
		context.Background(),
		option.WithCredentialsFile(options.GoogleCredential),
	)

	return NewGCSWithClient(client, options)
}

func NewGCSWithClient(client *storage.Client, options config.Options) StorageProvider {
	return &GCSProvider{
		Client:         client,
		ProjectID:      options.GCSProjectID,
		BucketName:     options.GCSBucket,
		UploadPath:     options.AppName,
		googleAccessID: options.GoogleAccessID,
		privateKey:     options.GooglePrivateKey,
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
			GoogleAccessID: p.googleAccessID,
			PrivateKey:     p.privateKey,
			Method:         http.MethodGet,
			Expires:        time.Now().Add(expired),
		},
	)
}
