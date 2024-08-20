package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var outputScale float64
	var webpQuality int

	fmt.Print("Output scale (ex: 0.5): ")
	fmt.Scanln(&outputScale)

	fmt.Print("WEBP Quality (ex: 80): ")
	fmt.Scanln(&webpQuality)

	files, err := getPNGFiles(".")
	if err != nil {
		fmt.Println("Error getting PNG files:", err)
		return
	}

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Use all available CPU cores
	numCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numCores)

	// Create a channel to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, numCores)

	for _, file := range files {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a semaphore slot

		go func(f string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the semaphore slot

			err := convertPNGToWebp(f, outputScale, webpQuality)
			if err != nil {
				fmt.Println("Error converting", f, ":", err)
			}
		}(file)
	}

	wg.Wait()
	fmt.Println("Conversion complete!")
}

func getPNGFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var pngFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
			pngFiles = append(pngFiles, file.Name())
		}
	}

	return pngFiles, nil
}

func convertPNGToWebp(file string, outputScale float64, webpQuality int) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", file,
		"-vf", fmt.Sprintf("scale=iw*%f:-1:flags=lanczos", outputScale),
		"-c:v", "libwebp",
		"-quality", strconv.Itoa(webpQuality),
		"-lossless", "0",
		"-map_metadata", "-1",
		fmt.Sprintf("%s.webp", file[:len(file)-4]),
	)

	// Redirect stderr to stdout to capture ffmpeg output
	cmd.Stderr = cmd.Stdout

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}

	fmt.Println("Converted:", file)
	return nil
}
