package scores

import (
	"github.com/gin-gonic/gin"
	"github.com/eahrend/papabob/common/nflapi"
	"net/http"
)


// Get today's score
func getScores(c *gin.Context){
	nflClient := c.MustGet("nflclient").(*nflapi.NFLClient)
	nflScores, err := nflClient.GetTodayScores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	basicReport, err := nflClient.BasicReport(nflScores)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.Writer.Write([]byte(basicReport))


}
