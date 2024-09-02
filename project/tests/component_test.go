package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"testing"
	"tickets/adapters"
	"tickets/adapters/mock_services"
	"tickets/app"
	"tickets/db"
	"tickets/domain"
	"tickets/ports"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var receiptsClient mock_services.ReceiptsServiceMock
var spreadsheetsClient mock_services.SpreadsheetsTrackerClient

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := watermill.NewStdLogger(false, false)
	dbInst, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	defer dbInst.Close()
	err = db.CreateDbSchema(dbInst)
	if err != nil {
		panic(err)
	}

	repo := adapters.NewPGTicketRepository(dbInst)
	receiptsClient = mock_services.NewReceiptsServiceMock()
	spreadsheetsClient = mock_services.NewSpreadsheetsClient()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
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
	printer := mock_services.NewMockPrinterClient()
	application := app.NewApplication(&receiptsClient, spreadsheetsClient, repo, printer, eventBus)
	router, err := ports.NewRouter(application, rdb)
	if err != nil {
		panic(err)
	}
	e := ports.NewHttpServer(eventBus, application)

	go func() {
		err := router.StartConsumers(ctx)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		logrus.Info("Server starting...")
		<-router.Running()
		logrus.Info("Server started...")

		err := e.Start(":8080")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}

		return
	}()
	m.Run()
	//_ = e.Shutdown(ctx)
	//_ = router.ShutDown()
}

func TestComponent(t *testing.T) {
	// place for your tests!
	waitForHttpServer(t)
	ticketID := uuid.NewString()
	st := domain.TicketStatus{
		TicketID: ticketID,
		Status:   "created",
		Price: domain.Money{
			Amount:   "56.13",
			Currency: "BRL",
		},
		CustomerEmail: "123q@abc.com",
	}
	payload, err := json.Marshal(domain.TicketsStatusRequest{Tickets: []domain.TicketStatus{st}})
	require.NoError(t, err)
	_, err = http.Post("http://localhost:8080/tickets-status", "application/json", bytes.NewBuffer(payload))
	require.NoError(t, err)
	st.Status = "confirmed"
	payload, err = json.Marshal(domain.TicketsStatusRequest{Tickets: []domain.TicketStatus{st}})
	require.NoError(t, err)
	_, err = http.Post("http://localhost:8080/tickets-status", "application/json", bytes.NewBuffer(payload))
	require.NoError(t, err)
	st.Status = "canceled"
	payload, err = json.Marshal(domain.TicketsStatusRequest{Tickets: []domain.TicketStatus{st}})
	require.NoError(t, err)
	_, err = http.Post("http://localhost:8080/tickets-status", "application/json", bytes.NewBuffer(payload))
	require.NoError(t, err)
	time.Sleep(400 * time.Millisecond)
	_, receiptIssued := receiptsClient.IssuedReceipts[ticketID]
	assert.Greater(t, len(receiptsClient.IssuedReceipts), 0)
	assert.True(t, receiptIssued)
	_, confirmed := spreadsheetsClient.RowsCreatedBySheet["tickets-to-print"]
	assert.True(t, confirmed)
	_, refunded := spreadsheetsClient.RowsCreatedBySheet["tickets-to-refund"]
	assert.True(t, refunded)
}

func waitForHttpServer(t *testing.T) {
	t.Helper()

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8080/health")
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()

			if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
				return
			}
		},
		time.Second*10,
		time.Millisecond*50,
	)
}
