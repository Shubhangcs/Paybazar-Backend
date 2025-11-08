package repositories

import (
	"fmt"

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

func (fr *fundRequestRepo) CreateFundRequest(e echo.Context) (string, error) {
	var req structures.FundRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request body: %w", err)
	}
	err := fr.query.CreateFundRequest(&req)
	if err != nil {
		return "", fmt.Errorf("failed to register fund request: %w", err)
	}
	return "fund request created successfully", err
}

func (fr *fundRequestRepo) RejectFundRequest(e echo.Context) (string, error) {
	var req string = e.Param("request_id")
	err := fr.query.RejectFundRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to reject fund request: %w", err)
	}
	return "fund request rejected successfully", nil
}

func (fr *fundRequestRepo) AcceptFundRequest(e echo.Context) (string, error) {
	var req structures.AcceptFundRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request body: %w", err)
	}
	err := fr.query.AcceptFundRequest(&req)
	if err != nil {
		return "", fmt.Errorf("failed to accept fund request: %w", err)
	}
	return "fund request accepted successfully", nil
}

func (fr *fundRequestRepo) GetFundRequestsById(e echo.Context) (*[]structures.FundRequest, error) {
	var req string = e.Param("requester_id")
	res, err := fr.query.GetFundRequestsById(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all fund request: %w", err)
	}
	return res, err
}

func (fr *fundRequestRepo) GetAllFundRequests(e echo.Context) (*[]structures.FundRequest, error) {
	var req = e.Param("admin_id")
	res, err := fr.query.GetAllFundRequests(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all fund request: %w", err)
	}
	return res, err
}
