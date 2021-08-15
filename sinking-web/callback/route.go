package callback

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

type Route interface {
	NotFound(c *sinking_web.Context)
}

func NotFound(c *sinking_web.Context) {
	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
}
