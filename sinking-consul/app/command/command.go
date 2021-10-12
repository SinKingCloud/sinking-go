package command

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/command/service"
)

func Run() {
	go service.Run()
}
