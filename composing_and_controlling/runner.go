package composing_and_controlling

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func Run() {
	f, err := os.Open("composing_and_controlling/Miami_Slice_-_04_-_Step_Into_Me.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Looping and Resampling. Showcases wrapping streamers.
	/*
	loop := beep.Loop(3, streamer)
	fast := beep.ResampleRatio(4, 10, loop)

	done := make(chan bool)
	speaker.Play(beep.Seq(fast, beep.Callback(func() {
		done <- true
	})))

	for {
		select {
		case <- done:
			return
		case <- time.After(time.Second):
			speaker.Lock()
			fmt.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
			speaker.Unlock()
		}
	}*/
	// Control with Ctrl.
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	// Volume. Base is multiplier to use. Volume is the power with which to multiple base. Base ^ Volume
	volume := &effects.Volume{
		Streamer: ctrl,
		Base: 2,
		Volume: 0,
		Silent: false,
	}
	speedy := beep.ResampleRatio(4, 1, volume)
	speaker.Play(speedy)
	for {
		nextAction := "pause"
		if ctrl.Paused {
			nextAction = "resume"
			volume.Volume += 1
			speedy.SetRatio(speedy.Ratio() + 0.2)
		}
		fmt.Printf("Press [ENTER] to %s.", nextAction)
		fmt.Scanln()

		speaker.Lock()
		ctrl.Paused = !ctrl.Paused
		speaker.Unlock()
	}
}