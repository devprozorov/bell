package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/api/user")
	{
		user.GET("/profile", GetUserProfile)
		user.PUT("/profile", UpdateUserProfile)
	}
}

func GetUserProfile(c *gin.Context) {
	userID := c.GetString("userID")
	c.JSON(http.StatusOK, gin.H{"message": "profile fetched", "userID": userID})
}

func UpdateUserProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
