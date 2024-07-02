package client

type MTEConfigDto struct {
	Routes    []RoutesDto   `json:"routes"`
	BasicAuth *BasicAuthDto `json:"basicAuth,omitempty"`
	Cache     []CacheDto    `json:"cache,omitempty"`
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
	AppendPathPrefix   string         `json:"appendPathPrefix,omitempty"`
	ShieldLocation     ShieldLocation `json:"shieldLocation,omitempty"`
}

type CacheDto struct {
	Keys      *CacheKeyDto `json:"keys,omitempty"`
	TtlSeconds *int64      `json:"ttlSeconds,omitempty"`
	PathRules  *MatcherDto  `json:"pathRules,omitempty"`
}

type MatcherDto struct {
	AnyMatch  []string `json:"anyMatch,omitempty"`
	NoneMatch []string `json:"noneMatch,omitempty"`
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
