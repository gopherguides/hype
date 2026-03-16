package hype

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type LinkValidator struct {
	Config LinkCheckConfig
	client *http.Client

	mu       sync.RWMutex
	cache    map[string]error
	limiters map[string]*rate.Limiter
}

func NewLinkValidator(cfg LinkCheckConfig) *LinkValidator {
	v := &LinkValidator{
		Config:   cfg,
		cache:    make(map[string]error),
		limiters: make(map[string]*rate.Limiter),
	}

	v.client = &http.Client{
		Timeout: cfg.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= cfg.MaxRedirects {
				return fmt.Errorf("too many redirects (%d)", len(via))
			}
			return nil
		},
	}

	return v
}

func (v *LinkValidator) Check(ctx context.Context, rawURL string) error {
	if v.shouldSkip(rawURL) {
		return nil
	}

	v.mu.RLock()
	if err, ok := v.cache[rawURL]; ok {
		v.mu.RUnlock()
		return err
	}
	v.mu.RUnlock()

	u, err := url.Parse(rawURL)
	if err != nil {
		lce := LinkCheckError{URL: rawURL, Err: err}
		v.cacheResult(rawURL, lce)
		return lce
	}

	limiter := v.limiterFor(u.Host)
	if err := limiter.Wait(ctx); err != nil {
		return LinkCheckError{URL: rawURL, Err: err}
	}

	checkErr := v.doCheck(ctx, rawURL)
	v.cacheResult(rawURL, checkErr)
	return checkErr
}

func (v *LinkValidator) shouldSkip(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	switch u.Scheme {
	case "http", "https":
	case "":
		return true
	default:
		return true
	}

	for _, pattern := range v.Config.ExcludePatterns {
		matched, err := path.Match(pattern, rawURL)
		if err == nil && matched {
			return true
		}
		if strings.Contains(rawURL, strings.ReplaceAll(pattern, "*", "")) {
			continue
		}
	}

	return false
}

func (v *LinkValidator) limiterFor(host string) *rate.Limiter {
	v.mu.Lock()
	defer v.mu.Unlock()

	if l, ok := v.limiters[host]; ok {
		return l
	}

	l := rate.NewLimiter(rate.Limit(v.Config.RatePerHost), v.Config.RateBurst)
	v.limiters[host] = l
	return l
}

func (v *LinkValidator) doCheck(ctx context.Context, rawURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, nil)
	if err != nil {
		return LinkCheckError{URL: rawURL, Err: err}
	}
	req.Header.Set("User-Agent", "hype-link-checker/1.0")

	resp, err := v.client.Do(req)
	if err != nil {
		return LinkCheckError{URL: rawURL, Err: err}
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode == http.StatusMethodNotAllowed {
		return v.doGet(ctx, rawURL)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return v.retryAfter(ctx, rawURL, resp)
	}

	if !v.isAccepted(resp.StatusCode) {
		return LinkCheckError{URL: rawURL, StatusCode: resp.StatusCode}
	}

	return nil
}

func (v *LinkValidator) doGet(ctx context.Context, rawURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return LinkCheckError{URL: rawURL, Err: err}
	}
	req.Header.Set("User-Agent", "hype-link-checker/1.0")

	resp, err := v.client.Do(req)
	if err != nil {
		return LinkCheckError{URL: rawURL, Err: err}
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return v.retryAfter(ctx, rawURL, resp)
	}

	if !v.isAccepted(resp.StatusCode) {
		return LinkCheckError{URL: rawURL, StatusCode: resp.StatusCode}
	}

	return nil
}

func (v *LinkValidator) retryAfter(ctx context.Context, rawURL string, resp *http.Response) error {
	wait := parseRetryAfter(resp.Header.Get("Retry-After"))
	wait = min(wait, 60*time.Second)
	if wait <= 0 {
		wait = 5 * time.Second
	}

	select {
	case <-ctx.Done():
		return LinkCheckError{URL: rawURL, Err: ctx.Err()}
	case <-time.After(wait):
	}

	return v.doGet(ctx, rawURL)
}

func parseRetryAfter(val string) time.Duration {
	if val == "" {
		return 0
	}

	if secs, err := strconv.Atoi(val); err == nil {
		return time.Duration(secs) * time.Second
	}

	if t, err := http.ParseTime(val); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}

	return 0
}

func (v *LinkValidator) isAccepted(code int) bool {
	return slices.Contains(v.Config.AcceptedCodes, code)
}

func (v *LinkValidator) cacheResult(rawURL string, err error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.cache[rawURL] = err
}
