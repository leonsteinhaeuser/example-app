package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/file-service/accessobjects"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type fileRouter struct {
	fileAccessObject accessobjects.File
	log              log.Logger
}

func NewFileRouter(aco accessobjects.File, log log.Logger) *fileRouter {
	return &fileRouter{
		fileAccessObject: aco,
		log:              log,
	}
}

func (a *fileRouter) Router(r chi.Router) {
	r.Route("/file", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", a.get)
			r.Put("/", a.put)
			r.Delete("/", a.delete)
		})
		r.Get("/", a.list)
		r.Post("/", a.post)
	})
}

func (a *fileRouter) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	file, err := a.fileAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get filefrom database", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, file)
	if err != nil {
		a.log.Error(err).Log("failed to write file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write file", err)
		return
	}
}

func (a *fileRouter) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	files, err := a.fileAccessObject.List(ctx)
	if err != nil {
		a.log.Error(err).Log("failed to list files")
		lib.WriteError(w, http.StatusInternalServerError, "failed to list files", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, files)
	if err != nil {
		a.log.Error(err).Log("failed to write files")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write files", err)
		return
	}
}

func (a *fileRouter) put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	file, err := a.fileAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get filefrom database", err)
		return
	}

	err = lib.ReadJSON(r, file)
	if err != nil {
		a.log.Error(err).Log("failed to read file")
		lib.WriteError(w, http.StatusBadRequest, "failed to read file", err)
		return
	}

	if file.ID.String() != id {
		a.log.Warn().
			Field("urlID", id).
			Field("body", file.ID.String()).
			Log("id in path and body do not match")
		lib.WriteError(w, http.StatusBadRequest, "id in path and body do not match", nil)
		return
	}

	err = a.fileAccessObject.Update(ctx, file)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to update file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to update file", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, file)
	if err != nil {
		a.log.Error(err).Log("failed to write file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write file", err)
		return
	}
}

func (a *fileRouter) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	uid, err := uuid.Parse(id)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to parse id")
		lib.WriteError(w, http.StatusBadRequest, "failed to parse id", err)
		return
	}

	file := &lib.File{ID: uid}
	err = a.fileAccessObject.Delete(ctx, file)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to delete file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to delete file", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *fileRouter) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file := &lib.File{}
	err := lib.ReadJSON(r, file)
	if err != nil {
		a.log.Error(err).Log("failed to read file")
		lib.WriteError(w, http.StatusBadRequest, "failed to read file", err)
		return
	}

	err = a.fileAccessObject.Create(ctx, file)
	if err != nil {
		a.log.Error(err).Log("failed to create file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to create file", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusCreated, file)
	if err != nil {
		a.log.Error(err).Log("failed to write file")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write file", err)
		return
	}
}
