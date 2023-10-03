package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/apoloa/bjournal/src/api"
	"github.com/apoloa/bjournal/src/service"
	"github.com/apoloa/bjournal/src/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
