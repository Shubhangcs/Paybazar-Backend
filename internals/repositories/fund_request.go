package repositories

import (
	"log"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type fundRequestRepo struct {
	query *queries.Query
}

func NewFundRequestRepository(query *queries.Query) *fundRequestRepo {
	return &fundRequestRepo{
		query: query,
	}
}

func (fr *fundRequestRepo) bindAndValidate(e echo.Context, v interface{}) error {
	if err := e.Bind(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request format")
	}
	if err := e.Validate(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request data")
	}
	return nil
}

func (fr *fundRequestRepo) CreateFundRequest(e echo.Context) (string, error) {
	var req structures.CreateFundRequestModel
	if err := fr.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	if err := fr.query.CreateFundRequest(&req); err != nil {
		log.Println("DB create fund request error:", err)
		return "", echo.NewHTTPError(500, "Failed to create fund request")
	}
	return "Fund request created successfully", nil
}

func (fr *fundRequestRepo) RejectFundRequest(e echo.Context) (string, error) {
	requestID := e.Param("request_id")
	if requestID == "" {
		return "", echo.NewHTTPError(400, "request_id is required")
	}
	if err := fr.query.RejectFundRequest(requestID); err != nil {
		log.Println("DB reject fund request error:", err)
		return "", echo.NewHTTPError(500, "Failed to reject fund request")
	}
	return "Fund request rejected successfully", nil
}

func (fr *fundRequestRepo) AcceptFundRequest(e echo.Context) (string, error) {
	var req structures.AcceptFundRequestModel
	if err := fr.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	if err := fr.query.AcceptFundRequest(&req); err != nil {
		log.Println("DB accept fund request error:", err)
		return "", echo.NewHTTPError(500, "Failed to accept fund request")
	}
	return "Fund request accepted successfully", nil
}

func (fr *fundRequestRepo) GetFundRequestsById(e echo.Context) (*[]structures.GetFundRequestModel, error) {
	requesterID := e.Param("requester_id")
	if requesterID == "" {
		return nil, echo.NewHTTPError(400, "requester_id is required")
	}
	res, err := fr.query.GetFundRequestsById(requesterID)
	if err != nil {
		log.Println("DB get fund requests by id error:", err)
		return nil, echo.NewHTTPError(500, "Failed to fetch fund requests")
	}
	if res == nil {
		empty := []structures.GetFundRequestModel{}
		return &empty, nil
	}
	return res, nil
}

func (fr *fundRequestRepo) GetAllFundRequests(e echo.Context) (*[]structures.GetFundRequestModel, error) {
	adminID := e.Param("admin_id")
	if adminID == "" {
		return nil, echo.NewHTTPError(400, "admin_id is required")
	}
	res, err := fr.query.GetAllFundRequests(adminID)
	if err != nil {
		log.Println("DB get all fund requests error:", err)
		return nil, echo.NewHTTPError(500, "Failed to fetch fund requests")
	}
	if res == nil {
		empty := []structures.GetFundRequestModel{}
		return &empty, nil
	}
	return res, nil
}
