package service

import (
	"fmt"
	utils "go-lang-web-api/server/service/utils"
)

// 国のISOコードを取得
func FetchRegionCode() (map[string]struct{}, error) {
	// 国のISOコード
	codes := make(map[string]struct{})
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
	// 国のISOコードを取得
	req := service.I18nRegions.List([]string{"snippet"})

	res, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response")
	}

	items := res.Items

	// 国のISOコードをcodesに格納
	for _, item := range items {
		codes[item.Snippet.Gl] = struct{}{}
	}
	return codes, nil
}
