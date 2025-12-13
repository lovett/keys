package device

import (
	"fmt"
	"keys/internal/config"
	"log"
	"net/http"
	"os/user"
	"path/filepath"
	"slices"
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
	result, err := canListen()
	if err != nil {
		log.Fatalf("Failed to confirm ability to listen: %s", err)
		return
	}

	if !result {
		log.Fatalf("Current user cannot listen to keyboard events")
		return
	}

	var wg sync.WaitGroup

	c := make(chan *EventPair)

	wg.Add(1)
	if cfg.Mode == config.KeyTestMode {
		go echo(c, cfg)
	} else {
		go trigger(c, cfg)
	}

	devices, err := ListDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		if cfg.Keymap.DesignatedKeyboard != "" && device != cfg.Keymap.DesignatedKeyboard {
			log.Printf("Skipping %s\n", filepath.Base(device))
			continue
		}

		cfg.KeyboardFound = true
		wg.Add(1)
		go listener(device, c, &wg, cfg)
	}

	if cfg.KeyboardFound {
		wg.Wait()
	} else {
		time.Sleep(time.Duration(10 * time.Second))
		Listen(cfg)
	}
}

func ListDevices() ([]string, error) {
	return filepath.Glob("/dev/input/by-id/*-event-kbd")
}

func echo(c chan *EventPair, cfg *config.Config) {
	for pair := range c {
		codeName := evdev.CodeName(pair.Event.Type, pair.Event.Code)
		format := "\nKey pressed on %s: %s \n"
		if cfg.Keymap.DesignatedKeyboard != "" && pair.Path == cfg.Keymap.DesignatedKeyboard {
			format = format[1:]
		}

		key := cfg.Keymap.FindKey(codeName)

		fmt.Printf(
			format,
			filepath.Base(pair.Path),
			key.PhysicalKey,
		)
	}
}

func trigger(c chan *EventPair, cfg *config.Config) {
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

	device, err := evdev.Open(path)

	if err != nil {
		log.Fatalf("Failed to open device %s: %v", filepath.Base(path), err)
	}
	defer func() {
		err := device.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	if path == cfg.Keymap.DesignatedKeyboard {
		err = device.Grab()
		if err != nil {
			log.Fatalf("Failed to grab device %s: %v", filepath.Base(path), err)
		}
		log.Printf("Grabbed %s for exclusive access\n", filepath.Base(path))
		defer func() {
			err := device.Ungrab()
			if err != nil {
				log.Fatalf("failed to ungrab deice %s: %v", filepath.Base(path), err)
			}
		}()
	} else {
		log.Printf("Listening for keyboard events on %s\n", filepath.Base(path))
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

func canListen() (bool, error) {
	u, err := user.Current()
	if err != nil {
		return false, err
	}

	g, err := user.LookupGroup("input")
	if err != nil {
		return false, err
	}

	uids, err := u.GroupIds()
	if err != nil {
		return false, err
	}

	if slices.Contains(uids, g.Gid) {
		return true, nil
	}

	return false, nil
}
