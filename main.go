package main

import (
	"fmt"
	"os"

	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/ViniZap4/devnook-tui/internal/config"
	"github.com/ViniZap4/devnook-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := config.Load()
	client := api.NewClient(cfg.ServerURL, cfg.Token)

	p := tea.NewProgram(ui.NewModel(client, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
