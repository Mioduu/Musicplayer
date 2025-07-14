package main

import (
	"fmt"
	"image/color"
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
var volume *effects.Volume
var stopTicker chan struct{}

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

func playSong(timeLabel *widget.Label, songLabel *widget.Label) {
	go func() {
		if stopTicker != nil {
			close(stopTicker)
		}

		stopTicker = make(chan struct{})

		f, err := os.Open(selectedTrack)
		if err != nil {
			fmt.Println("Błąd otwierania:", err)
			return
		}
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			fmt.Println("Błąd dekodowania:", err)
			return
		}

		currentStreamer = streamer

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
		volume = &effects.Volume{
			Streamer: ctrl,
			Base:     2,
			Volume:   -2,
			Silent:   false,
		}

		durationInSeconds := float64(streamer.Len()) / float64(format.SampleRate)
		totalMin := int(durationInSeconds) / 60
		totalSec := int(durationInSeconds) % 60

		songLabel.SetText(fmt.Sprintf("Now playing: %s", userSong))

		go func() {
			ticker := time.NewTicker(time.Millisecond * 500)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if currentStreamer == nil {
						return
					}
					pos := currentStreamer.Position()
					curSec := float64(pos) / float64(format.SampleRate)
					curMin := int(curSec) / 60
					curSecInt := int(curSec) % 60

					timeLabel.SetText(fmt.Sprintf("Time duration: %d:%02d - %d:%02d", curMin, curSecInt, totalMin, totalSec))
				case <-stopTicker:
					return
				}
			}
		}()
		speaker.Play(volume)
	}()
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
	a.Settings().SetTheme(&LofiTheme{})
	r, _ := fyne.LoadResourceFromPath("icons/ic_launcher.ico")
	a.SetIcon(r)
	w := a.NewWindow("Music player")
	w.Resize(fyne.NewSize(900, 700))
	data, err := os.ReadFile("background/background.png")
	if err != nil {
		fmt.Println("Błąd ładowania obrazu")

	}
	res := fyne.NewStaticResource("background.png", data)
	background := canvas.NewImageFromResource(res)
	background.FillMode = canvas.ImageFillStretch
	w.SetFixedSize(true)
	label := widget.NewLabelWithStyle("Music Player", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	title := container.NewCenter(label)
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Podaj patha do folderu z piosenkami")

	songList := widget.NewLabelWithStyle("Song list:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	songListCentered := container.NewCenter(songList)

	list := widget.NewList(
		func() int { return len(songs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(songs[i])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		userSong = songs[id]
		fmt.Println("Id piosenki: ", songs[id])
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

	songLabel := widget.NewLabel("")
	timeLabel := widget.NewLabel("")
	songLabel.Alignment = fyne.TextAlignCenter
	timeLabel.Alignment = fyne.TextAlignCenter

	songLabel.TextStyle = fyne.TextStyle{Bold: true}
	timeLabel.TextStyle = fyne.TextStyle{Bold: true}

	playButton := widget.NewButtonWithIcon("", iconPlay.Resource, func() {
		var songToPlay string
		if userSong != "" {
			songToPlay = userSong
		} else {
			fmt.Println("Nie wybrano żadnej piosenki")
			return
		}
		selectedTrack = filepath.Join(basePath, songToPlay+".mp3")
		playSong(timeLabel, songLabel)
	})
	playButton.Resize(fyne.NewSize(64, 64))
	stopButton := widget.NewButtonWithIcon("", iconStop.Resource, func() {
		cancelSong()
	})
	stopButton.Resize(fyne.NewSize(64, 64))
	pauseOrResumeButton := widget.NewButtonWithIcon("", iconPause.Resource, func() {
		pauseOrResume()
	})
	pauseOrResumeButton.Resize(fyne.NewSize(64, 64))

	volumeSlider := widget.NewSlider(-5, 0)
	volumeSlider.Orientation = widget.Vertical
	volumeSlider.Value = 0
	volumeSlider.Step = 0.1

	playWrapped := container.NewGridWrap(fyne.NewSize(64, 64), playButton)
	pauseWrapped := container.NewGridWrap(fyne.NewSize(64, 64), pauseOrResumeButton)
	stopWrapped := container.NewGridWrap(fyne.NewSize(64, 64), stopButton)
	volumeSliderWrapped := container.NewGridWrap(fyne.NewSize(23, 120), volumeSlider)

	sliderRect := canvas.NewRectangle(color.RGBA{255, 255, 255, 120})
	sliderRect.SetMinSize(fyne.NewSize(23, 120))

	volumeSlider.OnChanged = func(vol float64) {
		if volume == nil {
			return
		}
		speaker.Lock()
		volume.Volume = vol
		speaker.Unlock()
	}

	rect := canvas.NewRectangle(color.RGBA{255, 255, 255, 120})
	rect.SetMinSize(fyne.NewSize(250, 215))

	listScroll := container.NewVScroll(list)
	listScroll.SetMinSize(fyne.NewSize(250, 215))
	listBg := container.NewStack(rect, listScroll)

	labelsFixed := container.NewWithoutLayout(songLabel, timeLabel)

	volumeWithBg := container.NewStack(
		sliderRect,
		volumeSliderWrapped,
	)

	// Przyciski
	wWidth := float32(900)
	wHeight := float32(700)
	btnSize := float32(64)
	btnMarginBottom := float32(10)
	buttonsCount := 3
	buttonsSpacing := float32(36)
	totalButtonsWidth := btnSize*float32(buttonsCount) + buttonsSpacing*float32(buttonsCount-1)
	buttonsStartX := (wWidth - totalButtonsWidth) / 2
	buttonsY := wHeight - btnSize - btnMarginBottom

	//Volume Bar
	volWidth := float32(23)
	volHeight := float32(120)
	volMargin := float32(10)
	volX := wWidth - volWidth - volMargin
	volY := wHeight - volHeight - volMargin
	labelsHeight := float32(25)
	labelsWidth := float32(300)
	labelsX := (wWidth-labelsWidth)/2 + 150
	labelsSpacingBetween := float32(4)
	timeLabelY := buttonsY - labelsHeight*2 - labelsSpacingBetween
	songLabelY := buttonsY - labelsHeight - labelsSpacingBetween/2

	// GUI
	playWrapped.Resize(fyne.NewSize(btnSize, btnSize))
	playWrapped.Move(fyne.NewPos(buttonsStartX, buttonsY))

	pauseWrapped.Resize(fyne.NewSize(btnSize, btnSize))
	pauseWrapped.Move(fyne.NewPos(buttonsStartX+btnSize+buttonsSpacing, buttonsY))

	stopWrapped.Resize(fyne.NewSize(btnSize, btnSize))
	stopWrapped.Move(fyne.NewPos(buttonsStartX+(btnSize+buttonsSpacing)*2, buttonsY))

	volumeWithBg.Move(fyne.NewPos(volX, volY))
	volumeWithBg.Resize(fyne.NewSize(volWidth, volHeight))

	timeLabel.Move(fyne.NewPos(labelsX, timeLabelY))
	songLabel.Move(fyne.NewPos(labelsX, songLabelY))

	controlsFixed := container.NewWithoutLayout(
		playWrapped,
		pauseWrapped,
		stopWrapped,
		volumeWithBg,
	)

	grid := container.NewGridWrap(fyne.NewSize(890, 700), container.NewStack(
		background,
		container.NewVBox(
			title,
			entry,
			checkSavedDir,
			dirButton,
			songListCentered,
			listBg,
		),
	))

	w.SetContent(container.NewStack(
		background,
		grid,
		labelsFixed,
		controlsFixed,
	))
	fmt.Println("Rozmiar:", background.Size())
	w.ShowAndRun()
}
