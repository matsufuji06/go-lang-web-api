package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-lang-web-api/server/service"
)

func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	// CORS対応 (これしないとうまくいかない)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// URLからクエリの内容を取ってくる
	query := r.URL.Query()
	channel := query.Get("query")

	// チャン名が空の場合は、400を返す
	if channel == "" {
		writeJSONError(w, http.StatusBadRequest, "Need to enter channel name") // 400
		return
	} else {
		result, err := service.GetChannelInfo(channel)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrChannelNotFound):
				writeJSONError(w, http.StatusNotFound, "Channel not found") // 404
			case errors.Is(err, service.ErrRateLimitExceeded):
				writeJSONError(w, http.StatusTooManyRequests, "Request limit exceeded") // 429
			default:
				writeJSONError(w, http.StatusInternalServerError, "Internal server error") // 500
			}
			return
		}

		// レスポンスを返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// エラーのレスポンスを作る関数
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
	})
}

// エラー用のstruct
type ErrorResponse struct {
	Message string `json:"message"`
}
