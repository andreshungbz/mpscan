package scan

import (
	"testing"
)

func TestPortList(t *testing.T) {
	t.Run("correct print format", func(t *testing.T) {
		portList := PortList{22, 80, 443}

		got := portList.String()
		want := "22,80,443"

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("out of range port", func(t *testing.T) {
		var portList PortList

		for _, value := range []string{"0", "22", "80", "70000"} {
			portList.Set(value)
		}

		got := portList.String()
		want := "22,80"

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestTargetList(t *testing.T) {
	t.Run("correct print format", func(t *testing.T) {
		targetList := TargetList{"example.com", "localhost"}

		got := targetList.String()
		want := "example.com,localhost"

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
