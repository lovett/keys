package main

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Key struct {
	Name    string
	Label   string
	Command string
	Note    string
}

func NewKeyFromSection(s *ini.Section) *Key {
	k := &Key{
		Name:    s.Name(),
		Label:   s.Key("label").String(),
		Note:    s.Key("note").String(),
		Command: s.Key("command").String(),
	}

	if k.Name == "" {
		return nil
	}

	if k.Label == "" {
		return nil
	}

	if k.Command == "" {
		return nil
	}

	return k
}

func (k *Key) RunCommand() ([]byte, error) {
	log.Printf("Running command: %s", k.Command)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	commandParts := strings.Split(k.Command, " ")

	var cmd *exec.Cmd

	switch len(commandParts) {
	case 0:
		return nil, errors.New("Command not specified")
	case 1:
		cmd = exec.CommandContext(ctx, commandParts[0])
	default:
		cmd = exec.CommandContext(ctx, commandParts[0], commandParts[1:]...)
	}

	return cmd.Output()
}
