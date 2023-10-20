package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func main() {
	r := gin.Default()

	r.GET("/", handler)
	if err := r.Run(os.Args[1]); err != nil {
		panic(err)
	}
}

func handler(c *gin.Context) {
	c.Render(http.StatusOK, render.Data{
		ContentType: "text/plain; charset=utf-8",
		Data:        []byte("ok"),
	})
}
