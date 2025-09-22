package sinking_web

import (
	"log"
	"os"
	"time"
)

// Author 输出框架信息
func Author(engine *Engine, addr string) {
	log.Println("The framework is starting...\n ____    ______   __  __  __  __   ______   __  __  ____              ____    _____      \n/\\  _`\\ /\\__  _\\ /\\ \\/\\ \\/\\ \\/\\ \\ /\\__  _\\ /\\ \\/\\ \\/\\  _`\\           /\\  _`\\ /\\  __`\\    \n\\ \\,\\L\\_\\/_/\\ \\/ \\ \\ `\\\\ \\ \\ \\/'/'\\/_/\\ \\/ \\ \\ `\\\\ \\ \\ \\L\\_\\         \\ \\ \\L\\_\\ \\ \\/\\ \\   \n \\/_\\__ \\  \\ \\ \\  \\ \\ , ` \\ \\ , <    \\ \\ \\  \\ \\ , ` \\ \\ \\L_L   _______\\ \\ \\L_L\\ \\ \\ \\ \\  \n   /\\ \\L\\ \\ \\_\\ \\__\\ \\ \\`\\ \\ \\ \\\\`\\   \\_\\ \\__\\ \\ \\`\\ \\ \\ \\/, \\/\\______\\\\ \\ \\/, \\ \\ \\_\\ \\ \n   \\ `\\____\\/\\_____\\\\ \\_\\ \\_\\ \\_\\ \\_\\ /\\_____\\\\ \\_\\ \\_\\ \\____/\\/______/ \\ \\____/\\ \\_____\\\n    \\/_____/\\/_____/ \\/_/\\/_/\\/_/\\/_/ \\/_____/ \\/_/\\/_/\\/___/            \\/___/  \\/_____/\n                                                                                         ")
	log.Println("SinKing-Go Framework " + FrameWorkVersion)
	log.Println("Author:SinKingCloud")
	log.Println("Blog:www.clwl.online")
	for k := range engine.router.handlers {
		log.Println("RequestHandle:", k)
	}
	log.Printf("The total handle is %d\n", len(engine.router.handlers))
	log.Printf("The request url is %s\n", addr)
	if engine.debug {
		log.Printf("The run mode is debug\n")
	} else {
		log.Printf("The run mode is release\n")
	}
}

// Logger debug log
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[PID:%d] [IP:%s METHOD:%s STATUS:%d] [URL:%s TIME:%v] ", os.Getppid(), c.ClientIP(false), c.Method, c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
