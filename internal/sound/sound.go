package sound

import (
	"fmt"
	"keys/internal/asset"
	"keys/internal/config"
	"keys/internal/keymap"
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
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

func PlayConfirmationSound(cfg *config.Config) {
	if !cfg.SoundAllowed() {
		return
	}

	soundMap["confirmation"].Play()
}

func PlayErrorSound(cfg *config.Config) {
	if !cfg.SoundAllowed() {
		return
	}
	soundMap["error"].Play()
}

func PlayToggleSound(cfg *config.Config, key *keymap.Key) {
	if !cfg.SoundAllowed() {
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

func SoundBuffer(path string) *Sound {
	b, err := asset.AssetFS.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := vorbis.Decode(b)
	if err != nil {
		log.Fatal(err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	err = streamer.Close()
	if err != nil {
		log.Fatal(err)
	}

	return &Sound{
		Path:   path,
		Format: format,
		Buffer: *buffer,
	}
}
