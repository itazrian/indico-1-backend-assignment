package main

import (
    "backend-assignment/module/database"
    "backend-assignment/module/handlers"
    "backend-assignment/module/services"
    "github.com/gin-gonic/gin"
    "log"
    "os"
)

func main() {
    os.MkdirAll("/tmp/settlements", 0o755)
    
    database.Init()
    
    svc := services.NewSettlementService(database.DB)
    r := gin.Default()

    handlers.OrderRoutes(r)
    handlers.JobRoutes(r, svc)

    r.Static("/downloads", "/tmp/settlements")

    log.Println("listening on :8080")
    r.Run(":8080")
}
