package handlers

import (
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

func (ch *commonHandler) GetAllMasterDistributorsByAdminID(e echo.Context) error {
	adminID := e.Param("admin_id")
	res, err := ch.commonRepo.GetAllMasterDistributorsByAdminID(adminID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.CommonResponse{
			Message: err.Error(),
			Status:  "failed",
		})
	}
	return e.JSON(http.StatusOK, structures.CommonResponse{
		Message: "master distributors fetched successfully",
		Status:  "success",
		Data:    res,
	})
}

func (ch *commonHandler) GetAllDistributorsByMasterDistributorID(e echo.Context) error {
	masterDistributorID := e.Param("master_distributor_id")
	res, err := ch.commonRepo.GetAllDistributorsByMasterDistributorID(masterDistributorID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.CommonResponse{
			Message: err.Error(),
			Status:  "failed",
		})
	}
	return e.JSON(http.StatusOK, structures.CommonResponse{
		Message: "distributors fetched successfully",
		Status:  "success",
		Data:    res,
	})
}

func (ch *commonHandler) GetAllUsersByDistributorID(e echo.Context) error {
	distributorID := e.Param("distributor_id")
	res, err := ch.commonRepo.GetAllUsersByDistributorID(distributorID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.CommonResponse{
			Message: err.Error(),
			Status:  "failed",
		})
	}
	return e.JSON(http.StatusOK, structures.CommonResponse{
		Message: "users fetched successfully",
		Status:  "success",
		Data:    res,
	})
}
