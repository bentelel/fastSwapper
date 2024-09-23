package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"fastSwapper/tuiAssets"
	"fastSwapper/utils"
)

var (
	mainColor          = tuiAssets.GetDefaultColor()
	footerColor        = tuiAssets.GREY
	headerColor        = tuiAssets.GREY
	choicesColor       = tuiAssets.WHITE
	choiceStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(choicesColor))
	keywordStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(mainColor))
	cursorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(mainColor))
	footerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(footerColor))
	headerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(headerColor))
	boxStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color(mainColor))
	activeBox          = tuiAssets.GetDefaultBox()
	footerItems        = []string{"q: quit", "u: swap", "c: change colors", "b: change box"}
	numFooterRows      = 2
	cursorSymbol       = ">"
	checkmarkSymbol    = "x"
	leftbracketSymbol  = "["
	rightbracketSymbol = "]"
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
	dirsWithOutTgkFolder := DirectoriesInTgkDirExcludingTgkFolder()
	m.choices = dirsWithOutTgkFolder
	m.lastSelected = nil
	m.selected = make(map[int]struct{})
	m.active = GetActiveVersion()
	return m
}

func (m model) swapFolders() error {
	var err error = nil
	target := m.choices[m.cursor]
	err = SwapDirectories(target)
	if err != nil {
		return err
	}
	return nil
}

func killExcel(name string) error {
	return utils.KillProcessByName(name)
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

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "c":
			changeColors()
		case "b":
			changeBoxStyle()

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
	s := headerStyle.Render("Please chose which version to swap in.") + "\n"
	s += headerStyle.Render("Currently active: ") + keywordStyle.Render(m.active) + "\n\n"
	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = cursorStyle.Render(cursorSymbol) // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = cursorStyle.Render(checkmarkSymbol) // selected!
		}

		leftBracket := choiceStyle.Render(leftbracketSymbol)
		rightBracket := choiceStyle.Render(rightbracketSymbol)
		choice = choiceStyle.Render(choice)
		// Render the row
		s += fmt.Sprintf("%s %s%s%s %s\n", cursor, leftBracket, checked, rightBracket, choice)
	}

	// The footer
	s += "\n" + drawInGrid(footerItems, numFooterRows)
	s = drawInBox(s, activeBox) + "\n"
	// Send the UI for rendering
	return s
}

func drawInGrid(items []string, numRows int) string {
	padding := " "
	// split items into n slices, where n is number of rows
	// use ceiling division to ensure we always have enough room for all items
	itemsPerRow := (len(items) + numRows - 1) / numRows
	rows := make([][]string, numRows)
	for i := range rows {
		rows[i] = make([]string, itemsPerRow)
	}
	counter := 0
	maxLength := 0
	// Split items into rows and items per row
	// since we are already iterating over the slice, we find the maximum length among items at the same time
	for idx, item := range items {
		rows[counter][idx%itemsPerRow] = item
		if idx%itemsPerRow == itemsPerRow-1 {
			counter += 1
		}
		if len(item) > maxLength {
			maxLength = len(item)
		}
	}
	s := ""
	for rowNum, row := range rows {
		for idx, item := range row {
			if idx != itemsPerRow-1 {
				s += item + strings.Repeat(padding, maxLength-len(item)) + strings.Repeat(padding, 2)
			} else {
				s += item + strings.Repeat(padding, maxLength-len(item))
			}
		}
		if rowNum != numRows-1 {
			s += "\n"
		}
	}
	return s
}

func drawInBox(s string, b tuiAssets.Box) string {
	padRune := ' '
	numRunesLeftBar := utf8.RuneCountInString(b.LeftBar)
	numRunesRightBar := utf8.RuneCountInString(b.RightBar)
	// num of spaces to add between borders and content; vertical: in rows, horizontal: in spaces
	verticalPaddingCount := 1 // vertical padding currently isnt clean as the box characters at the sides are missing!
	horizontalPaddingCount := 4
	horizontalPadding := strings.Repeat(string(padRune), horizontalPaddingCount)
	// split string into slice to find longest row
	ss := strings.Split(s, "\n")
	// loop over all rows, add padding and find longest row
	maxLineLength := 0
	for _, l := range ss {
		// get hypothetical len to not change the line just now >> needs padding based on max len
		// len of line string + 2 for box chars + twice the horizontalPadding
		lineLength := len(utils.StripANSI(l)) + numRunesLeftBar + numRunesRightBar + horizontalPaddingCount*2
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
		l = l + utils.PadRight("", padRune, maxLineLength-len(utils.StripANSI(l))-horizontalPaddingCount+numRunesLeftBar) + horizontalPadding + boxStyle.Render(b.RightBar)
		padded_ss = append(padded_ss, l)
	}
	topLine := boxStyle.Render(b.TopLeftCorner)
	topLine += strings.Repeat(boxStyle.Render(b.TopBar), maxLineLength-numRunesLeftBar-numRunesRightBar)
	topLine += boxStyle.Render(b.TopRightCorner)
	bottomLine := boxStyle.Render(b.BottomLeftCorner)
	bottomLine += strings.Repeat(boxStyle.Render(b.BottomBar), maxLineLength-numRunesLeftBar-numRunesRightBar)
	bottomLine += boxStyle.Render(b.BottomRightCorner)
	verticalPad := boxStyle.Render(b.LeftBar) + strings.Repeat(" ", maxLineLength-numRunesLeftBar-numRunesRightBar) + boxStyle.Render(b.RightBar) + "\n"
	verticalPadding := strings.Repeat(verticalPad, verticalPaddingCount)
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
	dirsWithOutTgkFolder := DirectoriesInTgkDirExcludingTgkFolder()
	activeVersion := GetActiveVersion()
	p := tea.NewProgram(mainModel(dirsWithOutTgkFolder, activeVersion))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}
}
