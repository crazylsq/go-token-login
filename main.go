package main

import (
	"github.com/gin-gonic/gin"
	"opserver/controllers"
	"net/http"
	"opserver/middlewares"
)



func main() {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	do := router.Group("/do")
	{
		do.POST("/login", controllers.Login)
		do.GET("/verifytoken", controllers.LoginTokenVerify)
	}
	_user := router.Group("/user", middlewares.VerifyTokenMiddleware())
	{
		//_user.POST("/auth", user.Auth)
		_user.POST("/add",  controllers.Add)
		_user.PUT("/update", controllers.Update)
		_user.GET("/view", controllers.View)
		_user.DELETE("/delete" , controllers.Delete)
	}

	_test := router.Group("/test", middlewares.VerifyTokenMiddleware())
	{
		_test.POST("/xx", Xx)
	}
	router.Run(":8080")

}

func Xx(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
