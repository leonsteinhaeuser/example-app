package db

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPosgresConfigFromEnv(t *testing.T) {
	type env struct {
		Host            string
		Port            string
		User            string
		Password        string
		Database        string
		Options         string
		MaxOpenConns    string
		MaxIdleConns    string
		MaxIdleConnTime string
		MaxConnLifetime string
	}
	tests := []struct {
		name string
		env  *env
		want PostgresConfig
	}{
		{
			name: "Test Postgres config from env",
			env: &env{
				Host:            "localhost",
				Port:            "5432",
				User:            "postgres",
				Password:        "postgres",
				Database:        "postgres",
				Options:         "sslmode=disable",
				MaxOpenConns:    "10",
				MaxIdleConns:    "10",
				MaxIdleConnTime: "10",
				MaxConnLifetime: "10",
			},
			want: PostgresConfig{
				Host:            "localhost",
				Port:            "5432",
				Username:        "postgres",
				Password:        "postgres",
				Database:        "postgres",
				Options:         "sslmode=disable",
				MaxOpenConns:    10,
				MaxIdleConns:    10,
				MaxIdleConnTime: 10 * time.Second,
				MaxConnLifetime: 10 * time.Second,
			},
		},
		{
			name: "env not set",
			env:  nil,
			want: PostgresConfig{
				Host:            "localhost",
				Port:            "5432",
				Username:        "",
				Password:        "",
				Database:        "",
				Options:         "sslmode=disable",
				MaxOpenConns:    10,
				MaxIdleConns:    10,
				MaxIdleConnTime: 10 * time.Second,
				MaxConnLifetime: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set postgres host env variable
			if tt.env != nil && tt.env.Host != "" {
				os.Setenv("POSTGRES_HOST", tt.env.Host)
				defer os.Unsetenv("POSTGRES_HOST")
			}
			// set postgres port env variable
			if tt.env != nil && tt.env.Port != "" {
				os.Setenv("POSTGRES_PORT", tt.env.Port)
				defer os.Unsetenv("POSTGRES_PORT")
			}
			// set postgres username env variable
			if tt.env != nil && tt.env.User != "" {
				os.Setenv("POSTGRES_USERNAME", tt.env.User)
				defer os.Unsetenv("POSTGRES_USERNAME")
			}
			// set postgres password env variable
			if tt.env != nil && tt.env.Password != "" {
				os.Setenv("POSTGRES_PASSWORD", tt.env.Password)
				defer os.Unsetenv("POSTGRES_PASSWORD")
			}
			// set postgres database env variable
			if tt.env != nil && tt.env.Database != "" {
				os.Setenv("POSTGRES_DATABASE", tt.env.Database)
				defer os.Unsetenv("POSTGRES_DATABASE")
			}
			// set postgres options env variable
			if tt.env != nil && tt.env.Options != "" {
				os.Setenv("POSTGRES_OPTIONS", tt.env.Options)
				defer os.Unsetenv("POSTGRES_OPTIONS")
			}
			// set postgres max open conns env variable
			if tt.env != nil && tt.env.MaxOpenConns != "" {
				os.Setenv("POSTGRES_MAX_OPEN_CONNS", tt.env.MaxOpenConns)
				defer os.Unsetenv("POSTGRES_MAX_OPEN_CONNS")
			}
			// set postgres max idle conns env variable
			if tt.env != nil && tt.env.MaxIdleConns != "" {
				os.Setenv("POSTGRES_MAX_IDLE_CONNS", tt.env.MaxIdleConns)
				defer os.Unsetenv("POSTGRES_MAX_IDLE_CONNS")
			}
			// set postgres max idle conn time env variable
			if tt.env != nil && tt.env.MaxIdleConnTime != "" {
				os.Setenv("POSTGRES_MAX_IDLE_CONN_TIME_SEC", tt.env.MaxIdleConnTime)
				defer os.Unsetenv("POSTGRES_MAX_IDLE_CONN_TIME_SEC")
			}
			// set postgres max conn lifetime env variable
			if tt.env != nil && tt.env.MaxConnLifetime != "" {
				os.Setenv("POSTGRES_MAX_CONN_LIFETIME_SEC", tt.env.MaxConnLifetime)
				defer os.Unsetenv("POSTGRES_MAX_CONN_LIFETIME_SEC")
			}

			if got := PosgresConfigFromEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PosgresConfigFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresConfig_dsn(t *testing.T) {
	type fields struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
		Options  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Postgres dsn",
			fields: fields{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				Database: "postgres",
				Options:  "sslmode=disable",
			},
			want: "postgres://postgres@localhost:5432/postgres?password=postgres&sslmode=disable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresConfig{
				Host:     tt.fields.Host,
				Port:     tt.fields.Port,
				Username: tt.fields.User,
				Password: tt.fields.Password,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
			}
			if got := p.dsn(); got != tt.want {
				t.Errorf("Postgres.dsn() = %v, want %v", got, tt.want)
			}
		})
	}
}
