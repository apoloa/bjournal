package ui

import (
	"fmt"
	"github.com/apoloa/bjournal/model"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type IndexList struct {
	*tview.Box

	index *model.Index

	// The index of the currently selected item.
	currentItem int

	// The item main text style.
	mainTextStyle tcell.Style

	// The item secondary text style.
	secondaryTextStyle tcell.Style

	// The item shortcut text style.
	shortcutStyle tcell.Style

	// The style for selected items.
	selectedStyle tcell.Style

	selectedStyleSubLog tcell.Style

	// If true, the selection is only shown when the list has focus.
	selectedFocusOnly bool

	// If true, the entire row is highlighted when selected.
	highlightFullLine bool

	// Whether or not navigating the list will wrap around.
	wrapAround bool

	// The number of list items skipped at the top before the first item is
	// drawn.
	itemOffset int

	// The number of cells skipped on the left side of an item text. Shortcuts
	// are not affected.
	horizontalOffset int

	// Set to true if a currently visible item flows over the right border of
	// the box. This is set by the Draw() function. It determines the behaviour
	// of the right arrow key.
	overflowing bool

	// An optional function which is called when the user has navigated to a list
	// item.
	changed func(index int, log model.IndexItem)

	// An optional function which is called when a list item was selected. This
	// function will be called even if the list item defines its own callback.
	selected func(index int, log model.IndexItem)

	// An optional function which is called when the user presses the Escape key.
	done func()
}

// NewIndexList returns a new form.
func NewIndexList() *IndexList {
	return &IndexList{
		Box:                tview.NewBox(),
		wrapAround:         false,
		highlightFullLine:  true,
		currentItem:        -1,
		mainTextStyle:      tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor),
		secondaryTextStyle: tcell.StyleDefault.Foreground(tview.Styles.TertiaryTextColor),
		shortcutStyle:      tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor),
		selectedStyle: tcell.StyleDefault.Foreground(tview.Styles.PrimitiveBackgroundColor).
			Background(tview.Styles.PrimaryTextColor),
		selectedStyleSubLog: tcell.StyleDefault.Foreground(tview.Styles.PrimitiveBackgroundColor).
			Background(tview.Styles.PrimaryTextColor),
	}
}

func (i *IndexList) AddIndexModel(index *model.Index) *IndexList {
	i.index = index
	return i
}

func (i *IndexList) GetCurrentItem() *model.IndexItem {
	if i.currentItem == -1 || i.currentItem > len(i.index.Items) {
		return nil
	}
	return &i.index.Items[i.currentItem]
}

// Draw draws this primitive onto the screen.
func (i *IndexList) Draw(screen tcell.Screen) {
	i.Box.DrawForSubclass(screen, i)

	// Determine the dimensions.
	x, y, width, height := i.GetInnerRect()
	bottomLimit := y + height
	_, totalHeight := screen.Size()
	if bottomLimit > totalHeight {
		bottomLimit = totalHeight
	}

	x += 4
	width -= 4

	// Adjust offset to keep the current selection in view.
	if i.currentItem < i.itemOffset {
		i.itemOffset = i.currentItem
	} else if 2*(i.currentItem-i.itemOffset) >= height-1 {
		i.itemOffset = (2*i.currentItem + 3 - height) / 2
	} else {
		if i.currentItem-i.itemOffset >= height {
			i.itemOffset = i.currentItem + 1 - height
		}
	}
	if i.horizontalOffset < 0 {
		i.horizontalOffset = 0
	}

	// Draw the list items.
	var (
		maxWidth    int  // The maximum printed item width.
		overflowing bool // Whether a text's end exceeds the right border.
	)

	for index, item := range i.index.Items {
		if index < i.itemOffset {
			continue
		}

		if y >= bottomLimit {
			break
		}

		// Shortcuts.
		printWithStyle(screen, fmt.Sprint(" - "), x-5, y, 0, 4, AlignRight, i.mainTextStyle, true)

		// Main text.
		for _, wordWrap := range WordWrap(item.Name, width) {
			printWithStyle(screen, wordWrap, x, y, i.horizontalOffset, width, AlignLeft, i.mainTextStyle, true)
			// Background color of selected text.
			if index == i.currentItem && (!i.selectedFocusOnly || i.HasFocus()) {
				textWidth := width
				if !i.highlightFullLine {
					if w := TaggedStringWidth(item.Name); w < textWidth {
						textWidth = w
					}
				}

				mainTextColor, _, _ := i.mainTextStyle.Decompose()
				for bx := 0; bx < textWidth; bx++ {
					m, c, style, _ := screen.GetContent(x+bx, y)
					fg, _, _ := style.Decompose()
					style = i.selectedStyle
					if fg != mainTextColor {
						style = style.Foreground(fg)
					}
					screen.SetContent(x+bx, y, m, c, style)
				}
			}
			y++
		}

		if y >= bottomLimit {
			break
		}

	}

	// We don't want the item text to get out of view. If the horizontal offset
	// is too high, we reset it and redraw. (That should be about as efficient
	// as calculating everything up front.)
	if i.horizontalOffset > 0 && maxWidth < width {
		i.horizontalOffset -= width - maxWidth
		i.Draw(screen)
	}
	i.overflowing = overflowing
}

// InputHandler returns the handler for this primitive.
func (l *IndexList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return l.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		previousItem := l.currentItem

		switch key := event.Key(); key {
		case tcell.KeyTab, tcell.KeyDown:
			l.currentItem++
		case tcell.KeyBacktab, tcell.KeyUp:
			l.currentItem--
		case tcell.KeyRight:
			if l.overflowing {
				l.horizontalOffset += 2 // We shift by 2 to account for two-cell characters.
			} else {
				l.currentItem++
			}
		case tcell.KeyLeft:
			if l.horizontalOffset > 0 {
				l.horizontalOffset -= 2
			} else {
				l.currentItem--
			}
		case tcell.KeyHome:
			l.currentItem = 0
		case tcell.KeyEnd:
			l.currentItem = len(l.index.Items) - 1
		case tcell.KeyPgDn:
			_, _, _, height := l.GetInnerRect()
			l.currentItem += height
			if l.currentItem >= len(l.index.Items) {
				l.currentItem = len(l.index.Items) - 1
			}
		case tcell.KeyPgUp:
			_, _, _, height := l.GetInnerRect()
			l.currentItem -= height
			if l.currentItem < 0 {
				l.currentItem = 0
			}
		case tcell.KeyEnter:
			if l.currentItem >= 0 && l.currentItem < len(l.index.Items) {
				item := l.index.Items[l.currentItem]
				if l.selected != nil {
					l.selected(l.currentItem, item)
				}
			}
		}

		if l.currentItem < 0 {
			if l.wrapAround {
				l.currentItem = len(l.index.Items) - 1
			} else {
				l.currentItem = -1
			}
		} else if l.currentItem >= len(l.index.Items) {
			if l.wrapAround {
				l.currentItem = 0
			} else {
				l.currentItem = -1
			}
		}

		if l.currentItem != previousItem && l.currentItem < len(l.index.Items) && l.changed != nil {
			item := l.index.Items[l.currentItem]
			l.changed(l.currentItem, item)
		}
	})
}
