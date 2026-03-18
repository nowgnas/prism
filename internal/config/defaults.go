package config

import (
	"os"
)

// Defaults returns the factory-default configuration.
func Defaults() Config {
	return Config{
		Shell:        detectShell(),
		StartupTabs:  1,
		SidebarWidth: 20,
		PollInterval: 500,
		AIProcesses: []string{
			"claude", "aider", "llm", "ollama",
			"sgpt", "gpt4all", "chatgpt",
		},
		Theme:        "dark",
		MouseEnabled: false,
	}
}

func detectShell() string {
	if s := os.Getenv("SHELL"); s != "" {
		return s
	}
	return "/bin/sh"
}
