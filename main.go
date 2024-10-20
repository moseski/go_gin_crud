package main

import (
    "github.com/gin-gonic/gin"
    "github.com/moseski/go_gin_crud/initializers"
    "github.com/moseski/go_gin_crud/routes"
)

func init() {
    // Load environment variables
    initializers.LoadEnvVariables()

    // Connect to the database
    initializers.ConnectToDB()
}

func main() {
    r := gin.Default()

    routes.UserRoutes(r)

    r.Run(":8080")
}


