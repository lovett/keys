package sound

import (
	"testing"
)

func TestSpeakerInit(t *testing.T) {
	t.Cleanup(func() {
		clear(cache)
	})

	if len(sounds) == 0 {
		t.Fatal("Sound map is empty")
	}

	if len(cache) > 0 {
		t.Fatal("Sound cache is not empty")
	}
}

func TestPlay(t *testing.T) {
	clear(cache)

	if err := Play(Confirmation); err != nil {
		t.Fatal(err)
	}

	if len(cache) == 0 {
		t.Fatal("sound was not cached")
	}

	// Ensure speaker re-initialization is prevented.
	if err := Play(Confirmation); err != nil {
		t.Fatal(err)
	}
}
