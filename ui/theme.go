package ui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed font/Helvetica-Rounded-Bold.otf
var HelveticaTTF []byte

var HelveticaFont = fyne.NewStaticResource("assets/font/Helvetica-Rounded-Bold.otf", HelveticaTTF)

type LofiTheme struct{}

func (m *LofiTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{72, 12, 168, 255} //TŁO
	case theme.ColorNameButton:
		return color.RGBA{247, 37, 133, 255} //PRZYCISKI
	case theme.ColorNameDisabled:
		return color.RGBA{100, 100, 100, 255} //WYLACZONE
	case theme.ColorNameForeground:
		return color.RGBA{148, 230, 255, 255} //AKCENT
	case theme.ColorNamePrimary:
		return color.RGBA{181, 23, 158, 255} //SLIDERY ITP
	case theme.ColorNameInputBackground:
		return color.RGBA{181, 23, 158, 255} //INPUT
	case theme.ColorNamePlaceHolder:
		return color.RGBA{253, 128, 184, 255} //PLACEHOLDER
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m *LofiTheme) Font(style fyne.TextStyle) fyne.Resource {
	return HelveticaFont
}

func (m *LofiTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *LofiTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
