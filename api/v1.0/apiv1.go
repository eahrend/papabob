package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/eahrend/papabob/api/v1.0/scores"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	scores.ApplyRoutes(v1)
}