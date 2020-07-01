package fazzthrottle

// LimitType for limit type
type LimitType string

//noinspection GoSnakeCaseUsage
const (
	// IP IP
	IP LimitType = "IP"
	// ENDPOINT ENDPOINT
	ENDPOINT LimitType = "ENDPOINT"
	// IP_ENDPOINT IP_ENDPOINT
	IP_ENDPOINT LimitType = "IP_ENDPOINT"
)
