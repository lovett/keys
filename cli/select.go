package cli

import (
	"fmt"
	"keys/internal/config"
	"keys/internal/device"
	"log"
	"os"
)

func Select(cfg *config.Config, args []string) int {
	var target string
	if len(args) > 0 {
		target = args[0]
	}

	switch target {
	case "keyboard":
		keyboard, err := device.Prompt()

		if err != nil {
			log.Fatal(err)
		}

		err = cfg.Keymap.StoreKeyboard(*keyboard)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Updated %s\n", cfg.Keymap.Filename)
		os.Exit(0)
	}

	return 0
}
