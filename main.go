package main

import (
	"github.com/apoloa/bjournal/src/api"
	"github.com/apoloa/bjournal/src/service"
	"github.com/apoloa/bjournal/src/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
)

/*
func makeWeekFlex(m *model.Model) *tview.Flex {
	flex := tview.NewFlex()
	timeNow := time.Now()
	intWeekDay := int(timeNow.Weekday())
	for i := 1; i <= 5; i++ {
		dayTime := time.Now().Add(time.Duration((i-intWeekDay)*24) * time.Hour)
		dl, _ := m.ReadDay(dayTime)
		list := ui.NewList()
		for _, log := range dl.Logs {
			list.AddItem(log.Name, "", log.Mark.Print(), nil).ShowSecondaryText(false)
		}
		list.SetBorder(true)
		list.SetTitle(fmt.Sprintf("%v (%v/%v)", dayTime.Weekday(), dayTime.Day(), dayTime.Month()))
		flex.AddItem(list, 0, 1, false)
	}
	return flex
}*/

func main() {
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(executablePath)
	mainPath := path.Join(exePath, "main.log")
	file, err := os.OpenFile(mainPath, mod, 0777)
	if err != nil {
		log.Printf("Error %v \n", err)
		log.Fatal()
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})

	m := service.NewLogService("/Users/apoloalcaide/Developer/Journal")

	router := api.NewRouter(8778, m)
	router.Init()
	go router.Start()

	app := view.NewApp(m)
	app.Show()
}

/*
func main2() {
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile("main.log", mod, 0777)
	if err != nil {
		log.Printf("Error %v \n", err)
		log.Fatal()
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
	m := model.NewModel("./bjournal-files")
	prompt := ui.NewPrompt(false)

	app := tview.NewApplication()

	weekWindow := makeWeekFlex(m)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(weekWindow, 0, 1, false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if flex.ItemAt(0) != prompt {
				flex.Clear()
				flex.AddItemAtIndex(0, prompt, 3, 1, false).AddItem(weekWindow, 0, 1, false)
			} else {
				flex.Clear().AddItem(weekWindow, 0, 1, false)
			}

			return nil
		}
		return event
	})
	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		x, y := event.Position()
		log.Printf("Event %v, at position %v - %v ", action, x, y)
		return event, action
	})
	flex.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		x, y := event.Position()
		log.Printf("Event %v, at position %v - %v ", action, x, y)
		if action == tview.MouseLeftClick {
			log.Printf("Click at %v-%v \n", x, y)
			weekWindow.GetMouseCapture()(action, event)
		}
		return action, event
	})

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func mainPages() {
	const pageCount = 5

	app := tview.NewApplication()
	pages := tview.NewPages()
	for page := 0; page < pageCount; page++ {
		func(page int) {
			pages.AddPage(fmt.Sprintf("page-%d", page),
				ui.NewPrompt(false),
				false,
				page == 0)
		}(page)
	}
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
t
*/
