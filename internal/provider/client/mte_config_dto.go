package client

type MTEConfigDto struct {
	Routes    []RoutesDto   `json:"routes"`
	BasicAuth *BasicAuthDto `json:"basicAuth,omitempty"`
}

type BasicAuthDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoutesDto struct {
	Host               string         `json:"host"`
	Path               string         `json:"path"`
	EnableSsl          bool           `json:"enableSsl"`
	PreservePathPrefix bool           `json:"preservePathPrefix"`
	CacheKey           *CacheKeyDto   `json:"cacheKey,omitempty"`
	AppendPathPrefix   string         `json:"appendPathPrefix,omitempty"`
	ShieldLocation     ShieldLocation `json:"shieldLocation,omitempty"`
	CacheMaxAge        string         `json:"cacheMaxAge,omitempty"`
}

type CacheKeyDto struct {
	Header []string `json:"header"`
	Cookie []string `json:"cookie"`
}

type ShieldLocation string

const (
	London        ShieldLocation = "London"
	Manchester    ShieldLocation = "Manchester"
	Frankfurt                    = "Frankfurt"
	Madrid                       = "Madrid"
	New_York_City                = "New York City"
	Los_Angeles                  = "Los Angeles"
	Toronto                      = "Toronto"
	Johannesburg                 = "Johannesburg"
	Seoul                        = "Seoul"
	Sydney                       = "Sydney"
	Tokyo                        = "Tokyo"
	Hong_Kong                    = "Hong Kong"
	Mumbai                       = "Mumbai"
	Singapore                    = "Singapore"
)
