package player

import (
	"fmt"
	"os"
	"time"

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
)

var (
	UserSong      string
	SelectedTrack string
	BasePath      string
)

func PlaySong(timeLabel, songLabel *widget.Label) {
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

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		Ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
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
