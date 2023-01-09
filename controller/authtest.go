package controller

import (
	"github.com/gin-gonic/gin"
)

func GetAdmin(c *gin.Context) {
	c.String(200, "Hola soy el admin")
}
