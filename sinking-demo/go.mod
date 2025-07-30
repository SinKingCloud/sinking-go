module sinking-demo

go 1.11

require (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20250330150011-1690c6225578
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20250330150011-1690c6225578
)

replace (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211012114015-f249644cbf78 => ../sinking-web
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211012114015-f249644cbf78 => ../sinking-websocket
)
