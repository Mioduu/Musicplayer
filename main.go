package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var songs []string
var userPath string
var selectedTrack string
var userSong string
var basePath string
var userData []byte
var currentStreamer beep.StreamSeekCloser
var ctrl *beep.Ctrl
var iconPlay *canvas.Image
var iconPause *canvas.Image
var iconStop *canvas.Image

func loadResourceFromPath(path string) fyne.Resource {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Błąd przy ładowaniu ikony:", path, err)
		return theme.CancelIcon()
	}
	name := filepath.Base(path)
	return fyne.NewStaticResource(name, data)
}

func pauseOrResume() {
	if ctrl != nil {
		speaker.Lock()
		ctrl.Paused = !ctrl.Paused
		speaker.Unlock()
	}
}

func playSong() {
	f, err := os.Open(selectedTrack)
	if err != nil {
		panic(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}
	currentStreamer = streamer
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   -2,
		Silent:   false,
	}
	speaker.Play(volume)

	select {}

}

func cancelSong() {
	if currentStreamer != nil {
		currentStreamer.Close()
		currentStreamer = nil
	}

}

func checkDir(entry *widget.Entry) {
	_, err := os.Stat("cfg.txt")
	if os.IsNotExist(err) {
		os.WriteFile("cfg.txt", []byte(userPath), 0644)
		basePath = userPath
	}
	data, err := os.ReadFile("cfg.txt")
	if err != nil {
		panic(err)
	}
	userData = data
	entry.SetText(string(userData))
	basePath = string(userData)
}

func showSongs(list *widget.List) {
	songs = nil
	listDirectory, err := os.ReadDir(basePath)
	if err != nil {
		panic(err)
	}
	for _, file := range listDirectory {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.HasSuffix(fileName, ".mp3") {
				fileNameShort := strings.TrimSuffix(fileName, ".mp3")
				songs = append(songs, fileNameShort)
			}
		}
	}
	list.Refresh()
}

func main() {
	a := app.New()
	w := a.NewWindow("Music player")
	w.Resize(fyne.NewSize(600, 400))

	label := widget.NewLabel("Music player")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Podaj patha do folderu z piosenkami")

	list := widget.NewList(
		func() int { return len(songs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(songs[i])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		userSong = songs[id]
	}
	checkSavedDir := widget.NewButton("Check last dir", func() {
		checkDir(entry)
	})
	dirButton := widget.NewButton("Explore directory", func() {
		userPath = entry.Text
		if userPath == "" {
			fmt.Println("Podaj ścieżke przed sprawdzeniem")
			return
		}
		basePath = userPath
		os.WriteFile("cfg.txt", []byte(userPath), 0644)
		showSongs(list)
	})

	iconPlay = canvas.NewImageFromResource(loadResourceFromPath("icons/icons8-play-button-100.png"))
	iconPlay.FillMode = canvas.ImageFillContain
	iconPlay.SetMinSize(fyne.NewSize(64, 64))

	iconPause = canvas.NewImageFromResource(loadResourceFromPath("icons/icons8-pause-button-100.png"))
	iconPause.FillMode = canvas.ImageFillContain
	iconPause.SetMinSize(fyne.NewSize(64, 64))

	iconStop = canvas.NewImageFromResource(loadResourceFromPath("icons/icons8-stop-squared-100.png"))
	iconStop.FillMode = canvas.ImageFillContain
	iconStop.SetMinSize(fyne.NewSize(64, 64))

	playButton := widget.NewButtonWithIcon("", iconPlay.Resource, func() {
		var songToPlay string
		if userSong != "" {
			songToPlay = userSong
		} else {
			fmt.Println("Nie wybrano żadnej piosenki")
			return
		}
		selectedTrack = filepath.Join(basePath, songToPlay+".mp3")
		go playSong()
	})
	stopButton := widget.NewButtonWithIcon("", iconStop.Resource, func() {
		cancelSong()
	})
	pauseOrResumeButton := widget.NewButtonWithIcon("", iconPause.Resource, func() {
		pauseOrResume()
	})

	listScroll := container.NewVScroll(list)
	listScroll.SetMinSize(fyne.NewSize(300, 200))

	content := container.NewVBox(
		label,
		entry,
		checkSavedDir,
		dirButton,
		widget.NewLabel("Lista piosenek:"),
		listScroll,
		playButton,
		pauseOrResumeButton,
		stopButton,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
