package youtube

type YouTubeSearchResponse struct {
	Items []struct {
		ID struct {
			ChannelID string `json:"channelId"`
		} `json:"id"`
	} `json:"items"`
}
