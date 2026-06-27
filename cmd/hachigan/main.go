package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"

	"github.com/hachigan/hachigan/internal/app"
	"github.com/hachigan/hachigan/internal/config"
)

func main() {
	klog.SetLogger(logr.Discard())
	klog.SetOutput(io.Discard)

	configPath := flag.String("config", "", "path to Hachigan YAML config")
	kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}
	if *kubeconfig != "" {
		cfg.Kubeconfig = *kubeconfig
	}

	model, err := app.New(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start Hachigan: %v\n", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "hachigan exited with error: %v\n", err)
		os.Exit(1)
	}
}
