package interfaces

import (
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type Banks interface {
	GetBanks(echo.Context) (*[]structures.BankModel, error)
	AddNewBank(echo.Context) error
}