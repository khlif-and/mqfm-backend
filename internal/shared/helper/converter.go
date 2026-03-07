package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParamToUint(c *gin.Context, name string) (uint, bool) {
	val := c.Param(name)
	id, err := strconv.Atoi(val)
	if err != nil || id < 1 {
		return 0, false
	}
	return uint(id), true
}
