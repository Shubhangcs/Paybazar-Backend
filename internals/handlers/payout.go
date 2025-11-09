package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type payoutHandler struct {
	payoutRepo interfaces.PayoutInterface
}

func NewPayoutHandler(payoutRepo interfaces.PayoutInterface) *payoutHandler {
	return &payoutHandler{
		payoutRepo: payoutRepo,
	}
}

func (ph *payoutHandler) PayoutRequest(e echo.Context) error {
	res, err := ph.payoutRepo.PayoutRequest(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "all fund requests fetched successfully", Status: "success", Data: res})
}
