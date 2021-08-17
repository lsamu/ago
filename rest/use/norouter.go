package use

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// NoRouter NoRouter
func NoRouter() gin.HandlerFunc {
    return func(c *gin.Context){
        c.String(http.StatusNotFound, "%s", "Page Not Found")
    }
}
