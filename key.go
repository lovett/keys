package main

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Key struct {
	Name    string
	Label   string
	Command []string
	Note    string
}

func NewKeyFromSection(s *ini.Section) *Key {
	k := &Key{
		Name:  s.Name(),
		Label: s.Key("label").String(),
		Note:  s.Key("note").String(),
	}

	if k.Name == "" {
		return nil
	}

	if k.Label == "" {
		return nil
	}

	command := s.Key("command").String()
	if strings.HasPrefix(command, "get") {
		// HTTP GET
	} else if strings.HasPrefix(command, "post") {
		// HTTP POST
	} else {
		k.Command = strings.Split(command, " ")
	}

	return k
}

func (k *Key) RunCommand() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var cmd *exec.Cmd
	switch len(k.Command) {
	case 0:
		cmd = nil
	case 1:
		cmd = exec.CommandContext(ctx, k.Command[0])
	default:
		cmd = exec.CommandContext(ctx, k.Command[0], k.Command[1:]...)
	}

	if cmd == nil {
		return nil, errors.New("Command not specified")
	}

	return cmd.Output()
}
