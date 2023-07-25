package main

import (
	"Gee/gee"
	"fmt"
)

// 我实现的
func main() {
	r := gee.Default()
	//尝试一下
	a := r.Group("/a")
	a.GET("/:id/a", func(c *gee.Context) {
		c.JSON(200, gee.H{"test2": "test2"})
	})
	a.GET("/:id/b", func(c *gee.Context) {
		c.JSON(200, gee.H{"test3": "test3"})
	})
	a.Use(test1())
	a.GET("/test", Test)
	r.Run(":8080")
}
func test1() gee.HandleFunc {
	return func(c *gee.Context) {
		fmt.Println(111)
	}
}
func Test(c *gee.Context) {
	c.JSON(200, gee.H{"test": "test"})
}
