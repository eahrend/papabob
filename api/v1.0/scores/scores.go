package scores

import "github.com/gin-gonic/gin"

func ApplyRoutes(r *gin.RouterGroup) {
	scores := r.Group("/scores")
	{
		scores.GET("", getScores)
		scores.GET("/", getScores)
	}
}