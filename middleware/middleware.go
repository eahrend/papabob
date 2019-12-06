package middleware

import (
	"github.com/eahrend/papabob/common/nflapi"
	"github.com/gin-gonic/gin"
)

func NFLClientMW(nfl *nflapi.NFLClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("nflclient", nfl)
		c.Next()
	}
}

