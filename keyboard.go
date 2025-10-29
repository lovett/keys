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

type EventPair struct {
	Event *evdev.InputEvent
	Path  string
}

func StartKeyTest(config *Config) {
	log.Println("Running key test. Press a key to see its name.")
	StartKeyboardListener(config)
}

func StartKeyboardListener(config *Config) {

	if !userInGroup("input") {
		log.Fatal("Current user doesn't belong to input group")
	}

	var wg sync.WaitGroup

	c := make(chan *EventPair)

	wg.Add(1)
	if config.Mode == KeyTestMode {
		go testFire(c, config)
	} else {
		go fire(c, config)
	}

	for _, device := range ListDevices() {
		if config.DesignatedKeyboard() != "" && device != config.DesignatedKeyboard() {
			log.Printf("Skipping %s\n", deviceName(device))
			continue
		}

		config.KeyboardFound = true
		wg.Add(1)
		go listen(device, c, &wg, config)
	}

	if config.KeyboardFound {
		wg.Wait()
	} else {
		sleepDuration := time.Duration(10 * time.Second)
		time.Sleep(sleepDuration)
		StartKeyboardListener(config)
	}
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
		return nil, errors.New("no keyboards found")
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
		return nil, errors.New("invalid input")
	}

	if index < 1 || index > len(devices) {
		return nil, errors.New("invalid selection")
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

func testFire(c chan *EventPair, config *Config) {
	for pair := range c {
		key := evdev.CodeName(pair.Event.Type, pair.Event.Code)
		format := "\nKey pressed on %s: %s \n"
		if config.DesignatedKeyboard() != "" && pair.Path == config.DesignatedKeyboard() {
			format = format[1:]
		}

		fmt.Printf(
			format,
			deviceName(pair.Path),
			config.Keymap.KeyNameToSectionName(key),
		)
	}
}

func fire(c chan *EventPair, config *Config) {
	var timer *time.Timer
	keyBuffer := []string{}

	callback := func() {
		url := config.PublicTriggerUrl(strings.Join(keyBuffer, ","))
		resp, err := http.Post(url, "", nil)

		if err != nil {
			log.Print("Error reading POST response:", err)
			return
		}

		log.Printf("POST to %s returned %d", url, resp.StatusCode)
		keyBuffer = keyBuffer[:0]
	}

	for pair := range c {
		codeName := evdev.CodeName(pair.Event.Type, pair.Event.Code)
		if config.KeyboardLocked {
			log.Printf("Ignoring keypress of %s because keyboard is locked", codeName)
			continue
		}

		keyBuffer = append(keyBuffer, codeName)

		if timer != nil {
			timer.Stop()
		}

		if config.Keymap.IsPrefix(keyBuffer) {
			timer = time.AfterFunc(500*time.Millisecond, callback)
		} else {
			callback()
		}
	}
}

func listen(path string, c chan *EventPair, wg *sync.WaitGroup, config *Config) {
	defer wg.Done()

	deviceName := deviceName(path)
	device, err := evdev.Open(path)

	if err != nil {
		log.Fatalf("Failed to open device %s: %v", deviceName, err)
	}
	defer func() {
		err := device.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	if path == config.DesignatedKeyboard() {
		err = device.Grab()
		if err != nil {
			log.Fatalf("Failed to grab device %s: %v", deviceName, err)
		}
		log.Printf("Grabbed %s for exclusive access\n", deviceName)
		defer func() {
			err := device.Ungrab();
			if err != nil {
				log.Fatalf("failed to ungrab deice %s: %v", deviceName, err)
			}
		}()
	} else {
		log.Printf("Listening for keyboard events on %s\n", deviceName)
	}

	for {
		event, err := device.ReadOne()
		if err != nil {
			log.Printf("Failed to read keyboard input: %s", err)
			StartKeyboardListener(config)
			return
		}

		if event.Type == evdev.EV_KEY && event.Value == 0 {
			c <- &EventPair{event, path}
		}
	}
}

func deviceName(path string) string {
	return filepath.Base(path)
}
