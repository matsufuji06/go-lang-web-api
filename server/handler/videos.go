package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {
	// APIキー取得
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		return
	}

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

	// Youtube APIのURLを取得
	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&order=date&type=video&maxResults=%d&key=%s",
		keywordEscaped,
		limit,
		apiKey,
	)

	// 2回目以降の取得の場合、URLにpublishedAfterを付与
	if sinceTime != nil {
		apiURL += "&publishedAfter=" + url.QueryEscape(sinceTime.Format(time.RFC3339))
	}

	// YouTube API 呼び出し
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "failed to call youtube api", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// レスポンスをそのまま返す（まずはここまで）
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(json.RawMessage(mustRead(resp)))
}

// レスポンスを丸ごと読む用（一時的）
func mustRead(resp *http.Response) []byte {
	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	b, _ := json.Marshal(data)
	return b
}
