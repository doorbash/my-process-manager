package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func KillProcessByPID(pid int) error {
	switch runtime.GOOS {
	case "windows":
		log.Printf("Killing process %d (Windows)", pid)
		kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))
		err := kill.Run()
		if err != nil {
			log.Printf("Error killing process %d: %v", pid, err)
		} else {
			log.Printf("Process %d was killed", pid)
		}
		return err
	case "linux":
		log.Printf("Killing process %d (Linux)", pid)
		kill := exec.Command("kill", "-9", strconv.Itoa(pid))
		err := kill.Run()
		if err != nil {
			log.Printf("Error killing process %d: %v", pid, err)
		} else {
			log.Printf("Process %d was killed", pid)
		}
		return err
	}
	return errors.New("unsupported platform")
}

func KillProcessByName(name string) error {
	switch runtime.GOOS {
	case "windows":
		log.Printf("Killing process %s (Windows)", name)
		kill := exec.Command("TASKKILL", "/T", "/F", "/IM", name)
		err := kill.Run()
		if err != nil {
			log.Printf("Error killing process %s: %v", name, err)
		} else {
			log.Printf("Process %s was killed", name)
		}
		return err
	case "linux":
		log.Printf("Killing process %s (Linux)", name)
		kill := exec.Command("pkill", "-9", name)
		err := kill.Run()
		if err != nil {
			log.Printf("Error killing process %s: %v", name, err)
		} else {
			log.Printf("Process %s was killed", name)
		}
		return err
	}
	return errors.New("unsupported platform")
}
