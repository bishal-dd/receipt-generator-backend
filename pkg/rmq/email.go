package rmq

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/redis/go-redis/v9"
	"github.com/resend/resend-go/v2"
)

type EmailMessage struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

type EmailConsumer struct {
    resendClient *resend.Client
}

var (
    EmailQueue rmq.Queue
    rmqConnection rmq.Connection

)

func NewEmailConsumer() (*EmailConsumer, error) {
    resendClient := resend.NewClient(os.Getenv("RESEND_API_KEY"))
    return &EmailConsumer{
        resendClient: resendClient,
    }, nil
}

func (c *EmailConsumer) Consume(delivery rmq.Delivery) {
    var email EmailMessage
    if err := json.Unmarshal([]byte(delivery.Payload()), &email); err != nil {
        log.Printf("Failed to unmarshal email: %v", err)
        if err := delivery.Reject(); err != nil {
            log.Printf("Failed to reject delivery: %v", err)
        }
        return
    }

    params := &resend.SendEmailRequest{
        From:    "hello@updates.yangkoo.com",
        To:      []string{email.To},
        Subject: email.Subject,
        Html:    email.Body,
    }

    _, err := c.resendClient.Emails.Send(params)
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        if err := delivery.Reject(); err != nil {
            log.Printf("Failed to reject delivery: %v", err)
        }
        return
    }

    if err := delivery.Ack(); err != nil {
        log.Printf("Failed to ack delivery: %v", err)
    }
}

// InitEmailQueue initializes the RMQ connection and email queue
func InitEmailQueue(redisClient *redis.Client)  error {
    errChan := make(chan error, 100)
    // Start error handler
    go func() {
        for err := range errChan {
            log.Printf("RMQ error: %v", err)
        }
    }()

    connection, err := rmq.OpenConnectionWithRedisClient("email-service", redisClient, errChan)
    if err != nil {
        return  err
    }
    rmqConnection = connection
    queue, err := connection.OpenQueue("email-queue")
    if err != nil {
        return err
    }

	EmailQueue = queue

    // Start consuming with prefetch 10
    if err := queue.StartConsuming(10, time.Second); err != nil {
        return  err
    }

    // Add consumer
    consumer, err := NewEmailConsumer()
    if err != nil {
        return  err
    }

    _, err = queue.AddConsumer("email-consumer", consumer)
    if err != nil {
        return  err
    }
    PollQueueStats(5 * time.Second)
    return nil
}

