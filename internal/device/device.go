package device

import (
	"bufio"
	"errors"
	"fmt"
	"keys/internal/config"
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

func Listen(cfg *config.Config) {

	if !userInGroup("input") {
		log.Fatal("Current user doesn't belong to input group")
	}

	var wg sync.WaitGroup

	c := make(chan *EventPair)

	wg.Add(1)
	if cfg.Mode == config.KeyTestMode {
		go testFire(c, cfg)
	} else {
		go fire(c, cfg)
	}

	for _, device := range ListDevices() {
		if cfg.DesignatedKeyboard != "" && device != cfg.DesignatedKeyboard {
			log.Printf("Skipping %s\n", deviceName(device))
			continue
		}

		cfg.KeyboardFound = true
		wg.Add(1)
		go listener(device, c, &wg, cfg)
	}

	if cfg.KeyboardFound {
		wg.Wait()
	} else {
		sleepDuration := time.Duration(10 * time.Second)
		time.Sleep(sleepDuration)
		Listen(cfg)
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

func Prompt() (*string, error) {
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

func testFire(c chan *EventPair, cfg *config.Config) {
	for pair := range c {
		codeName := evdev.CodeName(pair.Event.Type, pair.Event.Code)
		format := "\nKey pressed on %s: %s \n"
		if cfg.DesignatedKeyboard != "" && pair.Path == cfg.DesignatedKeyboard {
			format = format[1:]
		}

		key := cfg.Keymap.FindKey(codeName)

		fmt.Printf(
			format,
			deviceName(pair.Path),
			key.PhysicalKey,
		)
	}
}

func fire(c chan *EventPair, cfg *config.Config) {
	var timer *time.Timer
	keyBuffer := []string{}

	callback := func() {
		url := fmt.Sprintf("%s/trigger/%s", cfg.PublicUrl, strings.Join(keyBuffer, ","))
		resp, err := http.Post(url, "", nil)

		if err != nil {
			log.Print("Error reading POST response:", err)
			return
		}

		log.Printf("POST to %s returned %d", url, resp.StatusCode)
		keyBuffer = keyBuffer[:0]
	}

	for pair := range c {
		// Because the mute key on an Apple keyboard showed up as "mute/min_interesting".
		codeName := strings.SplitN(evdev.CodeName(pair.Event.Type, pair.Event.Code), "/", 2)[0]

		if cfg.KeyboardLocked {
			log.Printf("Ignoring keypress of %s because keyboard is locked", codeName)
			continue
		}

		keyBuffer = append(keyBuffer, codeName)

		if timer != nil {
			timer.Stop()
		}

		bufferString := strings.Join(keyBuffer, "")
		if cfg.Keymap.IsPhysicalKeyPrefix(bufferString) {
			timer = time.AfterFunc(500*time.Millisecond, callback)
		} else {
			callback()
		}
	}
}

func listener(path string, c chan *EventPair, wg *sync.WaitGroup, cfg *config.Config) {
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

	if path == cfg.DesignatedKeyboard {
		err = device.Grab()
		if err != nil {
			log.Fatalf("Failed to grab device %s: %v", deviceName, err)
		}
		log.Printf("Grabbed %s for exclusive access\n", deviceName)
		defer func() {
			err := device.Ungrab()
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
			Listen(cfg)
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
