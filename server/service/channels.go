package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"go-lang-web-api/server/model/api"
	"go-lang-web-api/server/model/youtube"
)

// 404と429のメッセージ
var (
	ErrChannelNotFound   = errors.New("channel not found")   // 404
	ErrRateLimitExceeded = errors.New("rate limit exceeded") // 429
)

// チャンネルの情報を取得する関数
func GetChannelInfo(channel string) (*api.ChannelResponse, error) {
	// APIキーを取得
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// チャンネルIDを取得
	channelId, err := getChannelId(channel, apiKey)

	endpoint := "https://www.googleapis.com/youtube/v3/channels"

	params := url.Values{}
	params.Add("part", "snippet,statistics")
	params.Add("id", channelId)
	params.Add("key", apiKey)

	reqURL := endpoint + "?" + params.Encode()

	resp, err := http.Get((reqURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// ステータスコードを確認
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrChannelNotFound
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimitExceeded
	}

	var youTubeChannelsResponse youtube.YouTubeChannelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&youTubeChannelsResponse); err != nil {
		return nil, err
	}

	// 検索結果が空の場合
	if len(youTubeChannelsResponse.Items) == 0 {
		return nil, ErrChannelNotFound
	}

	// 情報を取ってくる
	title := youTubeChannelsResponse.Items[0].Snippet.Title
	subscriberCount, _ := strconv.ParseUint(youTubeChannelsResponse.Items[0].Statistics.SubscriberCount, 10, 64)
	viewCount, _ := strconv.ParseUint(youTubeChannelsResponse.Items[0].Statistics.ViewCount, 10, 64)
	videoCount, _ := strconv.ParseUint(youTubeChannelsResponse.Items[0].Statistics.VideoCount, 10, 64)
	thumbnail := youTubeChannelsResponse.Items[0].Snippet.Thumbnails.High.URL
	averageViews := calcAverageViews(viewCount, videoCount)

	return &api.ChannelResponse{
		Title:           title,
		SubscriberCount: subscriberCount,
		ViewCount:       viewCount,
		VideoCount:      videoCount,
		AverageViews:    averageViews,
		Thumbnail:       thumbnail,
	}, nil
}

// チャンネルIDを取得する関数(この中でしか使わないため、関数名の先頭は小文字)
func getChannelId(channel string, apiKey string) (string, error) {
	endpoint := "https://www.googleapis.com/youtube/v3/search"

	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("q", channel)
	params.Add("type", "channel")
	params.Add("maxResults", "1")
	params.Add("key", apiKey)

	reqURL := endpoint + "?" + params.Encode()

	resp, err := http.Get(reqURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// ステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to call YouTube Search API")
	}

	var searchResponse youtube.YouTubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return "", err
	}

	// 検索結果が空の場合
	if len(searchResponse.Items) == 0 {
		return "", errors.New("channel not found")
	}

	channelId := searchResponse.Items[0].ID.ChannelID
	if channelId == "" {
		return "", errors.New("channelId is empty")
	}

	return channelId, nil
}

// 平均再生数を計算する関数
func calcAverageViews(views uint64, videos uint64) uint64 {
	if videos == 0 {
		return 0
	}

	return views / videos
}
