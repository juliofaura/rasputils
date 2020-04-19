package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const (
	delay = 2 * time.Second
	on    = rpio.Low
	off   = rpio.High
)

var (
	outpin = rpio.Pin(17)
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func estadoDeLaCaldera() string {
	estado := outpin.Read()
	if estado == on {
		return "encendida"
	} else {
		return "apagada"
	}
}

func encender() {
	outpin.Write(on)
}

func apagar() {
	outpin.Write(off)
}

func main() {

	fmt.Println("Consola de la caldera\n")

	// Open and map memory to access gpio, check for errors
	fmt.Println("Opening rpio ...")
	check(rpio.Open())

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	fmt.Println("Configuring pins ...")
	outpin.Output()
	defer outpin.Input()

	for {
		fmt.Println("La caldera está", estadoDeLaCaldera())
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Introducir comando (encender / apagar / exit): ")
		s, _ := reader.ReadString('\n')
		command := strings.Fields(s)
		if len(command) >= 1 {
			switch command[0] {
			case "exit":
				fmt.Println("Have a nice day!")
				os.Exit(0)
			case "encender":
				// To do
				encender()
			case "apagar":
				// To do
				apagar()
			default:
				fmt.Printf("Unknown command %v\n", command)
			}
		}
	}

}
