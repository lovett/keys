package cli

import (
	"flag"
	"keys/internal/config"
	"keys/internal/device"
	"keys/internal/server"
	"keys/internal/sound"
	"log"
	"strings"
)

var flagSet *flag.FlagSet

func Start(cfg *config.Config, args []string) int {
	flagSet = flag.NewFlagSet("server", flag.ExitOnError)

	port := flagSet.Int("port", 4004, "Web server port")
	publicUrl := flagSet.String("url", "http://localhost:4004", "Web server URL")
	inputs := flagSet.String("inputs", "browser,keyboard", "Where to listen for input")

	flagSet.Usage = startUsage
	if err := flagSet.Parse(args); err != nil {
		log.Println(err)
	}

	cfg.PublicUrl = *publicUrl

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
	log.SetFlags(0)
	log.SetPrefix("")
	log.Print("Launch the web server and listen for keyboard input.\n\n")

	log.Print("Options\n")
	flagSet.PrintDefaults()
}
