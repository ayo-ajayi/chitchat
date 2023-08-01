package chat

type ChatGPTModel string
type ChatGPTModelRole string

const (
	ChatGPTModelRoleUser      ChatGPTModelRole = "user"
	ChatGPTModelRoleSystem    ChatGPTModelRole = "system"
	ChatGPTModelRoleAssistant ChatGPTModelRole = "assistant"
)

type ChatMessage struct {
	Role    ChatGPTModelRole `json:"role"`
	Content string           `json:"content"`
}
type RequestStruct struct {
	Model            ChatGPTModel  `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	Temperature      float64       `json:"temperature,omitempty"`
	Top_P            float64       `json:"top_p,omitempty"`
	N                int           `json:"n,omitempty"`
	MaxTokens        int           `json:"max_tokens,omitempty"`
	PresencePenalty  float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64       `json:"frequency_penalty,omitempty"`
	User             string        `json:"user,omitempty"`
}
type ResponseStruct struct {
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}
type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int64   `json:"index"`
	Message      Message `json:"message"`
}
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}
type Usage struct {
	CompletionTokens int64 `json:"completion_tokens"`
	PromptTokens     int64 `json:"prompt_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}
