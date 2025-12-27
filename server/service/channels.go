package service

import (
	"errors"
	"fmt"

	"go-lang-web-api/server/model/api"

	utils "go-lang-web-api/server/service/utils"
)

// 404と429のメッセージ
var (
	ErrChannelNotFound   = errors.New("channel not found")   // 404
	ErrRateLimitExceeded = errors.New("rate limit exceeded") // 429
)

// チャンネルの情報を取得する関数
func FetchChannelInfo(channel string) (*api.ChannelResponse, error) {
	// APIキーの作成
	apiKey, err := utils.GetApiKey()
	if err != nil {
		return nil, err
	}

	// サービスの作成
	service, err := utils.CreateService(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create service")
	}

	// search.list (チャンネルIDを取ってくる)
	searchCall := service.Search.List([]string{"snippet"}).
		Q(channel).
		Type("channel").
		MaxResults(1)

	searchResp, err := searchCall.Do()
	if err != nil {
		// err の中身を見て 429 / その他 を判定
		return nil, err
	}

	if len(searchResp.Items) == 0 {
		return nil, ErrChannelNotFound
	}

	channelId := searchResp.Items[0].Snippet.ChannelId

	// チャンネルの情報を取ってくる
	channelCall := service.Channels.List([]string{"snippet", "statistics"}).
		Id(channelId)

	channelResp, err := channelCall.Do()
	if err != nil {
		return nil, err
	}

	if len(channelResp.Items) == 0 {
		return nil, ErrChannelNotFound
	}

	item := channelResp.Items[0]

	// 取った情報を変数に入れる
	title := item.Snippet.Title
	subscriberCount := item.Statistics.SubscriberCount
	viewCount := item.Statistics.ViewCount
	videoCount := item.Statistics.VideoCount
	thumbnail := item.Snippet.Thumbnails.High.Url
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

// 平均再生数を計算する関数
func calcAverageViews(views uint64, videos uint64) uint64 {
	if videos == 0 {
		return 0
	}

	return views / videos
}
