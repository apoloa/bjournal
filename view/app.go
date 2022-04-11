package view

import (
	"bjournal/model"
	"bjournal/service"
	"bjournal/ui"
	"bjournal/utils"
	"fmt"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type App struct {
	logService       *service.LogService
	prompt           *ui.Prompt
	buffer           *model.CmdBuff
	app              *tview.Application
	mainFlex         *tview.Flex
	dailyFlex        *ui.List
	showingPrompt    bool
	selectedCategory *model.Category
}

func NewApp(logService *service.LogService) *App {
	prompt := ui.NewPrompt(false)
	buffer := model.NewCmdBuff('>')
	mainFlex := tview.NewFlex()
	mainFlex.SetDirection(tview.FlexRow)
	prompt.SetModel(buffer)
	app := &App{
		logService: logService,
		prompt:     prompt,
		buffer:     buffer,
		app:        tview.NewApplication(),
		mainFlex:   mainFlex,
	}
	buffer.AddListener(app)
	return app
}

func (a *App) makeDayFlex() *tview.Flex {
	flex := tview.NewFlex()
	timeNow := time.Now()
	dl, _ := a.logService.ReadDay(timeNow)
	list := ui.NewList()
	for index, _ := range dl.Logs {
		list.AddItem(&dl.Logs[index], nil)
	}
	list.SetBorder(true)
	list.SetTitle(fmt.Sprintf("%02d.%02d %v", timeNow.Day(), timeNow.Month(), utils.ToShortString(timeNow.Weekday())))
	flex.AddItem(list, 0, 1, false)
	a.dailyFlex = list
	return flex
}

func (a *App) showPrompt() {
	if a.mainFlex.ItemAt(0) != a.prompt {
		a.mainFlex.Clear()
		a.mainFlex.
			AddItemAtIndex(0, a.prompt, 3, 1, false).
			AddItem(a.dailyFlex, 0, 1, false)
	}
	a.showingPrompt = true
}

func (a *App) rebuild() {
	a.prompt = ui.NewPrompt(false)
	a.prompt.SetModel(a.buffer)
	a.makeDayFlex()
	a.mainFlex.Clear()
	a.mainFlex.AddItem(a.dailyFlex, 0, 1, false)
}

func (a *App) hidePrompt() {
	a.mainFlex.Clear()
	a.mainFlex.AddItem(a.dailyFlex, 0, 1, false)
	a.showingPrompt = false
}

func (a *App) Show() {

	a.rebuild()
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if a.showingPrompt {
			a.prompt.GetInputCapture()(event)
		} else {
			if event.Key() == tcell.KeyRune && event.Rune() == 't' {
				category := model.Task
				a.selectedCategory = &category
				a.showPrompt()
			} else if event.Key() == tcell.KeyRune && event.Rune() == 'n' {
				category := model.Note
				a.selectedCategory = &category
				a.showPrompt()
			} else if event.Key() == tcell.KeyRune && event.Rune() == 'e' {
				category := model.Event
				a.selectedCategory = &category
				a.showPrompt()
			} else if event.Key() == tcell.KeyRune && event.Rune() == 'c' {
				index := a.dailyFlex.GetCurrentItem()
				if index > 0 {
					a.dailyFlex.GetItem(index)
				}

			} else {
				handler := a.dailyFlex.InputHandler()
				handler(event, func(p tview.Primitive) {})
			}
		}
		return event
	})

	if err := a.app.SetRoot(a.mainFlex, true).SetFocus(a.mainFlex).Run(); err != nil {
		panic(err)
	}
}

func (a *App) BufferCompleted(text string) {}

func (a *App) BufferChanged(text string) {}

func (a *App) BufferActive(state bool) {
	if state == false {
		if a.selectedCategory == nil {
			log.Print("Buffer complete without selected category")
			os.Exit(101)
		}
		var selectedLog *model.Log
		index := a.dailyFlex.GetCurrentItem()
		if index >= 0 {
			selectedLog = a.dailyFlex.GetItem(index)
			if selectedLog.Parent != nil {
				selectedLog = selectedLog.Parent
			}
		}
		text := a.buffer.GetText()
		if len(text) != 0 {
			if selectedLog != nil {
				selectedLog.AppendNewSubLog(text, *a.selectedCategory)
				_, err := a.logService.SaveLog(time.Now())
				if err != nil {
					return
				}
			} else {
				_, err := a.logService.AddNewLog(time.Now(), text, *a.selectedCategory)
				if err != nil {
					return
				}
			}
		}
		a.showingPrompt = false
		a.buffer.ClearText(true)
		a.hidePrompt()
		a.rebuild()
	}
}
