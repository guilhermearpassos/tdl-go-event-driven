package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func CreateDbSchema(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tickets (
	ticket_id UUID PRIMARY KEY,
	price_amount NUMERIC(10, 2) NOT NULL,
	price_currency CHAR(3) NOT NULL,
	customer_email VARCHAR(255) NOT NULL
);`)
	if err != nil {
		return fmt.Errorf("creating table %s: %w", "tickets", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS shows (
	show_id UUID PRIMARY KEY,
	external_id UUID unique NOT NULL,
	start_time timestamptz NOT NULL,
	title varchar(100) NOT NULL,
	venue varchar(50) NOT NULL
);`)
	if err != nil {
		return fmt.Errorf("creating table %s: %w", "shows", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bookings (
	booking_id UUID PRIMARY KEY,
	customer_email varchar(50) NOT NULL,
	show_id uuid references shows(show_id) ON DELETE CASCADE,
	number_of_tickets integer
);`)
	if err != nil {
		return fmt.Errorf("creating table %s: %w", "bookings", err)
	}
	return err
}
