package util

import (
	"strings"

	"github.com/cyralinc/sidecar-failopen/internal/config"
)

func ParseOptString(cfg config.RepoConfig) string {
	if cfg.PGStringOptions == "" {
		return ""
	}
	opts := strings.Split(cfg.PGStringOptions, ",")
	return "?" + strings.Join(opts, "&")
}
