package cli

import (
	"flag"
	"fmt"
	"io"
	"keys/internal/asset"
	"keys/internal/config"
	"os"
	"path/filepath"
)

func Run(stdout io.Writer, stderr io.Writer) int {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(stderr, "Could not determine config dir. Giving up.")
		return 1
	}

	versionFlag := flag.Bool("version", false, "Application version")
	configFlag := flag.String("config", filepath.Join(configDir, "keys.ini"), "Configuration file")

	flag.Usage = topUsage
	flag.Parse()

	if *versionFlag {
		fmt.Fprintln(stdout, string(asset.ReadVersion()))
		return 0
	}

	cfg, err := config.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintln(stderr, "Could not parse config. Giving up.")
		return 1
	}

	command := flag.Arg(0)

	var args []string
	if flag.NArg() > 1 {
		args = flag.Args()[1:]
	}

	switch command {
	case "test":
		return Test(cfg, args)
	case "setup":
		return Setup(args)
	case "select":
		return Select(cfg, args)
	case "start":
		return Start(cfg, args)
	default:
		fmt.Fprintln(stderr, "Command not specified. Run keys --help for available commands.")
		return 1
	}
}

func topUsage() {
	fmt.Print(`Use a regular keyboard as a macro pad to run arbitrary commands headlessly.

Commands:
  select keyboard
        Choose which physical keyboard to use for input.
  start
        Launch the webserver and keyboard listener. Add --help for additional options.
  setup
        Install a startup service. Linux/systemd only.
  test key
        Test mode to see the name of a pressed key.
  test sound
        Test mode to see if sound output works.

Options:
`)

	flag.PrintDefaults()
}
