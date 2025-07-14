package player

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func LoadResourceFromPath(path string) fyne.Resource {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Błąd przy ładowaniu ikony:", path, err)
		return theme.CancelIcon()
	}
	name := filepath.Base(path)
	return fyne.NewStaticResource(name, data)
}
