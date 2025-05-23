package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Byteloader - Bytecharge Programming Language Installer")

	osName := runtime.GOOS
	var url, fileName string

	switch osName {
	case "linux":
		url = "https://bytecharger.42web.io/download/Bytecharge"
		fileName = "Bytecharge"
	case "windows":
		url = "https://bytecharger.42web.io/download/Bytecharge_Win64"
		fileName = "Bytecharge.exe"
	case "darwin":
		url = "https://bytecharger.42web.io/download/Bytecharge"
		fileName = "Bytecharge"
	default:
		fmt.Println("Unsupported operating system!")
		return
	}

	fmt.Printf("Detected OS: %s\n", osName)
	fmt.Println("Downloading Bytecharge...")

	err := downloadFileWithProgress(url, fileName)
	if err != nil {
		fmt.Printf("Download error: %v\n", err)
		return
	}

	fmt.Printf("\nDownload completed: %s\n", fileName)

	if osName == "linux" {
		fmt.Println("Making Bytecharge globally accessible...")
		err := globalizeLinux(fileName)
		if err != nil {
			fmt.Printf("Error making global: %v\n", err)
			return
		}
		fmt.Println("Done! You can now run 'Bytecharge' from anywhere.")
	} else {
		fmt.Println("To use Bytecharge, run the downloaded file directly.")
	}
}

// Download file with progress bar and estimated time
func downloadFileWithProgress(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	total := resp.ContentLength
	var downloaded int64 = 0
	start := time.Now()

	buf := make([]byte, 32*1024) // 32 KB buffer
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := out.Write(buf[0:nr])
			if nw > 0 {
				downloaded += int64(nw)
			}
			if ew != nil {
				return ew
			}
			if nr != nw {
				return io.ErrShortWrite
			}
			printProgress(downloaded, total, start)
		}
		if er != nil {
			if er == io.EOF {
				break
			}
			return er
		}
	}
	fmt.Println() // line break after download
	return nil
}

// Display progress bar with estimated time
func printProgress(downloaded, total int64, start time.Time) {
	percent := float64(downloaded) / float64(total) * 100
	elapsed := time.Since(start).Seconds()
	speed := float64(downloaded) / 1024 / elapsed // KB/s

	var eta float64
	if speed > 0 {
		eta = float64(total-downloaded) / 1024 / speed
	}

	barWidth := 30
	progress := int(float64(barWidth) * percent / 100)

	bar := ""
	for i := 0; i < progress; i++ {
		bar += "="
	}
	for i := progress; i < barWidth; i++ {
		bar += " "
	}

	fmt.Printf("\r[%s] %.2f%% | %.2f KB/s | ETA: %.0fs", bar, percent, speed, eta)
}

// On Linux, move the file to /usr/local/bin to make it globally accessible
func globalizeLinux(fileName string) error {
	// Make it executable
	err := os.Chmod(fileName, 0755)
	if err != nil {
		return err
	}

	// Move to /usr/local/bin (may require sudo)
	cmd := exec.Command("mv", fileName, "/usr/local/bin/"+fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
