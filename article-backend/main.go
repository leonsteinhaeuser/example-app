package main

import (
	"context"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/article-backend/api/v1/article"
	"github.com/leonsteinhaeuser/example-app/internal/db"
	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
)

var (
	logr = log.NewZerlog()

	httpServer = server.NewDefaultServer(logr, env.GetStringEnvOrDefault("LISTEN_ADDRESS", ":1200"))
	httpRouter = server.NewGenericRouter()

	dbr db.Repository
)

func init() {
	db, err := db.NewGormRepository(db.PostgresConfig{
		Host:     env.GetStringEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     env.GetStringEnvOrDefault("POSTGRES_PORT", "5432"),
		Username: env.GetStringEnvOrDefault("POSTGRES_USERNAME", "postgres"),
		Password: env.GetStringEnvOrDefault("POSTGRES_PASSWORD", "postgres"),
		Database: env.GetStringEnvOrDefault("POSTGRES_DATABASE", "articles"),
	})
	if err != nil {
		panic(err)
	}
	dbr = db

	dbr.Migrate(context.Background(), &article.Article{})
}

func main() {
	defer dbr.Close(context.Background())

	httpRouter.AddEndpoint("GET", "/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	httpServer.AddRouter(httpRouter)
	httpServer.AddRouter(article.NewArticleRouter(logr, dbr))
	err := httpServer.Start()
	if err != nil {
		panic(err)
	}
}
