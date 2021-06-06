package slack

type PlanTextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type HeaderBlock struct {
	Type string        `json:"type"`
	Text PlanTextBlock `json:"text"`
}

type ButtonBlock struct {
	Type  string        `json:"type"`
	Text  PlanTextBlock `json:"text"`
	Style *string       `json:"style,omitempty"`
	Url   string        `json:"url"`
}

type ButtonSection struct {
	Type     string        `json:"type"`
	Elements []ButtonBlock `json:"elements"`
}

type WebhookMessage struct {
	Blocks []interface{} `json:"blocks"`
}
