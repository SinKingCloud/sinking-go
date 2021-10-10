package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

func Index(s *sinking_web.Context) {
	s.JSON(http.StatusOK, "success")
}
