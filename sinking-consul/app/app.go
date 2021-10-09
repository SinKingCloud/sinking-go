package app

import (
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/str"
)

func Run() {
	fmt.Println(str.GetExternalIP())
}
