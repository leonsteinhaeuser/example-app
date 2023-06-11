package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/internal/utils"
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

// getThread gets a thread by ID.
//
//	@Summary      Get thread by ID
//	@Description  Get thread by ID
//	@Tags         thread
//	@Accept       json
//	@Produce      json
//	@Param        id    query     string  false  "find thread b ID"  Format(uuid)
//	@Success      200  {object}   thread.Thread
//	@Failure      400  {object}  utils.HTTPError
//	@Failure      404  {object}  utils.HTTPError
//	@Failure      500  {object}  utils.HTTPError
//	@Router       /thread/{id} [get]
func (t *ThreadRouter) getThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	t.log.Info().Field("id", id).Log("get thread by ID")
	th, err := t.ts.GetByID(ctx, id)
	if err != nil {
		t.log.Error(err).Field("id", id).Log("failed to get thread")
		utils.NewHTTPError(http.StatusInternalServerError, "failed find thread by id in database", err).Write(w)
		return
	}
	err = utils.WriteJSON(w, http.StatusOK, th)
	if err != nil {
		t.log.Error(err).Field("id", id).Log("failed to encode thread")
		utils.NewHTTPError(http.StatusInternalServerError, "failed to encode thread", err).Write(w)
		return
	}
}

// listThreads lists all threads.
//
//	@Summary      List all threads
//	@Description  List all threads
//	@Tags         thread
//	@Accept       json
//	@Produce      json
//	@Success      200  {array}   thread.Thread
//	@Failure      400  {object}  utils.HTTPError
//	@Failure      404  {object}  utils.HTTPError
//	@Failure      500  {object}  utils.HTTPError
//	@Router       /thread [get]
func (t *ThreadRouter) listThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	threads, err := t.ts.List(ctx)
	if err != nil {
		t.log.Error(err).Log("failed to list threads")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSON(w, http.StatusOK, threads)
	if err != nil {
		t.log.Error(err).Log("failed to encode thread")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// createThread creates a thread.
//
//	@Summary      Create a thread
//	@Description  Create a thread
//	@Tags         thread
//	@Accept       json
//	@Produce      json
//	@Success      200  {object}  thread.Thread
//	@Failure      400  {object}  utils.HTTPError
//	@Failure      404  {object}  utils.HTTPError
//	@Failure      500  {object}  utils.HTTPError
//	@Router       /thread [post]
func (t *ThreadRouter) createThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	th := &thread.Thread{}
	err := utils.ReadJSON(r, th)
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
	err = utils.WriteJSON(w, http.StatusCreated, th)
	if err != nil {
		t.log.Error(err).Log("failed to encode thread")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// updateThread updates a thread.
//
//	@Summary      Update a thread
//	@Description  Update a thread
//	@Tags         thread
//	@Accept       json
//	@Produce      json
//	@Success      200  {object}  thread.Thread
//	@Failure      400  {object}  utils.HTTPError
//	@Failure      404  {object}  utils.HTTPError
//	@Failure      500  {object}  utils.HTTPError
//	@Router       /thread/{id} [put]
func (t *ThreadRouter) updateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tid := chi.URLParam(r, "id")

	th := &thread.Thread{}
	err := utils.ReadJSON(r, th)
	if err != nil {
		t.log.Error(err).Field("id", tid).Log("failed to decode thread")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if th.ID.String() != tid {
		t.log.Error(err).Field("id", tid).Log("thread id mismatch")
		http.Error(w, "thread id mismatch", http.StatusBadRequest)
		return
	}

	err = t.ts.UpdateByID(ctx, th)
	if err != nil {
		t.log.Error(err).Field("id", tid).Log("failed to update thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, th)
	if err != nil {
		t.log.Error(err).Field("id", tid).Log("failed to encode thread")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// deleteThread deletes a thread.
//
//	@Summary      Deleteh a thread
//	@Description  Deleteh a thread
//	@Tags         thread
//	@Accept       json
//	@Produce      json
//	@Success      200  {object}  thread.Thread
//	@Failure      400  {object}  utils.HTTPError
//	@Failure      404  {object}  utils.HTTPError
//	@Failure      500  {object}  utils.HTTPError
//	@Router       /thread/{id} [delete]
func (t *ThreadRouter) deleteThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	thread := &thread.Thread{}
	err := utils.ReadJSON(r, thread)
	if err != nil {
		t.log.Error(err).Field("id", id).Log("failed to decode thread")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if thread.ID.String() != id {
		t.log.Error(err).Field("id", id).Log("thread id mismatch")
		http.Error(w, "thread id mismatch", http.StatusBadRequest)
		return
	}

	err = t.ts.DeleteByID(ctx, id)
	if err != nil {
		t.log.Error(err).Field("id", id).Log("failed to delete thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
