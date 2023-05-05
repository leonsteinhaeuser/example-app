package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonsteinhaeuser/example-app/article-consumer-service/comment"
	"github.com/leonsteinhaeuser/example-app/articlecomment-service/client"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
	"github.com/leonsteinhaeuser/example-app/lib/worker"
)

var (
	natsURL                  = os.Getenv("NATS_URL")
	articleCommentServiceURL = os.Getenv("ARTICLE_COMMENT_SERVICE_URL")

	clog log.Logger = log.NewZerlog()

	pl         lib.ProcessLifecycle = lib.NewProcessLifecycle([]os.Signal{os.Interrupt, os.Kill})
	natsClient pubsub.Client

	articleCommentServiceClient lib.CustomArticleClient[lib.ArticleComment] = client.NewArticleCommentClient(clog, articleCommentServiceURL)
)

func init() {
	nc, err := pubsub.NewNatsClient(clog, natsURL, "article")
	if err != nil {
		clog.Panic(err).Log("failed to initialize nats client")
		return
	}
	natsClient = nc
}

func main() {
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	pl.RegisterShutdownProcess(natsClient.Close)

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

	msgCh := make(<-chan pubsub.SubscribeDater)
	go func() {
		mCh, err := natsClient.Subscribe()
		if err != nil {
			clog.Panic(err).Log("failed to subscribe to article-created")
			return
		}
		msgCh = mCh
	}()

	wCtx, wcf := context.WithCancel(context.Background())
	defer wcf()

	worker.Worker(wCtx, 10, func(ctx context.Context, workerID int) {
		clog.Debug().Field("worker_id", workerID).Log("starting process...")
		for {
			select {
			case <-ctx.Done():
				clog.Warn().Field("worker_id", workerID).Log("context was canceled, shutting down...")
			case msg := <-msgCh:
				clog.Warn().Field("worker_id", workerID).Log("processing message")
				err := comment.DeleteByArticleID(ctx, msg, articleCommentServiceClient)
				if err != nil {
					clog.Error(err).Field("worker_id", workerID).Field("message", msg).Log("failed to delete comments by article id")
					// we are not returning here because we want to continue processing messages
					continue
				}
			}
		}
	})

	// log routes
	lib.WalkRoutes(mux, clog)

	clog.Info().Log("starting article-consumer-service with address: 0.0.0.0:7777")
	go func() {
		err := http.ListenAndServe(":7777", mux)
		if err != nil {
			clog.Panic(err).Log("something went wrong with the server")
		}
	}()

	pl.Wait()
	pl.Shutdown(ctx)
}
