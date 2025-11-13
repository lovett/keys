package cli

import (
	"flag"
	"keys/internal/asset"
	"keys/internal/config"
	"log"
	"os"
	"path/filepath"
)

func Run() int {
	log.SetFlags(0)
	log.SetPrefix("")

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Could not determine config dir. Giving up.")
		return 1
	}

	versionFlag := flag.Bool("version", false, "Application version")
	configFlag := flag.String("config", filepath.Join(configDir, "keys.ini"), "Configuration file")

	flag.Usage = topUsage
	flag.Parse()

	if *versionFlag {
		log.Println(string(asset.ReadVersion()))
		return 0
	}

	cfg, err := config.NewConfig(*configFlag)
	if err != nil {
		log.Println("Could not parse config. Giving up.")
		return 1
	}

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "test":
		return Test(cfg, os.Args[2:])
	case "setup":
		return Setup(os.Args[2:])
	case "select":
		return Select(cfg, os.Args[2:])
	case "start":
		return Start(cfg, os.Args[2:])
	default:
		log.Println("Command not specified. Run keys --help for available commands.")
		return 1
	}
}

func topUsage() {
	log.Print(`Use a regular keyboard as a macro pad to run arbitrary commands headlessly.

Commands:
  select keyboard
        Choose which keyboard to use for input
  start
        Launch the webserver and keyboard listener
  setup
        Install a startup service. Linux/systemd only
  test key
        Test mode to see the name of a pressed key
  test sound
        Test mode to see if sound output works

Options:
`)

	flag.PrintDefaults()
}
