package db

import (
	"testing"
)

func TestPostgres_dsn(t *testing.T) {
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
