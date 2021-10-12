package sinking_sdk_go

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	server := New("106.52.89.187:80", "sinking-token", "test_token", "cloud.api", "sinking.go", "dev", "sinking-go-order", "127.0.0.1:8888")
	server.Listen()
	time.Sleep(999 * time.Second)
}
