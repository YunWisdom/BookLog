

package console

import (
	"net/http"

	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// LoginCheck checks login or not.
func LoginCheck(c *gin.Context) {
	session := util.GetSession(c)
	if 0 == session.UID {
		result := util.NewResult()
		result.Code = -2
		result.Msg = "unauthenticated request"
		c.AbortWithStatusJSON(http.StatusOK, result)

		return
	}

	c.Next()
}
