package env

import (
	"os"
	"testing"
)

func TestGetStringEnvOrDefault(t *testing.T) {
	type args struct {
		key string
		def string
	}
	type testEnv struct {
		val string
	}
	tests := []struct {
		name string
		args args
		env  *testEnv
		want string
	}{
		{
			name: "default",
			args: args{
				key: "TEST",
				def: "default",
			},
			env:  nil,
			want: "default",
		},
		{
			name: "default",
			args: args{
				key: "TEST",
				def: "default",
			},
			env:  &testEnv{val: "test"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != nil {
				os.Setenv(tt.args.key, tt.env.val)
				defer os.Unsetenv(tt.args.key)
			}

			if got := GetStringEnvOrDefault(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetStringEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntEnvOrDefault(t *testing.T) {
	type args struct {
		key string
		def int
	}
	type testEnv struct {
		val string
	}
	tests := []struct {
		name string
		args args
		env  *testEnv
		want int
	}{
		{
			name: "env not set",
			args: args{
				key: "TEST",
				def: 10,
			},
			env:  nil,
			want: 10,
		},
		{
			name: "env set and valid",
			args: args{
				key: "TEST",
				def: 0,
			},
			env:  &testEnv{val: "1"},
			want: 1,
		},
		{
			name: "env set and not valid",
			args: args{
				key: "TEST",
				def: 10,
			},
			env:  &testEnv{val: "asd"},
			want: 10,
		},
	}
	for _, tt := range tests {
		if tt.env != nil {
			os.Setenv(tt.args.key, tt.env.val)
			defer os.Unsetenv(tt.args.key)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIntEnvOrDefault(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetIntEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
