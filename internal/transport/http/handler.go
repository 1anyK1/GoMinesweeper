package http

import (
	"encoding/json"
	"net/http"

	"minesweeper/internal/service"
)

type Handler struct {
	gameService *service.GameService
}

func NewHandler(gameService *service.GameService) *Handler {
	return &Handler{
		gameService: gameService,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.Health)
	mux.HandleFunc("/state", h.State)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) State(w http.ResponseWriter, r *http.Request) {
	board, err := h.gameService.State()
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"board": board,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
