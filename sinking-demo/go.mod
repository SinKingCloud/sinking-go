module sinking-demo

go 1.11

require (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20251210092505-e7ad146c8491
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20251210092505-e7ad146c8491
	github.com/gorilla/websocket v1.5.3 // indirect
)

replace (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211012114015-f249644cbf78 => ../sinking-web
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211012114015-f249644cbf78 => ../sinking-websocket
)
