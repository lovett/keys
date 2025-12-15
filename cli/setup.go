package cli

import (
	"errors"
	"fmt"
	"keys/internal/asset"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"text/template"
)

func Setup(args []string) int {
	log.SetFlags(0)
	log.SetPrefix("")

	var err error

	switch runtime.GOOS {
	case "linux":
		err = installSystemdUserService()
	default:
		err = errors.New("not supported on this OS")
	}

	if err != nil {
		log.Println(err.Error())
		return 1
	}

	return 0
}

func installSystemdUserService() error {
	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	destination := filepath.Clean(filepath.Join(home, ".config", "systemd", "user", "keys.service"))

	if err := os.MkdirAll(path.Dir(destination), 0750); err != nil {
		return err
	}

	f, err := os.Create(destination)
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
		ExecStart:        fmt.Sprintf("%s start", execPath),
		WorkingDirectory: home,
	}

	if err := template.Execute(f, templateVars); err != nil {
		return err
	}

	fmt.Printf("Wrote %s.\n\n", destination)
	fmt.Println("To enable: systemctl --user enable keys.service")
	fmt.Println(" To start: systemctl --user start keys.service")

	return nil

}
