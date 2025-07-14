package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type LofiTheme struct{}

func (m *LofiTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{230, 212, 190, 255}
	case theme.ColorNameButton:
		return color.RGBA{70, 120, 105, 255}
	case theme.ColorNameDisabled:
		return color.RGBA{160, 160, 160, 255}
	case theme.ColorNameForeground:
		return color.RGBA{60, 40, 30, 255}
	case theme.ColorNamePrimary:
		return color.RGBA{230, 147, 98, 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{210, 195, 170, 255}
	case theme.ColorNamePlaceHolder:
		return color.RGBA{150, 135, 115, 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m *LofiTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m *LofiTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *LofiTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
