package tui

import (
	"strings"

	"github.com/Drack112/go-youtube/internal/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	result models.SearchResult
}

func (i item) FilterValue() string {
	return i.result.Title
}
func (i item) Title() string {
	var badges []string

	if i.result.IsLive {
		badges = append(badges, "LIVE")
	}

	if i.result.IsShort {
		badges = append(badges, "SHORT")
	}

	if i.result.Duration != "" && !i.result.IsLive {
		badges = append(badges, i.result.Duration)
	}

	title := i.result.Title
	if len(badges) > 0 {
		title += " " + strings.Join(badges, " ")
	}

	return title
}

func (i item) Description() string {
	var parts []string
	parts = append(parts, i.result.URL)

	if i.result.ChannelName != "" {
		parts = append(parts, i.result.ChannelName)
	}

	return strings.Join(parts, " * ")
}

func NewItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// Customize the styles
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#8B5CF6")).
		BorderForeground(lipgloss.Color("#8B5CF6"))

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#94A3B8"))

	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(lipgloss.Color("#F8FAFC"))

	d.Styles.NormalDesc = d.Styles.NormalDesc.
		Foreground(lipgloss.Color("#64748B"))

	d.Styles.FilterMatch = lipgloss.NewStyle().Underline(true)

	return d
}
