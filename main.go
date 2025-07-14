package main

import (
	"fmt"

	"musicplayer/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&ui.LofiTheme{})
	r, _ := fyne.LoadResourceFromPath("assets/icons/ic_launcher.ico")
	a.SetIcon(r)

	w := a.NewWindow("Music player")
	w.Resize(fyne.NewSize(900, 700))
	w.SetFixedSize(true)

	background := ui.LoadBackground("assets/background/background.png")
	title := ui.MakeTitle()

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Podaj patha do folderu z piosenkami")

	songListLabel := ui.MakeSongListLabel()

	songList, songs := ui.MakeSongList()

	entryButtons := ui.MakeEntryButtons(entry, songList, songs)

	playerIcons := ui.LoadIcons()

	songLabel := widget.NewLabel("")
	timeLabel := widget.NewLabel("")
	ui.StyleLabels(songLabel, timeLabel)

	controls := ui.MakeControls(playerIcons, timeLabel, songLabel)

	labels := ui.MakeLabels(songLabel, timeLabel)

	listUI := ui.MakeSongListUI(songList)

	grid := container.NewGridWrap(fyne.NewSize(890, 700), container.NewStack(
		background,
		container.NewVBox(
			title,
			entry,
			entryButtons[0],
			entryButtons[1],
			songListLabel,
			listUI,
		),
	))

	w.SetContent(container.NewStack(
		background,
		grid,
		labels,
		controls,
	))

	fmt.Println("Rozmiar:", background.Size())
	w.ShowAndRun()
}
