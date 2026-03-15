package blog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caddyserver/caddy/v2"
)

type ProductionConfig struct {
	PublicDir string
	Addr      string
}

func StartProductionServer(cfg ProductionConfig) error {
	if cfg.PublicDir == "" {
		return fmt.Errorf("public directory is required")
	}
	if cfg.Addr == "" {
		cfg.Addr = ":3000"
	}

	caddyCfg, err := buildCaddyConfig(cfg)
	if err != nil {
		return fmt.Errorf("building caddy config: %w", err)
	}

	return caddy.Run(caddyCfg)
}

func StopProductionServer() error {
	return caddy.Stop()
}

func buildCaddyConfig(cfg ProductionConfig) (*caddy.Config, error) {
	securityHeaders := map[string][]string{
		"X-Content-Type-Options": {"nosniff"},
		"X-Frame-Options":       {"DENY"},
		"Referrer-Policy":       {"strict-origin-when-cross-origin"},
		"Permissions-Policy":    {"camera=(), microphone=(), geolocation=()"},
	}

	staticCacheHeaders := map[string][]string{
		"Cache-Control": {"public, max-age=31536000, immutable"},
	}

	htmlCacheHeaders := map[string][]string{
		"Cache-Control": {"public, max-age=3600, must-revalidate"},
	}

	encodeHandler := json.RawMessage(`{
		"handler": "encode",
		"encodings": {
			"zstd": {},
			"gzip": {}
		},
		"prefer": ["zstd", "gzip"]
	}`)

	securityHeadersJSON, err := json.Marshal(securityHeaders)
	if err != nil {
		return nil, err
	}
	securityHeaderHandler := json.RawMessage(fmt.Sprintf(`{
		"handler": "headers",
		"response": {
			"set": %s
		}
	}`, securityHeadersJSON))

	staticCacheJSON, err := json.Marshal(staticCacheHeaders)
	if err != nil {
		return nil, err
	}
	staticCacheHandler := json.RawMessage(fmt.Sprintf(`{
		"handler": "headers",
		"response": {
			"set": %s
		}
	}`, staticCacheJSON))

	htmlCacheJSON, err := json.Marshal(htmlCacheHeaders)
	if err != nil {
		return nil, err
	}
	htmlCacheHandler := json.RawMessage(fmt.Sprintf(`{
		"handler": "headers",
		"response": {
			"set": %s
		}
	}`, htmlCacheJSON))

	fileServerHandler := json.RawMessage(fmt.Sprintf(`{
		"handler": "file_server",
		"root": %q,
		"index_names": ["index.html"]
	}`, cfg.PublicDir))

	routes := []json.RawMessage{
		mustJSON(map[string]any{
			"handle": []json.RawMessage{encodeHandler, securityHeaderHandler},
		}),
		mustJSON(map[string]any{
			"match": []map[string]any{
				{"path": []string{"*.css", "*.js", "*.woff", "*.woff2", "*.ttf", "*.eot", "*.otf", "*.svg", "*.png", "*.jpg", "*.jpeg", "*.gif", "*.webp", "*.ico", "*.avif"}},
			},
			"handle": []json.RawMessage{staticCacheHandler},
		}),
		mustJSON(map[string]any{
			"match": []map[string]any{
				{"path": []string{"*.html", "/", "/*/"}},
			},
			"handle": []json.RawMessage{htmlCacheHandler},
		}),
		mustJSON(map[string]any{
			"handle":   []json.RawMessage{fileServerHandler},
			"terminal": true,
		}),
	}

	serverCfg := map[string]any{
		"listen": []string{cfg.Addr},
		"routes": routes,
	}

	notFoundPage := filepath.Join(cfg.PublicDir, "404.html")
	if _, err := os.Stat(notFoundPage); err == nil {
		serverCfg["errors"] = map[string]any{
			"routes": []map[string]any{
				{
					"handle": []json.RawMessage{
						json.RawMessage(`{
							"handler": "rewrite",
							"uri": "/404.html"
						}`),
						fileServerHandler,
					},
				},
			},
		}
	}

	httpApp := map[string]any{
		"servers": map[string]any{
			"srv0": serverCfg,
		},
	}

	httpAppJSON, err := json.Marshal(httpApp)
	if err != nil {
		return nil, err
	}

	return &caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: true,
		},
		Logging: &caddy.Logging{
			Logs: map[string]*caddy.CustomLog{
				"default": {
					BaseLog: caddy.BaseLog{
						Level: "WARN",
					},
				},
			},
		},
		AppsRaw: caddy.ModuleMap{
			"http": httpAppJSON,
		},
	}, nil
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
