package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wagoodman/go-partybus"

	"github.com/anchore/bubbly/bubbles/taskprogress"
	"github.com/anubhav06/sanrakshya-cli/internal/log"
	sanrakshyaEventParsers "github.com/anubhav06/sanrakshya-cli/sanrakshya/event/parsers"
)

func (m *Handler) handleFileIndexingStarted(e partybus.Event) []tea.Model {
	path, prog, err := sanrakshyaEventParsers.ParseFileIndexingStarted(e)
	if err != nil {
		log.WithFields("error", err).Warn("unable to parse event")
		return nil
	}

	tsk := m.newTaskProgress(
		taskprogress.Title{
			Default: "Index files system",
			Running: "Indexing file system",
			Success: "Indexed file system",
		},
		taskprogress.WithStagedProgressable(prog),
	)

	tsk.Context = []string{path}

	return []tea.Model{tsk}
}
