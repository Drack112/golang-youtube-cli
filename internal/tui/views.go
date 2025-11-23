package tui

import (
	"fmt"

	"github.com/Drack112/go-youtube/internal/api"
	"github.com/Drack112/go-youtube/internal/player"
	"github.com/Drack112/go-youtube/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) performSearch() tea.Cmd {
	return func() tea.Msg {
		resp, err := api.SearchVideosWithPagination(m.opts.Input, "")
		if err != nil {
			return searchResultsMsg{err: err}
		}
		return searchResultsMsg{
			results:           resp.Results,
			continuationToken: resp.ContinuationToken,
			hasMore:           resp.HasMore,
			err:               nil,
			isLoadMore:        false,
		}
	}
}

func (m *Model) loadMoreResults() tea.Cmd {
	return func() tea.Msg {
		resp, err := api.SearchVideosWithPagination(m.opts.Input, m.continuationToken)
		if err != nil {
			return searchResultsMsg{err: err, isLoadMore: true}
		}
		return searchResultsMsg{
			results:           resp.Results,
			continuationToken: resp.ContinuationToken,
			hasMore:           resp.HasMore,
			err:               nil,
			isLoadMore:        true,
		}
	}
}

func (m Model) loadingView() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.spinner.View(),
			"Searching for: "+m.opts.Input,
			"\n\nPress q to quit",
		),
	)
}

func (m Model) listView() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.list.View(),
	)
}

func (m Model) detailView() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.createDetailView(),
			"\n",
			m.createDetailControls(),
		),
	)
}

func (m Model) errorView() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		ui.CreateErrorBox("Search Error", m.err.Error()),
	)
}

func (m Model) createDetailView() string {
	if m.selectedVideo == nil {
		return "No video selected"
	}
	return ui.CreateDetailedVideoView(*m.selectedVideo)
}

func (m Model) createDetailControls() string {
	controls := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Padding(1).
		Width(60).
		Align(lipgloss.Center)

	var controlsText []string
	controlsText = append(controlsText, lipgloss.NewStyle().Bold(true).Render("Controls"))
	controlsText = append(controlsText, "")

	if m.playerType != "" {
		controlsText = append(controlsText, fmt.Sprintf("[>]  [p/enter/space] Play with %s", m.playerType))
	} else {
		controlsText = append(controlsText, "[!] No video player found (install mpv)")
	}

	controlsText = append(controlsText, "[<]  [esc] Back to list")
	controlsText = append(controlsText, "[x] [q] Quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		controlsText...,
	)

	return controls.Render(content)
}

func (m *Model) playVideo() tea.Cmd {
	return func() tea.Msg {
		if m.selectedVideo == nil {
			return nil
		}

		playerType := player.PlayerType(m.playerType)
		err := player.StreamVideo(m.selectedVideo.URL, playerType, m.opts.Quality, m.opts.WindowMode)
		if err != nil {
			// Return error message
			return tea.Println(fmt.Sprintf("Failed to play video: %v", err))
		}
		return nil
	}
}
