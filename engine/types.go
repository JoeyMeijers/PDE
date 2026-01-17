package engine

type Step struct {
	ID         string `json:"id"`
	Executor   string `json:"executor"`
	Language   string `json:"language"`
	Image      string `json:"image"`
	Function   string `json:"function,omitempty"`
	Input      string `json:"input"`
	Output     string `json:"output"`
	Type       string `json:"type,omitempty"`
	DurationMs int    `json:"duration_ms,omitempty"`
}

type Strategy struct {
	Source struct {
		Path string `json:"path"`
	} `json:"source"`
	Pipeline []Step `json:"pipeline"`
	Sink     struct {
		Path  string `json:"path"`
		Input string `json:"input"`
	} `json:"sink"`
}
