package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	dbDriver   = os.Getenv("DATABASE_DRIVER")
	dbHost     = os.Getenv("DATABASE_HOST")
	dbPort     = os.Getenv("DATABASE_PORT")
	dbUsername = os.Getenv("DATABASE_USERNAME")
	dbPassword = os.Getenv("DATABASE_PASSWORD")
	dbName     = os.Getenv("DATABASE_NAME")
	dbOptions  = os.Getenv("DATABASE_OPTIONS")

	gormDB *gorm.DB
)

type Article struct {
	// identifier and state fields
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// model fields

	// Title reprGormDBDataTypeesents the title of the article
	Title string `json:"title,omitempty"`
	// Description represents a short description what the article is about
	Description string `json:"description,omitempty"`
	// Author represents the author of the article
	AuthorID uuid.UUID `json:"author_id,omitempty" gorm:"type:uuid"`
	// Text is the actual article text
	Text string `json:"text,omitempty"`

	KeywordIDs pq.StringArray `json:"keyword_ids,omitempty" gorm:"type:uuid[]"`
}

func init() {
	log.Info().Msg("initializing database")

	var dialector gorm.Dialector
	switch dbDriver {
	case "postgres", "postgresql":
		dsn := fmt.Sprintf("postgres://%s@%s:%s/%s?password=%s", dbUsername, dbHost, dbPort, dbName, dbPassword)
		if dbOptions != "" {
			dsn += "&" + dbOptions
		}
		dialector = postgres.Open(dsn)
	default:
		log.Fatal().Msgf("unsupported database driver: %q", dbDriver)
	}

	log.Info().Msg("connecting to database...")
	db, err := gorm.Open(dialector)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
	gormDB = db

	err = gormDB.Raw(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		log.Fatal().Err(err).Msg("failed to enable uuid-ossp extension")
	}

	log.Info().Msg("migrating database tables...")
	err = gormDB.AutoMigrate(Article{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database table")
	}
}

func main() {
	// on shutdown close database
	defer func() {
		db, err := gormDB.DB()
		if err != nil {
			log.Error().Err(err).Msg("unable to close database")
			return
		}
		db.Close()
	}()

	log.Info().Msg("creating and initializing http router")
	mux := chi.NewRouter()

	// initialize middlewares
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.CleanPath)
	mux.Use(middleware.Logger)
	mux.Use(middleware.AllowContentType("application/json"))
	mux.Use(middleware.Recoverer)

	log.Info().Msg("defining http routes")
	mux.Route("/article", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getByID)
			r.Put("/", updateByID)
			r.Delete("/", deleteByID)
		})
		r.Get("/", list)
		r.Post("/", create)
	})

	log.Info().Msg("starting article-service with address: 0.0.0.0:3333")
	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal().Err(err).Msg("something went wrong with the server")
	}
}

func getURLID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func getByID(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	id := getURLID(r)

	log.Debug().Msg("requesting article by ID from database")
	article := &Article{}
	err := gormDB.Model(article).First(article, "id = ?", id).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to article by id")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Debug().Msg("parsing article to json")
	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		log.Error().Err(err).Msgf("failed to encode article with id %q", article.ID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	articles := &[]Article{}

	log.Debug().Msg("querying database for all articles")
	err := gormDB.Model(articles).Find(articles).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to get all articles from database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Debug().Msg("parsing list of articles to json")
	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		log.Error().Err(err).Msg("failed to encode list of articles")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateByID(w http.ResponseWriter, r *http.Request) {
	id := getURLID(r)

	log.Debug().Msg("parsing request from json to article")
	article := Article{}
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if article.ID.String() != id {
		http.Error(w, "request and URL id mismatch", http.StatusBadRequest)
		return
	}

	log.Debug().Msg("updating article in database")
	err = gormDB.Model(article).Save(article).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to update article in database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func deleteByID(w http.ResponseWriter, r *http.Request) {
	id := getURLID(r)

	log.Debug().Msg("unparsing request body to article")
	article := &Article{}
	err := json.NewDecoder(r.Body).Decode(article)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if article.ID.String() != id {
		http.Error(w, "request body and URL id mismatch", http.StatusBadRequest)
		return
	}

	log.Debug().Msg("deleting article from database")
	err = gormDB.Model(article).Delete(article).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to delete article form database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("unparsing request body to article")
	article := &Article{}
	err := json.NewDecoder(r.Body).Decode(article)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug().Msg("creating article in database")
	err = gormDB.Model(article).Create(article).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to create article in database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	log.Debug().Msg("encoding created article to json")
	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		log.Error().Err(err).Msg("failed to encode article")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
