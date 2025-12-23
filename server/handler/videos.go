package handler

import (
	"encoding/json"
	"go-lang-web-api/server/service"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {
	// クエリ取得
	// keyword
	keyword := r.URL.Query().Get("keyword")
	// 必須チェック
	if keyword == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}
	keywordEscaped := url.QueryEscape(keyword) // 日本語対応

	// limit
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // デフォルト値
	// 必須チェック
	if limitStr == "" {
		http.Error(w, "limit is required", http.StatusBadRequest)
		return
	}
	// 数値チェック
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "limit must be a number", http.StatusBadRequest)
			return
		}
		limit = l
	}
	// 最大限チェック（～50）
	if limit > 50 {
		http.Error(w, "limit must be 50 or less", http.StatusBadRequest)
		return
	}

	// since
	var sinceTime *time.Time
	sinceStr := r.URL.Query().Get("since")

	if sinceStr != "" { // 初回検索の時はnil
		t, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			http.Error(w, "since must be RFC3339 format", http.StatusBadRequest)
			return
		}
		sinceTime = &t
	}

	// service呼び出し
	result, err := service.SearchLatestVideos(keywordEscaped, limit, sinceTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// レスポンスを丸ごと読む用（一時的）
func mustRead(resp *http.Response) []byte {
	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	b, _ := json.Marshal(data)
	return b
}
