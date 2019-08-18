package queue

import "time"

// Connection is our connection string to be reused
const Connection = `amqp://guest:guest@localhost:5672/`

// TimeFormat formats our timestamps
const TimeFormat = "2006-01-02T15:04:05.999999-07:00"

// SenderRequest is the request body for our API
type SenderRequest struct {
	Username string    `json:"username"`
	Message  string    `json:"message"`
	ID       time.Time `json:"id,omitempty"`
}

// SenderResponse is the response body
type SenderResponse struct {
	ID time.Time `json:"id"`
}

// ReceiverResponse contains the full response array
type ReceiverResponse struct {
	Messages []SenderRequest `json:"messages"`
}
