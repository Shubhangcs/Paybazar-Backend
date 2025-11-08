package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type fundRequestHandler struct {
	fundRequestRepo interfaces.FundRequestInterface
}

func NewFundRequestHandler(fundRequestRepo interfaces.FundRequestInterface) *fundRequestHandler {
	return &fundRequestHandler{
		fundRequestRepo,
	}
}

func (fh *fundRequestHandler) CreateFundRequest(e echo.Context) error {
	res, err := fh.fundRequestRepo.CreateFundRequest(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: res, Status: "success"})
}

func (fh *fundRequestHandler) RejectFundRequest(e echo.Context) error {
	res, err := fh.fundRequestRepo.RejectFundRequest(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: res, Status: "success"})
}

func (fh *fundRequestHandler) GetFundRequestsById(e echo.Context) error {
	res, err := fh.fundRequestRepo.GetFundRequestsById(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "all fund requests fetched successfully", Status: "success", Data: res})
}

func (fh *fundRequestHandler) GetAllFundRequests(e echo.Context) error {
	res, err := fh.fundRequestRepo.GetAllFundRequests(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "all fund requests fetched successfully", Status: "success", Data: res})
}

func (fh *fundRequestHandler) AcceptFundRequest(e echo.Context) error {
	res, err := fh.fundRequestRepo.AcceptFundRequest(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: res, Status: "success"})
}
