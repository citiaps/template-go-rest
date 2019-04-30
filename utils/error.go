package utils

import (

	"github.com/gin-gonic/gin"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func FormatError(message string, e string) gin.H{
	return gin.H{"message":message, "error":e}
}
