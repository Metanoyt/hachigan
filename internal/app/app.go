package app

import (
	"context"

	"github.com/hachigan/hachigan/internal/config"
	"github.com/hachigan/hachigan/internal/orchestrator"
	"github.com/hachigan/hachigan/internal/providers/kubernetes"
	"github.com/hachigan/hachigan/internal/tui"
)

func New(_ context.Context, cfg config.Config) (tui.Model, error) {
	provider, err := kubernetes.New(cfg)
	if err != nil {
		return tui.Model{}, err
	}
	orch := orchestrator.New(provider)
	return tui.New(orch, cfg.RefreshInterval), nil
}
