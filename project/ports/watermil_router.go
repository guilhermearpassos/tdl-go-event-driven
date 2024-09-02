package ports

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/redis/go-redis/v9"
	"tickets/app"
	"tickets/app/command"
	"tickets/domain"
	"tickets/ports/watermill_midlewares"
	"time"
)

type Router struct {
	app         app.Application
	redisClient redis.UniversalClient
	router      *message.Router
	logger      watermill.LoggerAdapter
}

func NewRouter(app app.Application, redisClient redis.UniversalClient) (*Router, error) {

	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, fmt.Errorf("creating router: %w", err)
	}
	return &Router{
		app:         app,
		redisClient: redisClient,
		router:      router,
		logger:      logger,
	}, nil
}

func (r *Router) StartConsumers(ctx context.Context) error {
	r.router.AddMiddleware(watermill_midlewares.CustomLogginMiddleware)
	r.router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          r.logger,
	}.Middleware)
	ep, err := cqrs.NewEventProcessorWithConfig(
		r.router,
		cqrs.EventProcessorConfig{
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return redisstream.NewSubscriber(redisstream.SubscriberConfig{
					Client:        r.redisClient,
					ConsumerGroup: params.HandlerName,
				}, r.logger)
			},
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
			Logger: r.logger,
		},
	)
	if err != nil {
		return err
	}
	handlers := make([]cqrs.EventHandler, 0)
	handlers = append(handlers, cqrs.NewEventHandler("receipt-issuer", func(ctx context.Context, event *domain.TicketBookingConfirmed) error {
		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}
		request := domain.IssueReceiptRequest{
			TicketID: event.TicketID,
			Price:    event.Price,
		}
		err = r.app.Commands.IssueReceipt.Handle(ctx, request)
		return err
	}))
	handlers = append(handlers, cqrs.NewEventHandler("tracker-appender", func(ctx context.Context, event *domain.TicketBookingConfirmed) error {
		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}
		err = r.app.Commands.RecordTicket.Handle(ctx, *event)
		return err
	}))
	handlers = append(handlers, cqrs.NewEventHandler("printer", func(ctx context.Context, event *domain.TicketBookingConfirmed) error {
		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}
		err = r.app.Commands.PrintTicket.Handle(ctx, command.PrintTicket{Ticket: domain.TicketStatus{
			TicketID:      event.TicketID,
			Status:        "confirmed",
			Price:         event.Price,
			CustomerEmail: event.CustomerEmail,
		}})
		return err
	}))
	handlers = append(handlers, cqrs.NewEventHandler("tickets-to-refund", func(ctx context.Context, event *domain.TicketBookingCanceled) error {
		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}
		err = r.app.Commands.RecordRefund.Handle(ctx, *event)
		return err
	}))
	handlers = append(handlers, cqrs.NewEventHandler("tickets-to-delete", func(ctx context.Context, event *domain.TicketBookingCanceled) error {
		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}
		err = r.app.Commands.CancelTicket.Handle(ctx, event.TicketID)
		return err
	}))
	handlers = append(handlers, cqrs.NewEventHandler("ticket-created", func(ctx context.Context, event *domain.TicketBookingConfirmed) error {
		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}
		err = r.app.Commands.CreateTicket.Handle(ctx, command.CreateTicket{TicketStatus: domain.TicketStatus{
			TicketID:      event.TicketID,
			Status:        "confirmed",
			Price:         event.Price,
			CustomerEmail: event.CustomerEmail,
		}})
		return err
	}))
	err = ep.AddHandlers(handlers...)
	if err != nil {
		return err
	}
	err = r.router.Run(ctx)
	return err
}

func (r *Router) Running() chan struct{} {
	return r.router.Running()
}

func (r *Router) ShutDown() error {
	return r.router.Close()
}
