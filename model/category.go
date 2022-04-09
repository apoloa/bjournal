package model

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
