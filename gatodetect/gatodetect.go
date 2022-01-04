package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	samplesForCalibration = 10
	threshold             = 50
	pythonFile            = "/home/pi/Gasoleo/distance.py"
	alarmDir              = "/home/pi/Gasoleo/Alarms/"
)

func measure() (measurement float64) {
	cmd := exec.Command("python", pythonFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		errStr := string(stderr.Bytes())
		measurement, _ = strconv.ParseFloat(errStr, 64)
	}
	return
}

func main() {

	// First we calibrate
	samples := []float64{}
	for {
		thisSample := measure()
		if thisSample > 0 {
			samples = append(samples, thisSample)
			fmt.Println("Good sample:", thisSample)
		}
		if len(samples) > samplesForCalibration {
			break
		}
	}

	normalDistance := samples[len(samples)/2]
	fmt.Println("normal distance is", normalDistance)

	alarm := false

	for {
		thisMeasure := measure()
		if thisMeasure <= 0 {
			continue
		}
		if math.Abs(thisMeasure-normalDistance) > threshold {
			if !alarm {
				fmt.Println("Puerta abierta!!")
				now := time.Now()
				alarmFileName := fmt.Sprintf(alarmDir+"%0d-%d-%d-%d-%d-%d-%d-%d-%d",
					now.Unix(),
					now.Year(),
					now.Month(),
					now.Day(),
					now.Weekday(),
					now.Hour(),
					now.Minute(),
					now.Second(),
					now.Nanosecond(),
				)
				alarmFile, err := os.Create(alarmFileName)
				if err != nil {
					log.Println(err)
				}
				defer alarmFile.Close()
				alarmFile.WriteString(fmt.Sprintf("Puerta abierta!!\n"))

			}
			alarm = true
		} else {
			if alarm {
				fmt.Println("Puerta cerrada de nuevo")
			}
			alarm = false
		}
	}
}
