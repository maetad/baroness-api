package storageprovider

import (
	"time"
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

func NewProvider(storage Provider, options map[string]interface{}) StorageProvider {
	switch storage {
	case AWSS3:
		config := AWSS3Config{
			Region:          options["region"].(string),
			AccessKeyID:     options["access_key_id"].(string),
			SecretAccessKey: options["secret_access_key"].(string),
			Bucket:          options["bucket_name"].(string),
			AppName:         options["app_name"].(string),
		}
		session := NewAWSS3Session(config)
		return NewAWSS3(session, config)
	case GCS:
		config := GCSConfig{
			GoogleCredential: options["google_credential"].(string),
			GoogleAccessID:   options["google_access_id"].(string),
			ProjectID:        options["project_id"].(string),
			BucketName:       options["bucket_name"].(string),
			UploadPath:       options["upload_path"].(string),
			PrivateKey:       options["private_key"].([]byte),
		}

		return NewGCS(config)
	default:
		return nil
	}
}
