package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonsteinhaeuser/example-app/article-service/accessobjects"
	"github.com/leonsteinhaeuser/example-app/article-service/api"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
)

var (
	dbDriver   = os.Getenv("DATABASE_DRIVER")
	dbHost     = os.Getenv("DATABASE_HOST")
	dbPort     = os.Getenv("DATABASE_PORT")
	dbUsername = os.Getenv("DATABASE_USERNAME")
	dbPassword = os.Getenv("DATABASE_PASSWORD")
	dbName     = os.Getenv("DATABASE_NAME")
	dbOptions  = os.Getenv("DATABASE_OPTIONS")
	natsURL    = os.Getenv("NATS_URL")

	accessor db.Repository

	clog log.Logger = log.NewZerlog()

	pl         lib.ProcessLifecycle = lib.NewProcessLifecycle([]os.Signal{os.Interrupt, os.Kill})
	natsClient pubsub.Client
)

func init() {
	acsr, err := db.NewGormRepository(db.Config{
		Driver: dbDriver,

		Postgres: db.PostgresConfig{
			Host:     dbHost,
			Port:     dbPort,
			Password: dbPassword,
			Database: dbName,
			Username: dbUsername,
			Options:  dbOptions,
		},
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	})
	if err != nil {
		clog.Panic(err).Log("failed to initialize database accessor")
		return
	}
	accessor = acsr

	nc, err := pubsub.NewNatsClient(clog, natsURL, "general")
	if err != nil {
		clog.Panic(err).Log("failed to initialize nats client")
		return
	}
	natsClient = nc
}

func main() {
	ctx, cf := context.WithTimeout(context.Background(), 30*time.Second)
	defer cf()
	pl.RegisterShutdownProcess(accessor.Close, natsClient.Close)

	clog.Info().Log("creating and initializing http router")
	mux := chi.NewRouter()
	// initialize middlewares
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.CleanPath)
	mux.Use(log.LoggerMiddleware(clog))
	mux.Use(middleware.AllowContentType("application/json"))
	mux.Use(middleware.Recoverer)
	mux.Get("/healthz", lib.Healthz)

	artc := accessobjects.NewArticle(accessor, clog, natsClient)
	err := artc.Migrate(ctx)
	if err != nil {
		clog.Panic(err).Log("failed to migrate article table")
		return
	}

	clog.Info().Log("defining http routes")
	api.NewArticleRouter(*artc, clog).Router(mux)

	// log routes
	lib.WalkRoutes(mux, clog)

	clog.Info().Log("starting article-service with address: 0.0.0.0:3333")
	go func() {
		err = http.ListenAndServe(":3333", mux)
		if err != nil {
			clog.Panic(err).Log("something went wrong with the server")
		}
	}()

	pl.Wait()
	pl.Shutdown(ctx)
}
