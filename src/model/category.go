package model

import (
	"github.com/gdamore/tcell/v2"
)

type Category string

const (
	Task       Category = "task"
	Complete   Category = "complete"
	Irrelevant Category = "irrelevant"
	Migrated   Category = "migrated"
	Scheduled  Category = "scheduled"
	Note       Category = "note"
	Event      Category = "event"
)

func (c Category) Print() rune {
	switch {
	case c == Task:
		return '•'
	case c == Complete:
		return '✘'
	case c == Irrelevant:
		return ' '
	case c == Migrated:
		return '>'
	case c == Scheduled:
		return '<'
	case c == Note:
		return '-'
	case c == Event:
		return '○'
	}
	return ' '
}

func (c Category) Color() tcell.Color {
	switch {
	case c == Task:
		return tcell.ColorCadetBlue
	case c == Complete:
		return tcell.ColorYellowGreen
	case c == Irrelevant:
		return tcell.ColorYellow
	case c == Migrated:
		return tcell.ColorOrangeRed
	case c == Scheduled:
		return tcell.ColorYellow
	case c == Note:
		return tcell.ColorHotPink
	case c == Event:
		return tcell.ColorRebeccaPurple
	}
	return tcell.ColorBlack
}

func (c Category) Style() tcell.Style {
	return tcell.StyleDefault.Foreground(c.Color())
}
