package main

import (
	"fmt"

	"github.com/yryz/ds18b20"
)

func main() {

	sensors, err := ds18b20.Sensors()
	if err != nil {
		fmt.Println("No sensors")
	}

	for _, sensor := range sensors {
		t, err := ds18b20.Temperature(sensor)
		if err == nil {
			fmt.Printf("%.2f\n", t)
		} else {
			fmt.Println("NA")
		}
	}
}
