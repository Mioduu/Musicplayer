package player

import (
	"fmt"
	"os"
	"time"

	"path/filepath"

	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	CurrentStreamer beep.StreamSeekCloser
	Ctrl            *beep.Ctrl
	Volume          *effects.Volume
	StopTicker      chan struct{}
	IsLooping       = false
	SongListPointer *[]string
)

var (
	UserSong      string
	SelectedTrack string
	BasePath      string
	isSeeking     bool
)

func PlayNextSong(timeLabel, songLabel *widget.Label, seekSlider *widget.Slider) {
	currentIndex := -1
	for i, name := range *SongListPointer {
		if name == UserSong {
			currentIndex = i
			break
		}
	}
	if currentIndex >= 0 {
		nextIndex := (currentIndex + 1) % len(*SongListPointer)
		UserSong = (*SongListPointer)[nextIndex]
		SelectedTrack = filepath.Join(BasePath, UserSong+".mp3")
		PlaySong(timeLabel, songLabel, seekSlider)
	} else {
		fmt.Println("Nie znaleziono obecnej piosenki w liście")
	}

}

func PlaySong(timeLabel, songLabel *widget.Label, seekSlider *widget.Slider) {
	go func() {
		if StopTicker != nil {
			close(StopTicker)
		}
		StopTicker = make(chan struct{})

		f, err := os.Open(SelectedTrack)
		if err != nil {
			fmt.Println("Błąd otwierania:", err)
			return
		}
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			fmt.Println("Błąd dekodowania:", err)
			return
		}
		CurrentStreamer = streamer
		var mainStreamer beep.Streamer = streamer

		if IsLooping {
			mainStreamer = beep.Loop(-1, streamer)
		}

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		Ctrl = &beep.Ctrl{Streamer: mainStreamer, Paused: false}
		Volume = &effects.Volume{
			Streamer: Ctrl,
			Base:     2,
			Volume:   -2,
			Silent:   false,
		}

		duration := float64(streamer.Len()) / float64(format.SampleRate)
		totalMin := int(duration) / 60
		totalSec := int(duration) % 60

		songLabel.SetText(fmt.Sprintf("Now playing: %s", UserSong))

		seekSlider.Min = 0
		seekSlider.Max = float64(streamer.Len())
		seekSlider.Value = 0

		seekSlider.OnChanged = func(value float64) {
			if !isSeeking && CurrentStreamer != nil {
				speaker.Lock()
				_ = CurrentStreamer.Seek(int(value))
				speaker.Unlock()
			}
		}

		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if CurrentStreamer == nil {
						return
					}
					pos := CurrentStreamer.Position()
					cur := float64(pos) / float64(format.SampleRate)
					curMin := int(cur) / 60
					curSec := int(cur) % 60
					timeLabel.SetText(fmt.Sprintf("Time duration: %d:%02d - %d:%02d", curMin, curSec, totalMin, totalSec))
					isSeeking = true
					seekSlider.SetValue(float64(CurrentStreamer.Position()))
					isSeeking = false

					if pos >= CurrentStreamer.Len() {
						if IsLooping {
							PlaySong(timeLabel, songLabel, seekSlider)
						} else {
							PlayNextSong(timeLabel, songLabel, seekSlider)
						}
						return
					}

				case <-StopTicker:
					return
				}
			}
		}()
		speaker.Play(Volume)
	}()
}

func PauseOrResume() {
	if Ctrl != nil {
		speaker.Lock()
		Ctrl.Paused = !Ctrl.Paused
		speaker.Unlock()
	}
}

func CancelSong() {
	if CurrentStreamer != nil {
		CurrentStreamer.Close()
		CurrentStreamer = nil
	}
}

func LoopSong(timeLabel, songLabel *widget.Label, seekSlider *widget.Slider) {
	IsLooping = !IsLooping
	CancelSong()
	SelectedTrack = filepath.Join(BasePath, UserSong+".mp3")
	PlaySong(timeLabel, songLabel, seekSlider)
}
