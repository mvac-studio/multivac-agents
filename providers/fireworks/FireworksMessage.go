package fireworks

type FireworksRequestMessage struct {
	Role    string             `json:"role"`
	Content []FireworksContent `json:"content"`
}

type FireworksResponseMessage struct {
	Role    string           `json:"role"`
	Content FireworksContent `json:"content"`
}
