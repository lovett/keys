package cli

import (
	"flag"
	"fmt"
	"keys/internal/config"
	"keys/internal/keyboard"
	"keys/internal/server"
	"keys/internal/sound"
	"os"
	"strings"
)

var flagSet *flag.FlagSet

func Server(cfg *config.Config, args []string) int {
	flagSet = flag.NewFlagSet("server", flag.ExitOnError)

	port := flagSet.Int("port", 4004, "Server port")
	publicUrl := flagSet.String("url", "http://localhost:4004", "Server URL")
	inputs := flagSet.String("inputs", "browser,keyboard", "Where to listen for input")

	cfg.PublicUrl = *publicUrl

	flagSet.Usage = serverUsage
	flagSet.Parse(args)

	if strings.Contains(*inputs, "keyboard") {
		go keyboard.StartKeyboardListener(cfg)
	}

	if strings.Contains(*inputs, "browser") {
		sound.LoadSounds()
		server.StartServer(cfg, *port)
	}

	return 0
}

func serverUsage() {
	fmt.Fprint(os.Stderr, "Start the web server and listen for keyboard input.\n\n")

	fmt.Fprint(os.Stderr, "Options\n")
	flagSet.PrintDefaults()
}
