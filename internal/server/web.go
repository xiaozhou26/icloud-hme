package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed web/*
var webFiles embed.FS

func (s *Server) registerUI() {
	assets, err := fs.Sub(webFiles, "web")
	if err != nil {
		panic(err)
	}
	s.r.GET("/", func(c *gin.Context) {
		index, readErr := fs.ReadFile(assets, "index.html")
		if readErr != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	s.r.StaticFS("/assets", http.FS(assets))
}
