package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/avpetkun/jessy-go"
)

func JSON(w http.ResponseWriter, data any, statuscode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	err := jessy.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("JSON encode error:", slog.Any("msg", err.Error()))
		return
	}
}

func isInvalidCollectionName(name string) bool {
	return name == "" || strings.Contains(name, ".") || strings.Contains(name, "$") || len(name) > 120
}
