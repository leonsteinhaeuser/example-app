package db

import (
	"context"
	"fmt"
	"time"

	"github.com/leonsteinhaeuser/example-app/internal/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PosgresConfigFromEnv returns a PostgresConfig struct with values from the environment.
// If no environment variables are set, it will return a PostgresConfig with default values.
//
// The following environment variables are used:
// - POSTGRES_HOST (default: localhost)
// - POSTGRES_PORT (default: 5432)
// - POSTGRES_USERNAME (default: "")
// - POSTGRES_PASSWORD (default: "")
// - POSTGRES_DATABASE (default: "")
// - POSTGRES_OPTIONS (default: sslmode=disable)
// - POSTGRES_MAX_OPEN_CONNS (default: 10)
// - POSTGRES_MAX_IDLE_CONNS (default: 10)
// - POSTGRES_MAX_IDLE_CONN_TIME_SEC (default: 10)
// - POSTGRES_MAX_CONN_LIFETIME_SEC (default: 10)
func PosgresConfigFromEnv() PostgresConfig {
	return PostgresConfig{
		Host: env.GetStringEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port: env.GetStringEnvOrDefault("POSTGRES_PORT", "5432"),
		// Username, Password and Database are required, so we don't provide a default value.
		Username:     env.GetStringEnvOrDefault("POSTGRES_USERNAME", ""),
		Password:     env.GetStringEnvOrDefault("POSTGRES_PASSWORD", ""),
		Database:     env.GetStringEnvOrDefault("POSTGRES_DATABASE", ""),
		Options:      env.GetStringEnvOrDefault("POSTGRES_OPTIONS", "sslmode=disable"),
		MaxOpenConns: env.GetIntEnvOrDefault("POSTGRES_MAX_OPEN_CONNS", 10),
		MaxIdleConns: env.GetIntEnvOrDefault("POSTGRES_MAX_IDLE_CONNS", 10),
		// MaxIdleConnTime and MaxConnLifetime are time.Duration values, so we need to parse them.
		MaxIdleConnTime: time.Duration(env.GetIntEnvOrDefault("POSTGRES_MAX_IDLE_CONN_TIME_SEC", 10)) * time.Second,
		MaxConnLifetime: time.Duration(env.GetIntEnvOrDefault("POSTGRES_MAX_CONN_LIFETIME_SEC", 10)) * time.Second,
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Options  string

	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleConnTime time.Duration
	MaxConnLifetime time.Duration
}

func (p *PostgresConfig) dsn() string {
	dsn := fmt.Sprintf("postgres://%s@%s:%s/%s?password=%s",
		p.Username,
		p.Host,
		p.Port,
		p.Database,
		p.Password,
	)
	if p.Options != "" {
		dsn += "&" + p.Options
	}
	return dsn
}

// NewGormRepository opens a connection to the database using the GORM library.
func NewGormRepository(conf PostgresConfig) (Repository, error) {
	gormDB, err := gorm.Open(postgres.Open(conf.dsn()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxIdleTime(conf.MaxIdleConnTime)
	db.SetConnMaxLifetime(conf.MaxConnLifetime)
	return &gormRepository{
		DB: gormDB,
	}, nil
}

// gormRepository represents the interface between the application and the database.
// It uses the GORM library to interact with the database.
// gormRepository implements the Repository interface.
type gormRepository struct {
	DB *gorm.DB
}

func (p *gormRepository) Create(ctx context.Context, data any) error {
	return p.DB.Model(data).Create(data).Error
}

func (p *gormRepository) Find(ctx context.Context, data any, selectors ...Options) error {
	tx := p.DB.Model(data)
	limitCount := 0
	for _, selector := range selectors {
		fmt.Println("selector", selector)
		// if limit is set, use it only once (first selector)
		if selector.Limit != nil && limitCount == 0 {
			tx = tx.Limit(*selector.Limit)
			limitCount++
			continue
		}
		tx = tx.Where(selector.Field+" = ?", selector.Value)
	}
	return tx.Find(data).Error
}

func (p *gormRepository) Update(ctx context.Context, data any, selectors ...Options) error {
	tx := p.DB.Model(data)
	for _, selector := range selectors {
		tx = tx.Where(selector.Field, selector.Value)
	}
	return tx.Updates(data).Error
}

func (p *gormRepository) Delete(ctx context.Context, data any, selectors ...Options) error {
	tx := p.DB.Model(data)
	for _, selector := range selectors {
		tx = tx.Where(selector.Field, selector.Value)
	}
	return tx.Delete(data).Error
}

func (p *gormRepository) Migrate(ctx context.Context, model any) error {
	return p.DB.AutoMigrate(model)
}

func (p *gormRepository) Close(context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *gormRepository) Raw(ctx context.Context, query string, args ...any) error {
	return p.DB.Raw(query, args).Error
}
