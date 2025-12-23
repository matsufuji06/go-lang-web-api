package service

import (
	"encoding/json"
	"fmt"
	apiModel "go-lang-web-api/server/model/api"
	youtubeModel "go-lang-web-api/server/model/youtube"
	"net/http"
	"net/url"
	"os"
	"time"
)

func SearchLatestVideos(keyword string, limit int, since *time.Time) (*apiModel.VideoResponse, error) {
	// APIキー取得
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// YouTube API URLパラメータ組み立て
	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("q", keyword)
	params.Set("order", "date")
	params.Set("type", "video")
	params.Set("maxResults", fmt.Sprintf("%d", limit))
	params.Set("key", apiKey)

	if since != nil {
		params.Set("publishedAfter", since.Format(time.RFC3339))
	}

	apiURL := "https://www.googleapis.com/youtube/v3/search?" + params.Encode()

	// API呼び出し
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// YouTube API用structにデコードして詰める
	var ytResp youtubeModel.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&ytResp); err != nil {
		return nil, err
	}

	items := []apiModel.VideoItem{}
	var latest time.Time

	// 自分のAPI用struct に変換
	for _, item := range ytResp.Items {
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		// 前回取得した時刻（since）よりも後に公開された動画の場合
		isNew := false
		if since != nil && publishedAt.After(*since) {
			isNew = true
		}

		items = append(items, apiModel.VideoItem{
			VideoID:     item.ID.VideoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
			ChannelID:   item.Snippet.ChannelID,
			ChannelName: item.Snippet.ChannelTitle,
			URL:         "https://www.youtube.com/watch?v=" + item.ID.VideoID,
			IsNew:       isNew,
			Thumbnails: apiModel.Thumbnails{
				Default: apiModel.Thumbnail{
					URL:    item.Snippet.Thumbnails.Default.URL,
					Width:  item.Snippet.Thumbnails.Default.Width,
					Height: item.Snippet.Thumbnails.Default.Height,
				},
				Medium: apiModel.Thumbnail{
					URL:    item.Snippet.Thumbnails.Medium.URL,
					Width:  item.Snippet.Thumbnails.Medium.Width,
					Height: item.Snippet.Thumbnails.Medium.Height,
				},
				High: apiModel.Thumbnail{
					URL:    item.Snippet.Thumbnails.High.URL,
					Width:  item.Snippet.Thumbnails.High.Width,
					Height: item.Snippet.Thumbnails.High.Height,
				},
			},
		})

		// 最新の投稿日時を更新する
		if publishedAt.After(latest) {
			latest = publishedAt
		}
	}

	return &apiModel.VideoResponse{
		Items: items,
		Meta: apiModel.Meta{
			Count:           len(items),
			LastPublishedAt: latest.Format(time.RFC3339),
		},
	}, nil

}
