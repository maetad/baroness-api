package storageprovider_test

import (
	"reflect"
	"testing"

	"github.com/pakkaparn/no-idea-api/internal/config"
	"github.com/pakkaparn/no-idea-api/internal/services/fileservice/storageprovider"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		storage storageprovider.Provider
		options config.Options
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
			},
			want: storageprovider.NewGCS(storageprovider.GCSConfig{}),
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
