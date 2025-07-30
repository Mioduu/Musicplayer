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
	rect.SetMinSize(fyne.NewSize(SONG_BACKGROUND_WIDTH, SONG_BACKGROUND_HEIGHT))
	listScroll := container.NewVScroll(list)
	listScroll.SetMinSize(fyne.NewSize(SONG_LIST_WIDTH, SONG_LIST_HEIGHT))
	return container.NewStack(rect, listScroll)
}

func MakeSongListLabel() *fyne.Container {
	songList := widget.NewLabelWithStyle("Song list:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewCenter(songList)
}

func checkLastDir() string {
	_, err := os.Stat("cfg.txt")
	if os.IsNotExist(err) {
		os.WriteFile("cfg.txt", []byte(player.UserSong), 0644)
		player.BasePath = player.UserSong
	}
	data, err := os.ReadFile("cfg.txt")
	if err != nil {
		panic(err)
	}
	player.BasePath = string(data)
	return string(data)
}

func LoadSongs(list *widget.List, songs *[]string) {
	userPath := checkLastDir()
	fmt.Println("Userpath: ", userPath)
	if userPath == "" {
		fmt.Println("Podaj ścieżke przed sprawdzeniem")
		return
	}
	player.BasePath = userPath
	os.WriteFile("cfg.txt", []byte(userPath), 0644)
	listDirectory, err := os.ReadDir(userPath)
	if err != nil {
		panic(err)
	}
	for _, file := range listDirectory {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp3") {
			name := strings.TrimSuffix(file.Name(), ".mp3")
			*songs = append(*songs, name)
		}
	}
	fmt.Println("Songs: ", songs)
	list.Refresh()
}

func MakeControls(icons *Icons, timeLabel, songLabel *widget.Label, seekSlider *widget.Slider) *fyne.Container {
	volumeSlider := widget.NewSlider(VOLUME_SLIDER_MIN_DB, VOLUME_SLIDER_MAX_DB)
	volumeSlider.Value = 0
	volumeSlider.Step = 0.1
	volumeSlider.Orientation = widget.Vertical
	player.ChangeVolume(volumeSlider)

	playButton := widget.NewButtonWithIcon("", icons.Play.Resource, func() {
		player.PlaySong(timeLabel, songLabel, seekSlider, volumeSlider)
	})

	stopButton := widget.NewButtonWithIcon("", icons.Stop.Resource, func() {
		player.CancelSong()
	})

	pauseButton := widget.NewButtonWithIcon("", icons.Pause.Resource, func() {
		player.PauseOrResume()
	})

	loopStatus := widget.NewLabel("🔁 Loop: OFF")

	loopButton := widget.NewButtonWithIcon("", icons.Loop.Resource, func() {
		player.LoopSong(timeLabel, songLabel, seekSlider, volumeSlider, loopStatus)
	})

	playWrapped := container.NewGridWrap(fyne.NewSize(ICON_MIN_SIZE, ICON_MAX_SIZE), playButton)
	pauseWrapped := container.NewGridWrap(fyne.NewSize(ICON_MIN_SIZE, ICON_MAX_SIZE), pauseButton)
	stopWrapped := container.NewGridWrap(fyne.NewSize(ICON_MIN_SIZE, ICON_MAX_SIZE), stopButton)
	loopWrapped := container.NewGridWrap(fyne.NewSize(ICON_MIN_SIZE, ICON_MAX_SIZE), loopButton)

	sliderRect := canvas.NewRectangle(color.RGBA{255, 255, 255, 120})
	sliderRect.SetMinSize(fyne.NewSize(SLIDER_RECT_WIDTH, SLIDER_RECT_HEIGHT))
	volumeSliderWrapped := container.NewGridWrap(fyne.NewSize(VOLUME_SLIDER_WIDTH, VOLUME_SLIDER_HEIGHT), volumeSlider)
	volumeWithBg := container.NewStack(sliderRect, volumeSliderWrapped)

	playWrapped.Move(fyne.NewPos(BUTTON_START_X-20, BUTTONS_Y))
	pauseWrapped.Move(fyne.NewPos(BUTTON_START_X+BUTTON_SIZE+BUTTONS_SPACING-20, BUTTONS_Y))
	stopWrapped.Move(fyne.NewPos(BUTTON_START_X+(BUTTON_SIZE+BUTTONS_SPACING)*2-20, BUTTONS_Y))
	loopWrapped.Move(fyne.NewPos(BUTTON_START_X+(BUTTON_SIZE+BUTTONS_SPACING)*3-20, BUTTONS_Y))
	loopStatus.Move(fyne.NewPos(BUTTON_START_X+(BUTTON_SIZE+BUTTONS_SPACING)*3-33, BUTTONS_Y-25))
	volumeWithBg.Move(fyne.NewPos(VOLUME_X, VOLUME_Y))
	volumeWithBg.Resize(fyne.NewSize(VOLUME_SLIDER_WIDTH, VOLUME_SLIDER_HEIGHT))

	seekSlider.Resize(fyne.NewSize(SEEK_SLIDER_WIDTH, SEEK_SLIDER_HEIGHT))
	seekSlider.Move(fyne.NewPos((WINDOW_WIDTH-300)/2, BUTTONS_Y-20))

	return container.NewWithoutLayout(playWrapped, pauseWrapped, stopWrapped, loopWrapped, volumeWithBg, loopStatus, seekSlider)
}

func MakeLabels(songLabel, timeLabel *widget.Label) *fyne.Container {
	timeLabelY := LABEL_BUTTONS_Y - LABELS_HEIGHT*3 - LABELS_SPACING
	songLabelY := LABEL_BUTTONS_Y - LABELS_HEIGHT*2 - LABELS_SPACING/2

	timeLabel.Move(fyne.NewPos(LABELS_X, timeLabelY))
	songLabel.Move(fyne.NewPos(LABELS_X, songLabelY))

	return container.NewWithoutLayout(songLabel, timeLabel)
}

func StyleLabels(songLabel, timeLabel *widget.Label) {
	songLabel.Alignment = fyne.TextAlignCenter
	timeLabel.Alignment = fyne.TextAlignCenter
	songLabel.TextStyle = fyne.TextStyle{Bold: true}
	timeLabel.TextStyle = fyne.TextStyle{Bold: true}
}
