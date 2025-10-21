package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/holoplot/go-evdev"
)

func StartKeyboardListener(config *Config) {

	if !userInGroup("input") {
		log.Fatal("Current user doesn't belong to input group")
	}

	var wg sync.WaitGroup

	c := make(chan *evdev.InputEvent)

	wg.Add(1)
	go fire(c, config)

	for _, device := range ListDevices() {
		wg.Add(1)
		go listen(device, c, &wg)
	}

	wg.Wait()
}

func userInGroup(groupName string) bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Error getting current user:", err)
	}

	userGroupIds, err := currentUser.GroupIds()
	if err != nil {
		log.Fatal("Error getting user groups:", err)
	}

	inputGroup, err := user.LookupGroup(groupName)
	if err != nil {
		log.Fatalf("Error getting %s group: %s", groupName, err)
	}

	for _, id := range userGroupIds {
		if id == inputGroup.Gid {
			return true
		}
	}

	return false
}

func PromptForKeyboard() (*string, error) {
	devices := ListDevices()

	if len(devices) == 0 {
		return nil, errors.New("No keyboards found.")
	}

	if len(devices) == 1 {
		fmt.Printf("Only one keyboard found (%s) so using that.\n", devices[0])
		return &devices[0], nil
	}

	fmt.Println("\nSelect a keyboard by number:")
	for i, device := range devices {
		fmt.Printf("%2d. %s\n", i+1, device)
	}

	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	index, err := strconv.Atoi(strings.TrimSuffix(answer, "\n"))
	if err != nil {
		return nil, errors.New("Invalid input.")
	}

	if index < 1 || index > len(devices) {
		return nil, errors.New("Invalid selection.")
	}

	return &devices[index-1], nil
}

func ListDevices() []string {
	matches, err := filepath.Glob("/dev/input/by-id/*-event-kbd")
	if err != nil {
		log.Fatalf("Failed to find input devices: %v", err)
	}

	return matches

}

func fire(c chan *evdev.InputEvent, config *Config) {
	var timer *time.Timer

	keyBuffer := []string{}

	action := func() {
		url := config.PublicTriggerUrl(strings.Join(keyBuffer, ","))
		resp, err := http.Post(url, "", nil)

		if err != nil {
			log.Printf("Error reading POST response:", err)
			return
		}

		log.Printf("POST to %s returned %d", url, resp.StatusCode)
		keyBuffer = keyBuffer[:0]
	}

	for event := range c {
		keyBuffer = append(keyBuffer, evdev.CodeName(event.Type, event.Code))

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(500*time.Millisecond, action)
	}
}

func listen(path string, c chan *evdev.InputEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	device, err := evdev.Open(path)

	if err != nil {
		log.Fatalf("Failed to open device %s: %v", path, err)
	}
	defer device.Close()

	log.Printf("Listening for keyboard events on %s\n", device.Path())

	for {
		event, err := device.ReadOne()
		if err != nil {
			log.Fatalf("failed to read event: %v", err)
		}

		if event.Type == evdev.EV_KEY && event.Value == 1 {
			c <- event
		}
	}
}
