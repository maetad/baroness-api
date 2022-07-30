package storageprovider_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pakkaparn/no-idea-api/internal/services/fileservice/storageprovider"
)

func TestNewAWSS3(t *testing.T) {
	type args struct {
		s      *session.Session
		config storageprovider.AWSS3Config
	}
	s := &session.Session{}
	tests := []struct {
		name string
		args args
		want storageprovider.StorageProvider
	}{
		{
			name: "NewAWSS3",
			args: func() args {
				a := args{}
				a.s = s
				return a
			}(),
			want: &storageprovider.AWSS3Provider{
				Session: s,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storageprovider.NewAWSS3(tt.args.s, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAWSS3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAWSS3Session(t *testing.T) {
	type args struct {
		config storageprovider.AWSS3Config
	}
	tests := []struct {
		name string
		args args
		want *session.Session
	}{
		{
			name: "new session",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storageprovider.NewAWSS3Session(tt.args.config); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewAWSS3Session() = %v, want %v", got, tt.want)
			}
		})
	}
}
