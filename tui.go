package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getModel() (model, error) {
	tabs := []string{"Issues", "Details"}

	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "P", Width: 3},
		{Title: "Subject", Width: 20},
		{Title: "S", Width: 3},
		{Title: "Assignee", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	m := model{
		Tabs:  tabs,
		Table: t,
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	m.Table.SetStyles(s)

	err := m.setRows()
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	if m.activeTab == 0 {
		t, cmd := m.Table.Update(msg)
		m.Table = t
		issueIdString := m.Table.SelectedRow()[0]
		activeIssueId, err := strconv.Atoi(issueIdString)
		if err == nil {
			m.setActiveIssue(activeIssueId)
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				config, err := getConfig()
				if err != nil {
					return m, cmd
				}
				url := fmt.Sprintf("%s/issues/%s", strings.TrimSuffix(config.RedmineURL, "/"), issueIdString)
				openUrl(url)
			}
		}

		return m, cmd
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	var tabContent string
	if m.activeTab == 0 {
		tabContent = m.Table.View()
	} else {
		tabContent = m.getDetailsTabContent()
	}

	windowWidth := (lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()) + 60
	// @todo: fix alignment
	doc.WriteString(windowStyle.Width(windowWidth).Render(tabContent))

	return docStyle.Render(doc.String())
}

func (m model) getDetailsTabContent() string {
	heading := func(s string) string {
		style := lipgloss.NewStyle().
			Bold(true).
			Align(0)

		return style.Render(s + ": ")
	}
	text := func(s string) string {
		style := lipgloss.NewStyle().
			Italic(true)

		return style.Render(s + "\n")
	}

	content := ""
	if m.ActiveIssue != nil {
		content = heading("Project") + text(m.ActiveIssue.ProjectCode)
		content = content + heading("ID") + text("#"+strconv.Itoa(m.ActiveIssue.Id))
		content = content + heading("Subject") + text(m.ActiveIssue.Subject)
		content = content + heading("Assignee") + text(m.ActiveIssue.Assignee)
	}

	return content
}

func (m *model) setRows() error {
	is, err := getRedmineService()
	if err != nil {
		return err
	}

	m.Issues, err = is.GetIssuesByQueryId(is.Config.QueryId)
	if err != nil {
		return err
	}

	var rows []table.Row
	for i := 0; i < len(m.Issues); i++ {
		row := table.Row{
			strconv.Itoa(m.Issues[i].Id),
			m.Issues[i].ProjectCode,
			m.Issues[i].Subject,
			m.Issues[i].Status.getStatusString(),
			m.Issues[i].Assignee,
		}

		rows = append(rows, row)
	}

	m.Table.SetRows(rows)

	return nil
}

func (m *model) setActiveIssue(issueId int) {
	for _, issue := range m.Issues {
		if issue.Id == issueId {
			m.ActiveIssue = issue
			return
		}
	}
}

func (status status) getStatusString() string {
	switch status {
	case New:
		return "🆕"
	case Solved:
		return "✅"
	case OnHold:
		return "✋🏻"
	case InProgress:
		return "👷"
	case Closed:
		return "🔒"
	}

	return "❔"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
