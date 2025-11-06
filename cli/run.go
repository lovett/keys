package cli

import (
	"flag"
	"keys/internal/config"
	"os"
	"path/filepath"
)

func Run(appVersion string) int {
	configDir, err := os.UserConfigDir()
	if err != nil {
		os.Stderr.WriteString("Could not determine config dir. Giving up.")
		return 1
	}

	versionFlag := flag.Bool("version", false, "Application version")
	configFlag := flag.String("config", filepath.Join(configDir, "keys.ini"), "Configuration file")

	flag.Usage = topUsage
	flag.Parse()

	cfg := config.NewConfig(*configFlag, appVersion)

	if *versionFlag {
		os.Stdout.WriteString(appVersion)
		return 0
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
	case "server":
		return Server(cfg, os.Args[2:])
	default:
		topUsage()
		return 1
	}
}

func topUsage() {
	os.Stderr.WriteString(`Use a regular keyboard as a macro pad to run arbitrary commands headlessly.

Commands:
  select keyboard
        Choose which keyboard to use for input
  server
        Start the webserver and keyboard listener
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
