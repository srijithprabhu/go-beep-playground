package hello_beep

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	_ "github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func Run() {
	f, err := os.Open("hello_beep/Lame_Drivers_-_01_-_Frozen_Egg.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// to hang forever
	// speaker.Play(streamer)
	// select {}

	// Calls in sequence and uses channel to indicate that the stream is over.
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// use beep.Resample(quality_index int, old, new SampleRate, s Streamer) to resample a streamer... returns a new streamer.
	// higher quality_index means better quality, but also higher CPU consumption.
	// Used if your speaker has a different sample rate than your input audio.

	<- done
}