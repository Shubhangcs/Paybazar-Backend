package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type bankHandler struct {
	bankRepository interfaces.Banks
}

func NewBankHandler(bankRepo interfaces.Banks) *bankHandler {
	return &bankHandler{
		bankRepository: bankRepo,
	}
}

func (bh *bankHandler) CreateBank(e echo.Context) error {
	err := bh.bankRepository.AddNewBank(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "bank added successfully", Status: "success"})
}

func (bh *bankHandler) GetAllBanks(e echo.Context) error {
	res, err := bh.bankRepository.GetBanks(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "banks fetched successfully", Status: "success", Data: map[string]any{
		"banks": res,
	}})
}
