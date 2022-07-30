package storageprovider_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pakkaparn/no-idea-api/internal/config"
	"github.com/pakkaparn/no-idea-api/internal/services/fileservice/storageprovider"
	"google.golang.org/api/option"
)

func TestNewGCS(t *testing.T) {
	type args struct {
		options config.Options
	}
	tests := []struct {
		name string
		args args
		want storageprovider.StorageProvider
	}{
		{
			name: "NewGCS",
			args: args{
				options: config.Options{
					GCSProjectID: "project-id",
					GCSBucket:    "bucket",
					AppName:      "app-name",
				},
			},
			want: func() storageprovider.StorageProvider {
				client, _ := storage.NewClient(context.Background())
				return &storageprovider.GCSProvider{
					Client:     client,
					ProjectID:  "project-id",
					BucketName: "bucket",
					UploadPath: "app-name",
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storageprovider.NewGCS(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGCS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGCSWithClient(t *testing.T) {
	type args struct {
		client  *storage.Client
		options config.Options
	}
	tests := []struct {
		name string
		args args
		want storageprovider.StorageProvider
	}{
		{
			name: "NewGCSWithClient",
			args: func() args {
				client, _ := storage.NewClient(context.Background())
				return args{
					client: client,
					options: config.Options{
						GCSProjectID: "project-id",
						GCSBucket:    "bucket",
						AppName:      "app-name",
					},
				}
			}(),
			want: func() storageprovider.StorageProvider {
				client, _ := storage.NewClient(context.Background())
				return &storageprovider.GCSProvider{
					Client:     client,
					ProjectID:  "project-id",
					BucketName: "bucket",
					UploadPath: "app-name",
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storageprovider.NewGCSWithClient(tt.args.client, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGCSWithClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGCSProvider_Upload(t *testing.T) {
	type fields struct {
		Client  *storage.Client
		options config.Options
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
			name: "Upload",
			fields: func() fields {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{\"name\":\"test-file.jpg\"}"))
				}))

				client, _ := storage.NewClient(
					context.Background(),
					option.WithoutAuthentication(),
					option.WithEndpoint(ts.URL),
				)
				return fields{
					Client: client,
					options: config.Options{
						GCSProjectID: "project-id",
						GCSBucket:    "bucket",
						AppName:      "app-name",
					},
				}
			}(),
			args: args{
				file:     []byte("test"),
				filename: "test.txt",
			},
			want: "test.txt",
		},
		{
			name: "Upload fail",
			fields: func() fields {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "error", http.StatusInternalServerError)
				}))
				defer ts.Close()

				client, _ := storage.NewClient(
					context.Background(),
					option.WithoutAuthentication(),
					option.WithEndpoint(ts.URL),
				)
				return fields{
					Client: client,
					options: config.Options{
						GCSProjectID: "project-id",
						GCSBucket:    "bucket",
						AppName:      "app-name",
					},
				}
			}(),
			args: args{
				file:     []byte("test"),
				filename: "test.txt",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := storageprovider.NewGCSWithClient(tt.fields.Client, tt.fields.options)

			got, err := p.Upload(tt.args.file, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GCSProvider.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GCSProvider.Upload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGCSProvider_GetSigedURL(t *testing.T) {
	type fields struct {
		Client  *storage.Client
		options config.Options
	}
	type args struct {
		filename string
		expired  time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *regexp.Regexp
		wantErr bool
	}{
		{
			name: "GetSigedURL",
			fields: func() fields {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "error", http.StatusInternalServerError)
				}))
				defer ts.Close()

				client, _ := storage.NewClient(
					context.Background(),
					option.WithoutAuthentication(),
					option.WithEndpoint(ts.URL),
				)

				return fields{
					Client: client,
					options: config.Options{
						GCSProjectID:     "project-id",
						GCSBucket:        "bucket",
						AppName:          "app-name",
						GoogleAccessID:   "access-id",
						GooglePrivateKey: generatePrivateKey(),
					},
				}
			}(),
			args: args{
				filename: "test.txt",
				expired:  time.Hour,
			},
			want: regexp.MustCompile(
				`https:\/\/storage\.googleapis\.com\/bucket\/app-name\/test\.txt\?Expires=\d+&GoogleAccessId=access-id&Signature=.+`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := storageprovider.NewGCSWithClient(tt.fields.Client, tt.fields.options)
			got, err := p.GetSigedURL(tt.args.filename, tt.args.expired)
			if (err != nil) != tt.wantErr {
				t.Errorf("GCSProvider.GetSigedURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.want.MatchString(got) {
				t.Errorf("GCSProvider.GetSigedURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generatePrivateKey() []byte {
	type PKCS8Key struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}

	var pkey PKCS8Key

	key, _ := rsa.GenerateKey(rand.Reader, 4096)

	pkey.Version = 0
	pkey.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	pkey.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	pkey.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	bytes, _ := asn1.Marshal(pkey)

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: bytes,
		},
	)
}
