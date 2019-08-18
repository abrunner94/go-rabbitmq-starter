package queue

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

// Receive polls for messages to at a specific channel/queue
func Receive(c *gin.Context) {
	var messages []SenderRequest
	// Parse the request parameters
	channel := c.Param("channel")

	// Establish a connection to a local RabbitMQ
	conn, err := amqp.Dial(Connection)
	checkErrors(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	checkErrors(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		channel,
		true,
		false,
		false,
		false,
		nil,
	)
	checkErrors(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	checkErrors(err, "Failed to register a consumer")
	go handle(c, msgs, messages)

}

func handle(c *gin.Context, msgs <-chan amqp.Delivery, messages []SenderRequest) {
	// Corresponds to the timestamp
	lastID := c.Query("last_id")

	// Iterate over the messages in this queue
	for msg := range msgs {
		newMsg, _ := deserialize(msg.Body)

		if strings.TrimSpace(lastID) != "" {
			// Parse the query parameter from the URL into a timestamp
			timeQuery, _ := time.Parse(TimeFormat, lastID)
			// The timestamp for each message in the queue
			timeQueue := newMsg.ID
			diff := timeQueue.Sub(timeQuery)
			seconds := int(diff.Seconds())

			// Display only messages if the time difference is equal to or greater than 0
			if seconds >= 0 {
				messages = append(messages, newMsg)
			}

		} else {
			messages = append(messages, newMsg)
		}
		msg.Ack(false)
	}

	if len(messages) > 0 {
		newMessages := ReceiverResponse{Messages: messages}
		c.JSON(200, &newMessages)
	} else {
		c.Status(400)
	}

}

// Returns a deserialized version of a byte array as JSON
func deserialize(b []byte) (SenderRequest, error) {
	var msg SenderRequest
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
