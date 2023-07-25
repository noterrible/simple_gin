package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// gin框架
func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	a := r.Group("/a")
	a.GET("/:id/a", func(c *gin.Context) {
		c.JSON(200, gin.H{"test2": "test2"})
	})
	a.GET("/:id/b", func(c *gin.Context) {
		c.JSON(200, gin.H{"test3": "test3"})
	})
	a.Use(test1())
	a.GET("/test", Test)
	r.Run(":8080")
}
func test1() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(111)
	}
}
func Test(c *gin.Context) {
	c.JSON(200, gin.H{"test": "test"})
}
