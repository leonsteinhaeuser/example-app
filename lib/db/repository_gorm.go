package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Driver          string
	Postgres        PostgresConfig
	Mysql           MysqlConfig
	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleConnTime time.Duration
	MaxConnLifetime time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Options  string
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

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (m *MysqlConfig) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
	)
}

// NewGormRepository opens a connection to the database using the GORM library.
func NewGormRepository(conf Config) (Repository, error) {
	var dial gorm.Dialector
	switch conf.Driver {
	case "postgres":
		dial = postgres.Open(conf.Postgres.dsn())
	case "mysql":
		dial = mysql.Open(conf.Mysql.dsn())
	default:
		panic(fmt.Sprintf("unsupported driver: %s", conf.Driver))
	}
	gormDB, err := gorm.Open(dial, &gorm.Config{})
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

func (p *gormRepository) Find(ctx context.Context, data any, selectors ...Selector) error {
	tx := p.DB.Model(data)
	for _, selector := range selectors {
		fmt.Println("selector", selector)
		tx = tx.Where(selector.Field+" = ?", selector.Value)
	}
	return tx.Find(data).Error
}

func (p *gormRepository) Update(ctx context.Context, data any, selectors ...Selector) error {
	tx := p.DB.Model(data)
	for _, selector := range selectors {
		tx = tx.Where(selector.Field, selector.Value)
	}
	return tx.Updates(data).Error
}

func (p *gormRepository) Delete(ctx context.Context, data any, selectors ...Selector) error {
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
