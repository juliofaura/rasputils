package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	chart "github.com/wcharczuk/go-chart"
	"github.com/yryz/ds18b20"
)

type datapoint struct {
	Timestamp                                       int64
	Year, Month, Day, Weekday, Hour, Minute, Second int64
	Temperature                                     float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func weekday(wd time.Weekday) (s string) {
	switch wd {
	case 0:
		s = "Sun"
	case 1:
		s = "Mon"
	case 2:
		s = "Tue"
	case 3:
		s = "Wed"
	case 4:
		s = "Thu"
	case 5:
		s = "Fri"
	case 6:
		s = "Sat"
	default:
		s = "Error!"
	}
	return
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " <prefix>, where <prefix> is the full path + prexif that will be prepended to the .txt and .png files")
		os.Exit(1)
	}

	dataFile := os.Args[1] + ".txt"
	graphFile := os.Args[1] + ".png"

	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("sensor IDs: %v\n", sensors)

	csvFile, err := os.OpenFile(dataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	for _, sensor := range sensors {
		t, err := ds18b20.Temperature(sensor)
		if err == nil {
			// fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, t)
			now := time.Now()
			if _, err = fmt.Fprintf(csvFile, "%d,%04d,%02d,%02d,%01d,%02d,%02d,%02d,%.2f\n",
				now.Unix(), now.Year(), now.Month(), now.Day(), now.Weekday(), now.Hour(), now.Minute(), now.Second(), t); err != nil {
				panic(err)
			}
		}
	}

	csvFile.Close()

	csvFile, err = os.Open(dataFile)
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	var data []datapoint
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		var dataPoint datapoint
		dataPoint.Timestamp, err = strconv.ParseInt(line[0], 10, 64)
		check(err)
		dataPoint.Year, err = strconv.ParseInt(line[1], 10, 64)
		check(err)
		dataPoint.Month, err = strconv.ParseInt(line[2], 10, 64)
		check(err)
		dataPoint.Day, err = strconv.ParseInt(line[3], 10, 64)
		check(err)
		dataPoint.Weekday, err = strconv.ParseInt(line[4], 10, 64)
		check(err)
		dataPoint.Hour, err = strconv.ParseInt(line[5], 10, 64)
		check(err)
		dataPoint.Minute, err = strconv.ParseInt(line[6], 10, 64)
		check(err)
		dataPoint.Second, err = strconv.ParseInt(line[7], 10, 64)
		check(err)
		dataPoint.Temperature, err = strconv.ParseFloat(line[8], 64)
		check(err)
		data = append(data, dataPoint)
	}
	var XValues []float64
	var YValues []float64
	for _, v := range data {
		XValues = append(XValues, float64(v.Timestamp))
		YValues = append(YValues, v.Temperature)
	}

	LastX := XValues[len(XValues)-1]
	LastY := YValues[len(YValues)-1]

	labelColor := chart.ColorGreen

	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64) * 1e9
				typedDate := chart.TimeFromFloat64(typed)
				return fmt.Sprintf("%s %d:%d", weekday(typedDate.Weekday()), typedDate.Hour(), typedDate.Minute())
			},
			Style: chart.Style{
				TextRotationDegrees: 45,
			},
		},
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.1f", v.(float64))
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
					StrokeWidth: 5,
					DotWidth:    4,
				},
				XValues: XValues,
				YValues: YValues,
			},
			chart.AnnotationSeries{
				Annotations: []chart.Value2{
					{
						XValue: LastX,
						YValue: LastY,
						Label:  fmt.Sprintf("%.1f", LastY),
						Style: chart.Style{
							StrokeWidth: 10,
							FontSize:    chart.StyleTextDefaults().FontSize,
							StrokeColor: labelColor,
						},
					},
				},
			},
		},
		Title: "Room temperature vs time",
	}

	f, _ := os.Create(graphFile)
	defer f.Close()
	graph.Render(chart.PNG, f)
}
