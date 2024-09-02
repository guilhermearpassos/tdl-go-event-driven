package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
	"tickets/adapters"
	"tickets/domain"
)

var db *sqlx.DB
var getDbOnce sync.Once

func getDb() *sqlx.DB {
	getDbOnce.Do(func() {
		var err error
		db, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
	})
	return db
}

func TestMain(m *testing.M) {
	db = getDb()
	err := CreateDbSchema(db)
	if err != nil {
		panic(err)
	}
	m.Run()
}
func TestCreateTicket(t *testing.T) {
	repo := adapters.NewPGTicketRepository(db)
	ticketID := uuid.NewString()
	for i := 0; i < 3; i++ {
		err := repo.CreateTicket(context.Background(), domain.TicketStatus{
			TicketID: ticketID,
			Status:   "confirmed",
			Price: domain.Money{
				Amount:   "4.52",
				Currency: "USD",
			},
			CustomerEmail: "abc@123.com",
		})
		require.NoError(t, err)
	}
	tickets, err := repo.GetTickets(context.Background())
	require.NoError(t, err)
	count := 0
	for _, ticket := range tickets {
		if ticket.TicketID == ticketID {
			count += 1
		}
	}
	require.Equal(t, 1, count, "incorrect number of tickets with id %s: %d expected one", ticketID, count)
}
