package tuiAssets

import "sync"

const (
	ORANGE    string = "#ed832d"
	RED              = "#fc0303"
	ROSE             = "#f562c1"
	HOTPINK          = "#f5029f"
	LILAC            = "#cf02f7"
	PURPLE           = "#720787"
	LIGHTBLUE        = "#03b1fc"
	DARKBLUE         = "#020fc4"
	YELLOW           = "#f7e40a"
	LIME             = "#bcf70a"
	GREEN            = "#0be04b"
	DARKGREEN        = "#005e1c"
	WHITE            = "#ffffff"
	GREY             = "#807d7d"
)

// first color here is the default color.
var (
	availableColors = []string{ORANGE, RED, ROSE, HOTPINK, LILAC, PURPLE, LIGHTBLUE, DARKBLUE, YELLOW, LIME, GREEN, DARKGREEN, WHITE, GREY}
	onceColors      sync.Once
	colorIterator   *ColorIterator
)

type ColorIterator struct {
	colors []string
	index  int
}

func newColorIterator() *ColorIterator {
	return &ColorIterator{
		colors: availableColors,
		index:  0,
	}
}

func GetColorIterator() *ColorIterator {
	// Initialize the colorIterator only once
	onceColors.Do(func() {
		colorIterator = newColorIterator()
	})
	return colorIterator
}

func GetDefaultColor() string {
	return newColorIterator().colors[0]
}

func (ci *ColorIterator) Next() string {
	ci.index = (ci.index + 1) % len(ci.colors)
	color := ci.colors[ci.index]
	return color
}

func (ci *ColorIterator) Previous() string {
	ci.index = (ci.index - 1 + len(ci.colors)) % len(ci.colors)
	color := ci.colors[ci.index]
	return color
}
