package hype

import "time"

type LinkCheckConfig struct {
	Enabled         bool
	Timeout         time.Duration
	AcceptedCodes   []int
	ExcludePatterns []string
	RatePerHost     float64
	RateBurst       int
	MaxRedirects    int
}

func DefaultLinkCheckConfig() LinkCheckConfig {
	return LinkCheckConfig{
		Enabled:   false,
		Timeout:   10 * time.Second,
		AcceptedCodes: []int{
			200, 201, 202, 203, 204,
			301, 302, 307, 308,
		},
		RatePerHost:  2,
		RateBurst:    1,
		MaxRedirects: 10,
	}
}
