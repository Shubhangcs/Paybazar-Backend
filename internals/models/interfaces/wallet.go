package interfaces

import (
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type WalletInterface interface {
	GetAdminWalletBalance(echo.Context) (string, error)
	GetMasterDistributorWalletBalance(echo.Context) (string, error)
	GetDistributorWalletBalance(echo.Context) (string, error)
	GetUserWalletBalance(echo.Context) (string, error)
	AdminWalletTopup(echo.Context) (string, error)
	GetTransactions(echo.Context) (*[]structures.WalletTransaction, error)
	DistributorRefund(echo.Context) error
	MasterDistributorRefund(echo.Context) error
	UserRefund(echo.Context) error
}
