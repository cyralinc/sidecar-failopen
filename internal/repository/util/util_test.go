package util

import (
	"testing"

	"github.com/cyralinc/sidecar-failopen/internal/config"
)

func TestParseOptString(t *testing.T) {
	type args struct {
		cfg config.RepoConfig
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid single config",
			args: args{
				cfg: config.RepoConfig{
					ConnectionStringOptions: "this=that",
				},
			},
			want: "?this=that",
		},
		{

			name: "valid multiple config",
			args: args{
				cfg: config.RepoConfig{
					ConnectionStringOptions: "this=that,that=this",
				},
			},
			want: "?this=that&that=this",
		},
		{

			name: "valid empty config",
			args: args{
				cfg: config.RepoConfig{
					ConnectionStringOptions: "",
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseOptString(tt.args.cfg); got != tt.want {
				t.Errorf("ParseOptString() = %v, want %v", got, tt.want)
			}
		})
	}
}
