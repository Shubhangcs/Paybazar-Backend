package handlers

import (
	"net/http"

	"github.com/Srujankm12/paybazar-api/internals/models/interfaces"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type payoutHandler struct {
	payoutRepo interfaces.PayoutInterface
}

func NewPayoutHandler(payoutRepo interfaces.PayoutInterface) *payoutHandler {
	return &payoutHandler{
		payoutRepo: payoutRepo,
	}
}

// respondWithError inspects error; if it's an *echo.HTTPError use its code/message,
// otherwise return 500 and a generic message.
func payoutRespondWithError(e echo.Context, err error) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		msg := httpErr.Message
		if s, ok := msg.(string); ok {
			return e.JSON(httpErr.Code, structures.FundRequestResponse{Message: s, Status: "failed"})
		}
		return e.JSON(httpErr.Code, structures.FundRequestResponse{Message: "request failed", Status: "failed"})
	}
	return e.JSON(http.StatusInternalServerError, structures.FundRequestResponse{Message: "Internal server error", Status: "failed"})
}

func (ph *payoutHandler) PayoutRequest(e echo.Context) error {
	res, err := ph.payoutRepo.PayoutRequest(e)
	if err != nil {
		return payoutRespondWithError(e, err)
	}
	// res is a success message string from the repo (e.g., "Transaction successful")
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: res, Status: "success"})
}

func (ph *payoutHandler) GetPayoutTransactionRequest(e echo.Context) error {
	res, err := ph.payoutRepo.GetPayoutTransactions(e)
	if err != nil {
		return payoutRespondWithError(e, err)
	}
	// res is a success message string from the repo (e.g., "Transaction successful")
	return e.JSON(http.StatusOK, structures.FundRequestResponse{Message: "successfully fetched transactions", Status: "success", Data: map[string]any{"transactions": res}})
}

func (ph *payoutHandler) VerifyPayoutAccountNumber(e echo.Context) error {
	res, err := ph.payoutRepo.VerifyAccountNumber(e)
	if err != nil {
		return payoutRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, res)
}

func (ph *payoutHandler) PayoutTransactionRefund(e echo.Context) error {
	transactionId := e.Param("transaction_id")
	err := ph.payoutRepo.RefundPayoutTransaction(transactionId)
	if err != nil {
		return payoutRespondWithError(e, err)
	}
	return e.JSON(http.StatusOK, "refund successfull")
}
