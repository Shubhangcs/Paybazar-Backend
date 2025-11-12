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

// respondWithError inspects error; if it's an *echo.HTTPError use its code/message,
// otherwise return 500 and a generic message.
func walletRespondWithError(e echo.Context, err error) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		msg := httpErr.Message
		if s, ok := msg.(string); ok {
			return e.JSON(httpErr.Code, structures.WalletResponse{Message: s, Status: "failed"})
		}
		// fallback if message is not a string
		return e.JSON(httpErr.Code, structures.WalletResponse{Message: "request failed", Status: "failed"})
	}
	return e.JSON(http.StatusInternalServerError, structures.WalletResponse{Message: "Internal server error", Status: "failed"})
}

func (wh *walletHandler) GetAdminWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetAdminWalletBalance(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "admin wallet balance fetched successfully",
		Status:  "success",
		Data:    map[string]string{"balance": res},
	})
}

func (wh *walletHandler) GetMasterDistributorWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetMasterDistributorWalletBalance(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "master distributor wallet balance fetched successfully",
		Status:  "success",
		Data:    map[string]string{"balance": res},
	})
}

func (wh *walletHandler) GetDistributorWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetDistributorWalletBalance(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "distributor wallet balance fetched successfully",
		Status:  "success",
		Data:    map[string]string{"balance": res},
	})
}

func (wh *walletHandler) GetUserWalletBalanceRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetUserWalletBalance(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "user wallet balance fetched successfully",
		Status:  "success",
		Data:    map[string]string{"balance": res},
	})
}

func (wh *walletHandler) AdminWalletTopupRequest(e echo.Context) error {
	res, err := wh.walletRepo.AdminWalletTopup(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{Message: res, Status: "success"})
}

func (wh *walletHandler) GetTransactionsRequest(e echo.Context) error {
	res, err := wh.walletRepo.GetTransactions(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "wallet transactions fetched successfully",
		Status:  "success",
		Data:    res,
	})
}
