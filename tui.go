package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"fastSwapper/tuiAssets"
)

const (
	SETTINGSFILENAME string = "settings.json"
	SWAPFLAG                = "-sw"
)

// keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
var (
	mainColor    = tuiAssets.GetDefaultColor()
	helpColor    = tuiAssets.GREY
	choicesColor = tuiAssets.WHITE
	choiceStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(choicesColor))
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(mainColor))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(mainColor))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(helpColor)) //.Inline(true)
	boxStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(mainColor))
	activeBox    = tuiAssets.GetDefaultBox()
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

		// toggle colors of main UI elements
		case "c":
			changeColors()
		case "b":
			changeBoxStyle()

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

		leftBracket := choiceStyle.Render("[")
		rightBracket := choiceStyle.Render("]")
		choice = choiceStyle.Render(choice)
		// Render the row
		s += fmt.Sprintf("%s %s%s%s %s\n", cursor, leftBracket, checked, rightBracket, choice)
	}

	// The footer
	s += "\n" + helpStyle.Render("q: quit\tu: swap")
	s += "\n" + helpStyle.Render("c: change colors\tb: change box")
	s = drawInBox(s, activeBox) + "\n"
	// Send the UI for rendering
	return s
}

func drawInBox(s string, b tuiAssets.Box) string {
	padRune := ' '
	numRunesLeftBar := utf8.RuneCountInString(StripANSI(b.LeftBar))
	numRunesRightBar := utf8.RuneCountInString(StripANSI(b.RightBar))
	// num of spaces to add between borders and content; vertical: in rows, horizontal: in spaces
	verticalPaddingCount := 1 // vertical padding currently isnt clean as the box characters at the sides are missing!
	horizontalPaddingCount := 4
	horizontalPadding := strings.Repeat(" ", horizontalPaddingCount)
	// split string into slice to find longest row
	ss := strings.Split(s, "\n")
	// loop over all rows, add padding and find longest row
	maxLineLength := 0
	for _, l := range ss {
		// get hypothetical len to not change the line just now >> needs padding based on max len
		// len of line string + 2 for box chars + twice the horizontalPadding
		lineLength := len(StripANSI(l)) + numRunesLeftBar + numRunesRightBar + horizontalPaddingCount*2
		if lineLength > maxLineLength {
			maxLineLength = lineLength
		}
	}
	// need a second loop in order to apply padding based on maxLineLength!
	// we strip ANSI escape sequences from the line because those interfere with the padding.
	padded_ss := []string{}
	for _, l := range ss {
		l = boxStyle.Render(b.LeftBar) + horizontalPadding + l
		// the padding does not work cleanly if we swap out LeftBar for p.e. "x" instead of "\u2551"
		l = l + PadRight("", padRune, maxLineLength-len(StripANSI(l))-horizontalPaddingCount+numRunesLeftBar) + horizontalPadding + boxStyle.Render(b.RightBar)
		padded_ss = append(padded_ss, l)
	}
	topLine := boxStyle.Render(b.TopLeftCorner) + strings.Repeat(boxStyle.Render(b.TopBar), maxLineLength-numRunesLeftBar-numRunesRightBar) + boxStyle.Render(b.TopRightCorner)
	bottomLine := boxStyle.Render(b.BottomLeftCorner) + strings.Repeat(boxStyle.Render(b.BottomBar), maxLineLength-numRunesLeftBar-numRunesRightBar) + boxStyle.Render(b.BottomRightCorner)
	verticalPadding := strings.Repeat(boxStyle.Render(b.LeftBar)+strings.Repeat(" ", maxLineLength-numRunesLeftBar-numRunesRightBar)+boxStyle.Render(b.RightBar)+"\n", verticalPaddingCount)
	ret := topLine + "\n" +
		verticalPadding +
		strings.Join(padded_ss[:], "\n") + "\n" +
		verticalPadding +
		bottomLine
	return ret
}

func changeColors() {
	newColor := tuiAssets.GetColorIterator().Next()
	updateTextStyleColor(&keywordStyle, newColor)
	updateTextStyleColor(&boxStyle, newColor)
	updateTextStyleColor(&cursorStyle, newColor)
}

func changeBoxStyle() {
	activeBox = tuiAssets.GetBoxIterator().Next()
}

func updateTextStyleColor(toUpdate *lipgloss.Style, newColor string) {
	*toUpdate = toUpdate.Foreground(lipgloss.Color(newColor))
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
