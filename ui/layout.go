package ui

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"musicplayer/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep/speaker"
)

func LoadBackground(path string) *canvas.Image {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Błąd ładowania obrazu")
		return canvas.NewImageFromImage(nil)
	}
	res := fyne.NewStaticResource(filepath.Base(path), data)
	bg := canvas.NewImageFromResource(res)
	bg.FillMode = canvas.ImageFillStretch
	return bg
}

func MakeTitle() *fyne.Container {
	label := widget.NewLabelWithStyle("Music Player", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewCenter(label)
}

func MakeSongList() (*widget.List, *[]string) {
	songs := &[]string{}
	list := widget.NewList(
		func() int { return len(*songs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText((*songs)[i])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		player.UserSong = (*songs)[id]
		fmt.Println("Id piosenki:", (*songs)[id])
	}
	return list, songs
}

func MakeSongListUI(list *widget.List) *fyne.Container {
	rect := canvas.NewRectangle(color.RGBA{255, 255, 255, 120})
	rect.SetMinSize(fyne.NewSize(250, 215))
	listScroll := container.NewVScroll(list)
	listScroll.SetMinSize(fyne.NewSize(250, 215))
	return container.NewStack(rect, listScroll)
}

func MakeSongListLabel() *fyne.Container {
	songList := widget.NewLabelWithStyle("Song list:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewCenter(songList)
}

func MakeEntryButtons(entry *widget.Entry, list *widget.List, songs *[]string) [2]*widget.Button {
	checkSavedDir := widget.NewButton("Check last dir", func() {
		_, err := os.Stat("cfg.txt")
		if os.IsNotExist(err) {
			os.WriteFile("cfg.txt", []byte(player.UserSong), 0644)
			player.BasePath = player.UserSong
		}
		data, err := os.ReadFile("cfg.txt")
		if err != nil {
			panic(err)
		}
		entry.SetText(string(data))
		player.BasePath = string(data)
	})
	dirButton := widget.NewButton("Explore directory", func() {
		userPath := entry.Text
		if userPath == "" {
			fmt.Println("Podaj ścieżke przed sprawdzeniem")
			return
		}
		player.BasePath = userPath
		os.WriteFile("cfg.txt", []byte(userPath), 0644)
		// Refresh song list
		*songs = nil
		listDirectory, err := os.ReadDir(player.BasePath)
		if err != nil {
			panic(err)
		}
		for _, file := range listDirectory {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp3") {
				name := strings.TrimSuffix(file.Name(), ".mp3")
				*songs = append(*songs, name)
			}
		}
		list.Refresh()
	})
	return [2]*widget.Button{checkSavedDir, dirButton}
}

func MakeControls(icons *Icons, timeLabel, songLabel *widget.Label) *fyne.Container {
	playButton := widget.NewButtonWithIcon("", icons.Play.Resource, func() {
		if player.UserSong == "" {
			fmt.Println("Nie wybrano żadnej piosenki")
			return
		}
		player.SelectedTrack = filepath.Join(player.BasePath, player.UserSong+".mp3")
		player.PlaySong(timeLabel, songLabel)
	})
	stopButton := widget.NewButtonWithIcon("", icons.Stop.Resource, func() {
		player.CancelSong()
	})
	pauseButton := widget.NewButtonWithIcon("", icons.Pause.Resource, func() {
		player.PauseOrResume()
	})

	volumeSlider := widget.NewSlider(-5, 0)
	volumeSlider.Orientation = widget.Vertical
	volumeSlider.Value = 0
	volumeSlider.Step = 0.1
	volumeSlider.OnChanged = func(vol float64) {
		if player.Volume == nil {
			return
		}
		speaker.Lock()
		player.Volume.Volume = vol
		speaker.Unlock()
	}

	playWrapped := container.NewGridWrap(fyne.NewSize(64, 64), playButton)
	pauseWrapped := container.NewGridWrap(fyne.NewSize(64, 64), pauseButton)
	stopWrapped := container.NewGridWrap(fyne.NewSize(64, 64), stopButton)

	sliderRect := canvas.NewRectangle(color.RGBA{255, 255, 255, 120})
	sliderRect.SetMinSize(fyne.NewSize(23, 120))
	volumeSliderWrapped := container.NewGridWrap(fyne.NewSize(23, 120), volumeSlider)
	volumeWithBg := container.NewStack(sliderRect, volumeSliderWrapped)

	// Pozycjonowanie
	wWidth := float32(900)
	wHeight := float32(700)
	btnSize := float32(64)
	btnMarginBottom := float32(10)
	buttonsSpacing := float32(36)
	buttonsCount := 3
	totalButtonsWidth := btnSize*float32(buttonsCount) + buttonsSpacing*float32(buttonsCount-1)
	buttonsStartX := (wWidth - totalButtonsWidth) / 2
	buttonsY := wHeight - btnSize - btnMarginBottom
	volX := wWidth - 23 - 10
	volY := wHeight - 120 - 10

	playWrapped.Move(fyne.NewPos(buttonsStartX, buttonsY))
	pauseWrapped.Move(fyne.NewPos(buttonsStartX+btnSize+buttonsSpacing, buttonsY))
	stopWrapped.Move(fyne.NewPos(buttonsStartX+(btnSize+buttonsSpacing)*2, buttonsY))
	volumeWithBg.Move(fyne.NewPos(volX, volY))
	volumeWithBg.Resize(fyne.NewSize(23, 120))

	return container.NewWithoutLayout(playWrapped, pauseWrapped, stopWrapped, volumeWithBg)
}

func MakeLabels(songLabel, timeLabel *widget.Label) *fyne.Container {
	labelsHeight := float32(25)
	labelsSpacing := float32(4)
	wWidth := float32(900)
	btnSize := float32(64)
	btnMarginBottom := float32(10)
	buttonsY := float32(700) - btnSize - btnMarginBottom
	labelsWidth := float32(300)
	labelsX := (wWidth-labelsWidth)/2 + 150
	timeLabelY := buttonsY - labelsHeight*2 - labelsSpacing
	songLabelY := buttonsY - labelsHeight - labelsSpacing/2

	timeLabel.Move(fyne.NewPos(labelsX, timeLabelY))
	songLabel.Move(fyne.NewPos(labelsX, songLabelY))

	return container.NewWithoutLayout(songLabel, timeLabel)
}

func StyleLabels(songLabel, timeLabel *widget.Label) {
	songLabel.Alignment = fyne.TextAlignCenter
	timeLabel.Alignment = fyne.TextAlignCenter
	songLabel.TextStyle = fyne.TextStyle{Bold: true}
	timeLabel.TextStyle = fyne.TextStyle{Bold: true}
}
