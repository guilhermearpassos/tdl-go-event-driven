package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	main2 "tickets/domain"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type TicketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func main() {
	logger := watermill.NewStdLogger(false, false)

	log.Init(logrus.InfoLevel)

	clients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	receiptsClient := main2.NewReceiptsClient(clients)
	spreadsheetsClient := main2.NewSpreadsheetsClient(clients)
	w := NewWorker(receiptsClient, spreadsheetsClient)
	go w.Run()
	e := commonHTTP.NewEcho()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	receiptsSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, logger)
	if err != nil {
		panic(err)
	}
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	router.AddNoPublisherHandler("receipt-issuer", "issue-receipt", receiptsSub,
		func(msg *message.Message) error {
			err := receiptsClient.IssueReceipt(msg.Context(), string(msg.Payload))
			return err
		})
	spreadsheetSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, logger)
	if err != nil {
		panic(err)
	}
	router.AddNoPublisherHandler("tracker-appender", "append-to-tracker", spreadsheetSub,
		func(msg *message.Message) error {
			err := spreadsheetsClient.AppendRow(msg.Context(), "tickets-to-print", []string{string(msg.Payload)})
			return err
		})
	go func() {
		_ = router.Run(context.Background())
	}()

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}
	e.POST("/tickets-confirmation", func(c echo.Context) error {
		var request TicketsConfirmationRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		for _, ticket := range request.Tickets {
			err = publisher.Publish("issue-receipt", message.NewMessage(watermill.NewUUID(), []byte(ticket)))
			if err != nil {
				panic(err)
			}
			err = publisher.Publish("append-to-tracker", message.NewMessage(watermill.NewUUID(), []byte(ticket)))
			if err != nil {
				panic(err)
			}
		}

		return c.NoContent(http.StatusOK)
	})

	logrus.Info("Server starting...")

	err = e.Start(":8080")
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
