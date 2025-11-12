package handlers

import (
	"fmt"
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type commonHandler struct {
	commonRepo interfaces.CommonInterface
}

func NewCommonHandler(commonRepo interfaces.CommonInterface) *commonHandler {
	return &commonHandler{
		commonRepo: commonRepo,
	}
}

// Helper for standardized error responses
func commonRespondWithError(e echo.Context, err error) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		msg := fmt.Sprint(httpErr.Message)
		return e.JSON(httpErr.Code, structures.CommonResponse{
			Message: msg,
			Status:  "failed",
		})
	}
	// fallback: internal server error
	return e.JSON(http.StatusInternalServerError, structures.CommonResponse{
		Message: "Internal server error",
		Status:  "failed",
	})
}

// Get all master distributors under an admin
func (ch *commonHandler) GetAllMasterDistributorsByAdminID(e echo.Context) error {
	adminID := e.Param("admin_id")
	res, err := ch.commonRepo.GetAllMasterDistributorsByAdminID(adminID)
	if err != nil {
		return commonRespondWithError(e, err)
	}

	if len(*res) == 0 {
		return e.JSON(http.StatusOK, structures.CommonResponse{
			Message: "no master distributors found",
			Status:  "success",
			Data:    []interface{}{},
		})
	}

	return e.JSON(http.StatusOK, structures.CommonResponse{
		Message: "master distributors fetched successfully",
		Status:  "success",
		Data:    res,
	})
}

// Get all distributors under a master distributor
func (ch *commonHandler) GetAllDistributorsByMasterDistributorID(e echo.Context) error {
	masterDistributorID := e.Param("master_distributor_id")
	res, err := ch.commonRepo.GetAllDistributorsByMasterDistributorID(masterDistributorID)
	if err != nil {
		return commonRespondWithError(e, err)
	}

	if len(*res) == 0 {
		return e.JSON(http.StatusOK, structures.CommonResponse{
			Message: "no distributors found",
			Status:  "success",
			Data:    []interface{}{},
		})
	}

	return e.JSON(http.StatusOK, structures.CommonResponse{
		Message: "distributors fetched successfully",
		Status:  "success",
		Data:    res,
	})
}

// Get all users under a distributor
func (ch *commonHandler) GetAllUsersByDistributorID(e echo.Context) error {
	distributorID := e.Param("distributor_id")
	res, err := ch.commonRepo.GetAllUsersByDistributorID(distributorID)
	if err != nil {
		return commonRespondWithError(e, err)
	}

	if len(*res) == 0 {
		return e.JSON(http.StatusOK, structures.CommonResponse{
			Message: "no users found",
			Status:  "success",
			Data:    []interface{}{},
		})
	}

	return e.JSON(http.StatusOK, structures.CommonResponse{
		Message: "users fetched successfully",
		Status:  "success",
		Data:    res,
	})
}
