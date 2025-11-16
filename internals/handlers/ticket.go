package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type ticketHandler struct {
	ticketRepo interfaces.TicketInterface
}

func NewTicketHandler(ticketRepo interfaces.TicketInterface) *ticketHandler {
	return &ticketHandler{
		ticketRepo: ticketRepo,
	}
}

func (th *ticketHandler) AddNewTicket(e echo.Context) error {
	err := th.ticketRepo.AddNewTicket(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "ticket added successfully", Status: "success"})
}

func (th *ticketHandler) GetAllTickets(e echo.Context) error {
	res, err := th.ticketRepo.GetAllTickets(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "beneficiaries fetched successfully", Status: "success", Data: map[string]any{
		"tickets": res,
	}})
}