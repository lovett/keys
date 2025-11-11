package cli

import (
	"errors"
	"fmt"
	"keys/internal/asset"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
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
	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	destinationDir := filepath.Join(home, ".config", "systemd", "user")
	destinationPath := filepath.Join(destinationDir, "keys.service")

	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	template := template.Must(template.ParseFS(asset.AssetFS, "assets/keys.service"))

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	templateVars := struct {
		ExecStart        string
		WorkingDirectory string
	}{
		ExecStart:        fmt.Sprintf("%s server", execPath),
		WorkingDirectory: home,
	}

	if err := template.Execute(f, templateVars); err != nil {
		return err
	}

	fmt.Printf("Wrote %s.\n\n", destinationPath)
	fmt.Println("To enable: systemctl --user enable keys.service")
	fmt.Println(" To start: systemctl --user start keys.service")

	return nil

}
