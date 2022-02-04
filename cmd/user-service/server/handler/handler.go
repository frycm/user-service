package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	jsonErr := json.NewEncoder(w).Encode(struct {
		Err string `json:"error"`
	}{
		Err: err.Error(),
	})
	if jsonErr != nil {
		zap.L().Error("failed to write json error", zap.Error(jsonErr), zap.String("originalError", err.Error()))
	}
}
