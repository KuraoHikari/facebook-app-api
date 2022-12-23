package middleware

import (
	"net/http"

	"github.com/KuraoHikari/facebook-app-res-api/helper"
	"github.com/KuraoHikari/facebook-app-res-api/service"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(jwtServuce service.JWTService) gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	}
}