package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckRoute(route *gin.Engine) {
	route.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
}
