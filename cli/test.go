package cli

import (
	"fmt"
	"keys/internal/config"
	"keys/internal/device"
	"keys/internal/sound"
	"log"
)

func Test(cfg *config.Config, args []string) int {
	var testName string

	if len(args) > 0 {
		testName = args[0]
	}

	switch testName {
	case "sound":
		TestSound()
	case "key":
		TestKey(cfg)
	}

	return 0
}

func TestKey(cfg *config.Config) {
	log.Print("Press a key to see its details. Control-c to cancel.\n\n")
	cfg.EnableKeyTestMode()
	device.Listen(cfg)
}

func TestSound() {
	sounds := map[string]sound.Name{
		"Confirmation": sound.Confirmation,
		"Error":        sound.Error,
		"Tap":          sound.Tap,
		"Lock":         sound.Lock,
		"Unlock":       sound.Unlock,
	}

	for {
		for name, s := range sounds {
			fmt.Printf("Press ENTER to play the %s sound ", name)
			_, err := fmt.Scanln()
			if err != nil {
				log.Fatal(err)
			}

			if err := sound.Play(s); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("\nTest complete. Press Control-c to exit, or ENTER to test again.")
		_, err := fmt.Scanln()
		if err != nil {
			log.Fatal(err)
		}
	}
}
