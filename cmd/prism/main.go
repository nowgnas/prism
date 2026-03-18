package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nowgnas/prism/internal/app"
	"github.com/nowgnas/prism/internal/config"
)

var version = "dev"

func main() {
	var (
		configPath  = flag.String("config", config.DefaultPath(), "path to config file")
		showVersion = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("prism %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: config error: %v\n", err)
		cfg = config.Defaults()
	}

	model, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "prism: %v\n", err)
		os.Exit(1)
	}

	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}

	p := tea.NewProgram(model, opts...)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "prism: %v\n", err)
		os.Exit(1)
	}
}
