package util

import (
	"strings"

	"github.com/cyralinc/sidecar-failopen/internal/config"
)

func ParseOptString(cfg config.RepoConfig) string {
	if cfg.ConnectionStringOptions == "" {
		return ""
	}
	opts := strings.Split(cfg.ConnectionStringOptions, ",")
	return "?" + strings.Join(opts, "&")
}
