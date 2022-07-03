package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRouter(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
