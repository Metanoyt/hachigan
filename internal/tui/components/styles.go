package components

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	MutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	SelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62"))
	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	OKStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	WarnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	CritStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)
