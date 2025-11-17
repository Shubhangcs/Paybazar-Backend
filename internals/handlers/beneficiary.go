package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type beneficiaryHandler struct {
	repo interfaces.Beneficiary
}

func NewBeneficiaryHandler(repo interfaces.Beneficiary) *beneficiaryHandler {
	return &beneficiaryHandler{repo: repo}
}

func (bh *beneficiaryHandler) GetBeneficiaries(c echo.Context) error {
	res, err := bh.repo.GetBeneficiaries(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, structures.FundRequestResponse{Message: "beneficiaries fetched successfully", Status: "success", Data: map[string]any{
		"beneficieries": res,
	}})
}

func (bh *beneficiaryHandler) AddNewBeneficiary(c echo.Context) error {
	err := bh.repo.AddNewBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, structures.FundRequestResponse{Message: "beneficiaries added successfully", Status: "success"})
}

func (bh *beneficiaryHandler) VerifyBeneficiary(c echo.Context) error {
	err := bh.repo.VerifyBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, structures.FundRequestResponse{Message: "beneficiaries verification successfully", Status: "success"})
}

func (bh *beneficiaryHandler) DeleteBeneficiary(c echo.Context) error {
	err := bh.repo.DeleteBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, structures.FundRequestResponse{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, structures.FundRequestResponse{Message: "beneficiaries deleted successfully", Status: "success"})
}