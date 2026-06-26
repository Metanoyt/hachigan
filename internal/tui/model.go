package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hachigan/hachigan/internal/orchestrator"
	"github.com/hachigan/hachigan/internal/tui/screens"
	"github.com/hachigan/hachigan/internal/tui/state"
)

type Model struct {
	orchestrator orchestrator.Orchestrator
	refreshEvery time.Duration

	screen state.Screen
	cursor int
	err    error

	overview    orchestrator.OverviewView
	apps        orchestrator.ApplicationsView
	appDetail   orchestrator.ApplicationDetailView
	cluster     orchestrator.ClusterInventoryView
	initialized bool
}

type dataLoadedMsg struct {
	overview orchestrator.OverviewView
	apps     orchestrator.ApplicationsView
	cluster  orchestrator.ClusterInventoryView
}

type detailLoadedMsg struct {
	detail orchestrator.ApplicationDetailView
}

type errMsg struct{ err error }
type refreshMsg struct{}

func New(orchestrator orchestrator.Orchestrator, refreshEvery time.Duration) Model {
	return Model{orchestrator: orchestrator, refreshEvery: refreshEvery, screen: state.ScreenOverview}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadData, m.tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.screen = state.ScreenOverview
		case "2":
			m.screen = state.ScreenApplications
		case "3":
			m.screen = state.ScreenCluster
		case "up", "k":
			if m.screen == state.ScreenApplications && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.screen == state.ScreenApplications && m.cursor < len(m.apps.Applications)-1 {
				m.cursor++
			}
		case "enter":
			if m.screen == state.ScreenApplications && len(m.apps.Applications) > 0 {
				app := m.apps.Applications[m.cursor]
				m.screen = state.ScreenApplicationDetail
				return m, m.loadDetail(app.Namespace, app.Name)
			}
		case "esc", "backspace":
			if m.screen == state.ScreenApplicationDetail {
				m.screen = state.ScreenApplications
			}
		case "r":
			return m, m.loadData
		}
	case dataLoadedMsg:
		m.overview = msg.overview
		m.apps = msg.apps
		m.cluster = msg.cluster
		m.initialized = true
		m.err = nil
		if m.cursor >= len(m.apps.Applications) {
			m.cursor = max(0, len(m.apps.Applications)-1)
		}
	case detailLoadedMsg:
		m.appDetail = msg.detail
		m.err = nil
	case errMsg:
		m.err = msg.err
		m.initialized = true
	case refreshMsg:
		return m, tea.Batch(m.loadData, m.tick())
	}
	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return screens.Error(m.err)
	}
	if !m.initialized {
		return screens.Loading()
	}
	switch m.screen {
	case state.ScreenApplications:
		return screens.Applications(m.apps, m.cursor)
	case state.ScreenApplicationDetail:
		return screens.ApplicationDetail(m.appDetail)
	case state.ScreenCluster:
		return screens.Cluster(m.cluster)
	default:
		return screens.Overview(m.overview)
	}
}

func (m Model) loadData() tea.Msg {
	ctx := context.Background()
	overview, err := m.orchestrator.Overview(ctx)
	if err != nil {
		return errMsg{err: err}
	}
	apps, err := m.orchestrator.Applications(ctx)
	if err != nil {
		return errMsg{err: err}
	}
	cluster, err := m.orchestrator.ClusterInventory(ctx)
	if err != nil {
		return errMsg{err: err}
	}
	return dataLoadedMsg{overview: overview, apps: apps, cluster: cluster}
}

func (m Model) loadDetail(namespace, name string) tea.Cmd {
	return func() tea.Msg {
		detail, err := m.orchestrator.ApplicationDetail(context.Background(), namespace, name)
		if err != nil {
			return errMsg{err: err}
		}
		return detailLoadedMsg{detail: detail}
	}
}

func (m Model) tick() tea.Cmd {
	if m.refreshEvery <= 0 {
		return nil
	}
	return tea.Tick(m.refreshEvery, func(time.Time) tea.Msg { return refreshMsg{} })
}
