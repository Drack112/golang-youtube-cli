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
		controlsText = append(controlsText, "[d]  Download video")
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
			return tea.Println(fmt.Sprintf("Failed to play video: %v", err))
		}
		return nil
	}
}
func (m *Model) startDownloadCmd(url, container, quality string, withThumb bool) tea.Cmd {
	return func() tea.Msg {
		err := player.DownloadVideo(url, container, quality, withThumb)
		return downloadResultMsg{err: err}
	}
}

func (m Model) renderDownloadModal() string {
	width := 48
	lines := []string{}

	if m.opts.QualityProvided {
		lines = append(lines, "Quality (from flag): "+m.opts.Quality)
	} else {
		q := m.downloadQualities[m.selectedQuality]
		if m.downloadCursor == 0 {
			q = "> " + q
		}
		lines = append(lines, "Quality: "+q)
	}

	c := m.downloadContainers[m.selectedContainer]
	if m.downloadCursor == 1 {
		c = "> " + c
	}
	lines = append(lines, "Container: "+c)

	thumb := "No"
	if m.downloadWithThumb {
		thumb = "Yes"
	}
	if m.downloadCursor == 2 {
		thumb = "> " + thumb
	}
	lines = append(lines, "Download thumbnail: "+thumb)

	// actions / status
	if m.downloadInProgress {
		lines = append(lines, "\nDownloading...")
		lines = append(lines, " "+m.spinner.View())
	} else if m.downloadMessage != "" {
		lines = append(lines, "\n"+m.downloadMessage)
		lines = append(lines, "Press Enter/Esc to close")
	} else {
		lines = append(lines, "\nPress Enter to start download | Esc to cancel")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	box := lipgloss.NewStyle().Width(width).Padding(1, 2).Border(lipgloss.RoundedBorder()).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
