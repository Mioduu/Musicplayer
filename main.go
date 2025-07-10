package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	fmt.Println("Ze względu na prawa autorskie piosenek nie ma i nie będzie lolz dodajcie sobie sami")
	fmt.Println("Wybierz jedną z piosenek: 1: |ShadowWhisper| 2: |ChillLofi| 3: |Goodvibes| 4: |Dodaj piosenki|")
	var userInput int
	var selectedTrack string
	fmt.Scanf("%d\n", &userInput)
	switch userInput {
	case 1:
		selectedTrack = "ShadowWhisper.mp3"
	case 2:
		selectedTrack = "Lofitypeshit.mp3"
	case 3:
		selectedTrack = "Goodvibes.mp3"
	case 4:
		var userPath string
		var userSong string
		var songs []string
		_, err := os.Stat("cfg.txt")
		if os.IsNotExist(err) {
			fmt.Println("Podaj path do folderu z piosenkami: ")
			fmt.Scanf("%s\n", &userPath)
			os.WriteFile("cfg.txt", []byte(userPath), 0644)
		} else {
			fmt.Println("Znaleziono poprzedni path")
		}
		data, err := os.ReadFile("cfg.txt")
		if err != nil {
			panic(err)
		}
		listDirectory, err := os.ReadDir(string(data))
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
		fmt.Printf("Twoje piosenki: %v\n", songs)
		fmt.Println("Podaj nazwe piosenki: ")
		fmt.Scanf("%s\n", &userSong)
		selectedTrack = string(data) + userSong + ".mp3"

	}
	f, err := os.Open(selectedTrack)
	if err != nil {
		panic(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   -2,
		Silent:   false,
	}
	speaker.Play(volume)

	select {}
}
