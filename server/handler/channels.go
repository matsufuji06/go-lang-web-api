package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-lang-web-api/server/service"
	"go-lang-web-api/server/service/utils"
)

func GetChannelInfo(w http.ResponseWriter, r *http.Request) {
	// URLからクエリの内容を取ってくる
	query := r.URL.Query()
	channel := query.Get("query")

	// チャン名が空の場合は、400を返す
	if channel == "" {
		utils.ResponseErrorJson(w, http.StatusBadRequest, "Need to enter channel name") // 400
		return
	} else {
		result, err := service.FetchChannelInfo(channel)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrChannelNotFound):
				utils.ResponseErrorJson(w, http.StatusNotFound, "Channel not found") // 404
			case errors.Is(err, service.ErrRateLimitExceeded):
				utils.ResponseErrorJson(w, http.StatusTooManyRequests, "Request limit exceeded") // 429
			default:
				utils.ResponseErrorJson(w, http.StatusInternalServerError, "Internal server error") // 500
			}
			return
		}

		// レスポンスを返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
