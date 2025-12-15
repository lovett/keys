package device

import (
	"fmt"
	"keys/internal/config"
	"keys/internal/keymap"
	"log"
	"net/http"
	"net/url"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/holoplot/go-evdev"
)

type DeviceEvent struct {
	DevicePath string
	Event      *evdev.InputEvent
}

func Listen(cfg *config.Config) {
	var (
		u   *user.User
		g   *user.Group
		err error
	)

	if u, err = user.Current(); err != nil {
		log.Fatal(err)
	}

	if g, err = user.LookupGroup("input"); err != nil {
		log.Fatal(err)
	}

	if !canListen(u, g) {
		log.Fatal("current user cannot listen to keyboard events")
		return
	}

	c := make(chan *DeviceEvent)

	var wg sync.WaitGroup
	wg.Add(1)
	go worker(c, cfg)

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

func worker(deviceEvents <-chan *DeviceEvent, cfg *config.Config) {
	var timer *time.Timer
	keyBuffer := []string{}

	callback := func() {
		trigger(keyBuffer, cfg)
		keyBuffer = keyBuffer[:0]
	}

	for deviceEvent := range deviceEvents {
		if cfg.Mode == config.KeyTestMode {
			echo(deviceEvent, cfg)
			return
		}

		codeName := evdev.CodeName(deviceEvent.Event.Type, deviceEvent.Event.Code)

		if cfg.KeyboardLocked {
			log.Printf("Ignoring keypress of %s because the keyboard is locked", codeName)
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

func echo(deviceEvent *DeviceEvent, cfg *config.Config) string {
	codeName := evdev.CodeName(deviceEvent.Event.Type, deviceEvent.Event.Code)
	translatedName := keymap.Translate(codeName)

	key := cfg.Keymap.FindKey(translatedName)
	var mappedKey string
	if key != nil {
		mappedKey = key.Name
	} else {
		mappedKey = "none"
	}

	return fmt.Sprintf(`

Code: %s
From: %s
Translated to: %s
Mapped to: %s
`, codeName, deviceEvent.DevicePath, translatedName, mappedKey)
}

func trigger(keyBuffer []string, cfg *config.Config) {
	key := strings.Join(keyBuffer, ",")
	url := fmt.Sprintf("%s/trigger/%s", cfg.PublicUrl, url.PathEscape(key))
	resp, err := http.Post(url, "", nil)

	if err != nil {
		log.Print("Error reading POST response:", err)
		return
	}

	log.Printf("POST to %s returned %d", url, resp.StatusCode)
}

func open(path string, c chan *DeviceEvent, wg *sync.WaitGroup, cfg *config.Config) {
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
			c <- &DeviceEvent{path, event}
		}
	}
}

func canListen(u *user.User, group *user.Group) bool {
	if uids, err := u.GroupIds(); err == nil {
		return slices.Contains(uids, group.Gid)
	}
	return false
}
