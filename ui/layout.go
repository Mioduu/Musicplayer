package ui

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"musicplayer/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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

func MakeSliders() (*widget.Slider, *widget.Slider) {
	seekSlider := widget.NewSlider(0, 100)
	seekSlider.Step = 1
	seekSlider.Value = 0

	volumeSlider := widget.NewSlider(VOLUME_SLIDER_MIN_DB, VOLUME_SLIDER_MAX_DB)
	volumeSlider.Value = 0
	volumeSlider.Step = 0.1
	volumeSlider.Orientation = widget.Vertical
	player.ChangeVolume(volumeSlider)

	return seekSlider, volumeSlider
}

func MakeNewToolbar(icons *Icons, timeLabel, songLabel *widget.Label, seekSlider *widget.Slider, volumeSlider *widget.Slider) *fyne.Container {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(icons.Play.Resource, func() {
			player.PlaySong(timeLabel, songLabel, seekSlider, volumeSlider)
		}),
		widget.NewToolbarAction(icons.Pause.Resource, func() {
			player.PauseOrResume()
		}),
		widget.NewToolbarAction(icons.Stop.Resource, func() {
			player.CancelSong()
		}),
		widget.NewToolbarAction(icons.Loop.Resource, func() {
			player.LoopSong(timeLabel, songLabel, seekSlider, volumeSlider)
		}),
	)
	return container.NewCenter(toolbar)
}

func StyleLabels(songLabel, timeLabel *widget.Label) {
	songLabel.Alignment = fyne.TextAlignCenter
	timeLabel.Alignment = fyne.TextAlignCenter
	songLabel.TextStyle = fyne.TextStyle{Bold: true}
	timeLabel.TextStyle = fyne.TextStyle{Bold: true}
}
