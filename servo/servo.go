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
	T    = 20 * time.Millisecond
	Pos1 = 1250 * time.Microsecond
	Pos2 = 1500 * time.Microsecond
	Pos3 = 1750 * time.Microsecond
)

var (
	pin = rpio.Pin(14)
	pos = Pos1
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	fmt.Println("Controlador del servomotor\n")

	// Open and map memory to access gpio, check for errors
	fmt.Println("Opening rpio ...")
	check(rpio.Open())

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pins to output mode
	fmt.Println("Configuring pins ...")
	pin.Output()
	defer pin.Input()

	pin.Write(rpio.Low)

	// Motor loop
	// go func() {
	// 	for {
	// 		pin.High()
	// 		time.Sleep(pos)
	// 		pin.Low()
	// 		time.Sleep(T - pos)
	// 	}
	// }()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nControl console: ")

		s, _ := reader.ReadString('\n')
		command := strings.Fields(s)
		switch command[0] {
		case "1":
			pos = Pos1
			pin.High()
			time.Sleep(pos)
			pin.Low()
			fmt.Print("pos is now 1")
		case "2":
			pos = Pos2
			fmt.Print("pos is now 2")
			pin.High()
			time.Sleep(pos)
			pin.Low()
		case "3":
			pos = Pos3
			fmt.Print("pos is now 3")
			pin.High()
			time.Sleep(pos)
			pin.Low()
		}
	}

}
