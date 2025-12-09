package keymap

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/ini.v1"
)

func resetLogger() {
	log.SetOutput(os.Stdout)
}

func loadKeyFromFixture(t *testing.T, filename string) *Key {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(wd, "../../testdata", filename)

	options := ini.LoadOptions{
		SkipUnrecognizableLines: true,
		AllowShadows:            true,
	}

	ini, err := ini.LoadSources(options, path)
	if err != nil {
		t.Fatal(err)
	}

	s, err := ini.GetSection("test")
	if err != nil {
		t.Fatal(err)
	}

	return NewKeyFromSection(s)
}

func TestKey(t *testing.T) {
	key := loadKeyFromFixture(t, "key-single.ini")

	if key == nil {
		t.Fatal("Valid key was rejected")
	}

	if key.CanRoll() {
		t.Fatal("Single-command key cannot roll")
	}

	if key.CanLock() {
		t.Fatal("Misidentified lock key")
	}

	if key.State() != "" {
		t.Fatal("Single-command key should not have a state")
	}
}

func TestCommandRequired(t *testing.T) {
	key := loadKeyFromFixture(t, "key-invalid.ini")

	if key != nil {
		t.Fatal("Key without command was not rejected")
	}
}

func TestRoll(t *testing.T) {
	key := loadKeyFromFixture(t, "key-roll.ini")

	if key.CommandIndex != 0 {
		t.Fatal("Command index did not start at zero")
	}

	if !key.CanRoll() {
		t.Fatal("Multi-command key should be able to roll")
	}

	if key.CanLock() {
		t.Fatal("Misidentified lock key")
	}

	if key.CurrentCommand() != "echo hello" {
		t.Errorf("Unexpected first command: %s", key.CurrentCommand())
	}

	if key.State() != "state1" {
		t.Fatalf("Unexpected first state: %s", key.State())
	}

	key.RollForward()
	if key.CurrentCommand() != "echo hello 2" {
		t.Errorf("Unexpected second command: %s", key.CurrentCommand())
	}

	if key.State() != "state2" {
		t.Fatalf("Unexpected second state: %s", key.State())
	}

	key.RollForward()
	if key.CurrentCommand() != "echo hello 3" {
		t.Errorf("Unexpected second command: %s", key.CurrentCommand())
	}

	if key.State() != "state3" {
		t.Fatalf("Unexpected third state: %s", key.State())
	}

	key.RollForward()
	if key.CurrentCommand() != "echo hello" {
		t.Errorf("Failed to return to first command")
	}

	if key.State() != "state1" {
		t.Fatal("Failed to return to first state")
	}
}

func TestRollCommandStateMismatch(t *testing.T) {
	key := loadKeyFromFixture(t, "key-roll-invalid.ini")

	if key != nil {
		t.Fatal("Mismatch between commands and states was not caught")
	}
}

func TestLock(t *testing.T) {
	tests := []struct {
		fixture string
	}{
		{fixture: "key-single-lock.ini"},
		{fixture: "key-single-unlock.ini"},
		{fixture: "key-roll-lock.ini"},
	}
	for _, tt := range tests {
		key := loadKeyFromFixture(t, tt.fixture)
		if !key.CanLock() {
			t.Fatal("Misidentified lock key")
		}
	}
}

func TestMatchesCommand(t *testing.T) {
	key := loadKeyFromFixture(t, "key-roll.ini")

	if !key.MatchesCommand("hello") {
		t.Fatal("False negative during command match")
	}

	if !key.MatchesCommand("hello 3") {
		t.Fatal("Failed to exact match command")
	}

	if key.MatchesCommand("x") {
		t.Fatal("False positive during command match")
	}
}

func TestMatchesPhysicalKey(t *testing.T) {
	key := loadKeyFromFixture(t, "key-single.ini")

	if !key.MatchesPhysicalKey("i") {
		t.Fatal("False negative during physical key match")
	}

	if !key.MatchesPhysicalKey("hi") {
		t.Fatal("Failed to exact match physical key")
	}

	if key.MatchesCommand("x") {
		t.Fatal("False positive during physical key match")
	}
}

func TestMatchesName(t *testing.T) {
	key := loadKeyFromFixture(t, "key-single.ini")

	if !key.MatchesName("est") {
		t.Fatal("False negative during name match")
	}

	if !key.MatchesName("test") {
		t.Fatal("Failed to exact match name")
	}

	if key.MatchesName("x") {
		t.Fatal("False positive during name match")
	}
}

func TestCommandTimeout(t *testing.T) {
	t.Cleanup(resetLogger)
	log.SetOutput(io.Discard)

	key := loadKeyFromFixture(t, "key-timeout.ini")

	_, err := key.RunCommand()

	if err == nil {
		t.Fatalf("Command did not time out")
	}
}

func TestCommand(t *testing.T) {
	t.Cleanup(resetLogger)
	log.SetOutput(io.Discard)

	key := loadKeyFromFixture(t, "key-roll.ini")

	command := key.CurrentCommand()

	stdout, err := key.RunCommand()

	if err != nil {
		t.Fatalf("Command should not have timed out")
	}

	if string(stdout) != "hello\n" {
		t.Fatalf("Command did not return expected stdout. Got: \"%s\"", string(stdout))
	}

	if command == key.CurrentCommand() {
		t.Fatalf("Command did not advance after being run")
	}
}
