package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " <working dir>, where <working dir> is the directory where the .txt and .png files will be placed")
		os.Exit(1)
	}

	workingDir := os.Args[1] + "/"

	cmd := exec.Command("python", workingDir+"distance.py")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	errStr := string(stderr.Bytes())
	fmt.Printf("%s", errStr)

	// if err != nil {
	// 	panic(err)
	// }
	// thisDistance, err := strconv.ParseFloat(string(output), 64)
	// if err != nil {
	// 	panic(err)
	// }
	// distances := []float64{}
	// distances = append(distances, thisDistance)

}
