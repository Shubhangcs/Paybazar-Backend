package interfaces

import (
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type Beneficiary interface {
	GetBeneficiaries(echo.Context) (*[]structures.BeneficiaryModel, error)
	AddNewBeneficiary(echo.Context) error
	VerifyBeneficiary(echo.Context) error
}
