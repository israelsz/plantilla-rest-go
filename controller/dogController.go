package controller

import (
	"github.com/gin-gonic/gin"
)

func GetDog(c *gin.Context) {
	c.String(200, "Hola soy el get perro")
}
