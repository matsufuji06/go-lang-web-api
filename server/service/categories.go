package service

import (
	"fmt"
	utils "go-lang-web-api/server/service/utils"
)

// ビデオのカテゴリを取得
func FetchCategories(regionCode string) (map[string]string, error) {

	categoriesCodes := make(map[string]string)
// APIキー作成
	apiKey, err := utils.GetApiKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create API key")
	}
	// サービス作成
	service, err := utils.CreateService(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create service")
	}

	// ビデオカテゴリ取得
	req := service.VideoCategories.List([]string{"snippet"}).RegionCode(regionCode).Hl("ja")

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
