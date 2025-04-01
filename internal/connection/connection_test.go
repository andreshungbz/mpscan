package connection

import (
	"net"
	"testing"
	"time"
)

func TestCreateDialer(t *testing.T) {
	t.Run("correct dialer timeout", func(t *testing.T) {
		got := createDialer(5)
		want := net.Dialer{Timeout: 5 * time.Second}

		if got.Timeout != want.Timeout {
			t.Errorf("got %v want %v", got.Timeout, want.Timeout)
		}
	})
}