package ui

import (
	"context"
	"fmt"
	"github.com/apoloa/bjournal/src/model"
	"sync"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

const (
	defaultPrompt = "%c > [::b]%s"
	defaultSpacer = 4
)

// Suggester provides suggestions.
type Suggester interface {
	// CurrentSuggestion returns the current suggestion.
	CurrentSuggestion() (string, bool)

	// NextSuggestion returns the next suggestion.
	NextSuggestion() (string, bool)

	// PrevSuggestion returns the prev suggestion.
	PrevSuggestion() (string, bool)

	// ClearSuggestions clear out all suggestions.
	ClearSuggestions()
}

// PromptModel represents a prompt buffer.
type PromptModel interface {
	// SetText sets the model text.
	SetText(txt string)

	// GetText returns the current text.
	GetText() string

	// ClearText clears out model text.
	ClearText(fire bool)

	// AddListener registers a command listener.
	AddListener(model.BuffWatcher)

	// RemoveListener removes a listener.
	RemoveListener(model.BuffWatcher)

	// IsActive returns true if prompt is active.
	IsActive() bool

	// SetActive sets whether the prompt is active or not.
	SetActive(bool)

	// Add adds a new char to the prompt.
	Add(rune)

	// Delete deletes the last prompt character.
	Delete()
}

type CmdBuff struct {
	buff       []rune
	suggestion string
	hotKey     rune
	active     bool
	cancel     context.CancelFunc
	mx         sync.RWMutex
}

// Prompt captures users free from command input.
type Prompt struct {
	*tview.TextView

	noIcons bool
	icon    rune
	model   PromptModel
	spacer  int
}

// NewPrompt returns a new command view.
func NewPrompt(noIcons bool) *Prompt {
	p := Prompt{
		noIcons:  noIcons,
		TextView: tview.NewTextView(),
		spacer:   defaultSpacer,
	}
	if noIcons {
		p.spacer--
	}
	p.SetWordWrap(true)
	p.SetWrap(true)
	p.SetDynamicColors(true)
	p.SetBorder(true)
	p.SetBorderPadding(0, 0, 1, 1)
	p.SetBackgroundColor(tcell.ColorBlack)
	p.SetTextColor(tcell.ColorWhite)
	p.SetInputCapture(p.keyboard)
	return &p
}

// SendKey sends a keyboard event (testing only!).
func (p *Prompt) SendKey(evt *tcell.EventKey) {
	p.keyboard(evt)
}

// SendStrokes (testing only!)
func (p *Prompt) SendStrokes(s string) {
	for _, r := range s {
		p.keyboard(tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone))
	}
}

// SetModel sets the prompt buffer model.
func (p *Prompt) SetModel(m PromptModel) {
	if p.model != nil {
		p.model.RemoveListener(p)
	}
	p.model = m
	p.model.AddListener(p)
}

func (p *Prompt) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	// nolint:exhaustive
	switch evt.Key() {
	case tcell.KeyBackspace2, tcell.KeyBackspace, tcell.KeyDelete:
		p.model.Delete()
	case tcell.KeyRune:
		p.model.Add(evt.Rune())
	case tcell.KeyEscape:
		p.model.ClearText(true)
		p.model.SetActive(false)
	case tcell.KeyEnter, tcell.KeyCtrlE:
		p.model.SetText(p.model.GetText())
		p.model.SetActive(false)
	case tcell.KeyCtrlW, tcell.KeyCtrlU:
		p.model.ClearText(true)
	}

	return nil
}

// StylesChanged notifies skin changed.
/*
func (p *Prompt) StylesChanged(s *config.Styles) {
	p.styles = s
	p.SetBackgroundColor(s.K9s.Prompt.BgColor.Color())
	p.SetTextColor(s.K9s.Prompt.FgColor.Color())
}*/

// InCmdMode returns true if command is active, false otherwise.
func (p *Prompt) InCmdMode() bool {
	if p.model == nil {
		return false
	}
	return p.model.IsActive()
}

func (p *Prompt) activate() {
	p.SetCursorIndex(len(p.model.GetText()))
	p.write(p.model.GetText())
}

func (p *Prompt) update(text string) {
	p.Clear()
	p.write(text)
}

func (p *Prompt) suggest(text, suggestion string) {
	p.Clear()
	p.write(text)
}

func (p *Prompt) write(text string) {
	p.SetCursorIndex(p.spacer + len(text))
	txt := text
	fmt.Fprintf(p, defaultPrompt, p.icon, txt)
}

// ----------------------------------------------------------------------------
// Event Listener protocol...

// BufferCompleted indicates input was accepted.
func (p *Prompt) BufferCompleted(text string) {
	p.update(text)
}

// BufferChanged indicates the buffer was changed.
func (p *Prompt) BufferChanged(text string) {
	p.update(text)
}

// BufferActive indicates the buff activity changed.
func (p *Prompt) BufferActive(activate bool) {
	if activate {
		p.ShowCursor(true)
		p.SetBorder(true)
		p.SetTextColor(tcell.ColorWhite)
		p.SetBorderColor(tcell.ColorBlue)
		p.activate()
		return
	}

	p.ShowCursor(false)
	p.SetBorder(false)
	p.SetBackgroundColor(tcell.ColorBlack)
	p.Clear()
}

func (p *Prompt) SetIcon(icon rune) {
	p.icon = icon
}
