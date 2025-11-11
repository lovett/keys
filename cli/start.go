package cli

import (
	"flag"
	"fmt"
	"keys/internal/config"
	"keys/internal/device"
	"keys/internal/server"
	"keys/internal/sound"
	"os"
	"strings"
)

var flagSet *flag.FlagSet

func Start(cfg *config.Config, args []string) int {
	flagSet = flag.NewFlagSet("server", flag.ExitOnError)

	port := flagSet.Int("port", 4004, "Web server port")
	publicUrl := flagSet.String("url", "http://localhost:4004", "Web server URL")
	inputs := flagSet.String("inputs", "browser,keyboard", "Where to listen for input")

	cfg.PublicUrl = *publicUrl

	flagSet.Usage = startUsage
	flagSet.Parse(args)

	if strings.Contains(*inputs, "keyboard") {
		go device.Listen(cfg)
	}

	if strings.Contains(*inputs, "browser") {
		sound.LoadSounds()
		server.Serve(cfg, *port)
	}

	return 0
}

func startUsage() {
	fmt.Fprint(os.Stderr, "Launch the web server and listen for keyboard input.\n\n")

	fmt.Fprint(os.Stderr, "Options\n")
	flagSet.PrintDefaults()
}
