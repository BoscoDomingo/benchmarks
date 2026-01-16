package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dimStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	normalStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selectedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true)
	activeStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("215"))
	activeTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true)
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75")).Underline(true)
	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

// listSelector is a generic list selector with number input support.
type listSelector struct {
	title       string
	items       []string
	cursor      int
	numberInput string
	errorMsg    string
	selected    string
	quit        bool
}

func newListSelector(title string, items []string) listSelector {
	return listSelector{
		title: title,
		items: items,
	}
}

func (m listSelector) Init() tea.Cmd {
	return nil
}

func (m listSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMsg = ""

		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			m.numberInput = ""
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			m.numberInput = ""
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			if m.numberInput != "" {
				idx, err := strconv.Atoi(m.numberInput)
				if err != nil || idx < 1 || idx > len(m.items) {
					m.errorMsg = fmt.Sprintf("Invalid selection: %s (valid: 1-%d)", m.numberInput, len(m.items))
					m.numberInput = ""
					return m, nil
				}
				m.selected = m.items[idx-1]
				return m, tea.Quit
			}
			m.selected = m.items[m.cursor]
			return m, tea.Quit
		case "backspace":
			if len(m.numberInput) > 0 {
				m.numberInput = m.numberInput[:len(m.numberInput)-1]
			}
		case "esc":
			m.numberInput = ""
		default:
			if len(msg.String()) == 1 {
				ch := msg.String()[0]
				if ch >= '0' && ch <= '9' {
					m.numberInput += msg.String()
				}
			}
		}
	}
	return m, nil
}

func (m listSelector) View() string {
	var s string

	s += titleStyle.Render(m.title) + "\n\n"

	for i, item := range m.items {
		cursor := "  "
		style := normalStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedStyle
		}
		s += style.Render(fmt.Sprintf("%s%d. %s", cursor, i+1, item)) + "\n"
	}

	s += "\n"

	if m.numberInput != "" {
		s += normalStyle.Render(fmt.Sprintf("Selection: %s", m.numberInput)) + "\n"
	}

	if m.errorMsg != "" {
		s += errorStyle.Render(m.errorMsg) + "\n"
	}

	s += helpStyle.Render("↑/↓: navigate • [num]+enter: select • q: quit")

	return s
}

// selectLanguage prompts the user to select a language.
func selectLanguage(languages []string) (string, error) {
	model := newListSelector("Select a language", languages)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	m := finalModel.(listSelector)
	if m.quit {
		return "", fmt.Errorf("user quit")
	}

	return m.selected, nil
}

// nextActionSelector handles the post-benchmark action menu.
type nextActionSelector struct {
	cursor   int
	selected string
	quit     bool
}

var nextActionOptions = []struct {
	key    string
	label  string
	action string
}{
	{"r", "Run another benchmark (same language)", actionRunAnother},
	{"c", "Change language", actionChangeLang},
	{"q", "Exit", actionExit},
}

func (m nextActionSelector) Init() tea.Cmd {
	return nil
}

func (m nextActionSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(nextActionOptions)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = nextActionOptions[m.cursor].action
			return m, tea.Quit
		case "r":
			m.selected = actionRunAnother
			return m, tea.Quit
		case "c", "l":
			m.selected = actionChangeLang
			return m, tea.Quit
		case "q", "e":
			m.selected = actionExit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m nextActionSelector) View() string {
	var s string

	s += titleStyle.Render("What next?") + "\n\n"

	for i, opt := range nextActionOptions {
		cursor := "  "
		style := normalStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedStyle
		}
		s += style.Render(fmt.Sprintf("%s[%s] %s", cursor, opt.key, opt.label)) + "\n"
	}

	s += "\n"
	s += helpStyle.Render("↑/↓: navigate • enter: select • r/c/l/e/q: quick select")

	return s
}

// promptNextAction asks the user what to do after a benchmark completes.
func promptNextAction() (string, error) {
	model := nextActionSelector{}
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	m := finalModel.(nextActionSelector)
	if m.quit {
		return actionExit, nil
	}

	return m.selected, nil
}

// benchmarkSelector handles benchmark selection with preferences editing.
type benchmarkSelector struct {
	language    string
	benchmarks  []string
	prefs       *Preferences
	cursor      int
	editMode    bool   // true when editing preferences
	prefCursor  int    // which preference is being edited
	editBuffer  string
	numberInput string // buffer for typing benchmark number
	errorMsg    string // error message to display
	selected    string
	quit        bool
	err         error
}

func newBenchmarkSelector(language string, benchmarks []string, prefs *Preferences) benchmarkSelector {
	return benchmarkSelector{
		language:   language,
		benchmarks: benchmarks,
		prefs:      prefs,
	}
}

func (m benchmarkSelector) Init() tea.Cmd {
	return nil
}

func (m benchmarkSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editMode {
			return m.updateEditMode(msg)
		}
		return m.updateSelectMode(msg)
	}
	return m, nil
}

func (m benchmarkSelector) updateSelectMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.errorMsg = ""

	switch msg.String() {
	case "ctrl+c", "q":
		m.quit = true
		return m, tea.Quit
	case "up", "k":
		m.numberInput = ""
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		m.numberInput = ""
		if m.cursor < len(m.benchmarks)-1 {
			m.cursor++
		}
	case "enter":
		if m.numberInput != "" {
			idx, err := strconv.Atoi(m.numberInput)
			if err != nil || idx < 1 || idx > len(m.benchmarks) {
				m.errorMsg = fmt.Sprintf("Invalid selection: %s (valid: 1-%d)", m.numberInput, len(m.benchmarks))
				m.numberInput = ""
				return m, nil
			}
			m.selected = m.benchmarks[idx-1]
			return m, tea.Quit
		}
		m.selected = m.benchmarks[m.cursor]
		return m, tea.Quit
	case "p":
		m.editMode = true
		m.prefCursor = 0
		key := PreferenceKeys[m.prefCursor]
		m.editBuffer = m.prefs.Get(key)
	case "backspace":
		if len(m.numberInput) > 0 {
			m.numberInput = m.numberInput[:len(m.numberInput)-1]
		}
	case "esc":
		m.numberInput = ""
	default:
		if len(msg.String()) == 1 {
			ch := msg.String()[0]
			if ch >= '0' && ch <= '9' {
				m.numberInput += msg.String()
			}
		}
	}
	return m, nil
}

// isBoolPref returns true if the preference at the given key is a yes/no toggle.
func isBoolPref(key string) bool {
	def := PreferenceDefinitions[key]
	return def.DefaultValue == "yes" || def.DefaultValue == "no"
}

func (m benchmarkSelector) updateEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := PreferenceKeys[m.prefCursor]

	switch msg.String() {
	case "ctrl+c":
		m.quit = true
		return m, tea.Quit
	case "esc":
		m.editMode = false
	case "up", "k":
		if m.prefCursor > 0 {
			m.prefCursor--
			key := PreferenceKeys[m.prefCursor]
			m.editBuffer = m.prefs.Get(key)
		}
	case "down", "j":
		if m.prefCursor < len(PreferenceKeys)-1 {
			m.prefCursor++
			key := PreferenceKeys[m.prefCursor]
			m.editBuffer = m.prefs.Get(key)
		}
	case "enter":
		if m.editBuffer != "" {
			m.prefs.Set(key, m.editBuffer)
		}
		m.editMode = false
	case " ", "tab":
		// Toggle yes/no for boolean preferences
		if isBoolPref(key) {
			if m.editBuffer == "yes" {
				m.editBuffer = "no"
			} else {
				m.editBuffer = "yes"
			}
		}
	case "y":
		if isBoolPref(key) {
			m.editBuffer = "yes"
		}
	case "n":
		if isBoolPref(key) {
			m.editBuffer = "no"
		}
	case "backspace":
		if !isBoolPref(key) && len(m.editBuffer) > 0 {
			m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
		}
	default:
		if !isBoolPref(key) && len(msg.String()) == 1 {
			ch := msg.String()[0]
			if ch >= '0' && ch <= '9' {
				m.editBuffer += msg.String()
			}
		}
	}
	return m, nil
}

func (m benchmarkSelector) View() string {
	var s string

	s += titleStyle.Render(fmt.Sprintf("Select a benchmark (%s)", m.language)) + "\n\n"

	for i, bench := range m.benchmarks {
		cursor := "  "
		style := dimStyle
		if !m.editMode {
			if i == m.cursor {
				cursor = "> "
				style = selectedStyle
			} else {
				style = normalStyle
			}
		}
		s += style.Render(fmt.Sprintf("%s%d. %s", cursor, i+1, bench)) + "\n"
	}

	s += "\n"

	prefsTitle := "Preferences"
	if m.editMode {
		prefsTitle = activeTitleStyle.Render("Preferences (editing)")
	} else {
		prefsTitle = dimStyle.Render("Preferences")
	}
	s += prefsTitle + "\n"

	for i, key := range PreferenceKeys {
		def := PreferenceDefinitions[key]
		value := m.prefs.Get(key)

		cursor := "  "
		if m.editMode && i == m.prefCursor {
			cursor = "> "
			label := fmt.Sprintf("%s%s: %s█", cursor, def.DisplayName, m.editBuffer)
			s += activeStyle.Render(label)
			if def.Description != "" {
				s += dimStyle.Render(fmt.Sprintf("  (%s)", def.Description))
			}
		} else {
			label := fmt.Sprintf("%s%s: %s", cursor, def.DisplayName, value)
			s += dimStyle.Render(label)
			if def.Description != "" && m.editMode {
				s += dimStyle.Render(fmt.Sprintf("  (%s)", def.Description))
			}
		}
		s += "\n"
	}

	s += "\n"

	if !m.editMode && m.numberInput != "" {
		s += normalStyle.Render(fmt.Sprintf("Selection: %s", m.numberInput)) + "\n"
	}

	if m.errorMsg != "" {
		s += errorStyle.Render(m.errorMsg) + "\n"
	}

	if m.editMode {
		key := PreferenceKeys[m.prefCursor]
		if isBoolPref(key) {
			s += helpStyle.Render("↑/↓: navigate • space/y/n: toggle • enter: save • esc: cancel")
		} else {
			s += helpStyle.Render("↑/↓: navigate • enter: save • esc: cancel")
		}
	} else {
		s += helpStyle.Render("↑/↓: navigate • [num]+enter: select • p: preferences • q: quit")
	}

	return s
}

// selectBenchmark prompts the user to select a benchmark.
// Preferences are shown below and can be edited by pressing 'p'.
func selectBenchmark(language string, benchmarks []string, prefs *Preferences) (string, error) {
	model := newBenchmarkSelector(language, benchmarks, prefs)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	m := finalModel.(benchmarkSelector)
	if m.quit {
		return "", fmt.Errorf("user quit")
	}
	if m.err != nil {
		return "", m.err
	}

	return m.selected, nil
}
