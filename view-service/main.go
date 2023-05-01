package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonsteinhaeuser/example-app/article-service/client"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	numberClient "github.com/leonsteinhaeuser/example-app/number-service/client"
	userClient "github.com/leonsteinhaeuser/example-app/user-service/client"
	"github.com/leonsteinhaeuser/example-app/view-service/views"
)

var (
	clog log.Logger = log.NewZerlog()

	articleServiceURL    = os.Getenv("ARTICLE_SERVICE_URL")
	numberServiceAddress = os.Getenv("NUMBER_SERVICE_ADDRESS")
	userServiceURL       = os.Getenv("USER_SERVICE_URL")

	pl lib.ProcessLifecycle = lib.NewProcessLifecycle([]os.Signal{os.Interrupt, os.Kill})
)

func main() {
	ctx, cf := context.WithTimeout(context.Background(), 30*time.Second)
	defer cf()

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.CleanPath)
	mux.Use(log.LoggerMiddleware(clog))
	mux.Use(middleware.AllowContentType("application/json"))
	mux.Use(middleware.Recoverer)
	mux.Get("/healthz", lib.Healthz)

	views.NewArticleRouter(clog, client.NewArticleClient(clog, articleServiceURL)).Router(mux)
	views.NewNumberRouter(clog, numberClient.NewNumberClient(numberServiceAddress)).Router(mux)
	views.NewUserRouter(clog, userClient.NewUserClient(clog, userServiceURL)).Router(mux)

	routes := &[]string{}
	views.NewHomeRouter(clog, routes).Router(mux)
	chi.Walk(mux, func(_, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		*routes = append(*routes, route)
		return nil
	})
	lib.WalkRoutes(mux, clog)

	go func() {
		err := http.ListenAndServe(":2222", mux)
		if err != nil {
			clog.Panic(err).Log("something went wrong with the server")
		}
	}()

	pl.Wait()
	pl.Shutdown(ctx)
}
