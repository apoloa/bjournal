package ui

import (
	"context"
	"os"
	"sync"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

type ActionHandler func(*tcell.EventKey) *tcell.EventKey

type KeyAction struct {
	Description string
	Action      ActionHandler
	Visible     bool
	Shared      bool
}

func NewKeyAction(d string, a ActionHandler, display bool) KeyAction {
	return KeyAction{Description: d, Action: a, Visible: display}
}

type KeyActions map[tcell.Key]KeyAction

type App struct {
	*tview.Application
	showHeader bool
	running    bool
	actions    KeyActions
	cancelFn   context.CancelFunc
	Main       *Pages
	mx         sync.RWMutex
	views      map[string]tview.Primitive
}

func NewUi() *App {
	a := App{
		Application: tview.NewApplication(),
		showHeader:  false,
		Main:        NewPages(),
		actions:     make(KeyActions),
	}

	a.views = map[string]tview.Primitive{
		"prompt": NewPrompt(false),
	}

	return &a
}

// IsRunning checks if app is actually running.
func (a *App) IsRunning() bool {
	a.mx.RLock()
	defer a.mx.RUnlock()
	return a.running
}

// SetRunning sets the app run state.
func (a *App) SetRunning(f bool) {
	a.mx.Lock()
	defer a.mx.Unlock()
	a.running = f
}

func (a *App) quitCmd(evt *tcell.EventKey) *tcell.EventKey {
	if a.InCmdMode() {
		return evt
	}
	a.BailOut()

	return nil
}

// HasAction checks if key matches a registered binding.
func (a *App) HasAction(key tcell.Key) (KeyAction, bool) {
	act, ok := a.actions[key]
	return act, ok
}

// GetActions returns a collection of actions.
func (a *App) GetActions() KeyActions {
	return a.actions
}

// AddActions returns the application actions.
func (a *App) AddActions(aa KeyActions) {
	for k, v := range aa {
		a.actions[k] = v
	}
}

// Views return the application root views.
func (a *App) Views() map[string]tview.Primitive {
	return a.views
}

// BailOut exists the application.
func (a *App) BailOut() {
	a.Stop()
	os.Exit(0)
}

// InCmdMode check if command mode is active.
func (a *App) InCmdMode() bool {
	return a.Prompt().InCmdMode()
}

// RedrawCmd forces a redraw.
func (a *App) redrawCmd(evt *tcell.EventKey) *tcell.EventKey {
	a.QueueUpdateDraw(func() {})
	return evt
}

// ResetPrompt reset the prompt model and marks buffer as active.
func (a *App) ResetPrompt() {
	a.Prompt()
	a.SetFocus(a.Prompt())
}

func (a *App) activateCmd(evt *tcell.EventKey) *tcell.EventKey {
	if a.InCmdMode() {
		return evt
	}
	a.ResetPrompt()
	//a.cmdBuff.ClearText(true)

	return nil
}

func (a *App) bindKeys() {
	a.actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Cmd", a.activateCmd, false),
		tcell.KeyCtrlR:  NewKeyAction("Redraw", a.redrawCmd, false),
		tcell.KeyCtrlC:  NewKeyAction("Quit", a.quitCmd, false),
	}
}

// Init initializes the application.
func (a *App) Init() {
	a.bindKeys()
	a.Prompt()
	//a.cmdBuff.AddListener(a)
	//a.Styles.AddListener(a)

	a.SetRoot(a.Main, true).EnableMouse(true)
}

// Prompt returns command prompt.
func (a *App) Prompt() *Prompt {
	return a.views["prompt"].(*Prompt)
}

// Resume restarts the app event loop.
func (a *App) Resume() {

}
