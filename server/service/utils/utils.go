package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go-lang-web-api/server/model/api"
	errRes "go-lang-web-api/server/model/api"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// .envからAPIキーの取得
func GetApiKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set")
	}
	return apiKey, nil
}

func CreateService(apiKey string) (*youtube.Service, error) {
	svs, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create client")
	}
	return svs, nil
}

// 正常系レスポンス
func ResponseJson(w http.ResponseWriter, httpStatus int, data api.Response) {
	setCORS(w)
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// 異常系（エラー発生時）レスポンス
func ResponseErrorJson(w http.ResponseWriter, httpStatus int, err string) {
	setCORS(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(errRes.ErrorResponse{
		Message: err,
	})
}

// CORSを設定
func setCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}