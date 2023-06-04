package env

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadDynamicEnv(t *testing.T) {
	type args struct {
		prefix string
	}
	type setup struct {
		envs map[string]string
	}
	tests := []struct {
		name  string
		args  args
		setup setup
		want  map[string]string
	}{
		{
			name: "no envs",
			args: args{
				prefix: "XXYZ_",
			},
			setup: setup{},
			want:  map[string]string{},
		},
		{
			name: "1 envs",
			args: args{
				prefix: "XXYZ_",
			},
			setup: setup{
				envs: map[string]string{
					"XXYZ_TEST": "test",
				},
			},
			want: map[string]string{
				"XXYZ_TEST": "test",
			},
		},
		{
			name: "5 envs",
			args: args{
				prefix: "XXYZ_",
			},
			setup: setup{
				envs: map[string]string{
					"XXYZ_CLIENT_ID":           "test",
					"XXYZ_CLIENT_SECRET":       "test",
					"XXYZ_CLIENT_REDIRECT_URL": "test",
					"XXYZ_CLIENT_SCOPES":       "test",
					"XXYZ_CLIENT_ENDPOINT":     "test",
				},
			},
			want: map[string]string{
				"XXYZ_CLIENT_ID":           "test",
				"XXYZ_CLIENT_SECRET":       "test",
				"XXYZ_CLIENT_REDIRECT_URL": "test",
				"XXYZ_CLIENT_SCOPES":       "test",
				"XXYZ_CLIENT_ENDPOINT":     "test",
			},
		},
		{
			name: "5 envs but wrong prefix",
			args: args{
				prefix: "XXYZ_",
			},
			setup: setup{
				envs: map[string]string{
					"XXYZA_CLIENT_ID":           "test",
					"XXYZA_CLIENT_SECRET":       "test",
					"XXYZA_CLIENT_REDIRECT_URL": "test",
					"XXYZA_CLIENT_SCOPES":       "test",
					"XXYZA_CLIENT_ENDPOINT":     "test",
				},
			},
			want: map[string]string{},
		},
		{
			name: "env with > 2 =",
			args: args{
				prefix: "XXYZ_",
			},
			setup: setup{
				envs: map[string]string{
					"XXYZ_CLIENT_ID": "a=test",
				},
			},
			want: map[string]string{
				"XXYZ_CLIENT_ID": "a=test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup environment
			if tt.setup.envs != nil {
				for k, v := range tt.setup.envs {
					os.Setenv(k, v)
				}
			}
			// defer cleanup environment
			defer func() {
				for k := range tt.setup.envs {
					os.Unsetenv(k)
				}
			}()

			if got := LoadDynamicEnv(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadDynamicEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
