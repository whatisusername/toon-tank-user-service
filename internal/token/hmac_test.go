package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBase64HMAC(t *testing.T) {
	type args struct {
		key     string
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid Input",
			args: args{
				key:     "secret",
				message: "message",
			},
			want:    "i19IcCmVwVmMVz2x4hhmqbgl1KeU0WnXBgoDYFeWNgs=",
			wantErr: false,
		},
		{
			name: "Empty Key",
			args: args{
				key:     "",
				message: "message",
			},
			want:    "6wjB9W1d3uB/e9+ARoCD2ga2TPT6xk/jqQiD31/qyuQ=",
			wantErr: false,
		},
		{
			name: "Empty Message",
			args: args{
				key:     "key",
				message: "",
			},
			want:    "XV0TlWPJW1lnub2ajJsjOp3ttFByeUzSMtwbdIMmB9A=",
			wantErr: false,
		},
		{
			name: "Empty Key And Message",
			args: args{
				key:     "",
				message: "",
			},
			want:    "thNnmggU2ex3L5XXeMNfxf8Wl8STcVZTxscSFEKSxa0=",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateBase64HMAC(tt.args.key, tt.args.message)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error: %v", err)
			}
			assert.Equal(t, tt.want, got, "GenerateBase64HMAC() returned unexpected result")
		})
	}
}
