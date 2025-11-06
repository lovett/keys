package cli

import (
	"errors"
	"keys/internal/system"
	"os"
	"runtime"
)

func Setup(args []string) int {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = system.InstallSystemdUserService()
	default:
		err = errors.New("Not supported on this OS")
	}

	if err != nil {
		os.Stderr.WriteString(err.Error())
		return 1
	}

	return 0
}
