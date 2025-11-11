package cli

import (
	"fmt"
	"keys/internal/config"
	"keys/internal/device"
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
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = cfg.Keymap.StoreKeyboard(*keyboard)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Updated %s\n", cfg.Keymap.Filename)
		os.Exit(0)
	}

	return 0
}
