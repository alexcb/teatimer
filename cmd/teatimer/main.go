package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/nsf/termbox-go"
)

//go:embed ready.mp3
var soundData string

func main() {

	d := time.Minute * 5
	if len(os.Args) == 2 {
		var err error
		d, err = time.ParseDuration(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	r := io.NopCloser(strings.NewReader(soundData))

	streamer, format, err := mp3.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	fmt.Printf("Starting timer for %v\n", d)
	count := int(d / time.Millisecond)

	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	go printTeaIsReady()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
}

func printTeaIsReady() {
	//	termbox.Init()
	//.var colorRange []termbox.Attribute = []termbox.Attribute{
	//.	termbox.ColorDefault,
	//.	termbox.ColorBlack,
	//.	termbox.ColorRed,
	//.	termbox.ColorGreen,
	//.	termbox.ColorYellow,
	//.	termbox.ColorBlue,
	//.	termbox.ColorMagenta,
	//.	termbox.ColorCyan,
	//.	termbox.ColorWhite,
	//.	termbox.ColorDarkGray,
	//.	termbox.ColorLightRed,
	//.	termbox.ColorLightGreen,
	//.	termbox.ColorLightYellow,
	//.	termbox.ColorLightBlue,
	//.	termbox.ColorLightMagenta,
	//.	termbox.ColorLightCyan,
	//.	termbox.ColorLightGray,
	//.}

	//.col := 5
	//.row := 7
	//.fg := colorRange[3]
	//.bg := colorRange[6]

	//.text := "tea is ready"
	//.tbprint(col, row+0, fg, bg, text)

	//.termbox.Flush()
	//.termbox.PollEvent()

	//.termbox.Close()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	color := termbox.ColorDefault
	color_change_tick := time.NewTicker(1 * time.Second)
	draw_tick := time.NewTicker(30 * time.Millisecond)
loop:
	for {
		select {
		case ev := <-event_queue:
			_ = ev
			break loop
			//if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			//	break loop
			//}
		case <-color_change_tick.C:
			color++
			if color >= 8 {
				color = 0
			}
		case <-draw_tick.C:
			w, h := termbox.Size()
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

			posx := rand.Int() % w
			posy := rand.Int() % h
			text := "TEA IS READY"
			tbprint(posx, posy, termbox.ColorDefault, color, text)

			termbox.Flush()
		}
	}
}
