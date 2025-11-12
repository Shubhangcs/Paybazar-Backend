package interfaces

import "github.com/labstack/echo/v4"

type AuthInterface interface {
	RegisterAdmin(echo.Context) (string, error)
	RegisterMasterDistributor(echo.Context) (string, error)
	RegisterDistributor(echo.Context) (string, error)
	RegisterUser(echo.Context) (string, error)
	LoginAdmin(echo.Context) (string, error)
	LoginMasterDistributor(echo.Context) (string, error)
	LoginDistributor(echo.Context) (string, error)
	LoginUserSendOTP(echo.Context) (string, error)
	LoginUserValidateOTP(echo.Context) (string, error)
	SetUserMpin(echo.Context) (string, error)
}
