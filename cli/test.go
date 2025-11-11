package cli

import (
	"fmt"
	"keys/internal/config"
	"keys/internal/device"
	"keys/internal/sound"
	"log"
	"os"
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
	log.Println("Running key test. Press a key to see its name.")
	cfg.EnableKeyTestMode()
	device.Listen(cfg)
}

func TestSound() {
	sound.LoadSounds()

	prompt := func(name string) {
		fmt.Printf("Press ENTER to play the %s sound ", name)
		_, err := fmt.Scanln()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		sound.SoundMap[name].Play()
	}

	for {
		for name := range sound.SoundMap {
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
