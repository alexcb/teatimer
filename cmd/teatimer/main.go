package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed ready.mp3
var soundData string

func main() {
	r := io.NopCloser(strings.NewReader(soundData))

	streamer, format, err := mp3.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	d := time.Minute * 5
	fmt.Printf("Starting timer for %v\n", d)
	count := int(d / time.Millisecond)

	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
