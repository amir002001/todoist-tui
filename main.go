package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	todoistApiKey     = os.Getenv("API_TOKEN")
)

type item string

type model struct {
	list list.Model
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func main() {
	data := url.Values{}
	data.Set("sync_token", "*")
	data.Set("resource_types", "[\"all\"]")

	client := http.Client{}
	request, err := http.NewRequest(
		"POST",
		"https://api.todoist.com/sync/v9/sync",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(todoistApiKey)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistApiKey))

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	// items := []list.Item{
	// 	item("Ramen"),
	// 	item("Tomato Soup"),
	// 	item("Hamburgers"),
	// 	item("Cheeseburgers"),
	// 	item("Currywurst"),
	// 	item("Okonomiyaki"),
	// 	item("Pasta"),
	// 	item("Fillet Mignon"),
	// 	item("Caviar"),
	// 	item("Just Wine"),
	// }
	//
	// const defaultWidth = 20
	//
	// ls := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	// ls.Title = "What do you want for dinner?"
	// ls.SetShowStatusBar(false)
	// ls.SetFilteringEnabled(false)
	// ls.Styles.Title = titleStyle
	// ls.Styles.PaginationStyle = paginationStyle
	// ls.Styles.HelpStyle = helpStyle
	//
	// m := model{list: ls}
	// p := tea.NewProgram(m, tea.WithAltScreen())
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			// i, ok := m.list.SelectedItem().(item)
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(message)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
