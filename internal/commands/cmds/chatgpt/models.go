package chatgpt

type GPTResponse struct {
	ID          string    `json:"id"`
	Object      string    `json:"object"`
	Created     int       `json:"created"`
	Model       string    `json:"model"`
	Choices     []Choices `json:"choices"`
	Usage       Usage     `json:"usage"`
	ServiceTier string    `json:"service_tier"`
}

type Choices struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Logprobs     string  `json:"logprobs,omitempty"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role        string   `json:"role"`
	Content     string   `json:"content"`
	Refusal     string   `json:"refusal"`
	Annotations []string `json:"annotations"`
}

type Usage struct {
	PromptTokens            int                     `json:"prompt_tokens"`
	CompletionTokens        int                     `json:"completion_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	PromptTokenDetails      PromptTokenDetails      `json:"prompt_token_details"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type PromptTokenDetails struct {
	CachedTokens int `json:"cached_tokens"`
	AudioTokens  int `json:"audio_tokens"`
}

type CompletionTokensDetails struct {
	ReasoningTokens          int `json:"reasoning_tokens"`
	AudioTokens              int `json:"audio_tokens"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
}
