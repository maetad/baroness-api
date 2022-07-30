package fileservice_test

import (
	"errors"
	"testing"

	"github.com/maetad/baroness-api/internal/services/fileservice"
	"github.com/maetad/baroness-api/internal/services/fileservice/storageprovider"
	"github.com/maetad/baroness-api/mocks"
)

func TestFileService_Upload(t *testing.T) {
	type fields struct {
		provider storageprovider.StorageProvider
	}
	type args struct {
		file     []byte
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Upload success",
			fields: func() fields {
				provider := &mocks.StorageProvider{}
				provider.On("Upload", []byte("image"), "image.jpg").Return("image.jpg", nil)
				return fields{provider}
			}(),
			args: args{
				file:     []byte("image"),
				filename: "image.jpg",
			},
			want: "image.jpg",
		},
		{
			name: "Upload fail",
			fields: func() fields {
				provider := &mocks.StorageProvider{}
				provider.On("Upload", []byte("image"), "image.jpg").Return("", errors.New("error"))
				return fields{provider}
			}(),
			args: args{
				file:     []byte("image"),
				filename: "image.jpg",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fileservice.New(tt.fields.provider)
			got, err := s.Upload(tt.args.file, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileService.Upload() = %v, want %v", got, tt.want)
			}
		})
	}
}
