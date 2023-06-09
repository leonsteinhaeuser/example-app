package main

import (
	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/thread-service/client"
	"github.com/leonsteinhaeuser/example-app/thread-view-service/api"
	"github.com/rs/zerolog"
)

var (
	threadServiceAddr = env.GetStringEnvOrDefault("THREAD_SERVICE_ADDRESS", "http://localhost:8080")
	listenAddr        = env.GetStringEnvOrDefault("LISTEN_ADDR", ":8080")

	// logger settings
	logLevel            = env.GetIntEnvOrDefault("LOG_LEVEL", 0)
	logr     log.Logger = log.NewZerlog(func(l *zerolog.Logger) {
		*l = l.With().Logger().Level(zerolog.Level(logLevel))
	})
)

func main() {
	srvr := server.NewDefaultServer(logr, listenAddr, "application/x-www-form-urlencoded")

	client, err := client.NewDefaultClient(threadServiceAddr)
	if err != nil {
		logr.Panic(err).Log("error creating thread service client")
	}

	srvr.AddRouter(api.NewViewRouter(logr, client))

	defer srvr.Stop()
	err = srvr.Start()
	if err != nil {
		logr.Panic(err).Log("error starting server")
	}
}
