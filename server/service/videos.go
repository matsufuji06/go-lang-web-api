package service

import (
	"context"
	"fmt"
	apiModel "go-lang-web-api/server/model/api"
	utils "go-lang-web-api/server/service/utils"
	"time"

	"google.golang.org/api/youtube/v3"
)

/*
タイムアウト時間の定義
*/
const youtubeTimeout = 5 * time.Second

/*
YouTubeサービスを呼び出して動画情報を取得する関数
*/
func SearchLatestVideos(keyword string, limit int, since *time.Time) (*apiModel.VideoResponse, error) {
	// APIキーを取得する
	apiKey, err := utils.GetApiKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create API key")
	}

	// コンテキストを生成する（設定時間以内に処理が終わらなければタイムアウト）
	ctx, cancel := context.WithTimeout(context.Background(), youtubeTimeout)
	// 関数終了時にキャンセルを呼び出す
	defer cancel()

	// Youtubeクライアントを生成する
	youtubeService, err := utils.CreateServiceWithContext(apiKey, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create service")
	}

	// リクエスト内容を生成する
	call := youtubeService.Search.
		List([]string{"snippet"}).
		Q(keyword).
		Order("date").
		Type("video").
		MaxResults(int64(limit))

	// sinceパラメータが指定されていれば、PublishedAfterを設定する
	if since != nil {
		call = call.PublishedAfter(since.Format(time.RFC3339))
	}

	// APIコールを実行する（タイムアウトの場合、errが返る）
	// ※日本語対応のためエスケープ処理も内部的に行われる
	ytResp, err := call.Do()
	if err != nil {
		return nil, err
	}

	// YouTube APIのレスポンスを自分のAPI用のレスポンスに変換する
	items, latest := convertItems(ytResp.Items, since)

	return &apiModel.VideoResponse{
		Items: items,
		Meta: apiModel.Meta{
			Count:           len(items),
			LastPublishedAt: latest.Format(time.RFC3339),
		},
	}, nil
}

/*
YouTube APIのSearchResultから自分のAPI用のVideoItemに変換する関数
*/
func convertItems(ytItems []*youtube.SearchResult, since *time.Time) ([]apiModel.VideoItem, time.Time) {
	items := make([]apiModel.VideoItem, 0, len(ytItems))
	var latest time.Time

	// 各アイテムを変換する
	for _, item := range ytItems {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			continue
		}

		// sinceパラメータが指定されていれば、新着フラグを設定する
		isNew := since != nil && publishedAt.After(*since)

		// API用のVideoItemに変換
		video := apiModel.VideoItem{
			VideoID:     item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
			ChannelID:   item.Snippet.ChannelId,
			ChannelName: item.Snippet.ChannelTitle,
			URL:         "https://www.youtube.com/watch?v=" + item.Id.VideoId,
			IsNew:       isNew,
			Thumbnails:  convertThumbnails(item.Snippet.Thumbnails),
		}

		// 変換後のスライスに追加する
		items = append(items, video)

		// 最新の公開日時を更新する
		if publishedAt.After(latest) {
			latest = publishedAt
		}
	}

	return items, latest
}

/*
YouTubeのThumbnailDetailsを自分のAPI用のThumbnailsに変換する関数
*/
func convertThumbnails(t *youtube.ThumbnailDetails) apiModel.Thumbnails {
	// nilチェック
	if t == nil {
		return apiModel.Thumbnails{}
	}

	return apiModel.Thumbnails{
		Default: toThumbnail(t.Default),
		Medium:  toThumbnail(t.Medium),
		High:    toThumbnail(t.High),
	}
}

/*
YouTubeのThumbnailを自分のAPI用のThumbnailに変換する関数
*/
func toThumbnail(t *youtube.Thumbnail) apiModel.Thumbnail {
	// nilチェック
	if t == nil {
		return apiModel.Thumbnail{}
	}

	return apiModel.Thumbnail{
		URL:    t.Url,
		Width:  int(t.Width),
		Height: int(t.Height),
	}
}
