package cli

import (
	"flag"
	"keys/internal/asset"
	"keys/internal/config"
	"os"
	"path/filepath"
)

func Run() int {
	configDir, err := os.UserConfigDir()
	if err != nil {
		os.Stderr.WriteString("Could not determine config dir. Giving up.\n")
		return 1
	}

	versionFlag := flag.Bool("version", false, "Application version")
	configFlag := flag.String("config", filepath.Join(configDir, "keys.ini"), "Configuration file")

	flag.Usage = topUsage
	flag.Parse()

	if *versionFlag {
		os.Stdout.Write(asset.ReadVersion())
		return 0
	}

	cfg, err := config.NewConfig(*configFlag)
	if err != nil {
		os.Stderr.WriteString("Could not parse config. Giving up.")
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
		os.Stderr.WriteString("Command not specified. Run keys --help for available commands.\n")
		return 1
	}
}

func topUsage() {
	os.Stderr.WriteString(`Use a regular keyboard as a macro pad to run arbitrary commands headlessly.

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
