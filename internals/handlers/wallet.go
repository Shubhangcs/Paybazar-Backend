package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type walletHandler struct {
	walletRepo interfaces.WalletInterface
}

func NewWalletHandler(walletRepo interfaces.WalletInterface) *walletHandler {
	return &walletHandler{
		walletRepo: walletRepo,
	}
}

func (wh *walletHandler) GetAdminWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetAdminWalletBalance(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.WalletResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{Message: "admin wallet balance fetch successfull", Status: "success", Data: map[string]string{"balance": res}})
}

func (wh *walletHandler) GetMasterDistributorWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetMasterDistributorWalletBalance(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.WalletResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{Message: "master distributor wallet balance fetch successfull", Status: "success", Data: map[string]string{"balance": res}})
}

func (wh *walletHandler) GetDistributorWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetDistributorWalletBalance(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.WalletResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{Message: "distributor wallet balance fetch successfull", Status: "success", Data: map[string]string{"balance": res}})
}

func (wh *walletHandler) GetUserWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetUserWalletBalance(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.WalletResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{Message: "user wallet balance fetch successfull", Status: "success", Data: map[string]string{"balance": res}})
}

func (wh *walletHandler) AdminWalletTopupRequest(e echo.Context) error {
	res, err := wh.walletRepo.AdminWalletTopup(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.WalletResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{Message: res, Status: "success"})
}

func (wh *walletHandler) GetTransactionsRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetTransactions(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, structures.WalletResponse{Message: err.Error(), Status: "failed"})
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{Message: "wallet transactions fetch successfull", Status: "success", Data: res})
}
