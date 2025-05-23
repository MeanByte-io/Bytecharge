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
		handleLinux(fileName)
	} else if osName == "windows" {
		handleWindows(fileName)
	} else {
		fmt.Println("To use Bytecharge, run the downloaded file directly.")
	}
}

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

	buf := make([]byte, 32*1024)
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
	fmt.Println()
	return nil
}

func printProgress(downloaded, total int64, start time.Time) {
	percent := float64(downloaded) / float64(total) * 100
	elapsed := time.Since(start).Seconds()
	speed := float64(downloaded) / 1024 / elapsed

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

func handleLinux(fileName string) {
	fmt.Println("Checking for C compiler on Linux...")
	if !checkCCompiler() {
		fmt.Println("No C compiler found. Attempting to install one...")
		if err := installCCompiler(); err != nil {
			fmt.Printf("Failed to install C compiler: %v\n", err)
			return
		}
	} else {
		fmt.Println("C compiler found.")
	}

	fmt.Println("Making Bytecharge globally accessible...")
	err := globalizeLinux(fileName)
	if err != nil {
		fmt.Printf("Error making global: %v\n", err)
		return
	}
	fmt.Println("Done! You can now run 'Bytecharge' from anywhere.")
}

func globalizeLinux(fileName string) error {
	err := os.Chmod(fileName, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("mv", fileName, "/usr/local/bin/"+fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func handleWindows(fileName string) {
	fmt.Println("Checking for C compiler on Windows...")

	if !checkCCompilerWindows() {
		fmt.Println("No C compiler found on your system.")
		fmt.Println("Recommended: Install MinGW (https://www.mingw-w64.org/downloads/)")
		fmt.Println("or use Visual Studio with C++ tools.")
	} else {
		fmt.Println("C compiler found.")
	}

	fmt.Printf("You can run '%s' directly or add its directory to your PATH.\n", fileName)
}

func checkCCompiler() bool {
	commands := []string{"gcc", "clang", "cc"}

	for _, cmd := range commands {
		_, err := exec.LookPath(cmd)
		if err == nil {
			return true
		}
	}
	return false
}

func installCCompiler() error {
	pkgManagers := map[string][]string{
		"apt":        {"sudo", "apt", "update"},
		"aptInstall": {"sudo", "apt", "install", "-y", "build-essential"},
		"dnf":        {"sudo", "dnf", "install", "-y", "gcc"},
		"yum":        {"sudo", "yum", "install", "-y", "gcc"},
		"pacman":     {"sudo", "pacman", "-Syu", "--noconfirm", "base-devel"},
		"zypper":     {"sudo", "zypper", "install", "-y", "gcc"},
	}

	for pm := range pkgManagers {
		if _, err := exec.LookPath(pm); err == nil {
			fmt.Printf("Detected package manager: %s\n", pm)
			if pm == "apt" {
				if err := runCommand(pkgManagers["apt"]); err != nil {
					return err
				}
				return runCommand(pkgManagers["aptInstall"])
			}
			return runCommand(pkgManagers[pm])
		}
	}

	return fmt.Errorf("no supported package manager found, please install a C compiler manually")
}

func runCommand(cmdArgs []string) error {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func checkCCompilerWindows() bool {
	commands := []string{"gcc", "cl"}

	for _, cmd := range commands {
		_, err := exec.LookPath(cmd)
		if err == nil {
			return true
		}
	}
	return false
}
