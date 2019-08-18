package main

import (
	queue "rabbit-messaging/internal/queue"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/:channel/messages", queue.Send)
	r.GET("/:channel/messages", queue.Receive)

	r.Run("localhost:12312")
}
