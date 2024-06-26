package gee

import (
	"log"
	"time"
)

func Logger() HandleFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process Request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
