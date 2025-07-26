module sinking-demo

go 1.11

require (
	github.com/SinKingCloud/sinking-go/sinking-sdk/sinking-sdk-go v0.0.0-20241113113141-f2da4c041bf2
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211012114015-f249644cbf78
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211012114015-f249644cbf78
)

replace (
	github.com/SinKingCloud/sinking-go/sinking-sdk/sinking-sdk-go v0.0.0-20241113113141-f2da4c041bf2 => ../sinking-sdk/sinking-sdk-go
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211012114015-f249644cbf78 => ../sinking-web
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211012114015-f249644cbf78 => ../sinking-websocket
)
