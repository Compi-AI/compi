package llm

type Prompt struct {
	System string `json:"system"`
	User   string `json:"user"`
}

type Conversation struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}
