package util

import (
	"strings"

	"github.com/cyralinc/sidecar-failopen/internal/config"
)

func ParseOptString(cfg config.RepoConfig) string {
	if cfg.PGConfig.ConnectionStringOptions == "" {
		return ""
	}
	opts := strings.Split(cfg.PGConfig.ConnectionStringOptions, ",")
	return "?" + strings.Join(opts, "&")
}
