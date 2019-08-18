package queue

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

// Send sends a message to a specific channel/queue
func Send(c *gin.Context) {
	// Get the body of our POST request
	var req SenderRequest
	var res SenderResponse

	c.BindJSON(&req)
	// Parse the request body
	channel := c.Param("channel")

	// Check for empty fields and regular expression matching channel name [A-Z][a-z][0-9] and -
	user := strings.TrimSpace(req.Username)
	message := strings.TrimSpace(req.Message)
	newChannel := strings.TrimSpace(channel)
	channelAtMost := len(channel) >= 1
	channelMatch, _ := regexp.MatchString("^([a-z]|[A-Z]|[0-9]|-)+$", newChannel)

	// Check if regular expression matches, and no fields are empty
	if user == "" || message == "" || !channelMatch || !channelAtMost {
		c.Status(400)
	} else {
		// Establish a connection to a local RabbitMQ
		conn, err := amqp.Dial(Connection)
		checkErrors(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		checkErrors(err, "Failed to open a channel")
		defer ch.Close()

		// Declare a queue, which is our /{channel}/ for our API
		q, err := ch.QueueDeclare(
			newChannel,
			true,
			false,
			false,
			false,
			nil,
		)
		checkErrors(err, "Failed to declare a queue")

		// Add lastID to the request
		time := time.Now()
		req = SenderRequest{
			Username: req.Username,
			Message:  req.Message,
			ID:       time,
		}
		// Serialize request and publish it to the queue
		request, _ := serialize(req)

		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        request,
			})
		checkErrors(err, "Failed to publish a message")

		// Respond only with the ID
		res = SenderResponse{ID: time}
		c.JSON(http.StatusOK, &res)
	}
}

func checkErrors(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Returns a serialized version of the JSON
func serialize(msg SenderRequest) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}
