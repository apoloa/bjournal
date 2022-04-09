module bjournal

go 1.17

replace github.com/gdamore/tcell/v2 => github.com/derailed/tcell/v2 v2.3.1-rc.2

require (
	github.com/derailed/tview v0.6.6
	github.com/gdamore/tcell/v2 v2.4.0
	github.com/mattn/go-runewidth v0.0.13
	github.com/rivo/uniseg v0.2.0
	github.com/rs/zerolog v1.22.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/term v0.0.0-20210406210042-72f3dc4e9b72 // indirect
	golang.org/x/text v0.3.6 // indirect
)
