package tui

import (
	"fmt"

	"github.com/Drack112/go-youtube/internal/api"
	"github.com/Drack112/go-youtube/internal/flags"
	"github.com/Drack112/go-youtube/internal/models"
	"github.com/Drack112/go-youtube/internal/player"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateLoading state = iota
	stateList
	stateDetail
	stateError
)

type Model struct {
	state             state
	opts              *flags.Options
	results           []models.SearchResult
	selectedVideo     *models.SearchResult
	err               error
	continuationToken string
	hasMore           bool
	isLoadingMore     bool
	playerType        string

	list     list.Model
	viewport viewport.Model
	spinner  spinner.Model

	width  int
	height int
}

type searchResultsMsg struct {
	results           []models.SearchResult
	continuationToken string
	hasMore           bool
	err               error
	isLoadMore        bool
}

func NewModel(opts *flags.Options) Model {
	detectedPlayer, err := player.DetectAvailablePlayer()
	playerTypeStr := ""
	if err == nil {
		playerTypeStr = string(detectedPlayer)
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B5CF6"))

	// Initialize list with custom delegate
	delegate := NewItemDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "[?] Search Results"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B5CF6")).
		Bold(true).
		Padding(0, 1)

	// Add help text
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "view details"),
			),
			key.NewBinding(
				key.WithKeys("m"),
				key.WithHelp("m", "load more"),
			),
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "view details"),
			),
			key.NewBinding(
				key.WithKeys("m"),
				key.WithHelp("m", "load more results"),
			),
		}
	}

	return Model{
		state:      stateLoading,
		opts:       opts,
		spinner:    s,
		list:       l,
		playerType: playerTypeStr,
	}
}

func NewProgram(model Model) *tea.Program {
	return tea.NewProgram(model, tea.WithAltScreen())
}

func (m Model) Init() tea.Cmd {
	if m.opts.InputKind == flags.InputYoutubeURL {
		return tea.Batch(
			m.spinner.Tick,
			m.fetchVideoDetails(),
		)
	}
	return tea.Batch(
		m.spinner.Tick,
		m.performSearch(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-4)
		m.viewport = viewport.New(msg.Width-4, msg.Height-4)
		m.viewport.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state == stateList || m.state == stateError {
				return m, tea.Quit
			}
		case "esc":
			switch m.state {
			case stateDetail:
				m.state = stateList
				m.selectedVideo = nil
				return m, nil
			case stateError:
				return m, tea.Quit
			}
		}

	case searchResultsMsg:
		m.isLoadingMore = false
		if msg.err != nil {
			m.state = stateError
			m.err = msg.err
			return m, nil
		}

		if msg.isLoadMore {
			m.results = append(m.results, msg.results...)
		} else {
			m.results = msg.results
		}

		m.continuationToken = msg.continuationToken
		m.hasMore = msg.hasMore
		m.state = stateList

		items := make([]list.Item, len(m.results))
		for i, result := range m.results {
			items[i] = item{result: result}
		}
		m.list.SetItems(items)

		if m.hasMore {
			m.list.Title = fmt.Sprintf("[?] Search Results (%d results - Press 'm' for more)", len(m.results))
		} else {
			m.list.Title = fmt.Sprintf("[?] Search Results (%d results)", len(m.results))
		}
	}

	var cmd tea.Cmd
	switch m.state {
	case stateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
	case stateList:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				if selected, ok := m.list.SelectedItem().(item); ok {
					m.selectedVideo = &selected.result
					m.state = stateDetail
					m.viewport.SetContent(m.createDetailView())
					return m, nil
				}
			case "m", "M":
				if m.hasMore && !m.isLoadingMore {
					m.isLoadingMore = true
					return m, m.loadMoreResults()
				}
			}
		}
		m.list, cmd = m.list.Update(msg)
	case stateDetail:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "p", "P", "enter", " ":
				if m.selectedVideo != nil && m.playerType != "" {
					return m, tea.Sequence(
						tea.ExitAltScreen,
						m.playVideo(),
						tea.EnterAltScreen,
					)
				}
			}
		}
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.loadingView()
	case stateList:
		return m.listView()
	case stateDetail:
		return m.detailView()
	case stateError:
		return m.errorView()
	default:
		return "Unknown state"
	}
}

func (m *Model) fetchVideoDetails() tea.Cmd {
	return func() tea.Msg {
		results, err := api.SearchVideos(m.opts.Input)
		if err != nil {
			return searchResultsMsg{err: err}
		}

		if len(results) == 0 {
			return searchResultsMsg{err: fmt.Errorf("no video found for URL: %s", m.opts.Input)}
		}

		return searchResultsMsg{
			results: results,
			hasMore: false,
		}
	}
}
