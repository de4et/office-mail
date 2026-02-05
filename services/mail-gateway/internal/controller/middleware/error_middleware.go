package middleware

import (
	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Msg string `json:"msg"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		c.JSON(-1, errResponse{
			Msg: "something went wrong",
		})
	}
}
