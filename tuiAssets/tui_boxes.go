package tuiAssets

import (
	"sync"
)

type Box struct {
	TopLeftCorner     string
	TopRightCorner    string
	TopBar            string
	BottomLeftCorner  string
	BottomRightCorner string
	BottomBar         string
	LeftBar           string
	RightBar          string
}

func DoublePiped() Box {
	return Box{
		TopLeftCorner:     "\u2554",
		TopRightCorner:    "\u2557",
		BottomLeftCorner:  "\u255A",
		BottomRightCorner: "\u255D",
		TopBar:            "\u2550",
		BottomBar:         "\u2550",
		LeftBar:           "\u2551",
		RightBar:          "\u2551",
	}
}

func SinglePiped() Box {
	return Box{
		TopLeftCorner:     "\u250F",
		TopRightCorner:    "\u2513",
		BottomLeftCorner:  "\u2517",
		BottomRightCorner: "\u251B",
		TopBar:            "\u2501",
		BottomBar:         "\u2501",
		LeftBar:           "\u2503",
		RightBar:          "\u2503",
	}
}

func SingleRounded() Box {
	return Box{
		TopLeftCorner:     "\u256D",
		TopRightCorner:    "\u256E",
		BottomLeftCorner:  "\u2570",
		BottomRightCorner: "\u256F",
		TopBar:            "\u2500",
		BottomBar:         "\u2500",
		LeftBar:           "\u2502",
		RightBar:          "\u2502",
	}
}

var (
	onceBoxes      sync.Once
	boxIterator    *BoxIterator
	availableBoxes = []interface{}{SingleRounded, DoublePiped, SinglePiped}
)

type BoxIterator struct {
	boxes []interface{}
	index int
}

func newBoxIterator() *BoxIterator {
	return &BoxIterator{
		boxes: availableBoxes,
		index: 0,
	}
}

func GetBoxIterator() *BoxIterator {
	onceBoxes.Do(func() {
		boxIterator = newBoxIterator()
	})
	return boxIterator
}

func (bi *BoxIterator) Next() Box {
	bi.index = (bi.index + 1) % len(bi.boxes)
	boxConstructor := bi.boxes[bi.index]
	if f, ok := boxConstructor.(func() Box); ok {
		return f()
	}
	return Box{}
}

func (bi *BoxIterator) Previous() Box {
	bi.index = (bi.index - 1 + len(bi.boxes)) % len(bi.boxes)
	boxConstructor := bi.boxes[bi.index]
	if f, ok := boxConstructor.(func() Box); ok {
		return f()
	}
	return Box{}
}

func GetDefaultBox() Box {
	// currently this could just be an indexing into availableBoxes[0], but if we change the box constructors in the future to include more logic, we want to get the constructors and run them.
	if f, ok := newBoxIterator().boxes[0].(func() Box); ok {
		return f()
	}
	return Box{}
}
