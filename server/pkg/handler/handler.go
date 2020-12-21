package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error(err)
				var (
					errMsg string
					ok     bool
				)
				if errMsg, ok = err.(string); ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "system error, " + errMsg,
					})
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "system error",
					})
					return
				}
			}
		}()
		c.Next()
	}
}
