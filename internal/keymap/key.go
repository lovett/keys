package keymap

import (
	"context"
	"log"
	"math"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Key struct {
	Name         string
	PhysicalKey  string
	Commands     []string
	States       []string
	CommandIndex uint8
	ShowOutput   bool
	Timeout      time.Duration
	Confirmation bool
	Row          string
}

func NewKeyFromSection(s *ini.Section, row string) *Key {
	k := &Key{
		Name:         s.Name(),
		PhysicalKey:  s.Key("physical_key").MustString(""),
		Commands:     s.Key("command").ValueWithShadows(),
		States:       s.Key("state").ValueWithShadows(),
		CommandIndex: 0,
		ShowOutput:   s.Key("output").MustBool(true),
		Timeout:      time.Duration(s.Key("timeout").MustFloat64(10.0)) * time.Second,
		Confirmation: s.Key("confirmation").MustBool(true),
		Row:          row,
	}

	if k.CurrentCommand() == "" {
		return nil
	}

	if len(k.Commands) > 1 && len(k.Commands) != len(k.States) {
		return nil
	}

	return k
}

func (k *Key) CanLock() bool {
	return k.CurrentCommand() == "lock" || k.CurrentCommand() == "unlock"
}

func (k *Key) CanToggle() bool {
	count := len(k.Commands)
	return count > 1 && count < math.MaxUint8
}

func (k *Key) MatchesCommand(command string) bool {
	lcCommand := strings.ToLower(command)

	for _, c := range k.Commands {
		if strings.Contains(strings.ToLower(c), lcCommand) {
			return true
		}
	}

	return false
}

func (k *Key) MatchesPhysicalKey(physicalKey string) bool {
	return strings.Contains(strings.ToLower(k.PhysicalKey), strings.ToLower(physicalKey))
}

func (k *Key) MatchesName(name string) bool {
	return strings.Contains(strings.ToLower(k.Name), strings.ToLower(name))
}

func (k *Key) State() string {
	if !k.CanToggle() {
		return ""
	}

	return k.States[k.CommandIndex]
}

func (k *Key) LastCommand() string {
	count := len(k.Commands)
	if count == 0 {
		return ""
	}

	return k.Commands[count-1]
}

func (k *Key) CurrentCommand() string {
	if len(k.Commands) == 0 {
		return ""
	}

	return k.Commands[k.CommandIndex]
}

func (k *Key) Toggle() {
	if !k.CanToggle() {
		return
	}

	if k.CurrentCommand() == k.LastCommand() {
		k.CommandIndex = 0
	} else {
		k.CommandIndex += 1
	}
}

func (k *Key) RunCommand() ([]byte, error) {
	log.Printf("Running command: %s", k.CurrentCommand())

	ctx, cancel := context.WithTimeout(context.Background(), k.Timeout)
	defer cancel()

	// #nosec [204] [-- The command being run intentionally comes from a user-supplied value.]
	cmd := exec.CommandContext(ctx, "sh", "-c", k.CurrentCommand())

	k.Toggle()

	return cmd.Output()
}
