package interfaces

import (
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type TicketInterface interface {
	AddNewTicket(echo.Context) error
	GetAllTickets(echo.Context) (*[]structures.Ticket, error)
}
