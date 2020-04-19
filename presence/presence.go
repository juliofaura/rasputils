package main

import (
	"fmt"
	"os"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	prespin = rpio.Pin(15)
	buzzpin = rpio.Pin(14)
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	fmt.Println("Taking a measurement, date/time is", time.Now())

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
	buzzpin.Output()
	prespin.Input()

	for {
		// Wait until no presence is detected
		for prespin.Read() != 0 {
			fmt.Println("Waiting to reset")
			time.Sleep(time.Second)
		}

		for prespin.Read() == 0 {
			fmt.Println("No presence yet")
			time.Sleep(time.Second)
		}

		fmt.Println("Intruder detected!!")

		for i := 0; i < 10; i++ {
			buzzpin.High()
			time.Sleep(150 * time.Millisecond)
			buzzpin.Low()
			time.Sleep(150 * time.Millisecond)
		}

	}

}
