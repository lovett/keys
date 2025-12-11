package sound

import (
	"errors"
	"keys/internal/asset"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
)

type Name int

const (
	Confirmation Name = iota
	Error
	Up
	Down
	Lock
	Unlock
	Tap
)

var (
	sounds = make(map[Name]string)
	cache  = make(map[Name]*beep.Buffer)
)

func init() {
	sounds[Confirmation] = "assets/hero_simple-celebration-02.ogg"
	sounds[Error] = "assets/alert_error-03.ogg"
	sounds[Tap] = "assets/navigation_forward-selection-minimal.ogg"
	sounds[Lock] = "assets/ui_lock.ogg"
	sounds[Unlock] = "assets/ui_unlock.ogg"
}

func load(name Name) error {
	if _, found := cache[name]; found {
		return nil
	}

	path, found := sounds[name]
	if !found {
		return errors.New("unknown sound")
	}

	b, err := asset.AssetFS.Open(path)
	if err != nil {
		return err
	}

	streamer, format, err := vorbis.Decode(b)
	if err != nil {
		return err
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	err = streamer.Close()
	if err != nil {
		return err
	}

	if len(cache) == 0 {
		if err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30)); err != nil {
			return err
		}
	}

	cache[name] = buffer
	return nil
}

func Play(name Name) error {
	if buffer, found := cache[name]; found {
		stream := buffer.Streamer(0, buffer.Len())
		speaker.Play(stream)
		return nil
	}

	if err := load(name); err != nil {
		return err
	}
	return Play(name)

}
