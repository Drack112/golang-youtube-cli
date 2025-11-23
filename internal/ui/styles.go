package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryPurple   = lipgloss.Color("#8B5CF6")
	SecondaryPurple = lipgloss.Color("#7C3AED")
	AccentBlue      = lipgloss.Color("#3B82F6")
	AccentGreen     = lipgloss.Color("#10B981")
	Background      = lipgloss.Color("#1E1B2E")
	Surface         = lipgloss.Color("#2A2540")
	TextPrimary     = lipgloss.Color("#F8FAFC")
	TextSecondary   = lipgloss.Color("#94A3B8")
	TextMuted       = lipgloss.Color("#64748B")
	Warning         = lipgloss.Color("#F59E0B")
	Error           = lipgloss.Color("#EF4444")
	Success         = lipgloss.Color("#22C55E")
)

var (
	ContainerStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderBottomForeground(PrimaryPurple).Background(Background).Foreground(TextPrimary)
	TitleStyle     = lipgloss.NewStyle().Foreground(TextPrimary).Bold(true).MarginBottom(1)
	MainTitleStyle = lipgloss.NewStyle().Foreground(PrimaryPurple).Bold(true).MarginBottom(1).Padding(0, 1).Background(Surface).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(AccentBlue)

	NormalTextStyle  = lipgloss.NewStyle().Foreground(TextPrimary)
	MutedTextStyle   = lipgloss.NewStyle().Foreground(TextSecondary)
	AccentTextStyle  = lipgloss.NewStyle().Foreground(AccentBlue)
	SuccessTextStyle = lipgloss.NewStyle().Foreground(Warning)
	WarningTextStyle = lipgloss.NewStyle().Foreground(Warning)

	LiveBadgeStyle     = lipgloss.NewStyle().Foreground(TextPrimary).Background(Error).Bold(true).Padding(0, 1).MarginRight(1)
	ShortBadgeStyle    = lipgloss.NewStyle().Foreground(TextPrimary).Background(AccentBlue).Bold(true).Padding(0, 1).MarginRight(1)
	DurationBadgeStyle = lipgloss.NewStyle().Foreground(TextPrimary).Background(AccentGreen).Bold(true).Padding(0, 1).MarginRight(1)
	ChannelBadgeStyle  = lipgloss.NewStyle().Foreground(TextPrimary).Background(SecondaryPurple).Bold(true).Padding(0, 1).MarginRight(1)

	SectionStyle   = lipgloss.NewStyle().MarginTop(1).PaddingLeft(1).Border(lipgloss.NormalBorder(), false, false, false, true)
	MetadataStyle  = lipgloss.NewStyle().Foreground(TextMuted).Italic(true)
	SeparatorStyle = lipgloss.NewStyle().Foreground(Surface)
)

func FormatDuration(seconds int) string {
	if seconds == 0 {
		return "LIVE"
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}

func CreateBadge(text string, style lipgloss.Style) string {
	return style.Render(text)
}
