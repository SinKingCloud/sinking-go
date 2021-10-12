package sinking_sdk_go

import "testing"

func Test_main(t *testing.T) {
	test := &RequestServer{
		Server:    "106.52.89.187:80",
		TokenName: "sinking-token",
		Token:     "test_token",
	}
	test.registerServer("cloud.api", "sinking.go", "dev", "sinking-go-order", "127.0.0.1:8888")
	test.getServerList()
}
