package main

import (
	"context"
	"time"

	"github.com/leonsteinhaeuser/example-app/internal/db"
	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/pubsub"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	v1 "github.com/leonsteinhaeuser/example-app/thread-service/api/v1"
	"github.com/leonsteinhaeuser/example-app/thread-service/thread"
	"github.com/rs/zerolog"
)

const (
	pubsubTopic = "thread"
)

var (
	// server settings
	listenAddr string = env.GetStringEnvOrDefault("LISTEN_ADDR", ":8080")
	// logger settings
	logLevel            = env.GetIntEnvOrDefault("LOG_LEVEL", 0)
	logr     log.Logger = log.NewZerlog(func(l *zerolog.Logger) {
		*l = l.With().Logger().Level(zerolog.Level(logLevel))
	})
	// nats settings
	natsURL      = env.GetStringEnvOrDefault("NATS_CLIENTS", "nats://localhost:6222")
	pubsubClient pubsub.Publisher
	// database settings
	dbr db.Repository
)

func init() {
	// initialize pubsub client
	clnt, err := pubsub.NewNatsClient(natsURL, pubsubTopic)
	if err != nil {
		logr.Panic(err).Log("failed to initialize pubsub client")
	}
	pubsubClient = clnt
	// initialize database
	dbr, err = db.NewGormRepository(db.PosgresConfigFromEnv())
	if err != nil {
		logr.Panic(err).Log("failed to initialize database")
	}

	ctx, cf := context.WithTimeout(context.Background(), time.Second*5)
	defer cf()

	err = dbr.Migrate(ctx, &thread.Thread{})
	if err != nil {
		logr.Panic(err).Log("failed to migrate database")
	}
}

func main() {
	srvr := server.NewDefaultServer(logr, listenAddr)

	srvr.AddRouter(v1.NewThreadRouter(logr, thread.Store{
		DB: dbr,
		PS: pubsubClient,
	}))

	defer srvr.Stop()
	logr.Info().Logf("starting server on %s", listenAddr)
	err := srvr.Start()
	if err != nil {
		logr.Panic(err).Log("failed to start server")
		return
	}
}
