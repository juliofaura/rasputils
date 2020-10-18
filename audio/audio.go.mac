package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	termbox "github.com/nsf/termbox-go"
)

var initialized bool

func play(note string) {
	fmt.Println("Playing note: ", note)

	var filename string
	switch note {
	case "c":
		filename = "./Celtic-harp-c2.wav"
	case "d":
		filename = "./Celtic-harp-d5.wav"
	case "g":
		filename = "./Celtic-harp-g4.wav"
	case "a":
		filename = "./Celtic-harp-a5.wav"
	default:
		log.Panicln("Wrong note")
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	if !initialized {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		initialized = true
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	fmt.Println("Played note: ", note)
}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch key := ev.Ch; key {
			case 'd':
				go play("d")
			case 'g':
				go play("g")
			case 'a':
				go play("a")
			case 'q':
				fmt.Println("Exiting now. Have a nice day!")
				os.Exit(0)
			default:
				fmt.Printf("Unknown note %v\n", key)
			}
		default:
			fmt.Println("Unknown event: ", ev)
		}
	}

}
