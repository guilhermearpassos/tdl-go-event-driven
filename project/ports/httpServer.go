package ports

import (
	"context"
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
	"net/http"
	"tickets/app"
	"tickets/domain"
)

type HttpServer struct {
	*echo.Echo
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

		for _, ticket := range request.Tickets {
			var event domain.Event
			if ticket.Status == "confirmed" {
				event = domain.TicketBookingConfirmed{
					Header:        domain.NewHeader(),
					TicketID:      ticket.TicketID,
					CustomerEmail: ticket.CustomerEmail,
					Price:         ticket.Price,
				}
			} else if ticket.Status == "canceled" {
				event = domain.TicketBookingCanceled{
					Header:        domain.NewHeader(),
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

	server := HttpServer{e}
	return &server
}
