package sound

import (
	"keys/internal/asset"
	"keys/internal/config"
	"keys/internal/keymap"
	"log"
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
	SoundMap           = make(map[string]*Sound)
	speakerInitialized = false
)

func initializeSpeaker(s *Sound) {
	if speakerInitialized {
		return
	}

	if err := speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/30)); err != nil {
		log.Print(err)
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
	if len(SoundMap) != 0 {
		return
	}

	SoundMap["confirmation"] = SoundBuffer("assets/hero_simple-celebration-02.ogg")
	SoundMap["error"] = SoundBuffer("assets/alert_error-03.ogg")
	SoundMap["up"] = SoundBuffer("assets/state-change_confirm-up.ogg")
	SoundMap["down"] = SoundBuffer("assets/state-change_confirm-down.ogg")
	SoundMap["tap"] = SoundBuffer("assets/navigation_forward-selection-minimal.ogg")
}

func PlayConfirmationSound(cfg *config.Config, key *keymap.Key) {
	if !cfg.SoundAllowed {
		return
	}

	if !key.Confirmation {
		return
	}

	SoundMap["confirmation"].Play()
}

func PlayErrorSound(cfg *config.Config) {
	if !cfg.SoundAllowed {
		return
	}
	SoundMap["error"].Play()
}

func PlayToggleSound(cfg *config.Config, key *keymap.Key) {
	if !cfg.SoundAllowed {
		return
	}

	if !key.CanRoll() {
		return
	}

	if key.CommandIndex == 0 {
		SoundMap["down"].Play()
	} else {
		SoundMap["up"].Play()
	}
}

func PlayTapSound(cfg *config.Config, key *keymap.Key) {
	if !cfg.SoundAllowed {
		return
	}

	SoundMap["tap"].Play()
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
