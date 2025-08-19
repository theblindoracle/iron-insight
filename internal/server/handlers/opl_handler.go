package handlers

import (
	"context"
	"encoding/json"
	"iron-insight/internal/database"
	"log/slog"
	"net/http"
)

type OplHandler struct {
	queries *database.Queries
}

func NewOplHandler(q *database.Queries) *OplHandler {
	return &OplHandler{
		queries: q,
	}
}

func (h *OplHandler) GetLifterRecords(w http.ResponseWriter, r *http.Request) {

	data, err := h.queries.GetMeetDataForLifterName(context.Background(), "Travis Napier")
	if err != nil {
		slog.Error("Could not get lifter data", "err", err)
		http.Error(w, "Failed to get lifter data", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"lifters": data}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Could not marshall json", "data", data)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
}
