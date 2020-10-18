package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	initialized bool
	inputPin1   = rpio.Pin(17)
	inputPin2   = rpio.Pin(27)
	c           = []byte{}
	d           = []byte{}
	g           = []byte{}
	a           = []byte{}
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func play(note string) {
	fmt.Println("Playing note: ", note)

	var n []byte
	switch note {
	case "c":
		n = c
	case "d":
		n = d
	case "g":
		n = g
	case "a":
		n = a
	default:
		log.Panicln("Wrong note")
	}

	f := bytes.NewReader(n)
	streamer, format, err := wav.Decode(f)
	check(err)
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

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	inputPin1.Input()
	inputPin2.Input()

	fmt.Println("Hola, ya he configurado los pines 17 y 27 como entrada")

	var err error
	c, err = ioutil.ReadFile("./Celtic-harp-c2.wav")
	check(err)
	d, err = ioutil.ReadFile("./Celtic-harp-d5.wav")
	check(err)
	g, err = ioutil.ReadFile("./Celtic-harp-g4.wav")
	check(err)
	a, err = ioutil.ReadFile("./Celtic-harp-a5.wav")
	check(err)

	var dedo1puesto bool
	var dedo2puesto bool
	var dedo1estaba_puesto bool
	var dedo2estaba_puesto bool

	for {
		// fmt.Printf("En el bucle, dedo 1 = %v [antes %v], dedo 2 = %v [antes %v\n", dedo1puesto, dedo1estaba_puesto, dedo2puesto, dedo2estaba_puesto)
		if inputPin1.Read() == rpio.High {
			dedo1puesto = true
		} else {
			dedo1puesto = false
		}

		if inputPin2.Read() == rpio.High {
			dedo2puesto = true
		} else {
			dedo2puesto = false
		}

		if dedo1estaba_puesto == false && dedo1puesto == true {
			go play("g")
		}

		if dedo2estaba_puesto == false && dedo2puesto == true {
			go play("d")
		}

		dedo1estaba_puesto = dedo1puesto
		dedo2estaba_puesto = dedo2puesto
		time.Sleep(5 * time.Millisecond)
	}

}
