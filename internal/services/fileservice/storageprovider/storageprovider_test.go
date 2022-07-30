package storageprovider_test

import (
	"reflect"
	"testing"

	"github.com/maetad/baroness-api/internal/services/fileservice/storageprovider"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		storage storageprovider.Provider
		options map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want storageprovider.StorageProvider
	}{
		{
			name: "GCS",
			args: args{
				storage: storageprovider.GCS,
				options: map[string]interface{}{
					"google_credential": "google-credential",
					"google_access_id":  "google-access-id",
					"project_id":        "project-id",
					"bucket_name":       "bucket-name",
					"upload_path":       "upload-path",
					"private_key":       []byte("private-key"),
				},
			},
			want: storageprovider.NewGCS(storageprovider.GCSConfig{
				GoogleCredential: "google-credential",
				GoogleAccessID:   "google-access-id",
				ProjectID:        "project-id",
				BucketName:       "bucket-name",
				UploadPath:       "upload-path",
				PrivateKey:       []byte("private-key"),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storageprovider.NewProvider(tt.args.storage, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
