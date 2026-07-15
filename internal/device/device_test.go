package device

import (
	"os/user"
	"testing"
)

func TestCanListen(t *testing.T) {
	tests := []struct {
		user   string
		group  string
		result bool
	}{
		{"nobody", "root", false},
		{"nobody", "nobody", true},
	}

	for _, tt := range tests {
		u, err := user.Lookup(tt.user)
		if err != nil {
			t.Fatal(err)
		}

		g, err := user.LookupGroup(tt.group)
		if err != nil {
			t.Fatal(err)
		}

		result := canListen(u, g)

		if result != tt.result {
			t.Errorf("Expected %t for user %s in group %s, got %t", tt.result, tt.user, tt.group, result)
		}
	}
}
