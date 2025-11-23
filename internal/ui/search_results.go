package ui

import (
	"fmt"
	"strings"

	"github.com/Drack112/go-youtube/internal/models"
	"github.com/charmbracelet/lipgloss"
)

func CreateSearchResultsView(results []models.SearchResult) string {
	var content strings.Builder

	content.WriteString(MainTitleStyle.Render(fmt.Sprintf("Search results (%d found)", len(results))))
	content.WriteString("\n\n")

	for i, result := range results {
		content.WriteString(createResultItem(result, i+1))
		if i < len(results)-1 {
			content.WriteString("\n")
			content.WriteString(createSeparator())
			content.WriteString("\n")
		}
	}

	return ContainerStyle.Render(content.String())
}

func createResultItem(result models.SearchResult, index int) string {
	var item strings.Builder

	indexStyle := lipgloss.NewStyle().Foreground(PrimaryPurple).Bold(true).Width(3).Align(lipgloss.Right)
	titleStyle := NormalTextStyle.Copy().Bold(true).MaxWidth(60)

	item.WriteString(indexStyle.Render(fmt.Sprintf("%d.", index)))
	item.WriteString(" ")
	item.WriteString(titleStyle.Render(TruncateText(result.Title, 58)))

	var metaParts []string

	if result.IsLive {
		metaParts = append(metaParts, LiveBadgeStyle.Render("LIVE"))
	}

	if result.IsShort {
		metaParts = append(metaParts, ShortBadgeStyle.Render("SHORT"))
	}

	if result.Duration != "" && !result.IsLive {
		metaParts = append(metaParts, DurationBadgeStyle.Render(FormatDuration(result.DurationSec)))
	}

	if result.ChannelName != "" {
		channelText := ChannelBadgeStyle.Render(TruncateText(result.ChannelName, 20))
		metaParts = append(metaParts, channelText)
	}

	if len(metaParts) > 0 {
		metaLine := lipgloss.NewStyle().MarginLeft(4).Render(strings.Join(metaParts, " "))
		item.WriteString(metaLine)
		item.WriteString("\n")
	}

	urlStyle := MutedTextStyle.Copy().Italic(true).MarginLeft(4)
	item.WriteString(urlStyle.Render(TruncateText(result.URL, 63)))

	return item.String()
}

func createSeparator() string {
	return SeparatorStyle.MarginLeft(4).Render("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈")
}
