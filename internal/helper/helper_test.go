package helper

import (
	"slices"
	"testing"
)

func TestCreateTargets(t *testing.T) {
	t.Run("aggregated targets", func(t *testing.T) {
		target := "localhost"
		targets := []string{"example.com", "scanme.nmap.org"}

		got := CreateTargets(target, targets)
		want := []string{"localhost", "example.com", "scanme.nmap.org"}

		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("missing target", func(t *testing.T) {
		target := ""
		targets := []string{"example.com", "scanme.nmap.org"}

		got := CreateTargets(target, targets)
		want := []string{"example.com", "scanme.nmap.org"}

		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("missing targets", func(t *testing.T) {
		target := "localhost"
		targets := []string{}

		got := CreateTargets(target, targets)
		want := []string{"localhost"}

		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestValidateSEPorts(t *testing.T) {
	t.Run("ports below minimum", func(t *testing.T) {
		start := -1
		end := -1

		ValidateSEPorts(&start, &end)

		if start != 1 && end != 1024 {
			t.Errorf("got start %v end %v wanted start 1 end 1024", start, end)
		}
	})

	t.Run("ports above maximum", func(t *testing.T) {
		start := 70000
		end := 70000

		ValidateSEPorts(&start, &end)

		if start != 1 && end != 1024 {
			t.Errorf("got start %v end %v wanted start 1 end 1024", start, end)
		}
	})
}
