package ports

import (
	"context"
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"tickets/app"
	"tickets/app/command"
	"tickets/domain"
)

type HttpServer struct {
	*echo.Echo
}
type TicketBookingRequest struct {
	ShowId          string `json:"show_id"`
	NumberOfTickets int    `json:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email"`
}
type TicketBookingResponse struct {
	BookingId string `json:"booking_id"`
}

func NewHttpServer(bus *cqrs.EventBus, application app.Application) *HttpServer {
	e := commonHTTP.NewEcho()
	e.GET("/tickets", func(c echo.Context) error {
		tkts, err := application.Queries.ListTickets.Handle(context.Background())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, tkts)
	})
	e.POST("/tickets-status", func(c echo.Context) error {
		var request domain.TicketsStatusRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}
		idempotencyKey := c.Request().Header.Get("Idempotency-Key")
		if idempotencyKey == "" {
			return c.JSON(http.StatusBadRequest, "Idempotency-Key is required")
		}

		for _, ticket := range request.Tickets {
			var event domain.Event
			if ticket.Status == "confirmed" {
				event = domain.TicketBookingConfirmed{
					Header:        domain.NewHeader(idempotencyKey),
					TicketID:      ticket.TicketID,
					CustomerEmail: ticket.CustomerEmail,
					Price:         ticket.Price,
				}
			} else if ticket.Status == "canceled" {
				event = domain.TicketBookingCanceled{
					Header:        domain.NewHeader(idempotencyKey),
					TicketID:      ticket.TicketID,
					CustomerEmail: ticket.CustomerEmail,
					Price:         ticket.Price,
				}
			} else {
				continue
			}

			err = bus.Publish(context.Background(), event)
			if err != nil {
				return err
			}
		}

		return c.NoContent(http.StatusOK)
	})
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.POST("/shows", func(c echo.Context) error {
		var request domain.ShowRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}
		show := domain.Show{
			Id:              uuid.NewString(),
			ExternalID:      request.DeadNationId,
			NumberOfTickets: request.NumberOfTickets,
			StartTime:       request.StartTime,
			Title:           request.Title,
			Venue:           request.Venue,
		}
		err = application.Commands.CreateShow.Handle(context.Background(), show)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, show)

	})
	e.POST("/book-tickets", func(c echo.Context) error {
		var request TicketBookingRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}
		bookingID := uuid.NewString()
		err = application.Commands.BookTickets.Handle(context.Background(), command.BookTickets{
			ShowId:          request.ShowId,
			BookingId:       bookingID,
			CustomerEmail:   request.CustomerEmail,
			NumberOfTickets: request.NumberOfTickets,
		})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, TicketBookingResponse{BookingId: bookingID})
	})
	server := HttpServer{e}
	return &server
}
