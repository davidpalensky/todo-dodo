package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Generate index file and fill in the data
func Test(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "test.html", nil)
}
