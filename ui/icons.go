package ui

import (
	"musicplayer/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Icons struct {
	Play  *canvas.Image
	Pause *canvas.Image
	Stop  *canvas.Image
	Loop  *canvas.Image
	Mute  *canvas.Image
}

func LoadIcons() *Icons {
	return &Icons{
		Play:  load("assets/icons/icons8-play-button-100.png"),
		Pause: load("assets/icons/icons8-pause-button-100.png"),
		Stop:  load("assets/icons/icons8-stop-squared-100.png"),
		Loop:  load("assets/icons/icons8-loop-100.png"),
		Mute:  load("assets/icons/icons8-mute-100.png"),
	}
}

func load(path string) *canvas.Image {
	icon := canvas.NewImageFromResource(player.LoadResourceFromPath(path))
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(ICON_MIN_SIZE, ICON_MAX_SIZE))
	return icon
}
