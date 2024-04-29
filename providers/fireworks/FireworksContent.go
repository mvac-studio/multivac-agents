package fireworks

type FireworksImageUrl struct {
	Url string `json:"url,omitempty"`
}

type FireworksContent struct {
	Type     string             `json:"type,omitempty"`
	Text     string             `json:"text,omitempty"`
	ImageUrl *FireworksImageUrl `json:"image_url,omitempty"`
}
