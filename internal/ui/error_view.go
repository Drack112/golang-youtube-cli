package ui

import "github.com/charmbracelet/lipgloss"

func CreateErrorBox(title, message string) string {
	errorStyle := ContainerStyle.Copy().BorderForeground(Error).Background(lipgloss.Color("#2A1A2E"))
	titleStyle := TitleStyle.Copy().Foreground(Error)

	content := titleStyle.Render("[ERROR] " + title + "\n\n" + NormalTextStyle.Render(message))
	return errorStyle.Render(content)
}
