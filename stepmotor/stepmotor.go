package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const (
	PULSE        = 100 * time.Microsecond
	timeInterval = time.Microsecond
)

var (
	ms1pin   = rpio.Pin(14)
	ms2pin   = rpio.Pin(15)
	ms3pin   = rpio.Pin(18)
	dirpin   = rpio.Pin(2)
	steppin  = rpio.Pin(3)
	sleeppin = rpio.Pin(23)
	resetpin = rpio.Pin(24)

	period  = 1 * time.Millisecond
	motorOn = false
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func step() {
	steppin.Write(rpio.High)
	time.Sleep(PULSE)
	steppin.Write(rpio.Low)
	time.Sleep(PULSE)
}

func main() {

	fmt.Println("Controlador del motor de paso\n")

	// Open and map memory to access gpio, check for errors
	fmt.Println("Opening rpio ...")
	check(rpio.Open())

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pins to output mode
	fmt.Println("Configuring pins ...")
	ms1pin.Output()
	ms2pin.Output()
	ms3pin.Output()
	dirpin.Output()
	steppin.Output()
	sleeppin.Output()
	resetpin.Output()
	defer ms1pin.Input()
	defer ms2pin.Input()
	defer ms3pin.Input()
	defer dirpin.Input()
	defer steppin.Input()
	defer sleeppin.Input()
	defer resetpin.Input()

	ms1pin.Write(rpio.Low)
	ms2pin.Write(rpio.Low)
	ms3pin.Write(rpio.Low)
	resetpin.Write(rpio.High)

	// Motor loop
	go func() {
		pulseStarts := time.Now()
		for {
			if !motorOn {
				if sleeppin.Read() == rpio.High {
					sleeppin.Write(rpio.Low)
					fmt.Println("Sleep set to Low")
					time.Sleep(500 * time.Millisecond)
					pulseStarts = time.Now()
				}
			} else {
				if sleeppin.Read() == rpio.Low {
					sleeppin.Write(rpio.High)
					fmt.Println("Sleep set to High")
					time.Sleep(500 * time.Millisecond)
					pulseStarts = time.Now()
				}
				for time.Since(pulseStarts) < period {
				}
				step()
				pulseStarts = pulseStarts.Add(period)
			}
			time.Sleep(timeInterval)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nControl console: ")

		s, _ := reader.ReadString('\n')
		command := strings.Fields(s)
		if len(command) >= 1 {
			switch command[0] {
			case "exit":
				fmt.Println("Have a nice day!")
				log.Println("Ending program, closing log\n\n")
				sleeppin.Write(rpio.Low)
				os.Exit(0)
			case "status":
				fmt.Printf("Period is %d ms\n", period/time.Microsecond)
				fmt.Println("MotorOn is", motorOn, " and pin is", sleeppin.Read())
			case "changePeriod":
				if len(command) != 2 {
					fmt.Println("Missing target period, syntax is: changePeriod <period in us>")
					continue
				}
				newPeriod, err := strconv.ParseInt(command[1], 10, 64)
				if err != nil {
					fmt.Println("Wrong target period: ", command[1])
					continue
				}
				period = time.Duration(newPeriod) * time.Microsecond
				fmt.Printf("Target period changed to %d", newPeriod)
			case "motorOff":
				motorOn = false
				time.Sleep(100 * time.Millisecond)
				fmt.Println("Motor disconnected")
			case "motorOn":
				motorOn = true
				time.Sleep(100 * time.Millisecond)
				fmt.Println("Motor connected")
			case "help":
				fmt.Println("COMMANDS:")
				fmt.Println("status - prints current status")
				fmt.Println("changePeriod <period in us> - sets a new target period in us")
				fmt.Println("motorOff - disconnects the motor")
				fmt.Println("motorOn - connects the motor")
				fmt.Println("help - prints help ;-)")
				fmt.Println("exit - exists program")
			default:
				fmt.Printf("Unknown command %v\n", command)
			}
		}
	}

}
