package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/football-api/pkg/response"
)

// APIKeyAuth middleware — protects write endpoints (POST/PUT/DELETE)
// Client harus kirim header: X-API-Key: <key>
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" || key != apiKey {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
