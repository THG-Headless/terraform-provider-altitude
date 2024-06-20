package client

type MTELoggingEndpointsDto struct {
	Type			string						`json:"type"`
	EnvironmentId	string						`json:"environmentId"`
	Config 			MTELoggingEndpointsConfig	`json:"Config"`
}

type MTELoggingEndpointsConfig struct {
}
