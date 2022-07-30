package storageprovider

import (
	"time"

	"github.com/pakkaparn/no-idea-api/internal/config"
)

type StorageProvider interface {
	Upload(file []byte, filename string) (string, error)
	GetSigedURL(filename string, expired time.Duration) (string, error)
}

type Provider string

var (
	AWSS3 Provider = "AWSS3"
	GCS   Provider = "GCS"
)

func NewProvider(storage Provider, options config.Options) StorageProvider {
	switch storage {
	case AWSS3:
		config := AWSS3Config{}
		session := NewAWSS3Session(config)
		return NewAWSS3(session, config)
	case GCS:
		config := GCSConfig{}
		return NewGCS(config)
	default:
		return nil
	}
}
