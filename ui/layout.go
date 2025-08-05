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

var volBox *VolumeBox

type VolumeBox struct {
	VolumeSlider *widget.Slider
	VolumeLabel  *widget.Label
	MuteButton   *widget.Toolbar
	Container    *fyne.Container
	Muted        bool
}

func MakeTitle() *fyne.Container {
	label := canvas.NewText("Music player", color.RGBA{76, 201, 240, 255})
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.TextSize = 20
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
	rect := canvas.NewRectangle(color.RGBA{58, 12, 163, 255})
	rect.SetMinSize(fyne.NewSize(SONG_BACKGROUND_WIDTH, SONG_BACKGROUND_HEIGHT))
	listScroll := container.NewVScroll(list)
	listScroll.SetMinSize(fyne.NewSize(SONG_LIST_WIDTH, SONG_LIST_HEIGHT))
	return container.NewStack(rect, listScroll)
}

func MakeSongListLabel() *fyne.Container {
	songList := canvas.NewText("Song list: ", color.RGBA{76, 201, 240, 255})
	songList.Alignment = fyne.TextAlignCenter
	songList.TextStyle = fyne.TextStyle{Bold: true}
	songList.TextSize = 20
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

func MakeSliders(icons *Icons) (*widget.Slider, *VolumeBox) {
	seekSlider := widget.NewSlider(0, 100)
	seekSlider.Step = 1
	seekSlider.Value = 0

	volumeSlider := widget.NewSlider(VOLUME_SLIDER_MIN_DB, VOLUME_SLIDER_MAX_DB)
	volumeSlider.Value = 0
	volumeSlider.Step = 0.1
	volumeSlider.Orientation = widget.Vertical

	volumeLabel := widget.NewLabel("0 dB")

	ChangeVolumeUI(volumeSlider, volumeLabel)

	muteButton := widget.NewToolbar(
		widget.NewToolbarAction(icons.Mute.Resource, func() {
			if volBox.Muted {
				player.UnmuteSong()
				volBox.Muted = false
				vol := volBox.VolumeSlider.Value
				volBox.VolumeLabel.SetText(fmt.Sprintf("%.0f dB", vol))
			} else {
				player.MuteSong()
				volBox.Muted = true
				volumeLabel.SetText("🔇 mute")
			}
		}),
	)

	volumeBox := container.NewHBox(
		container.NewVBox(
			muteButton,
			volumeLabel,
		),
		volumeSlider,
	)

	volBox = &VolumeBox{
		VolumeSlider: volumeSlider,
		VolumeLabel:  volumeLabel,
		MuteButton:   muteButton,
		Container:    volumeBox,
		Muted:        false,
	}

	return seekSlider, volBox
}

func ChangeVolumeUI(volumeSlider *widget.Slider, volumeLabel *widget.Label) {
	volumeSlider.OnChanged = func(vol float64) {
		if volBox.Muted {
			player.UnmuteSong()
			volBox.Muted = false
		}
		volumeLabel.SetText(fmt.Sprintf("%.0f dB", vol))

		player.ChangeVolume(vol)
	}

}

func MakeNewToolbar(icons *Icons, timeLabel, songLabel *widget.Label, seekSlider *widget.Slider, volBox *VolumeBox) *fyne.Container {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(icons.Play.Resource, func() {
			player.PlaySong(timeLabel, songLabel, seekSlider, volBox.VolumeSlider.Value)
		}),
		widget.NewToolbarAction(icons.Pause.Resource, func() {
			player.PauseOrResume()
		}),
		widget.NewToolbarAction(icons.Stop.Resource, func() {
			player.CancelSong()
		}),
		widget.NewToolbarAction(icons.Loop.Resource, func() {
			player.LoopSong(timeLabel, songLabel, seekSlider, volBox.VolumeSlider.Value)
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
