package interfaces

import "github.com/labstack/echo/v4"

type PayoutInterface interface {
	PayoutRequest(echo.Context) (string, error)
}
