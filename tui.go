package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	SETTINGSFILENAME string = "settings.json"
	SWAPFLAG                = "-sw"
)

// keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Inline(true)
)

// model holds the state
type model struct {
	choices      []string
	cursor       int
	selected     map[int]struct{}
	lastSelected *int
	active       string
}

// initialization of a new model
func mainModel(dirs []string, activeVersion string) model {
	return model{
		// choices:  []string{"Buy carrots", "Buy celery", "Do somthing else"},
		choices:      dirs,
		selected:     make(map[int]struct{}),
		lastSelected: nil,
		active:       activeVersion,
	}
}

// Init is used when we want to do IO, for now we dont need it so it returns nil
func (m model) Init() tea.Cmd {
	return nil
}

// this should be used to update the model > when we swap folders the list of choices needs to be refreshed
// this currently keeps the cursor position and selection! if the order of choices would change this would lead to wrong highlighting
func (m model) UpdateChoices() tea.Model {
	settings := GetCompleteSettings(SETTINGSFILENAME)
	tgkDir := settings.Defaults.Tgkdir
	tgkFolder := settings.Defaults.Tgkfolder
	dirs := GetDirsInDir(tgkDir)
	dirsWithOutTgkFolder := Remove(dirs, tgkFolder)
	m.choices = dirsWithOutTgkFolder
	m.lastSelected = nil
	m.selected = make(map[int]struct{})
	m.active = settings.ActiveSettings.OldDirectory
	return m
}

func (m model) swapFolders() error {
	var err error = nil
	// for now lets basically build the CLI args in here as string
	bogusProgramName := "prog"
	swapFlag := SWAPFLAG
	target := m.choices[m.cursor]
	args := []string{bogusProgramName, swapFlag, target}
	err = RunSwapper(args)
	if err != nil {
		return err
	}
	return nil
}

func killExcel(name string) error {
	return KillProcessByName(name)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "u":
			// diplay pop window to check if user really wants to proceed, if not then restart mainModel
			// TO DO: add the logic so this does not directly kill excel but informs the user first.
			// if any entry is selected, make the swapping.
			if len(m.selected) > 0 {
				m.swapFolders()
			}
			return m.UpdateChoices(), nil

			// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
				m.lastSelected = nil
			} else {
				m.selected[m.cursor] = struct{}{}
				// this makes it so that only every one entry in the model is selected.
				if m.lastSelected != nil {
					delete(m.selected, *m.lastSelected)
				}
				m.lastSelected = &m.cursor
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := helpStyle.Render("Please chose which version to swap in.") + "\n"
	s += helpStyle.Render("Currently active: ") + keywordStyle.Render(m.active) + "\n\n"
	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = cursorStyle.Render(">") // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = cursorStyle.Render("x") // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\n" + helpStyle.Render("Press q to quit.\t Press u to update.") + "\n"
	s = drawInBox(s)
	// Send the UI for rendering
	return s
}

func drawInBox(s string) string {
	// num of spaces to add between borders and content; vertical: in rows, horizontal: in spaces
	// verticalPaddingCount := 1
	horizontalPaddingCount := 4
	// verticalPadding := strings.Repeat("\n", verticalPaddingCount)
	horizontalPadding := strings.Repeat(" ", horizontalPaddingCount)
	// split string into slice to find longest row
	ss := strings.Split(s, "\n")
	// loop over all rows, add padding and find longest row
	maxLineLength := 0
	new_ss := []string{}
	for _, l := range ss {
		l = "\u2551" + horizontalPadding + l
		l = l + horizontalPadding + "\u2551" // add padding here so that the right border is always at the same position. use PadRight()
		lineLength := len(l)
		if lineLength > maxLineLength {
			maxLineLength = lineLength
		}
		new_ss = append(new_ss, l)
	}
	topLine := "\u2554" + strings.Repeat("\u2550", maxLineLength-2) + "\u2557"
	bottomLine := "\u255A" + strings.Repeat("\u2550", maxLineLength-2) + "\u255D"
	ret := topLine +
		"\n" +
		strings.Join(new_ss[:], "\n") +
		"\n" +
		bottomLine
	return ret
}

func runTui() {
	// at runtime, create default settings JSON if it doesnt exists.
	InitSettingsJSON()

	settings := GetCompleteSettings(SETTINGSFILENAME)
	tgkDir := settings.Defaults.Tgkdir
	tgkFolder := settings.Defaults.Tgkfolder
	dirs := GetDirsInDir(tgkDir)
	dirsWithOutTgkFolder := Remove(dirs, tgkFolder)
	activeVersion := settings.ActiveSettings.OldDirectory
	p := tea.NewProgram(mainModel(dirsWithOutTgkFolder, activeVersion))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}
}
