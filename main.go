package main

import (
    "github.com/gin-gonic/gin"
    "github.com/moseski/go_gin_crud/initializers"
    "github.com/moseski/go_gin_crud/routes"
)

func init() {
    initializers.LoadEnvVariables()
    initializers.ConnectToDB()
}

func main() {
    r := gin.Default()

    // Setup routes
    routes.UserRoutes(r)

    // Run server
    r.Run(":8080")
}