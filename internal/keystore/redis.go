package keystore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUnsupportedRedisDriver = errors.New("unsupported redis driver")

	_ KeyStore = (*redisKeyStore)(nil)
)

type RedisDriver string

const (
	RedisDriverRedis           RedisDriver = "redis"
	RedisDriverCluster         RedisDriver = "cluster"
	RedisDriverSentinel        RedisDriver = "sentinel"
	RedisDriverSentinelCluster RedisDriver = "sentinel-cluster"
)

type RedisConfig struct {
	// Driver is the name of the redis driver to use
	//
	// Valid values are:
	// - "redis"
	// - "cluster"
	// - "sentinel"
	// - "sentinel-cluster"
	Driver string
	// ClientName is the name of the redis client
	ClientName string

	// Address is the host:port address of the redis server
	Address string
	// Addresses is a list of host:port addresses of redis servers
	Addresses []string
	// Password is the password to use when connecting to the redis server
	Password string
	// Username is the username to use when connecting to the redis server
	Username string
	// DB is the redis database to use
	DB int
}

type redisKeyStore struct {
	setFunc func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	getFunc func(ctx context.Context, key string) *redis.StringCmd

	closeFunc func() error
}

func NewRedisKeyStore(cfg RedisConfig) (KeyStore, error) {
	rks := &redisKeyStore{}

	switch RedisDriver(cfg.Driver) {
	case RedisDriverRedis:
		rc := redis.NewClient(&redis.Options{
			Addr:       cfg.Address,
			ClientName: cfg.ClientName,
			Password:   cfg.Password,
			Username:   cfg.Username,
			DB:         cfg.DB,
		})
		rks.setFunc = rc.Set
		rks.getFunc = rc.Get
		rks.closeFunc = rc.Close
	case RedisDriverCluster:
		rcc := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:      cfg.Addresses,
			ClientName: cfg.ClientName,
			Password:   cfg.Password,
			Username:   cfg.Username,
			NewClient: func(opt *redis.Options) *redis.Client {
				opt.DB = cfg.DB
				return redis.NewClient(opt)
			},
		})
		rks.setFunc = rcc.Set
		rks.getFunc = rcc.Get
		rks.closeFunc = rcc.Close
	case RedisDriverSentinel:
		rfc := redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.Address,
			SentinelAddrs: cfg.Addresses,
			ClientName:    cfg.ClientName,
			Password:      cfg.Password,
			Username:      cfg.Username,
			DB:            cfg.DB,
		})
		rks.setFunc = rfc.Set
		rks.getFunc = rfc.Get
		rks.closeFunc = rfc.Close
	case RedisDriverSentinelCluster:
		rfcc := redis.NewFailoverClusterClient(&redis.FailoverOptions{
			MasterName:    cfg.Address,
			SentinelAddrs: cfg.Addresses,
			ClientName:    cfg.ClientName,
			Password:      cfg.Password,
			Username:      cfg.Username,
			DB:            cfg.DB,
		})
		rks.setFunc = rfcc.Set
		rks.getFunc = rfcc.Get
		rks.closeFunc = rfcc.Close
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedRedisDriver, cfg.Driver)
	}
	return rks, nil
}

func (rks redisKeyStore) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return rks.setFunc(ctx, key, value, expiration).Err()
}

func (rks redisKeyStore) Get(ctx context.Context, key string) ([]byte, error) {
	return rks.getFunc(ctx, key).Bytes()
}

func (rks redisKeyStore) Close() error {
	return rks.closeFunc()
}
