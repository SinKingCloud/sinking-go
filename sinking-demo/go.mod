module sinking-demo

go 1.11

require (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008065154-eb6948aeb02a
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211008065154-eb6948aeb02a
	github.com/gorilla/websocket v1.4.2
)

replace (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008065154-eb6948aeb02a => ../sinking-web
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211008065154-eb6948aeb02a => ../sinking-websocket
)
