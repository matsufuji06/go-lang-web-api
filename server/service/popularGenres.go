package service

import (
	"fmt"
	api "go-lang-web-api/server/model/api"
	util "go-lang-web-api/server/service/util"
	"sort"
)

// 人気のビデオのジャンルを取得
func FetchPopularGenres(regionCode string) ([]api.GenreRatio, error) {

	// APIキー作成
	apiKey, err := util.GetApiKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create API key")
	}
	// サービス作成
	service, err := util.CreateService(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create service")
	}

	// 人気のビデオ最大200件取得
	req := service.Videos.List([]string{"snippet"}).Chart("mostPopular").
		RegionCode(regionCode).MaxResults(200)

	res, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response")
	}

	totalVideos := len(res.Items)


	categoryVideos := make(map[string][]api.VideoURL)
	
	// カテゴリごとに動画をためる
	for _, item := range res.Items {
		categoryID := item.Snippet.CategoryId
		categoryVideos[categoryID] = append(categoryVideos[categoryID], api.VideoURL{
			VideoURL: "https://www.youtube.com/watch?v=" + item.Id,
		})
	}

	// ジャンルの割合計算
	genres := make([]api.GenreRatio, 0, len(categoryVideos))
	for categoryID, videos := range categoryVideos {
		limit := 2
		if len(videos) < limit {
			limit = len(videos)
		}
		genres = append(genres, api.GenreRatio{
			CategoryID: categoryID,
			Ratio:      float64(len(videos)) / float64(totalVideos) * 100,
			VideoURL:   videos[:limit],
		})
	}

	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Ratio > genres[j].Ratio
	})

	// 上位5件の人気ジャンルを返す
	topN := 5
	if len(genres) < topN {
		topN = len(genres)
	}

	return genres[:topN], nil
}
