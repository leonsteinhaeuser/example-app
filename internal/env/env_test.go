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
