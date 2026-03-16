package hype

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_DefaultLinkCheckConfig(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cfg := DefaultLinkCheckConfig()

	r.False(cfg.Enabled)
	r.Equal(10*time.Second, cfg.Timeout)
	r.Contains(cfg.AcceptedCodes, 200)
	r.Contains(cfg.AcceptedCodes, 301)
	r.Contains(cfg.AcceptedCodes, 302)
	r.Equal(float64(2), cfg.RatePerHost)
	r.Equal(1, cfg.RateBurst)
	r.Equal(10, cfg.MaxRedirects)
	r.Empty(cfg.ExcludePatterns)
}
