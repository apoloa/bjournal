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
	dailyList        *ui.List
	showingPrompt    bool
	showPreviousDay  bool
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

func (a *App) makeDayFlex(fetchFromCache bool) *tview.Flex {
	flex := tview.NewFlex()
	timeNow := time.Now()
	if a.showPreviousDay {
		previousDate, _ := a.logService.GetPreviousDate()
		previousList := ui.NewList().AddDailyLog(&previousDate)
		previousList.
			SetBorder(true).
			SetTitle(fmt.Sprintf("%02d.%02d %v", timeNow.Day(), timeNow.Month(), utils.ToShortString(timeNow.Weekday())))
		flex.AddItem(previousList, 0, 1, false)
	}
	if fetchFromCache {
		dl, _ := a.logService.ReadDay(timeNow)
		list := ui.NewList().
			AddDailyLog(&dl)
		list.
			SetBorder(true).
			SetTitle(fmt.Sprintf("%02d.%02d %v", timeNow.Day(), timeNow.Month(), utils.ToShortString(timeNow.Weekday())))
		a.dailyList = list
	}
	flex.AddItem(a.dailyList, 0, 1, false)
	return flex
}

func (a *App) showPrompt() {
	a.showingPrompt = true
	a.rebuild(false)
}

func (a *App) rebuild(fetchFromCache bool) {
	itemsFlex := a.makeDayFlex(fetchFromCache)
	a.mainFlex.Clear()
	if a.showingPrompt {
		a.prompt = ui.NewPrompt(false)
		a.prompt.SetModel(a.buffer)
		a.prompt.SetIcon(a.selectedCategory.Print())
		a.mainFlex.
			AddItemAtIndex(0, a.prompt, 3, 1, false)
	}
	a.mainFlex.AddItem(itemsFlex, 0, 1, false)
}

func (a *App) hidePrompt() {
	a.showingPrompt = false
	a.rebuild(false)
}

func (a *App) Show() {

	a.rebuild(true)
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if a.showingPrompt {
			a.prompt.GetInputCapture()(event)
		} else {
			switch {
			case event.Key() == tcell.KeyRune && event.Rune() == 't':
				category := model.Task
				a.selectedCategory = &category
				a.showPrompt()
			case event.Key() == tcell.KeyRune && event.Rune() == 'n':
				category := model.Note
				a.selectedCategory = &category
				a.showPrompt()
			case event.Key() == tcell.KeyRune && event.Rune() == 'e':
				category := model.Event
				a.selectedCategory = &category
				a.showPrompt()
			case event.Key() == tcell.KeyRune && event.Rune() == 'c':
				actualLog := a.dailyList.GetCurrentLog()
				if actualLog != nil {
					actualLog.MarkAsComplete()
					_, err := a.logService.SaveLog(time.Now())
					if err != nil {
						log.Print("Error saving log", err)
					}
				}
			case event.Key() == tcell.KeyRune && event.Rune() == 'i':
				actualLog := a.dailyList.GetCurrentLog()
				if actualLog != nil {
					actualLog.MarkAsIrrelevant()
					_, err := a.logService.SaveLog(time.Now())
					if err != nil {
						log.Print("Error saving log", err)
					}
				}
			case event.Key() == tcell.KeyCtrlI:
				a.showPreviousDay = !a.showPreviousDay
				a.rebuild(false)
			default:
				handler := a.dailyList.InputHandler()
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
		index := a.dailyList.GetCurrentItem()
		if index >= 0 {
			selectedLog = a.dailyList.GetItem(index)
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
		a.rebuild(true)
	}
}
