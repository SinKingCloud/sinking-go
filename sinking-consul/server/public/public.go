package public

import (
	"embed"
	"net/http"
)

var (
	//go:embed dist/*
	Static embed.FS
	//go:embed sql/install.sql
	Sql        string
	FileServer = http.FileServer(http.FS(Static))
)

func Path() string {
	return "/dist"
}

func ReadDistFile(file string) []byte {
	bytes, err := Static.ReadFile("dist/" + file)
	if err != nil || bytes == nil {
		return nil
	}
	return bytes
}
