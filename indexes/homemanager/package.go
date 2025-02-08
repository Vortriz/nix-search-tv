package homemanager

type Package struct {
	Key string `json:"_key"`
	// Subs    []string
	Example map[string]any `json:"example"`

	Type         string         `json:"type"`
	Description  string         `json:"description"`
	Declarations []Declarations `json:"declarations"`
	Default      Default        `json:"default"`
}

type Example struct {
}

type Declarations struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Default struct {
	Text string `json:"text"`
}
