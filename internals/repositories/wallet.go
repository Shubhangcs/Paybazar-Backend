package repositories

import (
	"fmt"

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

func (wr *walletRepo) GetAdminWalletBalance(e echo.Context) (string, error) {
	var req = e.Param("admin_id")
	res, err := wr.query.GetAdminWalletBalance(req)
	if err != nil {
		return "", fmt.Errorf("failed to retrive admin balance: %w", err)
	}
	return res, nil
}

func (wr *walletRepo) GetAdminWalletTransactions(e echo.Context) (*[]structures.AdminWalletTransactions, error) {
	var req = e.Param("admin_id")
	res, err := wr.query.GetAdminWalletTransactions(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive admin wallet transactions: %w", err)
	}
	return res, nil
}

func (wr *walletRepo) GetMasterDistributorWalletBalance(e echo.Context) (string, error) {
	var req = e.Param("master_distributor_id")
	res, err := wr.query.GetMasterDistributorWalletBalance(req)
	if err != nil {
		return "", fmt.Errorf("failed to retrive master distributor balance: %w", err)
	}
	return res, nil
}

func (wr *walletRepo) GetMasterDistributorWalletTransactions(e echo.Context) (*[]structures.MasterDistributorWalletTransactions, error) {
	var req = e.Param("master_distributor_id")
	res, err := wr.query.GetMasterDistributorWalletTransactions(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive master distributor wallet transactions: %w", err)
	}
	return res, nil
}

func (wr *walletRepo) GetDistributorWalletBalance(e echo.Context) (string, error) {
	var req = e.Param("distributor_id")
	res, err := wr.query.GetDistributorWalletBalance(req)
	if err != nil {
		return "", fmt.Errorf("failed to retrive distributor balance: %w", err)
	}
	return res, nil
}

func (wr *walletRepo) GetDistributorWalletTransactions(e echo.Context) (*[]structures.DistributorWalletTransactions, error) {
	var req = e.Param("distributor_id")
	res, err := wr.query.GetDistributorWalletTransactions(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive distributor wallet transaction: %w", err)
	}
	return res, nil
}

func (wr *walletRepo) GetUserWalletBalance(e echo.Context) (string, error) {
	var req = e.Param("user_id")
	res, err := wr.query.GetUserWalletBalance(req)
	if err != nil {
		return "", fmt.Errorf("failed to retrive user balance: %w", err)
	}
	return res, nil
}

func (wr *walletRepo) GetUserWalletTransactions(e echo.Context) (*[]structures.UserWalletTransactions, error) {
	var req = e.Param("user_id")
	res, err := wr.query.GetUserWalletTransactions(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive user transactions: %w", err)
	}
	return res, err
}

func (wr *walletRepo) AdminWalletTopup(e echo.Context) (string, error) {
	var req structures.AdminWalletTopupRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request body: %w", err)
	}
	if err := wr.query.AdminWalletTopup(&req); err != nil {
		return "", fmt.Errorf("faild to topup admin wallet: %w", err)
	}
	return "admin wallet topup successfull", nil
}
