package repositories

import (
	"fmt"
	"log"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type ticketRepo struct {
	query *queries.Query
}

func NewTicketRepo(query *queries.Query) *ticketRepo {
	return &ticketRepo{
		query: query,
	}
}

func (tr *ticketRepo) AddNewTicket(e echo.Context) error {
	var req structures.Ticket
	if err := e.Bind(&req); err != nil {
		return fmt.Errorf("failed to create new ticket")
	}
	if err := tr.query.AddNewTicket(req); err != nil {
		log.Println(err)
		return fmt.Errorf("failed to create new ticket")
	}
	return nil
}

func (tr *ticketRepo) GetAllTickets(e echo.Context) (*[]structures.Ticket, error) {
	var adminId = e.Param("admin_id")
	if adminId == "" {
		return nil, fmt.Errorf("admin id not found")
	}
	res, err := tr.query.GetAllTickets(adminId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all tickets")
	}
	return res, nil
}
