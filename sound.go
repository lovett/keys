package main

import (
	"fmt"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

type Sound struct {
	Path   string
	Format beep.Format
	Buffer beep.Buffer
}

var soundMap = make(map[string]*Sound)

func (s *Sound) Play() {
	speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/30))

	buffer := s.Buffer.Streamer(0, s.Buffer.Len())
	speaker.Play(buffer)
}

func LoadSounds() {
	if len(soundMap) != 0 {
		return
	}

	soundMap["confirmation"] = SoundBuffer("assets/hero_simple-celebration-02.ogg")
	soundMap["error"] = SoundBuffer("assets/alert_error-03.ogg")
	soundMap["up"] = SoundBuffer("assets/state-change_confirm-up.ogg")
	soundMap["down"] = SoundBuffer("assets/state-change_confirm-down.ogg")
}

func PlayConfirmationSound(config *Config) {
	if !config.SoundAllowed() {
		return
	}

	soundMap["confirmation"].Play()
}

func PlayErrorSound(config *Config) {
	if !config.SoundAllowed() {
		return
	}
	soundMap["error"].Play()
}

func PlayToggleSound(config *Config, key *Key) {
	if !config.SoundAllowed() {
		return
	}

	if !key.Toggle {
		return
	}

	if key.CommandIndex == 0 {
		soundMap["down"].Play()
	} else {
		soundMap["up"].Play()
	}
}

func StartSoundTest() {
	LoadSounds()

	prompt := func(name string) {
		fmt.Printf("Press ENTER to play the %s sound ", name)
		fmt.Scanln()
		soundMap[name].Play()
	}

	for {
		for name := range soundMap {
			prompt(name)
		}

		fmt.Println("\nTest complete. Press Control-c to exit, or ENTER to test again.")
		fmt.Scanln()
	}
}
