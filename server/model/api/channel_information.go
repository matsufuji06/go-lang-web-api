package api

type ChannelResponse struct {
	Title           string `json:"title"`
	SubscriberCount uint64 `json:"subscriberCount"`
	ViewCount       uint64 `json:"viewCount"`
	VideoCount      uint64 `json:"videoCount"`
	AverageViews    uint64 `json:"averageViews"`
	Thumbnail       string `json:"thumbnail"`
}
