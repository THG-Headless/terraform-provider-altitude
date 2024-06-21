package client

type MTELoggingEndpointsDto struct {
	Type          string                    `json:"type"`
	EnvironmentId string                    `json:"environmentid"`
	Config        MTELoggingEndpointsConfig `json:"config"`
}

type MTELoggingEndpointsConfig struct {
	NonSensititve NonSensitiveBQLoggingConfig `json:"nonsensitive"`
	Sensitive     *SensitiveBQLoggingConfig   `json:"sensitive,omitempty"`
}

type NonSensitiveBQLoggingConfig struct {
	Dataset   string            `json:"dataset"`
	ProjectId string            `json:"projectid"`
	Table     string            `json:"table"`
	Email     string            `json:"email"`
	Headers   []BQLoggingHeader `json:"headers"`
}

type BQLoggingHeader struct {
	Col     string `json:"col"`
	Header  string `json:"header"`
	Default string `json:"default"`
}

type SensitiveBQLoggingConfig struct {
	SecretKey string `json:"secretkey"`
}
