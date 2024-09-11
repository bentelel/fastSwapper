package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const SETTINGSFILENAME string = "settings.json"

// model holds the state
type model struct {
	choices      []string
	cursor       int
	selected     map[int]struct{}
	lastSelected *int
}

// initialization of a new model
func initialModel(dirs []string) model {
	return model{
		// choices:  []string{"Buy carrots", "Buy celery", "Do somthing else"},
		choices:      dirs,
		selected:     make(map[int]struct{}),
		lastSelected: nil,
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
	return m
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
	s := "Please chose which version to swap in.\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\t Press u to update.\n"

	// Send the UI for rendering
	return s
}

func runTui() {
	settings := GetCompleteSettings(SETTINGSFILENAME)
	tgkDir := settings.Defaults.Tgkdir
	tgkFolder := settings.Defaults.Tgkfolder
	dirs := GetDirsInDir(tgkDir)
	dirsWithOutTgkFolder := Remove(dirs, tgkFolder)
	p := tea.NewProgram(initialModel(dirsWithOutTgkFolder))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}
}
