package view

import (
	"fmt"
	"os"
	"time"

	"github.com/apoloa/bjournal/src/model"
	"github.com/apoloa/bjournal/src/service"
	"github.com/apoloa/bjournal/src/ui"
	"github.com/apoloa/bjournal/src/utils"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
)

type SelectedView int

const (
	Today SelectedView = iota
	PreviousDate
	Index
)

type App struct {
	logService       *service.LogService
	prompt           *ui.Prompt
	buffer           *model.CmdBuff
	app              *tview.Application
	mainFlex         *tview.Flex
	dailyList        *ui.List
	previousDayList  *ui.List
	indexList        *ui.IndexList
	showingPrompt    bool
	showPreviousDay  bool
	showIndex        bool
	selectedView     SelectedView
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

func (a *App) buildPreviousDay(timeNow time.Time) {
	previousDate, _ := a.logService.GetPreviousDate(time.Now())
	previousList := ui.NewList().AddDailyLog(&previousDate)
	previousList.
		SetBorder(true).
		SetTitle(fmt.Sprintf("%02d.%02d %v", previousDate.Date.Day(), previousDate.Date.Month(), utils.ToShortString(timeNow.Weekday())))
	if a.selectedView == PreviousDate {
		previousList.SetBorderColor(tcell.ColorBlue)
	} else {
		previousList.SetBorderColor(tcell.ColorWhite)
	}
	a.previousDayList = previousList
}

func (a *App) makeDayFlex(fetchFromCache bool) *tview.Flex {
	flex := tview.NewFlex()
	timeNow := time.Now()
	if a.showPreviousDay {
		a.buildPreviousDay(timeNow)
		flex.AddItem(a.previousDayList, 0, 1, false)
	}
	if a.showIndex {
		indexList := ui.NewIndexList().AddIndexModel(&a.logService.Index)
		indexList.
			SetBorder(true).
			SetTitle("Index")
		if a.selectedView == Index {
			indexList.SetBorderColor(tcell.ColorBlue)
		} else {
			indexList.SetBorderColor(tcell.ColorWhite)
		}
		a.indexList = indexList
		flex.AddItem(indexList, 0, 1, false)
	}
	if fetchFromCache {
		dl, _ := a.logService.ReadDay(timeNow)
		list := ui.NewList().
			AddDailyLog(&dl)
		list.
			SetBorder(true).
			SetTitle(fmt.Sprintf("%02d.%02d %v", dl.Date.Day(), dl.Date.Month(), utils.ToShortString(dl.Date.Weekday())))
		a.dailyList = list
	}
	if a.selectedView == Today {
		a.dailyList.SetBorderColor(tcell.ColorBlue)
	} else {
		a.dailyList.SetBorderColor(tcell.ColorWhite)
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
			case event.Key() == tcell.KeyRune && event.Rune() == 't': // Create Task
				category := model.Task
				a.selectedCategory = &category
				a.showPrompt()
			case event.Key() == tcell.KeyRune && event.Rune() == 'n': // Create Note
				category := model.Note
				a.selectedCategory = &category
				a.showPrompt()
			case event.Key() == tcell.KeyRune && event.Rune() == 'e': //Create Event
				category := model.Event
				a.selectedCategory = &category
				a.showPrompt()
			case event.Key() == tcell.KeyRune && event.Rune() == 'c': // Complete
				var actualLog *model.Log
				var dateTime time.Time
				switch a.selectedView {
				case PreviousDate:
					actualLog = a.previousDayList.GetCurrentLog()
					dateTime = a.previousDayList.GetDaily().Date
				case Today:
					actualLog = a.dailyList.GetCurrentLog()
					dateTime = a.dailyList.GetDaily().Date
				}
				if actualLog != nil {
					actualLog.MarkAsComplete()
					_, err := a.logService.SaveLog(dateTime)
					if err != nil {
						log.Print("Error saving log", err)
					}
				}
			case event.Key() == tcell.KeyRune && event.Rune() == 'i': // Irrelevant
				var actualLog *model.Log
				var dateTime time.Time
				switch a.selectedView {
				case PreviousDate:
					actualLog = a.previousDayList.GetCurrentLog()
					dateTime = a.previousDayList.GetDaily().Date
				case Today:
					actualLog = a.dailyList.GetCurrentLog()
					dateTime = a.dailyList.GetDaily().Date
				}
				if actualLog != nil {
					actualLog.MarkAsIrrelevant()
					_, err := a.logService.SaveLog(dateTime)
					if err != nil {
						log.Print("Error saving log", err)
					}
				}
			case event.Key() == tcell.KeyRune && event.Rune() == 'm': // Migrate
				if a.selectedView == PreviousDate {
					previousLog := a.previousDayList.GetCurrentLog()
					if previousLog != nil {
						_, err := a.logService.MoveExistingLog(time.Now(), *previousLog)
						if err != nil {
							log.Print("Error saving log", err)
						}
						previousLog.MarkAsMigrated()
						_, err = a.logService.SaveLog(a.previousDayList.GetDaily().Date)
						if err != nil {
							log.Print("Error saving log", err)
						}
					}
				}
				a.rebuild(true)
			case event.Key() == tcell.KeyCtrlM:
				a.buildPreviousDay(time.Now())
				previousLog := a.previousDayList.GetDaily()
				if previousLog != nil {
					for i, _ := range previousLog.Logs {
						_, err := a.logService.MoveExistingLog(time.Now(), previousLog.Logs[i])
						if err != nil {
							log.Print("Error saving log", err)
						}
						previousLog.Logs[i].MarkAsMigrated()
						_, err = a.logService.SaveLog(a.previousDayList.GetDaily().Date)
						if err != nil {
							log.Print("Error saving log", err)
						}
					}
				}
				a.rebuild(true)
			case event.Key() == tcell.KeyEnter:
				if a.selectedView == Index {
					indexItem := a.indexList.GetCurrentItem()
					if indexItem != nil {
						a.app.Stop()
						a.logService.OpenIndexItem(*indexItem)
					}
				}
			case event.Key() == tcell.KeyCtrlP: // Show Previous Day
				a.showPreviousDay = !a.showPreviousDay
				if a.showIndex && a.showPreviousDay {
					a.showIndex = false
				}
				a.rebuild(false)
				if !a.showPreviousDay {
					a.selectedView = Today
				}
			case event.Key() == tcell.KeyCtrlI: // Show Index
				a.showIndex = !a.showIndex
				if a.showPreviousDay && a.showIndex {
					a.showPreviousDay = false
				}
				a.rebuild(false)
			case event.Key() == tcell.KeyCtrlJ: // Jump between views
				if a.selectedView != Today {
					a.selectedView = Today
				} else {
					switch {
					case a.showPreviousDay:
						a.selectedView = PreviousDate
					case a.showIndex:
						a.selectedView = Index
					}
				}
				a.rebuild(true)
			default:
				switch a.selectedView {
				case PreviousDate:
					handler := a.previousDayList.InputHandler()
					handler(event, func(p tview.Primitive) {})
				case Today:
					handler := a.dailyList.InputHandler()
					handler(event, func(p tview.Primitive) {})
				case Index:
					handler := a.indexList.InputHandler()
					handler(event, func(p tview.Primitive) {})
				}
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

		if a.selectedView == Index && *a.selectedCategory == model.Note {
			a.app.Stop()
			text := a.buffer.GetText()
			a.logService.CreateIndexItem(text)
			return
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
