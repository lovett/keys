package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

type Sound struct {
	Path               string
	Format             beep.Format
	Buffer             beep.Buffer
	SpeakerInitialized bool
}

var (
	soundMap           = make(map[string]*Sound)
	speakerInitialized = false
)

func initializeSpeaker(s *Sound) {
	if speakerInitialized {
		return
	}

	if err := speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/30)); err != nil {
		fmt.Fprint(os.Stderr, err)
	} else {
		speakerInitialized = true
	}
}

func (s *Sound) Play() {
	initializeSpeaker(s)
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
		_, err := fmt.Scanln()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		soundMap[name].Play()
	}

	for {
		for name := range soundMap {
			prompt(name)
		}

		fmt.Println("\nTest complete. Press Control-c to exit, or ENTER to test again.")
		_, err := fmt.Scanln()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
