package main

import (
	"log"
	"net/http"
	"os/user"
	"path/filepath"
	"sync"

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
		log.Fatalf("Error getting %s group:", groupName, err)
	}

	for _, id := range userGroupIds {
		if id == inputGroup.Gid {
			return true
		}
	}

	return false
}

func ListDevices() []string {
	matches, err := filepath.Glob("/dev/input/by-id/*-event-kbd")
	if err != nil {
		log.Fatalf("Failed to find input devices: %v", err)
	}

	return matches

}

func fire(c chan *evdev.InputEvent, config *Config) {
	for event := range c {
		keyName := evdev.CodeName(event.Type, event.Code)
		if !config.Keymap.IsMappedKey(keyName) {
			log.Printf("Ignoring unmapped key %s", keyName)
			continue
		}

		url := config.PublicTriggerUrl(keyName)
		resp, err := http.Post(url, "", nil)

		if err != nil {
			log.Printf("Error reading POST response:", err)
			return
		}

		log.Printf("POST to %s returned %d", url, resp.StatusCode)
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
