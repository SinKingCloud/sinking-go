package command

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/command/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/command/synchronized"
)

func Run() {
	go service.Run()
	go synchronized.Run()
}
