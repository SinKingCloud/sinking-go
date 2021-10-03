package sinking_web

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Author(engine *Engine, addr string) {
	fmt.Println(" ____    ______   __  __  __  __   ______   __  __  ____              ____    _____      \n/\\  _`\\ /\\__  _\\ /\\ \\/\\ \\/\\ \\/\\ \\ /\\__  _\\ /\\ \\/\\ \\/\\  _`\\           /\\  _`\\ /\\  __`\\    \n\\ \\,\\L\\_\\/_/\\ \\/ \\ \\ `\\\\ \\ \\ \\/'/'\\/_/\\ \\/ \\ \\ `\\\\ \\ \\ \\L\\_\\         \\ \\ \\L\\_\\ \\ \\/\\ \\   \n \\/_\\__ \\  \\ \\ \\  \\ \\ , ` \\ \\ , <    \\ \\ \\  \\ \\ , ` \\ \\ \\L_L   _______\\ \\ \\L_L\\ \\ \\ \\ \\  \n   /\\ \\L\\ \\ \\_\\ \\__\\ \\ \\`\\ \\ \\ \\\\`\\   \\_\\ \\__\\ \\ \\`\\ \\ \\ \\/, \\/\\______\\\\ \\ \\/, \\ \\ \\_\\ \\ \n   \\ `\\____\\/\\_____\\\\ \\_\\ \\_\\ \\_\\ \\_\\ /\\_____\\\\ \\_\\ \\_\\ \\____/\\/______/ \\ \\____/\\ \\_____\\\n    \\/_____/\\/_____/ \\/_/\\/_/\\/_/\\/_/ \\/_____/ \\/_/\\/_/\\/___/            \\/___/  \\/_____/\n                                                                                         ")
	fmt.Println("SinKing-Go Framework " + FrameWorkVersion)
	fmt.Println("Author:SinKingCloud")
	fmt.Println("Blog:www.clwl.online")
	for k := range engine.router.handlers {
		fmt.Println("RequestHandle:", k)
	}
	fmt.Printf("The total handle is %d", len(engine.router.handlers))
	fmt.Printf("The request url is %s \n", addr)
	if debug {
		fmt.Printf("The run mode is debug")
	} else {
		fmt.Printf("The run mode is release")
	}
}

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[PID:%d] [STATUS:%d] [URL:%s TIME:%v] ", os.Getppid(), c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
