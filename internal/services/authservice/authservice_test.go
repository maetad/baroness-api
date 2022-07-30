package authservice_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pakkaparn/no-idea-api/internal/services/authservice"
	"github.com/pakkaparn/no-idea-api/mocks"
)

var PEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----
`

var claimer = &mocks.Claimer{}
var jwtPattern = `^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-\+\/=]*)`
var jwtRegex = regexp.MustCompile(jwtPattern)
var RSAPublicKey, _ = jwt.ParseRSAPublicKeyFromPEM([]byte(PEM))

func TestAuthService_GenerateToken(t *testing.T) {
	type fields struct {
		signingMethod      jwt.SigningMethod
		signingKey         interface{}
		allowSigningMethod authservice.AllowSigningMethod
	}
	type args struct {
		c authservice.Claimer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "token generated",
			fields: fields{
				signingMethod: jwt.SigningMethodHS256,
				signingKey:    []byte("signing-key"),
				allowSigningMethod: authservice.AllowSigningMethod{
					HMAC: true,
				},
			},
			args: args{claimer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claimer.Mock.ExpectedCalls = nil
			claimer.On("GetClaims").Return(map[string]interface{}{
				"username": "admin",
			})

			s := authservice.New(tt.fields.signingMethod, tt.fields.signingKey, tt.fields.allowSigningMethod)
			got, err := s.GenerateToken(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !jwtRegex.Match([]byte(got)) {
				t.Errorf("AuthService.GenerateToken() = %v, want %v", got, jwtPattern)
			}

			claimer.AssertNumberOfCalls(t, "GetClaims", 1)
		})
	}
}

func TestAuthService_ParseToken(t *testing.T) {
	type fields struct {
		signingMethod      jwt.SigningMethod
		signingKey         interface{}
		allowSigningMethod authservice.AllowSigningMethod
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    jwt.MapClaims
		wantErr bool
	}{
		{
			name: "parsed complete",
			fields: fields{
				signingMethod: jwt.SigningMethodHS256,
				signingKey:    []byte("signing-key"),
				allowSigningMethod: authservice.AllowSigningMethod{
					HMAC: true,
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIn0.6nXKowsHAmLw7NDHQatY_WK6TVee4qMeN4Mm6wRMokA",
			},
			want: jwt.MapClaims{
				"username": "admin",
			},
		},
		{
			name: "wrong signing key",
			fields: fields{
				signingMethod: jwt.SigningMethodHS256,
				signingKey:    []byte("signing-key"),
				allowSigningMethod: authservice.AllowSigningMethod{
					HMAC: true,
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIn0.woPwc2H-gBMrpHPsCjR2vaQRkv8j24jf4B67ACKBAnA",
			},
			wantErr: true,
		},
		{
			name: "expired token",
			fields: fields{
				signingMethod: jwt.SigningMethodHS256,
				signingKey:    []byte("signing-key"),
				allowSigningMethod: authservice.AllowSigningMethod{
					HMAC: true,
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxMDAwMDAwMDAwfQ.Ea3e-IE_qOhFuqjhyj0JhZcFrH4rQVEuANETUnSYOyU",
			},
			wantErr: true,
		},
		{
			name: "RSA alg is not allow",
			fields: fields{
				signingMethod: jwt.SigningMethodRS256,
				signingKey:    RSAPublicKey,
				allowSigningMethod: authservice.AllowSigningMethod{
					HMAC: true,
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ",
			},
			wantErr: true,
		},
		{
			name: "ECDSA alg is not allow",
			fields: fields{
				signingMethod: jwt.SigningMethodES256,
				signingKey:    RSAPublicKey,
				allowSigningMethod: authservice.AllowSigningMethod{
					HMAC: true,
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA",
			},
			wantErr: true,
		},
		{
			name: "RSAPSS alg is not allow",
			fields: fields{
				signingMethod: jwt.SigningMethodPS256,
				signingKey:    RSAPublicKey,
				allowSigningMethod: authservice.AllowSigningMethod{
					HMAC: true,
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.iOeNU4dAFFeBwNj6qdhdvm-IvDQrTa6R22lQVJVuWJxorJfeQww5Nwsra0PjaOYhAMj9jNMO5YLmud8U7iQ5gJK2zYyepeSuXhfSi8yjFZfRiSkelqSkU19I-Ja8aQBDbqXf2SAWA8mHF8VS3F08rgEaLCyv98fLLH4vSvsJGf6ueZSLKDVXz24rZRXGWtYYk_OYYTVgR1cg0BLCsuCvqZvHleImJKiWmtS0-CymMO4MMjCy_FIl6I56NqLE9C87tUVpo1mT-kbg5cHDD8I7MjCW5Iii5dethB4Vid3mZ6emKjVYgXrtkOQ-JyGMh6fnQxEFN1ft33GX2eRHluK9eg",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := authservice.New(tt.fields.signingMethod, tt.fields.signingKey, tt.fields.allowSigningMethod)
			got, err := s.ParseToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.ParseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllowSigningMethod_Allowed(t *testing.T) {
	type fields struct {
		ECDSA   bool
		Ed25519 bool
		HMAC    bool
		RSA     bool
		RSAPSS  bool
	}
	type args struct {
		k string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *authservice.AllowSigningMethod
	}{
		{
			name: "key exists",
			args: args{
				k: "ECDSA",
			},
			want: &authservice.AllowSigningMethod{
				ECDSA:   true,
				Ed25519: false,
				HMAC:    false,
				RSA:     false,
				RSAPSS:  false,
			},
		},
		{
			name: "key not exists",
			args: args{
				k: "KEY",
			},
			want: &authservice.AllowSigningMethod{
				ECDSA:   false,
				Ed25519: false,
				HMAC:    false,
				RSA:     false,
				RSAPSS:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authservice.AllowSigningMethod{
				ECDSA:   tt.fields.ECDSA,
				Ed25519: tt.fields.Ed25519,
				HMAC:    tt.fields.HMAC,
				RSA:     tt.fields.RSA,
				RSAPSS:  tt.fields.RSAPSS,
			}
			a.Allowed(tt.args.k)
			if !reflect.DeepEqual(a, tt.want) {
				t.Errorf("AllowSigningMethod.Allow() got = %v, want %v", a, tt.want)
			}
		})
	}
}
