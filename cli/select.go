package cli

import (
	"bufio"
	"errors"
	"fmt"
	"keys/internal/config"
	"keys/internal/device"
	"log"
	"os"
	"strconv"
	"strings"
)

func Select(cfg *config.Config, args []string) int {
	var target string
	if len(args) > 0 {
		target = args[0]
	}

	switch target {
	case "keyboard":
		keyboard, err := prompt()

		if err != nil {
			log.Fatal(err)
		}

		cfg.Keymap.SetKeyboard(*keyboard)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Updated %s\n", cfg.Keymap.Filename)
		os.Exit(0)
	}

	return 0
}

func prompt() (*string, error) {
	devices, err := device.ListDevices()
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, errors.New("no keyboards found")
	}

	if len(devices) == 1 {
		fmt.Printf("Only one keyboard found (%s) so using that.\n", devices[0])
		return &devices[0], nil
	}

	fmt.Println("\nSelect a keyboard by number:")
	for i, device := range devices {
		fmt.Printf("%2d. %s\n", i+1, device)
	}

	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	index, err := strconv.Atoi(strings.TrimSuffix(answer, "\n"))
	if err != nil {
		return nil, errors.New("invalid input")
	}

	if index < 1 || index > len(devices) {
		return nil, errors.New("invalid selection")
	}

	return &devices[index-1], nil
}
