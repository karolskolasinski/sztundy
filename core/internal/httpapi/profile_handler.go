package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"core/internal/domain"
)

type ProfileService interface {
	Create(ctx context.Context, input domain.CreateProfileInput) (domain.Profile, error)
	Get(ctx context.Context, id int64) (domain.Profile, error)
	Update(ctx context.Context, id int64, input domain.UpdateProfileInput) (domain.Profile, error)
}

type ProfileHandler struct {
	manager ProfileService
}

func NewProfileHandler(manager ProfileService) *ProfileHandler {
	return &ProfileHandler{manager: manager}
}

func (h *ProfileHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *ProfileHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateProfileInput
	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "Nieprawidłowy JSON w żądaniu.")
		return
	}

	created, err := h.manager.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	found, err := h.manager.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, found)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var input domain.UpdateProfileInput
	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "Nieprawidłowy JSON w żądaniu.")
		return
	}

	updated, err := h.manager.Update(r.Context(), id, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *ProfileHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, "Nie znaleziono profilu.")
	case errors.Is(err, domain.ErrEmailAlreadyTaken):
		writeError(w, http.StatusConflict, "Konto z takim adresem e-mail już istnieje.")
	default:
		writeError(w, http.StatusInternalServerError, "Wystąpił błąd serwera.")
	}
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Identyfikator profilu musi być dodatnią liczbą całkowitą.")
		return 0, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("żądanie musi zawierać dokładnie jeden obiekt JSON")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
