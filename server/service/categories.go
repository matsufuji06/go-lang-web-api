package service

import (
	"fmt"
	util "go-lang-web-api/server/service/util"
)

// ビデオのカテゴリを取得
func FetchCategories(regionCode string) (map[string]string, error) {

	categoriesCodes := make(map[string]string)
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

	// ビデオカテゴリ取得
	req := service.VideoCategories.List([]string{"snippet"}).RegionCode(regionCode)

	res, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response")
	}

	items := res.Items

	// ビデオカテゴリをcategoriesCodesに格納
	for _, item := range items {
		categoriesCodes[item.Id] = item.Snippet.Title
	}
	return categoriesCodes, nil
}
