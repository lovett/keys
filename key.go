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
	Name         string
	Label        string
	Command      []string
	States       []string
	Toggle       bool
	CommandIndex int8
}

func NewKeyFromSection(s *ini.Section) *Key {
	commands := s.Key("command").ValueWithShadows()
	states := s.Key("state").ValueWithShadows()

	k := &Key{
		Name:    s.Name(),
		Label:   s.Key("label").MustString(""),
		Command: commands,
		States:  states,
		Toggle:  len(commands) > 1,
	}

	if k.Name == "" {
		return nil
	}

	if k.Label == "" {
		return nil
	}

	if k.CurrentCommand() == "" {
		return nil
	}

	return k
}

func (k *Key) CurrentCommand() string {
	return k.Command[k.CommandIndex]
}

func (k *Key) CurrentState() string {
	if !k.Toggle {
		return ""
	}

	return k.States[k.CommandIndex]
}

func (k *Key) updateCommandIndex() {
	if !k.Toggle {
		return
	}

	if k.CommandIndex == 0 {
		k.CommandIndex = 1
	} else {
		k.CommandIndex = 0
	}
}

func (k *Key) RunCommand() ([]byte, error) {
	log.Printf("Running command: %s", k.CurrentCommand())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	commandParts := strings.Split(k.CurrentCommand(), " ")

	var cmd *exec.Cmd

	switch len(commandParts) {
	case 0:
		return nil, errors.New("Command not specified")
	case 1:
		cmd = exec.CommandContext(ctx, commandParts[0])
	default:
		cmd = exec.CommandContext(ctx, commandParts[0], commandParts[1:]...)
	}

	k.updateCommandIndex()

	return cmd.Output()
}

func (k *Key) IsLockKey() bool {
	return k.CurrentCommand() == "lock" || k.CurrentCommand() == "unlock"
}
