package sinking_sdk_go

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	server := New("106.52.89.187:80", "sinking-token", "test_token", "sinking-go-api-order", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8888")
	server.Listen()
	time.Sleep(999 * time.Second)
}
