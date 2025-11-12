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

// respondWithError inspects error; if it's an *echo.HTTPError use its code/message,
// otherwise return 500 and a generic message.
func fundRequestRespondWithError(e echo.Context, err error) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		msg := httpErr.Message
		return e.JSON(httpErr.Code, structures.FundRequestResponse{Message: msg.(string), Status: "failed"})
	}
	return e.JSON(http.StatusInternalServerError, structures.FundRequestResponse{Message: "Internal server error", Status: "failed"})
}

func (fh *fundRequestHandler) CreateFundRequest(e echo.Context) error {
	res, err := fh.fundRequestRepo.CreateFundRequest(e)
	if err != nil {
		return fundRequestRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: res, Status: "success"})
}

func (fh *fundRequestHandler) RejectFundRequest(e echo.Context) error {
	res, err := fh.fundRequestRepo.RejectFundRequest(e)
	if err != nil {
		return fundRequestRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: res, Status: "success"})
}

func (fh *fundRequestHandler) GetFundRequestsById(e echo.Context) error {
	res, err := fh.fundRequestRepo.GetFundRequestsById(e)
	if err != nil {
		return fundRequestRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{
		Message: "fund requests fetched successfully",
		Status:  "success",
		Data:    res,
	})
}

func (fh *fundRequestHandler) GetAllFundRequests(e echo.Context) error {
	res, err := fh.fundRequestRepo.GetAllFundRequests(e)
	if err != nil {
		return fundRequestRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{
		Message: "all fund requests fetched successfully",
		Status:  "success",
		Data:    res,
	})
}

func (fh *fundRequestHandler) AcceptFundRequest(e echo.Context) error {
	res, err := fh.fundRequestRepo.AcceptFundRequest(e)
	if err != nil {
		return fundRequestRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: res, Status: "success"})
}
