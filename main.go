package main

import (
	"fmt"
	"musicplayer/player"
	"musicplayer/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&ui.LofiTheme{})
	r, _ := fyne.LoadResourceFromPath("assets/icons/ic_launcher.ico")
	a.SetIcon(r)

	w := a.NewWindow("Music player")
	w.Resize(fyne.NewSize(900, 700))

	title := ui.MakeTitle()

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Podaj patha do folderu z piosenkami")

	songListLabel := ui.MakeSongListLabel()

	songList, songs := ui.MakeSongList()
	player.SongListPointer = songs
	ui.LoadSongs(songList, songs)
	// searchEntry := ui.MakeSearchEntry(filtered, songs, songList)

	// entryButtons := ui.MakeEntryButtons(entry, songList, songs, filtered)

	playerIcons := ui.LoadIcons()

	songLabel := widget.NewLabel("")
	timeLabel := widget.NewLabel("")
	ui.StyleLabels(songLabel, timeLabel)

	seekSlider, volumeSlider := ui.MakeSliders()
	toolbarcontrol := ui.MakeNewToolbar(playerIcons, timeLabel, songLabel, seekSlider, volumeSlider)

	toolbarContainer := container.NewCenter(toolbarcontrol)
	volumeContainer := container.NewStack(volumeSlider)
	listUI := ui.MakeSongListUI(songList)
	wrappedSeekSlider := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(ui.SEEK_SLIDER_WIDTH, ui.SEEK_SLIDER_HEIGHT)),
		seekSlider,
	)
	centeredSeek := container.NewCenter(wrappedSeekSlider)
	// spacer := layout.NewSpacer()

	bottomBar := container.NewBorder(
		nil,
		nil,
		nil,
		volumeContainer,
		container.NewVBox(
			songLabel,
			timeLabel,
			centeredSeek,
			toolbarContainer,
		),
	)

	grid := container.NewBorder(
		nil,       // Top
		bottomBar, // Bottom
		nil,       // Left
		nil,       // Right
		container.NewVBox( // Center
			title,
			entry,
			songListLabel,
			listUI,
		),
	)

	w.SetContent(container.NewStack(
		grid,
	))

	fmt.Println("seekSlider MinSize:", seekSlider.MinSize())
	fmt.Println("toolbar MinSize:", toolbarContainer.MinSize())
	fmt.Println("volumeSlider MinSize:", volumeSlider.MinSize())

	// fmt.Println("Rozmiar:", background.Size())
	w.ShowAndRun()
}
