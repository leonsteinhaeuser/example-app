package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/thread-service/thread"
)

var (
	_ server.Router = (*ThreadRouter)(nil)
)

type ThreadRouter struct {
	log log.Logger
	ts  thread.Store
}

func NewThreadRouter(log log.Logger, ts thread.Store) *ThreadRouter {
	return &ThreadRouter{
		log: log,
		ts:  ts,
	}
}

func (t *ThreadRouter) Router(rt chi.Router) {
	rt.Route("/thread", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", t.getThread)
			r.Put("/", t.updateThread)
			r.Delete("/", t.deleteThread)
		})
		r.Get("/", t.listThreads)
		r.Post("/", t.createThread)
	})
}

func (t *ThreadRouter) getThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	th, err := t.ts.GetByID(ctx, id)
	if err != nil {
		t.log.Error(err).Log("failed to get thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(th)
	if err != nil {
		t.log.Error(err).Log("failed to encode thread")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *ThreadRouter) listThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	threads, err := t.ts.List(ctx)
	if err != nil {
		t.log.Error(err).Log("failed to list threads")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(threads)
	if err != nil {
		t.log.Error(err).Log("failed to encode threads")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *ThreadRouter) createThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	th := &thread.Thread{}
	err := json.NewDecoder(r.Body).Decode(th)
	if err != nil {
		t.log.Error(err).Log("failed to decode thread")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.ts.Create(ctx, th)
	if err != nil {
		t.log.Error(err).Log("failed to create thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (t *ThreadRouter) updateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tid := chi.URLParam(r, "id")

	th := &thread.Thread{}
	err := json.NewDecoder(r.Body).Decode(th)
	if err != nil {
		t.log.Error(err).Log("failed to decode thread")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if th.ID.String() != tid {
		t.log.Error(err).Log("thread id mismatch")
		http.Error(w, "thread id mismatch", http.StatusBadRequest)
		return
	}

	err = t.ts.UpdateByID(ctx, th)
	if err != nil {
		t.log.Error(err).Log("failed to update thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *ThreadRouter) deleteThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	thread := &thread.Thread{}
	err := json.NewDecoder(r.Body).Decode(thread)
	if err != nil {
		t.log.Error(err).Log("failed to decode thread")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if thread.ID.String() != id {
		t.log.Error(err).Log("thread id mismatch")
		http.Error(w, "thread id mismatch", http.StatusBadRequest)
		return
	}

	err = t.ts.DeleteByID(ctx, id)
	if err != nil {
		t.log.Error(err).Log("failed to delete thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
