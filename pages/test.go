package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Generate test file and fill in the data
func Test(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "test.html", nil)
}
