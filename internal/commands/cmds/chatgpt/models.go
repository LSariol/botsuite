package chatgpt

type GPTResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	CreatedAt         int64          `json:"created_at"`
	Status            string         `json:"status"`
	Error             *ErrorObject   `json:"error"`
	Incomplete        *Incomplete    `json:"incomplete_details"`
	Instructions      any            `json:"instructions"`
	MaxOutputTokens   any            `json:"max_output_tokens"`
	Model             string         `json:"model"`
	Output            []OutputItem   `json:"output"`
	ParallelToolCalls bool           `json:"parallel_tool_calls"`
	PreviousResponse  any            `json:"previous_response_id"`
	Reasoning         Reasoning      `json:"reasoning"`
	Store             bool           `json:"store"`
	Temperature       float64        `json:"temperature"`
	Text              TextFormat     `json:"text"`
	ToolChoice        string         `json:"tool_choice"`
	Tools             []any          `json:"tools"`
	TopP              float64        `json:"top_p"`
	Truncation        string         `json:"truncation"`
	Usage             Usage          `json:"usage"`
	User              any            `json:"user"`
	Metadata          map[string]any `json:"metadata"`
}

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Incomplete struct {
	Reason string `json:"reason"`
}

type OutputItem struct {
	Type    string        `json:"type"`
	ID      string        `json:"id"`
	Status  string        `json:"status"`
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
	// annotations is an array, but empty → keep flexible
	Annotations []any `json:"annotations"`
}

type Reasoning struct {
	Effort  any    `json:"effort"`
	Summary string `json:"summary"`
}

type TextFormat struct {
	Format TextFormatInner `json:"format"`
}

type TextFormatInner struct {
	Type string `json:"type"`
}

type Usage struct {
	InputTokens         int                `json:"input_tokens"`
	InputTokensDetails  TokenDetails       `json:"input_tokens_details"`
	OutputTokens        int                `json:"output_tokens"`
	OutputTokensDetails OutputTokenDetails `json:"output_tokens_details"`
	TotalTokens         int                `json:"total_tokens"`
}

type TokenDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type OutputTokenDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}
