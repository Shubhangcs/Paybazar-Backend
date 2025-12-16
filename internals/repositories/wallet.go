package repositories

import (
	"fmt"
	"log"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type walletRepo struct {
	query *queries.Query
}

func NewWalletRepository(query *queries.Query) *walletRepo {
	return &walletRepo{
		query: query,
	}
}

// Helper for binding + validation
func (wr *walletRepo) bindAndValidate(e echo.Context, v interface{}) error {
	if err := e.Bind(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request format")
	}
	if err := e.Validate(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request data")
	}
	return nil
}

// ---------------------------
// Wallet Balance Queries
// ---------------------------

func (wr *walletRepo) GetAdminWalletBalance(e echo.Context) (string, error) {
	adminID := e.Param("admin_id")
	if adminID == "" {
		return "", echo.NewHTTPError(400, "admin_id is required")
	}
	res, err := wr.query.GetAdminWalletBalance(adminID)
	if err != nil {
		log.Println("DB get admin wallet balance error:", err)
		return "", echo.NewHTTPError(500, "Failed to retrieve admin wallet balance")
	}
	return res, nil
}

func (wr *walletRepo) GetMasterDistributorWalletBalance(e echo.Context) (string, error) {
	masterDistributorID := e.Param("master_distributor_id")
	if masterDistributorID == "" {
		return "", echo.NewHTTPError(400, "master_distributor_id is required")
	}
	res, err := wr.query.GetMasterDistributorWalletBalance(masterDistributorID)
	if err != nil {
		log.Println("DB get master distributor wallet balance error:", err)
		return "", echo.NewHTTPError(500, "Failed to retrieve master distributor wallet balance")
	}
	return res, nil
}

func (wr *walletRepo) GetDistributorWalletBalance(e echo.Context) (string, error) {
	distributorID := e.Param("distributor_id")
	if distributorID == "" {
		return "", echo.NewHTTPError(400, "distributor_id is required")
	}
	res, err := wr.query.GetDistributorWalletBalance(distributorID)
	if err != nil {
		log.Println("DB get distributor wallet balance error:", err)
		return "", echo.NewHTTPError(500, "Failed to retrieve distributor wallet balance")
	}
	return res, nil
}

func (wr *walletRepo) GetUserWalletBalance(e echo.Context) (string, error) {
	userID := e.Param("user_id")
	if userID == "" {
		return "", echo.NewHTTPError(400, "user_id is required")
	}
	res, err := wr.query.GetUserWalletBalance(userID)
	if err != nil {
		log.Println("DB get user wallet balance error:", err)
		return "", echo.NewHTTPError(500, "Failed to retrieve user wallet balance")
	}
	return res, nil
}

// ---------------------------
// Admin Wallet Topup
// ---------------------------

func (wr *walletRepo) AdminWalletTopup(e echo.Context) (string, error) {
	var req structures.AdminWalletTopupRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	if err := wr.query.AdminWalletTopup(&req); err != nil {
		log.Println("DB admin wallet topup error:", err)
		return "", echo.NewHTTPError(500, "Failed to top up admin wallet")
	}
	return "Admin wallet top-up successful", nil
}

// ---------------------------
// Wallet Transaction History
// ---------------------------

func (wr *walletRepo) GetTransactions(e echo.Context) (*[]structures.WalletTransaction, error) {
	id := e.Param("id")
	if id == "" {
		return nil, echo.NewHTTPError(400, "id is required")
	}
	res, err := wr.query.GetTransactions(id)
	if err != nil {
		log.Println("DB get transactions error:", err)
		return nil, echo.NewHTTPError(500, "Failed to fetch transactions")
	}
	if res == nil {
		empty := []structures.WalletTransaction{}
		return &empty, nil
	}
	return res, nil
}

func (wr *walletRepo) UserRefund(e echo.Context) error {
	var req structures.RefundRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.UserRefund(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) MasterDistributorRefund(e echo.Context) error {
	var req structures.RefundRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.MasterDistributorRefund(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) DistributorRefund(e echo.Context) error {
	var req structures.RefundRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.DistributorRefund(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) MasterDistributorFundRetailer(e echo.Context) error {
	var req structures.MasterDistributorFundRetailerRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.MasterDistributorFundRetailer(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) MasterDistributorFundDistributor(e echo.Context) error {
	var req structures.MasterDistributorFundRetailerRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.MasterDistributorFundDistributor(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) DistributorFundRetailer(e echo.Context) error {
	var req structures.DistributorFundRetailerRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.DistributorFundRetailer(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) GetRevertHistoryPhone(phoneNumber string) (*[]structures.GetRevertHistory, error) {
	return wr.query.GetRevertHistoryPhone(phoneNumber)
}

func (wr *walletRepo) GetRevertHistory() (*[]structures.GetRevertHistory, error) {
	return wr.query.GetRevertHistory()
}

func (wr *walletRepo) MasterDistributorRefundUser(e echo.Context) error {
	var req structures.MasterDistributorFundRetailerRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.MDUserRefund(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) MasterDistributorRefundDistributor(e echo.Context) error {
	var req structures.MasterDistributorFundRetailerRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.MDDistributorRefund(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) DistributorRefundRetailer(e echo.Context) error {
	var req structures.DistributorFundRetailerRequest
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.DistributorUserRefund(&req); err != nil {
		return fmt.Errorf("low balance")
	}
	return nil
}

func (wr *walletRepo) UpdatePayoutTransaction(e echo.Context) error {
	var req structures.UpdatePayoutTransaction
	if err := wr.bindAndValidate(e, &req); err != nil {
		return err
	}
	if err := wr.query.UpdatePayoutTransaction(&req); err != nil {
		return err
	}
	return nil
}
