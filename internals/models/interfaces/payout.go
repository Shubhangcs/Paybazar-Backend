package interfaces

import (
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type PayoutInterface interface {
	PayoutRequest(echo.Context) (string, error)
	GetPayoutTransactions(echo.Context) (*[]structures.GetPayoutLogs, error)
}
