package keymap

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"time"

	"gopkg.in/ini.v1"
)

type Key struct {
	Name           string
	PhysicalKey    string
	Command        []string
	States         []string
	ShowOutput     bool
	CommandIndex   int8
	TimeoutSeconds int
	Confirmation   bool
	Row            string
}

func NewKeyFromSection(s *ini.Section) *Key {
	k := &Key{
		Name:           s.Name(),
		PhysicalKey:    s.Key("physical_key").MustString(""),
		Command:        s.Key("command").ValueWithShadows(),
		States:         s.Key("state").ValueWithShadows(),
		CommandIndex:   0,
		ShowOutput:     s.Key("output").MustBool(true),
		TimeoutSeconds: s.Key("timeout").MustInt(10),
		Confirmation:   s.Key("confirmation").MustBool(true),
		Row:            s.Key("row").MustString(""),
	}

	if k.Name == "" {
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
	if !k.IsToggle() {
		return ""
	}

	return k.States[k.CommandIndex]
}

func (k *Key) UpdateCommandIndex() {
	if !k.IsToggle() {
		return
	}

	k.CommandIndex += 1
	if k.CommandIndex >= int8(len(k.Command)) {
		k.CommandIndex = 0
	}
}

func (k *Key) RunCommand() ([]byte, error) {
	if len(k.CurrentCommand()) == 0 {
		return nil, errors.New("key command is empty, nothing to run")
	}

	log.Printf("Running command: %s", k.CurrentCommand())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(k.TimeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", k.CurrentCommand())

	k.UpdateCommandIndex()

	return cmd.Output()
}

func (k *Key) IsLockKey() bool {
	return k.CurrentCommand() == "lock" || k.CurrentCommand() == "unlock"
}

func (k *Key) IsToggle() bool {
	return len(k.Command) > 1
}
