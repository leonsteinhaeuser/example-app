package db

import (
	"context"
	"reflect"
	"testing"

	"gorm.io/gorm"
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

func TestMysql_dsn(t *testing.T) {
	type fields struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Mysql dsn",
			fields: fields{
				Host:     "localhost",
				Port:     3306,
				User:     "mysql",
				Password: "mysql",
				Database: "mysql",
			},
			want: "mysql:mysql@tcp(localhost:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MysqlConfig{
				Host:     tt.fields.Host,
				Port:     tt.fields.Port,
				User:     tt.fields.User,
				Password: tt.fields.Password,
				Database: tt.fields.Database,
			}
			if got := m.dsn(); got != tt.want {
				t.Errorf("Mysql.dsn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenGORM(t *testing.T) {
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGormRepository(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenGORM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenGORM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormRepository_Create(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx  context.Context
		data any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &gormRepository{
				DB: tt.fields.DB,
			}
			if err := p.Create(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GormRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormRepository_Find(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx       context.Context
		data      any
		selectors []Selector
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &gormRepository{
				DB: tt.fields.DB,
			}
			if err := p.Find(tt.args.ctx, tt.args.data, tt.args.selectors...); (err != nil) != tt.wantErr {
				t.Errorf("GormRepository.Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormRepository_Update(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx       context.Context
		data      any
		selectors []Selector
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &gormRepository{
				DB: tt.fields.DB,
			}
			if err := p.Update(tt.args.ctx, tt.args.data, tt.args.selectors...); (err != nil) != tt.wantErr {
				t.Errorf("GormRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormRepository_Delete(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx       context.Context
		data      any
		selectors []Selector
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &gormRepository{
				DB: tt.fields.DB,
			}
			if err := p.Delete(tt.args.ctx, tt.args.data, tt.args.selectors...); (err != nil) != tt.wantErr {
				t.Errorf("GormRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormRepository_Migrate(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx   context.Context
		model any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &gormRepository{
				DB: tt.fields.DB,
			}
			if err := p.Migrate(tt.args.ctx, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("GormRepository.Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormRepository_Close(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &gormRepository{
				DB: tt.fields.DB,
			}
			if err := p.Close(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("GormRepository.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
