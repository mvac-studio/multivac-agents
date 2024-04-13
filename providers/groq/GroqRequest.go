package groq

import "multivac.network/services/agents/providers"

type GroqRequest struct {
	Messages []providers.Message `json:"messages"`
	Model    string              `json:"model"`
}
