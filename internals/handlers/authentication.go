package handlers

import (
	"fmt"
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

// respondWithError inspects error; if it's an *echo.HTTPError use its code/message,
// otherwise return 500 and a generic message.
func authRespondWithError(e echo.Context, err error) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		msg := fmt.Sprint(httpErr.Message)
		return e.JSON(httpErr.Code, structures.AuthResponse{Message: msg, Status: "failed"})
	}
	// fallback
	return e.JSON(http.StatusInternalServerError, structures.AuthResponse{Message: "Internal server error", Status: "failed"})
}

func (ah *authHandler) RegisterAdminRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterAdmin(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "admin registration successful", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) RegisterMasterDistributorRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterMasterDistributor(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "master distributor registration successful", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) RegisterDistributorRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterDistributor(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "distributor registration successful", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) RegisterUserRequest(e echo.Context) error {
	token, err := ah.authRepo.RegisterUser(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "user registration successful", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) LoginAdminRequest(e echo.Context) error {
	token, err := ah.authRepo.LoginAdmin(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "admin login successful", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) LoginMasterDistributorRequest(e echo.Context) error {
	token, err := ah.authRepo.LoginMasterDistributor(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "master distributor login successful", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) LoginDistributorRequest(e echo.Context) error {
	token, err := ah.authRepo.LoginDistributor(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "distributor login successful", Status: "success", Data: map[string]string{"token": token}})
}

func (ah *authHandler) LoginUserSendOTPRequest(e echo.Context) error {
	message, err := ah.authRepo.LoginUserSendOTP(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: message, Status: "success"})
}

func (ah *authHandler) LoginUserValidateOTPRequest(e echo.Context) error {
	token, err := ah.authRepo.LoginUserValidateOTP(e)
	if err != nil {
		return authRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.AuthResponse{Message: "user login successful", Status: "success", Data: map[string]string{"token": token}})
}
