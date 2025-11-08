package interfaces

import (
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type WalletInterface interface {
	GetAdminWalletBalance(echo.Context) (string, error)
	GetAdminWalletTransactions(echo.Context) (*[]structures.AdminWalletTransactions, error)
	GetMasterDistributorWalletBalance(echo.Context) (string, error)
	GetMasterDistributorWalletTransactions(echo.Context) (*[]structures.MasterDistributorWalletTransactions, error)
	GetDistributorWalletBalance(echo.Context) (string, error)
	GetDistributorWalletTransactions(echo.Context) (*[]structures.DistributorWalletTransactions, error)
	GetUserWalletBalance(echo.Context) (string, error)
	GetUserWalletTransactions(echo.Context) (*[]structures.UserWalletTransactions, error)
	AdminWalletTopup(echo.Context) (string, error)
}
