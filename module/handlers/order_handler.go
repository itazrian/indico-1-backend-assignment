package handlers

import (
    "backend-assignment/module/models"
    "backend-assignment/module/services"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

func OrderRoutes(r *gin.Engine) {
    r.POST("/orders", createOrder)
    r.GET("/orders/:id", getOrder)
}

func createOrder(c *gin.Context) {
    var req models.Order
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    svc := services.NewOrderService()
    if err := svc.CreateOrder(&req); err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, req)
}

func getOrder(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    svc := services.NewOrderService()
    o, err := svc.GetOrder(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, o)
}
