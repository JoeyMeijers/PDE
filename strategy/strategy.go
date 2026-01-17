package strategy

type Strategy struct {
	Source   Source         `json:"source"`
	Pipeline []PipelineStep `json:"pipeline"`
	Sink     Sink           `json:"sink"`
}

type Source struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type Sink struct {
	Type  string `json:"type"`
	Path  string `json:"path"`
	Input string `json:"input"`
}

type PipelineStep struct {
	ID         string `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Executor   string `json:"executor,omitempty"`
	Image      string `json:"image,omitempty"`
	Function   string `json:"function,omitempty"` // nieuw
	Input      string `json:"input,omitempty"`
	Output     string `json:"output,omitempty"`
	DurationMs int    `json:"duration_ms,omitempty"`
}
