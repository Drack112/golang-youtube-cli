package ui

import (
	"strconv"
	"strings"

	"github.com/Drack112/go-youtube/internal/models"
	"github.com/Drack112/go-youtube/pkg/utils"
	"github.com/charmbracelet/lipgloss"
)

func CreateDetailedVideoView(video models.SearchResult) string {
	var content strings.Builder

	content.WriteString(createVideoHeader(video))
	content.WriteString("\n\n")

	content.WriteString(createQuickInfo(video))
	content.WriteString("\n\n")

	content.WriteString(createMetadataSection(video))
	content.WriteString("\n\n")

	if video.ChannelName != "" {
		content.WriteString(createChannelSection(video))
		content.WriteString("\n\n")
	}

	content.WriteString(createURLSection(video))

	return ContainerStyle.Render(content.String())
}

func createVideoHeader(video models.SearchResult) string {
	var header strings.Builder

	title := MainTitleStyle.
		Width(70).
		Render(utils.TruncateText(video.Title, 66))
	header.WriteString(title)
	header.WriteString("\n")

	var badges []string

	if video.IsLive {
		badges = append(badges, LiveBadgeStyle.Render("LIVE"))
	}
	if video.IsShort {
		badges = append(badges, ShortBadgeStyle.Render("SHORT"))
	}
	if video.Duration != "" && !video.IsLive {
		durationText := utils.FormatDuration(video.DurationSec)
		badges = append(badges, DurationBadgeStyle.Render(durationText))
	}

	if len(badges) > 0 {
		badgeLine := lipgloss.NewStyle().MarginTop(1).Render(strings.Join(badges, " "))
		header.WriteString(badgeLine)
	}

	return header.String()
}

func createQuickInfo(video models.SearchResult) string {
	var info strings.Builder

	info.WriteString(TitleStyle.Render("[i] Video Info"))
	info.WriteString("\n\n")

	info.WriteString(NormalTextStyle.Render("  ID: "))
	info.WriteString(AccentTextStyle.Render(video.ID))
	info.WriteString("\n")

	if !video.IsLive && !video.IsShort {
		info.WriteString(NormalTextStyle.Render("  Duration: "))
		info.WriteString(SuccessTextStyle.Render(utils.FormatDuration(video.DurationSec)))
		info.WriteString("\n")
	}

	info.WriteString(NormalTextStyle.Render("  Type: "))
	if video.IsLive {
		info.WriteString(WarningTextStyle.Render("[LIVE] Stream"))
	} else if video.IsShort {
		info.WriteString(AccentTextStyle.Render("[SHORT] Video"))
	} else {
		info.WriteString(SuccessTextStyle.Render("[VIDEO] Regular"))
	}

	return SectionStyle.Render(info.String())
}

func createMetadataSection(video models.SearchResult) string {
	var metadata strings.Builder

	metadata.WriteString(TitleStyle.Render("[+] Video Details"))
	metadata.WriteString("\n")

	if video.DurationSec > 0 && !video.IsLive {
		metadata.WriteString(NormalTextStyle.Render("  Length: "))
		metadata.WriteString(SuccessTextStyle.Render(utils.FormatDuration(video.DurationSec)))
		metadata.WriteString(MutedTextStyle.Render(" (" + strconv.Itoa(video.DurationSec) + " seconds)"))
		metadata.WriteString("\n")
	}

	metadata.WriteString(NormalTextStyle.Render("  Status: "))
	if video.IsLive {
		metadata.WriteString(WarningTextStyle.Render("[!] Live Now"))
	} else if video.IsShort {
		metadata.WriteString(AccentTextStyle.Render("[*] Short Video"))
	} else {
		metadata.WriteString(SuccessTextStyle.Render("[ok] Available"))
	}

	return SectionStyle.Render(metadata.String())
}

func createChannelSection(video models.SearchResult) string {
	var channel strings.Builder

	channel.WriteString(TitleStyle.Render("[@] Channel"))
	channel.WriteString("\n")

	channel.WriteString(NormalTextStyle.Render("  Name: "))
	channel.WriteString(AccentTextStyle.Render(video.ChannelName))
	channel.WriteString("\n")

	if video.ChannelID != "" {
		channel.WriteString(NormalTextStyle.Render("  ID: "))
		channel.WriteString(MutedTextStyle.Render(video.ChannelID))
		channel.WriteString("\n")
	}

	return SectionStyle.Render(channel.String())
}

func createURLSection(video models.SearchResult) string {
	urlStyle := lipgloss.NewStyle().
		Foreground(AccentBlue).
		Italic(true)

	return SectionStyle.Render(
		TitleStyle.Render("[>] URL") + "\n" +
			urlStyle.Render(video.URL),
	)
}
