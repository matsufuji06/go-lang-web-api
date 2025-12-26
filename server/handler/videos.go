package handler

import (
	"encoding/json"
	"go-lang-web-api/server/service"
	"net/http"
	"strconv"
	"time"
)

/*
limitパラメータのデフォルト値と最大値
*/
const (
	defaultLimit = 10
	maxLimit     = 50
)

/*
動画情報を取得するハンドラ関数
*/
func GetVideos(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータを取得する
	query := r.URL.Query()

	// keywordパラメータを取得する
	keyword := query.Get("keyword")
	// 必須チェックを行う
	if keyword == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}

	// limitパラメータを取得する
	limit, err := parseLimit(query.Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sinceパラメータを取得する（ポーリング時のみ取得される想定）
	since, err := parseSince(query.Get("since"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 動画情報を取得する
	result, err := service.SearchLatestVideos(keyword, limit, since)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 動画情報をJSON形式でレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

/*
limitパラメータのバリデーション関数
*/
func parseLimit(limitStr string) (int, error) {
	// パラメータが空の場合、デフォルト値を返す
	if limitStr == "" {
		return defaultLimit, nil
	}

	// 文字列を整数に変換する
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, err
	}

	// 範囲チェックを行う
	if limit > maxLimit {
		return 0, strconv.ErrRange
	}

	return limit, nil
}

/*
sinceパラメータのバリデーション関数
*/
func parseSince(sinceStr string) (*time.Time, error) {
	// パラメータが空の場合、nilを返す
	if sinceStr == "" {
		return nil, nil
	}

	// RFC3339形式でパースする
	t, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
