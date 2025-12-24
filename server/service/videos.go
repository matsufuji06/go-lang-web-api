package service

import (
	"context"
	"fmt"
	apiModel "go-lang-web-api/server/model/api"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func SearchLatestVideos(keyword string, limit int, since *time.Time) (*apiModel.VideoResponse, error) {
	// APIキー取得
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	ctx := context.Background()

	// Youtubeクライアント生成
	youtubeService, err := youtube.NewService(
		ctx,
		option.WithAPIKey(apiKey),
	)
	if err != nil {
		return nil, err
	}

	// リクエスト内容を生成
	call := youtubeService.Search.
		List([]string{"snippet"}).
		Q(keyword).
		Order("date").
		Type("video").
		MaxResults(int64(limit))

	if since != nil {
		call = call.PublishedAfter(since.Format(time.RFC3339))
	}

	ytResp, err := call.Do()
	if err != nil {
		return nil, err
	}

	// 自分のAPI用struct に変換
	items := []apiModel.VideoItem{}
	var latest time.Time

	for _, item := range ytResp.Items {
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		// 前回取得した時刻（since）よりも後に公開された動画の場合
		isNew := false
		if since != nil && publishedAt.After(*since) {
			isNew = true
		}

		items = append(items, apiModel.VideoItem{
			VideoID:     item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
			ChannelID:   item.Snippet.ChannelId,
			ChannelName: item.Snippet.ChannelTitle,
			URL:         "https://www.youtube.com/watch?v=" + item.Id.VideoId,
			IsNew:       isNew,
			Thumbnails: apiModel.Thumbnails{
				Default: apiModel.Thumbnail{
					URL:    item.Snippet.Thumbnails.Default.Url,
					Width:  int(item.Snippet.Thumbnails.Default.Width),
					Height: int(item.Snippet.Thumbnails.Default.Height),
				},
				Medium: apiModel.Thumbnail{
					URL:    item.Snippet.Thumbnails.Medium.Url,
					Width:  int(item.Snippet.Thumbnails.Medium.Width),
					Height: int(item.Snippet.Thumbnails.Medium.Height),
				},
				High: apiModel.Thumbnail{
					URL:    item.Snippet.Thumbnails.High.Url,
					Width:  int(item.Snippet.Thumbnails.High.Width),
					Height: int(item.Snippet.Thumbnails.High.Height),
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
