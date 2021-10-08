module sinking-demo

go 1.11

require (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008042124-8e8ec72bf463
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211008042124-8e8ec72bf463
	github.com/gorilla/websocket v1.4.2
)

replace (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008042124-8e8ec72bf463 => ../sinking-web
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211008042124-8e8ec72bf463 => ../sinking-websocket
)
