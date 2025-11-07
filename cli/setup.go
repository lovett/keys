package cli

import (
	"errors"
	"fmt"
	"keys/internal/asset"
	"os"
	"path/filepath"
	"runtime"
)

func Setup(args []string) int {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = installSystemdUserService()
	default:
		err = errors.New("Not supported on this OS")
	}

	if err != nil {
		os.Stderr.WriteString(err.Error())
		return 1
	}

	return 0
}

func installSystemdUserService() error {
	asset, err := asset.ReadAsset("assets/keys.service")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	destinationDir := filepath.Join(home, ".config", "systemd", "user")
	destinationPath := filepath.Join(destinationDir, filepath.Base(asset.Path))

	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(destinationPath, asset.Bytes, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(home, "Documents", "Keys"), os.ModePerm); err != nil {
		return err
	}

	fmt.Printf("Wrote %s.\n\n", destinationPath)
	fmt.Println("To enable: systemctl --user enable keys.service")
	fmt.Println(" To start: systemctl --user start keys.service")

	return nil

}
