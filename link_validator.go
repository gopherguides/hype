package hype

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	inflight map[string]chan struct{}
	limiters map[string]*rate.Limiter
}

func NewLinkValidator(cfg LinkCheckConfig) *LinkValidator {
	if cfg.RatePerHost <= 0 {
		cfg.RatePerHost = 2
	}
	if cfg.RateBurst <= 0 {
		cfg.RateBurst = 1
	}

	v := &LinkValidator{
		Config:   cfg,
		cache:    make(map[string]error),
		inflight: make(map[string]chan struct{}),
		limiters: make(map[string]*rate.Limiter),
	}

	v.client = &http.Client{
		Timeout: cfg.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > cfg.MaxRedirects {
				return fmt.Errorf("too many redirects (%d)", len(via))
			}
			return nil
		},
	}

	return v
}

func (v *LinkValidator) Check(ctx context.Context, rawURL string) error {
	if strings.HasPrefix(rawURL, "//") {
		rawURL = "https:" + rawURL
	}

	if v.shouldSkip(rawURL) {
		return nil
	}

	v.mu.Lock()
	if err, ok := v.cache[rawURL]; ok {
		v.mu.Unlock()
		return err
	}

	if ch, ok := v.inflight[rawURL]; ok {
		v.mu.Unlock()
		select {
		case <-ch:
		case <-ctx.Done():
			return LinkCheckError{URL: rawURL, Err: ctx.Err()}
		}
		v.mu.RLock()
		err := v.cache[rawURL]
		v.mu.RUnlock()
		return err
	}

	ch := make(chan struct{})
	v.inflight[rawURL] = ch
	v.mu.Unlock()

	u, err := url.Parse(rawURL)
	if err != nil {
		lce := LinkCheckError{URL: rawURL, Err: err}
		v.finishCheck(rawURL, lce, ch)
		return lce
	}

	limiter := v.limiterFor(u.Host)
	if err := limiter.Wait(ctx); err != nil {
		lce := LinkCheckError{URL: rawURL, Err: err}
		v.finishCheck(rawURL, lce, ch)
		return lce
	}

	checkErr := v.doCheck(ctx, rawURL)
	v.finishCheck(rawURL, checkErr, ch)
	return checkErr
}

func (v *LinkValidator) finishCheck(rawURL string, err error, ch chan struct{}) {
	v.mu.Lock()
	v.cache[rawURL] = err
	delete(v.inflight, rawURL)
	v.mu.Unlock()
	close(ch)
}

func (v *LinkValidator) shouldSkip(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	switch u.Scheme {
	case "http", "https":
	case "":
		if strings.HasPrefix(rawURL, "//") {
			return false
		}
		return true
	default:
		return true
	}

	for _, pattern := range v.Config.ExcludePatterns {
		if prefix, ok := strings.CutSuffix(pattern, "*"); ok {
			if strings.HasPrefix(rawURL, prefix) {
				return true
			}
		} else if pattern == rawURL {
			return true
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

	if resp.StatusCode == http.StatusTooManyRequests {
		return v.retryAfter(ctx, rawURL, resp)
	}

	if v.isAccepted(resp.StatusCode) {
		return nil
	}

	return v.doGet(ctx, rawURL)
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
