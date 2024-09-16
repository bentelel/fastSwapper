package main

import "github.com/charmbracelet/lipgloss"

type box struct {
	topLeftCorner     string
	topRightCorner    string
	topBar            string
	bottomLeftCorner  string
	bottomRightCorner string
	bottomBar         string
	leftBar           string
	rightBar          string
}

func DoublePiped(style lipgloss.Style) box {
	return box{
		topLeftCorner:     style.Render("\u2554"),
		topRightCorner:    style.Render("\u2557"),
		bottomLeftCorner:  style.Render("\u255A"),
		bottomRightCorner: style.Render("\u255D"),
		topBar:            style.Render("\u2550"),
		bottomBar:         style.Render("\u2550"),
		leftBar:           style.Render("\u2551"),
		rightBar:          style.Render("\u2551"),
	}
}

func SinglePiped(style lipgloss.Style) box {
	return box{
		topLeftCorner:     style.Render("\u250F"),
		topRightCorner:    style.Render("\u2513"),
		bottomLeftCorner:  style.Render("\u2517"),
		bottomRightCorner: style.Render("\u251B"),
		topBar:            style.Render("\u2501"),
		bottomBar:         style.Render("\u2501"),
		leftBar:           style.Render("\u2503"),
		rightBar:          style.Render("\u2503"),
	}
}

func SingleRounded(style lipgloss.Style) box {
	return box{
		topLeftCorner:     style.Render("\u256D"),
		topRightCorner:    style.Render("\u256E"),
		bottomLeftCorner:  style.Render("\u2570"),
		bottomRightCorner: style.Render("\u256F"),
		topBar:            style.Render("\u2500"),
		bottomBar:         style.Render("\u2500"),
		leftBar:           style.Render("\u2502"),
		rightBar:          style.Render("\u2502"),
	}
}
