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

type deviceEvent struct {
	DevicePath string
	Event      *evdev.InputEvent
}

func Listen(cfg *config.Config) {
	result, err := canListen()
	if err != nil {
		log.Fatal(err)
		return
	}

	if !result {
		log.Fatalf("Current user cannot listen to keyboard events")
		return
	}

	c := make(chan *deviceEvent)

	var wg sync.WaitGroup
	wg.Add(1)

	if cfg.Mode == config.KeyTestMode {
		go echo(c, cfg)
	} else {
		go trigger(c, cfg)
	}

	devices, err := ListKeyboards()
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
		go open(device, c, &wg, cfg)
	}

	if cfg.KeyboardFound {
		wg.Wait()
	} else {
		time.Sleep(time.Duration(10 * time.Second))
		Listen(cfg)
	}
}

func ListKeyboards() ([]string, error) {
	return filepath.Glob("/dev/input/by-id/*-event-kbd")
}

func echo(pairs <-chan *deviceEvent, cfg *config.Config) {
	for pair := range pairs {
		codeName := evdev.CodeName(pair.Event.Type, pair.Event.Code)
		translatedName := cfg.Keymap.Translate(codeName)

		fmt.Println("")
		fmt.Printf("Code: %s\n", codeName)
		fmt.Printf("From: %s\n", pair.DevicePath)
		fmt.Printf("Translates to: %s\n", translatedName)

		key := cfg.Keymap.FindKey(translatedName)
		if key != nil {
			fmt.Printf("Mapped to: %s\n", key.Name)
		}
		fmt.Print("\n")
	}
}

func trigger(c chan *deviceEvent, cfg *config.Config) {
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

func open(path string, c chan *deviceEvent, wg *sync.WaitGroup, cfg *config.Config) {
	defer wg.Done()

	deviceName := filepath.Base(path)

	device, err := evdev.Open(path)
	if err != nil {
		log.Fatalf("unable to open %s: %s", deviceName, err)
	}

	defer func() {
		err := device.Close()
		if err != nil {
			log.Fatalf("unable to close %s: %s", deviceName, err)
		}
	}()

	if path == cfg.Keymap.DesignatedKeyboard {
		err = device.Grab()
		if err != nil {
			log.Fatalf("unable to grab %s: %s", deviceName, err)
		}

		log.Printf("Grabbed %s for exclusive access\n", deviceName)

		defer func() {
			err := device.Ungrab()
			if err != nil {
				log.Fatalf("unable to ungrab %s: %s", deviceName, err)
			}
		}()
	}

	for {
		event, err := device.ReadOne()
		if err != nil {
			log.Fatalf("unable to read input: %s", err)
			Listen(cfg)
			return
		}

		if event.Type == evdev.EV_KEY && event.Value == 0 {
			c <- &deviceEvent{path, event}
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
