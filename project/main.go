package main

import (
	"context"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"tickets/adapters"
	"tickets/app"
	"tickets/db"
	"tickets/ports"
)

func main() {
	dbInst, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	defer dbInst.Close()
	err = db.CreateDbSchema(dbInst)
	if err != nil {
		panic(err)
	}
	logger := watermill.NewStdLogger(false, false)
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	log.Init(logrus.InfoLevel)

	clientsImpl, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
		return nil
	})
	if err != nil {
		panic(err)
	}

	receiptsClient := adapters.NewReceiptsClient(clientsImpl)
	spreadsheetsClient := adapters.NewSpreadsheetsClient(clientsImpl)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}
	eventBus, err := ports.NewEventBus(publisher, logger)
	if err != nil {
		panic(err)
	}
	repo := adapters.NewPGTicketRepository(dbInst)
	printer := adapters.NewPrinterClientAdapter(clientsImpl)
	application := app.NewApplication(receiptsClient, spreadsheetsClient, repo, printer, eventBus)
	e := ports.NewHttpServer(eventBus, application)
	router, err := ports.NewRouter(application, rdb)
	if err != nil {
		panic(err)
	}
	g.Go(func() error {
		logrus.Info("Server starting...")
		<-router.Running()
		logrus.Info("Server started...")

		err := e.Start(":8080")
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})
	g.Go(func() error {
		err := router.StartConsumers(ctx)
		return err
	})
	err = g.Wait()
	if err != nil {
		panic(err)
	}
}
