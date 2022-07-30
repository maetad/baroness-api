package storageprovider

import (
	"reflect"
	"testing"

	"github.com/pakkaparn/no-idea-api/internal/config"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		storage Provider
		options config.Options
	}
	tests := []struct {
		name string
		args args
		want StorageProvider
	}{
		{
			name: "GCS",
			args: args{
				storage: GCS,
			},
			want: NewGCS(config.Options{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProvider(tt.args.storage, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
