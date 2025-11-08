package interfaces

import (
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type FundRequestInterface interface {
	CreateFundRequest(echo.Context) (string, error)
	RejectFundRequest(echo.Context) (string, error)
	AcceptFundRequest(echo.Context) (string, error)
	GetAllFundRequests(echo.Context) (*[]structures.FundRequest, error)
}
