package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	authRepo interfaces.AuthInterface
}

func NewAuthHandler(ar interfaces.AuthInterface) *authHandler {
	return &authHandler{
		authRepo: ar,
	}
}

func (ah *authHandler) RegisterAdminRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterAdmin(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "admin registration successfull", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) RegisterMasterDistributorRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterMasterDistributor(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "master distributor registration successfull", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) RegisterDistributorRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterDistributor(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "distributor registration successfull", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) RegisterUserRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterUser(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.AuthResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "user registration successfull", Status: "success", Data: map[string]string{"token": token}})
}
