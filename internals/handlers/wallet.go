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

func (wh *walletHandler) UserRefundRequest(e echo.Context) error {
	err := wh.walletRepo.UserRefund(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "refund success",
		Status:  "success",
	})
}

func (wh *walletHandler) MasterDistributorRefundRequest(e echo.Context) error {
	err := wh.walletRepo.MasterDistributorRefund(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "refund success",
		Status:  "success",
	})
}

func (wh *walletHandler) DistributorRefundRequest(e echo.Context) error {
	err := wh.walletRepo.DistributorRefund(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "refund success",
		Status:  "success",
	})
}

func (wh *walletHandler) MasterDistributorFundRetailerRequest(e echo.Context) error {
	err := wh.walletRepo.MasterDistributorFundRetailer(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "fund success",
		Status:  "success",
	})
}

func (wh *walletHandler) MasterDistributorFundDistributorRequest(e echo.Context) error {
	err := wh.walletRepo.MasterDistributorFundDistributor(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "fund success",
		Status:  "success",
	})
}

func (wh *walletHandler) DistributorFundRetailerRequest(e echo.Context) error {
	err := wh.walletRepo.DistributorFundRetailer(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "fund success",
		Status:  "success",
	})
}

func (wh *walletHandler) GetRevertHistory(e echo.Context) error {
	res, err := wh.walletRepo.GetRevertHistory()
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "revert history fetched successfully",
		Status:  "success",
		Data:    map[string]any{"revert_history": res},
	})
}

func (wh *walletHandler) GetRevertHistoryPhone(e echo.Context) error {
	res, err := wh.walletRepo.GetRevertHistoryPhone(e.Param("phone_number"))
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "revert history fetched successfully",
		Status:  "success",
		Data:    map[string]any{"revert_history": res},
	})
}

func (wh *walletHandler) MasterDistributorRefundDistributorRequest(e echo.Context) error {
	err := wh.walletRepo.MasterDistributorRefundDistributor(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "refund success",
		Status:  "success",
	})
}

func (wh *walletHandler) MasterDistributorRefundUserRequest(e echo.Context) error {
	err := wh.walletRepo.MasterDistributorRefundUser(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "refund success",
		Status:  "success",
	})
}


func (wh *walletHandler) DistributorRefundUserRequest(e echo.Context) error {
	err := wh.walletRepo.DistributorRefundRetailer(e)
	if err != nil {
		return walletRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, structures.WalletResponse{
		Message: "refund success",
		Status:  "success",
	})
}


