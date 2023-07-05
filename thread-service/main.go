package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/leonsteinhaeuser/example-app/internal/db"
	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/pubsub"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/internal/utils"
	v1 "github.com/leonsteinhaeuser/example-app/thread-service/api/v1"
	"github.com/leonsteinhaeuser/example-app/thread-service/docs"
	"github.com/leonsteinhaeuser/example-app/thread-service/thread"
	"github.com/rs/zerolog"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	pubsubTopic = "thread"
)

var (
	// server settings
	listenAddr string = env.GetStringEnvOrDefault("LISTEN_ADDR", ":1112")
	// externalURL
	externalURL string = env.GetStringEnvOrDefault("EXTERNAL_URL", "http://localhost:1112")
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

	err = dbr.Migrate(ctx, &thread.ThreadModel{})
	if err != nil {
		logr.Panic(err).Log("failed to migrate database")
	}
}

// @title Thread Service API
// @version 1.0
// @description This is the thread-service.

// @contact.name API Support
// @contact.url http://github.com/leonsteinhaeuser/example-app
// @contact.email mail@example.localhost

// @license.name GNU GENERAL PUBLIC LICENSE 3.0
// @license.url https://github.com/leonsteinhaeuser/example-app/blob/main/LICENSE

// @host thread-service
// @BasePath /
func main() {
	srvr := server.NewDefaultServer(logr, listenAddr)

	srvr.AddRouter(v1.NewThreadRouter(logr, thread.Store{
		DB: dbr,
		PS: pubsubClient,
	}))

	// override swagger docs host
	docs.SwaggerInfo.Host = utils.ReplaceAllToReplace(externalURL, "", "http://", "https://")
	// define swagger endpoint to serve swagger docs
	dfltrt := server.NewGenericRouter()
	dfltrt.AddEndpoint(http.MethodGet, "/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", externalURL)),
	))
	srvr.AddRouter(dfltrt)

	defer srvr.Stop()
	logr.Info().Logf("starting server on %s", listenAddr)
	err := srvr.Start()
	if err != nil {
		logr.Panic(err).Log("failed to start server")
		return
	}
}
