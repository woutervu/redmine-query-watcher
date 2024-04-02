package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Tabs        []string
	Table       table.Model
	Issues      []*Issue
	ActiveIssue *Issue
	activeTab   int
}

func main() {
	ec, err := appRun()
	if err != nil {
		fmt.Printf("Exit code: %d\nMessage: %s", ec, err)
		os.Exit(ec)
	}

	os.Exit(ec)
}

func appRun() (int, error) {
	m, err := getModel()
	if err != nil {
		return 1, err
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return 1, err
	}

	return 0, nil
}
