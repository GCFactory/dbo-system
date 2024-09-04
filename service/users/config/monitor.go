package config

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Jaeger config
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}
