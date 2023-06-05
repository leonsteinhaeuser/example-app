package keystore

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNewRedisKeyStore(t *testing.T) {
	type args struct {
		cfg RedisConfig
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "redis",
			args: args{
				cfg: RedisConfig{
					Driver:     "redis",
					ClientName: "test",
					Address:    "localhost:6379",
					Password:   "",
					Username:   "",
					DB:         0,
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "redis-cluster",
			args: args{
				cfg: RedisConfig{
					Driver:     "cluster",
					ClientName: "test",
					Addresses:  []string{"localhost:6379"},
					Password:   "",
					Username:   "",
					DB:         0,
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sentinel",
			args: args{
				cfg: RedisConfig{
					Driver:     "sentinel",
					ClientName: "test",
					Addresses:  []string{"localhost:6379"},
					Password:   "",
					Username:   "",
					DB:         0,
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sentinel-cluster",
			args: args{
				cfg: RedisConfig{
					Driver:     "sentinel-cluster",
					ClientName: "test",
					Addresses:  []string{"localhost:6379"},
					Password:   "",
					Username:   "",
					DB:         0,
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "unsupported driver",
			args: args{
				cfg: RedisConfig{
					Driver:     "oops",
					ClientName: "test",
					Addresses:  []string{"localhost:6379"},
					Password:   "",
					Username:   "",
					DB:         0,
				},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisKeyStore(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisKeyStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil == tt.wantNil {
				t.Errorf("NewRedisKeyStore() = %v, wantNil %v", got, tt.wantNil)
				return
			}
		})
	}
}

func Test_redisKeyStore_Set(t *testing.T) {
	type fields struct {
		setFunc   func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		getFunc   func(ctx context.Context, key string) *redis.StringCmd
		closeFunc func() error
	}
	type args struct {
		ctx        context.Context
		key        string
		value      any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "with empty value",
			fields: fields{
				setFunc: func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
					return &redis.StatusCmd{}
				},
			},
			args: args{
				ctx:   context.Background(),
				key:   "test",
				value: []byte("test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rks := redisKeyStore{
				setFunc:   tt.fields.setFunc,
				getFunc:   tt.fields.getFunc,
				closeFunc: tt.fields.closeFunc,
			}
			if err := rks.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("redisKeyStore.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisKeyStore_Get(t *testing.T) {
	type fields struct {
		setFunc   func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		getFunc   func(ctx context.Context, key string) *redis.StringCmd
		closeFunc func() error
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "with value",
			fields: fields{
				getFunc: func(ctx context.Context, key string) *redis.StringCmd {
					return redis.NewStringResult("test", nil)
				},
			},
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			want: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rks := redisKeyStore{
				setFunc:   tt.fields.setFunc,
				getFunc:   tt.fields.getFunc,
				closeFunc: tt.fields.closeFunc,
			}
			got, err := rks.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisKeyStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisKeyStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisKeyStore_Delete(t *testing.T) {
	type fields struct {
		deleteFunc func(ctx context.Context, key ...string) *redis.IntCmd
		closeFunc  func() error
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "with value",
			fields: fields{
				deleteFunc: func(ctx context.Context, key ...string) *redis.IntCmd {
					return redis.NewIntResult(1, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rks := redisKeyStore{
				deleteFunc: tt.fields.deleteFunc,
				closeFunc:  tt.fields.closeFunc,
			}
			err := rks.Delete(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisKeyStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_redisKeyStore_Close(t *testing.T) {
	type fields struct {
		setFunc   func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		getFunc   func(ctx context.Context, key string) *redis.StringCmd
		closeFunc func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "with value",
			fields: fields{
				closeFunc: func() error {
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rks := redisKeyStore{
				setFunc:   tt.fields.setFunc,
				getFunc:   tt.fields.getFunc,
				closeFunc: tt.fields.closeFunc,
			}
			if err := rks.Close(); (err != nil) != tt.wantErr {
				t.Errorf("redisKeyStore.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
