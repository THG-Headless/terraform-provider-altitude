package client

type MTELoggingEndpointsDto struct {
	Endpoints []MTELoggingEndpoint `json:"endpoints"`
}

type MTELoggingEndpoint struct {
	Type          string                    `json:"type"`
	EnvironmentId string                    `json:"environmentId"`
	Config        MTELoggingEndpointsConfig `json:"config"`
}

type MTELoggingEndpointsConfig struct {
	Dataset   string            `json:"dataset"`
	ProjectId string            `json:"projectId"`
	Table     string            `json:"table"`
	Email     string            `json:"email"`
	Headers   []BQLoggingHeader `json:"headers"`
	SecretKey string            `json:"secretKey"`
}

type BQLoggingHeader struct {
	ColumnName   string `json:"columnName"`
	HeaderName   string `json:"headerName"`
	DefaultValue string `json:"defaultValue"`
}
