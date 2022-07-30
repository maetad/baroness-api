package fileservice

import (
	"github.com/maetad/baroness-api/internal/services/fileservice/storageprovider"
)

type FileServiceInterface interface {
	Upload(file []byte, filename string) (string, error)
}

type FileService struct {
	provider storageprovider.StorageProvider
}

func New(provider storageprovider.StorageProvider) FileServiceInterface {
	return &FileService{provider}
}

func (s *FileService) Upload(file []byte, filename string) (string, error) {
	return s.provider.Upload(file, filename)
}
