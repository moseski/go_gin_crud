package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/moseski/go_gin_crud/controllers"
)

func UserRoutes(r *gin.Engine) {
    r.POST("/users", controllers.CreateUser)
    r.GET("/users/:id", controllers.GetUser)
    r.GET("/users", controllers.GetAllUsers)  // New route to get all users
    r.PUT("/users/:id", controllers.UpdateUser)
    r.DELETE("/users/:id", controllers.DeleteUser)
}
