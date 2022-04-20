package ui

import (
	"bjournal/model"
	"fmt"
	"github.com/derailed/tview"

	"github.com/gdamore/tcell/v2"
)

// List displays rows of items, each of which can be selected.
//
// See https://github.com/rivo/tview/wiki/List for an example.
type List struct {
	*tview.Box

	// The items of the list.
	items []*model.Log

	daily *model.DailyLog

	// The index of the currently selected item.
	currentItem int

	secondaryIndex int

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
	changed func(index int, log *model.Log)

	// An optional function which is called when a list item was selected. This
	// function will be called even if the list item defines its own callback.
	selected func(index int, log *model.Log)

	// An optional function which is called when the user presses the Escape key.
	done func()
}

// NewList returns a new form.
func NewList() *List {
	return &List{
		Box:                tview.NewBox(),
		wrapAround:         false,
		highlightFullLine:  true,
		currentItem:        -1,
		secondaryIndex:     -1,
		mainTextStyle:      tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor),
		secondaryTextStyle: tcell.StyleDefault.Foreground(tview.Styles.TertiaryTextColor),
		shortcutStyle:      tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor),
		selectedStyle: tcell.StyleDefault.Foreground(tview.Styles.PrimitiveBackgroundColor).
			Background(tview.Styles.PrimaryTextColor),
		selectedStyleSubLog: tcell.StyleDefault.Foreground(tview.Styles.PrimitiveBackgroundColor).
			Background(tview.Styles.PrimaryTextColor),
	}
}

// SetCurrentItem sets the currently selected item by its index, starting at 0
// for the first item. If a negative index is provided, items are referred to
// from the back (-1 = last item, -2 = second-to-last item, and so on). Out of
// range indices are clamped to the beginning/end.
//
// Calling this function triggers a "changed" event if the selection changes.
func (l *List) SetCurrentItem(index int) *List {
	if index < 0 {
		index = len(l.items) + index
	}
	if index >= len(l.items) {
		index = len(l.items) - 1
	}
	if index < 0 {
		index = -1
	}

	if index != l.currentItem && l.changed != nil {
		item := l.items[index]
		l.changed(index, item)
	}

	l.currentItem = index

	return l
}

// GetCurrentItem returns the index of the currently selected list item,
// starting at 0 for the first item.
func (l *List) GetCurrentItem() int {
	return l.currentItem
}

// SetOffset sets the number of items to be skipped (vertically) as well as the
// number of cells skipped horizontally when the list is drawn. Note that one
// item corresponds to two rows when there are secondary texts. Shortcuts are
// always drawn.
//
// These values may change when the list is drawn to ensure the currently
// selected item is visible and item texts move out of view. Users can also
// modify these values by interacting with the list.
func (l *List) SetOffset(items, horizontal int) *List {
	l.itemOffset = items
	l.horizontalOffset = horizontal
	return l
}

// GetOffset returns the number of items skipped while drawing, as well as the
// number of cells item text is moved to the left. See also SetOffset() for more
// information on these values.
func (l *List) GetOffset() (int, int) {
	return l.itemOffset, l.horizontalOffset
}

// RemoveItem removes the item with the given index (starting at 0) from the
// list. If a negative index is provided, items are referred to from the back
// (-1 = last item, -2 = second-to-last item, and so on). Out of range indices
// are clamped to the beginning/end, i.e. unless the list is empty, an item is
// always removed.
//
// The currently selected item is shifted accordingly. If it is the one that is
// removed, a "changed" event is fired.
func (l *List) RemoveItem(index int) *List {
	if len(l.items) == 0 {
		return l
	}

	// Adjust index.
	if index < 0 {
		index = len(l.items) + index
	}
	if index >= len(l.items) {
		index = len(l.items) - 1
	}
	if index < 0 {
		index = -1
	}

	// Remove item.
	l.items = append(l.items[:index], l.items[index+1:]...)

	// If there is nothing left, we're done.
	if len(l.items) == 0 {
		return l
	}

	// Shift current item.
	previousCurrentItem := l.currentItem
	if l.currentItem >= index {
		l.currentItem--
	}

	// Fire "changed" event for removed items.
	if previousCurrentItem == index && l.changed != nil {
		item := l.items[l.currentItem]
		l.changed(l.currentItem, item)
	}

	return l
}

// SetMainTextColor sets the color of the items' main text.
func (l *List) SetMainTextColor(color tcell.Color) *List {
	l.mainTextStyle = l.mainTextStyle.Foreground(color)
	return l
}

// SetMainTextStyle sets the style of the items' main text. Note that the
// background color is ignored in order not to override the background color of
// the list itself.
func (l *List) SetMainTextStyle(style tcell.Style) *List {
	l.mainTextStyle = style
	return l
}

// SetSecondaryTextColor sets the color of the items' secondary text.
func (l *List) SetSecondaryTextColor(color tcell.Color) *List {
	l.secondaryTextStyle = l.secondaryTextStyle.Foreground(color)
	return l
}

// SetSecondaryTextStyle sets the style of the items' secondary text. Note that
// the background color is ignored in order not to override the background color
// of the list itself.
func (l *List) SetSecondaryTextStyle(style tcell.Style) *List {
	l.secondaryTextStyle = style
	return l
}

// SetShortcutColor sets the color of the items' shortcut.
func (l *List) SetShortcutColor(color tcell.Color) *List {
	l.shortcutStyle = l.shortcutStyle.Foreground(color)
	return l
}

// SetShortcutStyle sets the style of the items' shortcut. Note that the
// background color is ignored in order not to override the background color of
// the list itself.
func (l *List) SetShortcutStyle(style tcell.Style) *List {
	l.shortcutStyle = style
	return l
}

// SetSelectedTextColor sets the text color of selected items. Note that the
// color of main text characters that are different from the main text color
// (e.g. color tags) is maintained.
func (l *List) SetSelectedTextColor(color tcell.Color) *List {
	l.selectedStyle = l.selectedStyle.Foreground(color)
	return l
}

// SetSelectedBackgroundColor sets the background color of selected items.
func (l *List) SetSelectedBackgroundColor(color tcell.Color) *List {
	l.selectedStyle = l.selectedStyle.Background(color)
	return l
}

// SetSelectedStyle sets the style of the selected items. Note that the color of
// main text characters that are different from the main text color (e.g. color
// tags) is maintained.
func (l *List) SetSelectedStyle(style tcell.Style) *List {
	l.selectedStyle = style
	return l
}

// SetSelectedFocusOnly sets a flag which determines when the currently selected
// list item is highlighted. If set to true, selected items are only highlighted
// when the list has focus. If set to false, they are always highlighted.
func (l *List) SetSelectedFocusOnly(focusOnly bool) *List {
	l.selectedFocusOnly = focusOnly
	return l
}

// SetHighlightFullLine sets a flag which determines whether the colored
// background of selected items spans the entire width of the view. If set to
// true, the highlight spans the entire view. If set to false, only the text of
// the selected item from beginning to end is highlighted.
func (l *List) SetHighlightFullLine(highlight bool) *List {
	l.highlightFullLine = highlight
	return l
}

// SetWrapAround sets the flag that determines whether navigating the list will
// wrap around. That is, navigating downwards on the last item will move the
// selection to the first item (similarly in the other direction). If set to
// false, the selection won't change when navigating downwards on the last item
// or navigating upwards on the first item.
func (l *List) SetWrapAround(wrapAround bool) *List {
	l.wrapAround = wrapAround
	return l
}

// SetChangedFunc sets the function which is called when the user navigates to
// a list item. The function receives the item's index in the list of items
// (starting with 0), its main text, secondary text, and its shortcut rune.
//
// This function is also called when the first item is added or when
// SetCurrentItem() is called.
func (l *List) SetChangedFunc(handler func(index int, log *model.Log)) *List {
	l.changed = handler
	return l
}

// SetSelectedFunc sets the function which is called when the user selects a
// list item by pressing Enter on the current selection. The function receives
// the item's index in the list of items (starting with 0), its main text,
// secondary text, and its shortcut rune.
func (l *List) SetSelectedFunc(handler func(int, *model.Log)) *List {
	l.selected = handler
	return l
}

// SetDoneFunc sets a function which is called when the user presses the Escape
// key.
func (l *List) SetDoneFunc(handler func()) *List {
	l.done = handler
	return l
}

func (l *List) AddDailyLog(dailyLog *model.DailyLog) *List {
	for index, _ := range dailyLog.Logs {
		l.InsertItem(-1, &dailyLog.Logs[index], nil)
	}
	l.daily = dailyLog
	return l
}

// AddItem calls InsertItem() with an index of -1.
func (l *List) AddItem(log *model.Log, selected func()) *List {
	l.InsertItem(-1, log, selected)
	return l
}

// InsertItem adds a new item to the list at the specified index. An index of 0
// will insert the item at the beginning, an index of 1 before the second item,
// and so on. An index of GetItemCount() or higher will insert the item at the
// end of the list. Negative indices are also allowed: An index of -1 will
// insert the item at the end of the list, an index of -2 before the last item,
// and so on. An index of -GetItemCount()-1 or lower will insert the item at the
// beginning.
//
// An item has a main text which will be highlighted when selected. It also has
// a secondary text which is shown underneath the main text (if it is set to
// visible) but which may remain empty.
//
// The shortcut is a key binding. If the specified rune is entered, the item
// is selected immediately. Set to 0 for no binding.
//
// The "selected" callback will be invoked when the user selects the item. You
// may provide nil if no such callback is needed or if all events are handled
// through the selected callback set with SetSelectedFunc().
//
// The currently selected item will shift its position accordingly. If the list
// was previously empty, a "changed" event is fired because the new item becomes
// selected.
func (l *List) InsertItem(index int, item *model.Log, selected func()) *List {
	// Shift index to range.
	if index < 0 {
		index = len(l.items) + index + 1
	}
	if index < 0 {
		index = 0
	} else if index > len(l.items) {
		index = len(l.items)
	}

	// Shift current item.
	if l.currentItem < len(l.items) && l.currentItem >= index {
		l.currentItem++
	}

	// Insert item (make space for the new item, then shift and insert).
	l.items = append(l.items, nil)
	if index < len(l.items)-1 { // -1 because l.items has already grown by one item.
		copy(l.items[index+1:], l.items[index:])
	}
	l.items[index] = item

	return l
}

// GetItemCount returns the number of items in the list.
func (l *List) GetItemCount() int {
	return len(l.items)
}

// GetItem returns an item's. Panics if the index
// is out of range.
func (l *List) GetItem(index int) *model.Log {
	if index < 0 {
		return nil
	} else if index > len(l.items) {
		return nil
	}
	return l.items[index]
}

func (l *List) GetDaily() *model.DailyLog {
	return l.daily
}

func (l *List) GetCurrentLog() *model.Log {
	if l.currentItem < 0 {
		return nil
	} else if l.currentItem > len(l.items) {
		return nil
	}

	item := l.items[l.currentItem]

	if item.SubLogs == nil {
		return item
	}

	if l.secondaryIndex < 0 {
		return item
	} else if l.secondaryIndex > len(*item.SubLogs) {
		return item
	}

	return &(*item.SubLogs)[l.secondaryIndex]
}

// SetItemText sets an item's main and secondary text. Panics if the index is
// out of range.
func (l *List) SetItemText(index int, log *model.Log) *List {
	l.items[index] = log
	return l
}

// Clear removes all items from the list.
func (l *List) Clear() *List {
	l.items = nil
	l.currentItem = 0
	return l
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

// Draw draws this primitive onto the screen.
func (l *List) Draw(screen tcell.Screen) {
	l.Box.DrawForSubclass(screen, l)

	// Determine the dimensions.
	x, y, width, height := l.GetInnerRect()
	bottomLimit := y + height
	_, totalHeight := screen.Size()
	if bottomLimit > totalHeight {
		bottomLimit = totalHeight
	}

	x += 4
	width -= 4

	// Adjust offset to keep the current selection in view.
	if l.currentItem < l.itemOffset {
		l.itemOffset = l.currentItem
	} else if 2*(l.currentItem-l.itemOffset) >= height-1 {
		l.itemOffset = (2*l.currentItem + 3 - height) / 2
	} else {
		if l.currentItem-l.itemOffset >= height {
			l.itemOffset = l.currentItem + 1 - height
		}
	}
	if l.horizontalOffset < 0 {
		l.horizontalOffset = 0
	}

	// Draw the list items.
	var (
		maxWidth    int  // The maximum printed item width.
		overflowing bool // Whether a text's end exceeds the right border.
	)

	for index, item := range l.items {
		if index < l.itemOffset {
			continue
		}

		if y >= bottomLimit {
			break
		}

		// Shortcuts.
		printWithStyle(screen, fmt.Sprintf("(%s)", string(item.Mark.Print())), x-5, y, 0, 4, AlignRight, item.Mark.Style(), true)

		// Main text.
		for _, wordWrap := range WordWrap(item.Name, width) {
			printWithStyle(screen, wordWrap, x, y, l.horizontalOffset, width, AlignLeft, l.mainTextStyle, true)
			// Background color of selected text.
			if index == l.currentItem && (!l.selectedFocusOnly || l.HasFocus()) {
				textWidth := width
				if !l.highlightFullLine {
					if w := TaggedStringWidth(item.Name); w < textWidth {
						textWidth = w
					}
				}

				mainTextColor, _, _ := l.mainTextStyle.Decompose()
				for bx := 0; bx < textWidth; bx++ {
					m, c, style, _ := screen.GetContent(x+bx, y)
					fg, _, _ := style.Decompose()
					style = l.selectedStyle
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

		if item.SubLogs == nil {
			continue
		}

		for subIndex, sublog := range *item.SubLogs {
			printWithStyle(screen, fmt.Sprintf("(%s)", string(sublog.Mark.Print())), x-2, y, 0, 4, AlignRight, sublog.Mark.Style(), true)
			for _, wordWrap := range WordWrap(sublog.Name, width) {
				printWithStyle(screen, wordWrap, x+3, y, l.horizontalOffset, width, AlignLeft, l.mainTextStyle, true)
				// Background color of selected text.
				if (index == l.currentItem && (!l.selectedFocusOnly || l.HasFocus())) && subIndex == l.secondaryIndex {
					textWidth := width
					if !l.highlightFullLine {
						if w := TaggedStringWidth(item.Name); w < textWidth {
							textWidth = w
						}
					}

					mainTextColor, _, _ := l.mainTextStyle.Decompose()
					for bx := 0; bx < textWidth; bx++ {
						m, c, style, _ := screen.GetContent(x+3+bx, y)
						fg, _, _ := style.Decompose()
						style = l.selectedStyle
						if fg != mainTextColor {
							style = style.Foreground(fg)
						}
						screen.SetContent(x+3+bx, y, m, c, style)
					}
				}
				y++
			}
		}
	}

	// We don't want the item text to get out of view. If the horizontal offset
	// is too high, we reset it and redraw. (That should be about as efficient
	// as calculating everything up front.)
	if l.horizontalOffset > 0 && maxWidth < width {
		l.horizontalOffset -= width - maxWidth
		l.Draw(screen)
	}
	l.overflowing = overflowing
}

func (l *List) increaseIndex() {
	if l.currentItem == len(l.items) {
		l.currentItem = -1
		return
	}
	if l.currentItem >= 0 && l.currentItem < len(l.items) {
		if l.GetItem(l.currentItem).SubLogs != nil {
			sublogs := l.GetItem(l.currentItem).SubLogs
			if len(*sublogs) > 0 && l.secondaryIndex < len(*sublogs)-1 {
				l.secondaryIndex++
			} else {
				l.secondaryIndex = -1
				if l.currentItem > len(l.items)-1 {
					l.currentItem = -1
				} else {
					l.currentItem++
				}

			}

		} else {
			if l.currentItem > len(l.items)-1 {
				l.currentItem = -1
				l.secondaryIndex = -1
			} else {
				l.currentItem++
				l.secondaryIndex = -1
			}

		}
	} else {
		l.currentItem++
		l.secondaryIndex = -1
	}
}

func (l *List) decreaseIndex() {
	if l.currentItem == -1 {
		l.currentItem = len(l.items) - 1
		if l.GetItem(l.currentItem).SubLogs != nil {
			l.secondaryIndex = len(*(l.GetItem(l.currentItem)).SubLogs) - 1
		}
	} else if l.GetItem(l.currentItem).SubLogs != nil {
		sublogs := l.GetItem(l.currentItem).SubLogs
		if len(*sublogs) > 0 && l.secondaryIndex < len(*sublogs) {
			if l.secondaryIndex == -1 {
				l.currentItem--
				if l.currentItem == -1 {
					return
				}
				if l.GetItem(l.currentItem).SubLogs != nil {
					l.secondaryIndex = len(*(l.GetItem(l.currentItem)).SubLogs) - 1
				}
				return
			}
			l.secondaryIndex--
		}
	} else {
		l.currentItem--
		if l.currentItem == -1 {
			return
		}
		if l.GetItem(l.currentItem).SubLogs != nil {
			l.secondaryIndex = len(*(l.GetItem(l.currentItem)).SubLogs) - 1
		}
	}
}

// InputHandler returns the handler for this primitive.
func (l *List) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return l.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyEscape {
			if l.done != nil {
				l.done()
			}
			return
		} else if len(l.items) == 0 {
			return
		}

		previousItem := l.currentItem

		switch key := event.Key(); key {
		case tcell.KeyTab, tcell.KeyDown:
			l.increaseIndex()
		case tcell.KeyBacktab, tcell.KeyUp:
			l.decreaseIndex()
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
			l.currentItem = len(l.items) - 1
		case tcell.KeyPgDn:
			_, _, _, height := l.GetInnerRect()
			l.currentItem += height
			if l.currentItem >= len(l.items) {
				l.currentItem = len(l.items) - 1
			}
		case tcell.KeyPgUp:
			_, _, _, height := l.GetInnerRect()
			l.currentItem -= height
			if l.currentItem < 0 {
				l.currentItem = 0
			}
		case tcell.KeyEnter:
			if l.currentItem >= 0 && l.currentItem < len(l.items) {
				item := l.items[l.currentItem]
				if l.selected != nil {
					l.selected(l.currentItem, item)
				}
			}
		}

		if l.currentItem < 0 {
			if l.wrapAround {
				l.currentItem = len(l.items) - 1
			} else {
				l.currentItem = -1
			}
		} else if l.currentItem >= len(l.items) {
			if l.wrapAround {
				l.currentItem = 0
			} else {
				l.currentItem = -1
			}
		}

		if l.currentItem != previousItem && l.currentItem < len(l.items) && l.changed != nil {
			item := l.items[l.currentItem]
			l.changed(l.currentItem, item)
		}
	})
}

// indexAtPoint returns the index of the list item found at the given position
// or a negative value if there is no such list item.
func (l *List) indexAtPoint(x, y int) int {
	rectX, rectY, width, height := l.GetInnerRect()
	if rectX < 0 || rectX >= rectX+width || y < rectY || y >= rectY+height {
		return -1
	}

	index := y - rectY

	index /= 2

	index += l.itemOffset

	if index >= len(l.items) {
		return -1
	}
	return index
}

// MouseHandler returns the mouse handler for this primitive.
func (l *List) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return l.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !l.InRect(event.Position()) {
			return false, nil
		}

		// Process mouse event.
		switch action {
		case tview.MouseLeftClick:
			setFocus(l)
			index := l.indexAtPoint(event.Position())
			if index != -1 {
				item := l.items[index]
				if l.selected != nil {
					l.selected(index, item)
				}
				if index != l.currentItem && l.changed != nil {
					l.changed(index, item)
				}
				l.currentItem = index
			}
			consumed = true
		case tview.MouseScrollUp:
			if l.itemOffset > 0 {
				l.itemOffset--
			}
			consumed = true
		case tview.MouseScrollDown:
			lines := len(l.items) - l.itemOffset
			lines *= 2
			if _, _, _, height := l.GetInnerRect(); lines > height {
				l.itemOffset++
			}
			consumed = true
		}

		return
	})
}
