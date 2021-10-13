package sinking_sdk_go

import (
	"fmt"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	server := New("127.0.0.1:8888", "sinking-token", "test_token", "sinking-go-api-order", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8888")
	server.Listen()
	server2 := New("127.0.0.1:8888", "sinking-token", "test_token", "sinking-go-api-order", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8887")
	server2.Listen()
	server3 := New("127.0.0.1:8888", "sinking-token", "test_token", "sinking-go-api-order", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8886")
	server3.Listen()
	go func() {
		time.Sleep(5 * time.Second)
		for {
			fmt.Println(server3.GetService("sinking-go-api-order"))
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(999 * time.Second)
}
