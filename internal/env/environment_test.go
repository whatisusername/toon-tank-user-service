package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name     string
		args     args
		setupEnv func(t *testing.T, key string)
		want     string
		wantErr  bool
	}{
		{
			name: "Value Exists",
			args: args{key: "EXISTING_KEY"},
			setupEnv: func(t *testing.T, key string) {
				t.Setenv(key, "value")
			},
			want:    "value",
			wantErr: false,
		},
		{
			name:     "Value Not Exists",
			args:     args{key: "NON_EXISTING_KEY"},
			setupEnv: func(t *testing.T, key string) {},
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t, tt.args.key)

			got, err := GetValue(tt.args.key)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error: %v", err)
			}
			assert.Equal(t, tt.want, got, "GetValue() returned unexpected result")
		})
	}
}

func TestGetValueOrDefault(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name     string
		args     args
		setupEnv func(t *testing.T, key string)
		want     string
	}{
		{
			name: "Value Exists",
			args: args{
				key:          "EXISTING_KEY",
				defaultValue: "default",
			},
			setupEnv: func(t *testing.T, key string) {
				t.Setenv(key, "value")
			},
			want: "value",
		},
		{
			name: "Default Value",
			args: args{
				key:          "NON_EXISTING_KEY",
				defaultValue: "default",
			},
			setupEnv: func(t *testing.T, key string) {},
			want:     "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t, tt.args.key)

			got := GetValueOrDefault(tt.args.key, tt.args.defaultValue)
			assert.Equal(t, tt.want, got, "GetValueOrDefault() returned unexpected result")
		})
	}
}
