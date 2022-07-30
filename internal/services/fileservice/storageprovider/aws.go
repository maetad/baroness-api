package storageprovider

import (
	"bytes"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSS3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	AppName         string
}

type AWSS3Provider struct {
	Session    *session.Session
	Bucket     string
	UploadPath string
}

func NewAWSS3Session(config AWSS3Config) *session.Session {
	session, _ := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyID,
			config.SecretAccessKey,
			"", // a token will be created when the session it's used.
		),
	})

	return session
}

func NewAWSS3(s *session.Session, config AWSS3Config) StorageProvider {
	return &AWSS3Provider{
		s,
		config.Bucket,
		config.AppName,
	}
}

func (p *AWSS3Provider) Upload(file []byte, filename string) (string, error) {
	up, err := s3manager.
		NewUploader(p.Session).
		Upload(&s3manager.UploadInput{
			Bucket: aws.String(p.Bucket),
			ACL:    aws.String("public-read"),
			Key:    aws.String(p.UploadPath + "/" + filename),
			Body:   bytes.NewReader(file),
		})

	if err != nil {
		return "", err
	}

	return up.Location, nil
}

func (p *AWSS3Provider) GetSigedURL(filename string, expired time.Duration) (string, error) {
	req, _ := s3.New(p.Session).
		GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(p.Bucket),
			Key:    aws.String(filename),
		})

	url, err := req.Presign(expired)

	if err != nil {
		return "", err
	}

	return url, nil
}
