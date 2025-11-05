package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type authHandler struct{
	authRepo interfaces.AuthInterface
}

func NewAuthHandler(ar interfaces.AuthInterface) *authHandler {
	return &authHandler{
		authRepo: ar,
	}
}

func(ah *authHandler) RegisterAdminRequest(e echo.Context) error {
	token , err := ah.authRepo.RegisterAdmin(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest , structures.AuthResponse{Message: err.Error() , Status: "failed"})
	}
	return e.JSON(http.StatusOK , structures.AuthResponse{Message: "registration successfull" , Status: "success" , Data: map[string]string{"token": token}})
}