package groq

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
