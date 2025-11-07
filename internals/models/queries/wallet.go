package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) GetAdminWalletBalance(adminId string) (string, error) {
	var res string
	query := `SELECT balance FROM admin_wallets WHERE admin_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, adminId).Scan(&res)
	return res, err
}

func (q *Query) GetMasterDistributorWalletBalance(masterDistributorId string) (string, error) {
	var res string
	query := `SELECT balance FROM master_distributor_wallets WHERE master_distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, masterDistributorId).Scan(&res)
	return res, err
}

func (q *Query) GetDistributorWalletBalance(distributorId string) (string, error) {
	var res string
	query := `SELECT balance FROM distributor_wallets WHERE distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, distributorId).Scan(&res)
	return res, err
}

func (q *Query) GetUserWalletBalance(userId string) (string, error) {
	var res string
	query := `SELECT balance FROM user_wallets WHERE user_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, userId).Scan(&res)
	return res, err
}

func (q *Query) AddAmountToAdminWallet(req *structures.AdminAddAmountRequest) (*structures.AdminWallet, error) {
	var res structures.AdminWallet

	query := `
	WITH upd AS (
		UPDATE admin_wallets
		SET balance = balance + $2
		WHERE admin_id = $1
		RETURNING wallet_id, admin_id, balance
	),
	ins AS (
		INSERT INTO admin_wallet_transactions (
			admin_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks
		)
		SELECT
			admin_id,
			$2,
			$3,
			$4,
			$5,
			$6
		FROM upd
		RETURNING 1
	)
	SELECT 
		wallet_id::TEXT AS wallet_id,
		admin_id::TEXT AS admin_id,
		balance AS amount
	FROM upd;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminID,
		req.Amount,
		req.TransactionType,
		req.TransactionService,
		req.RefrenceID,
		req.Remarks,
	).Scan(
		&res.WalletID,
		&res.AdminID,
		&res.Balance,
	)

	return &res, err
}

func (q *Query) SendAmountToUserWalletFromAdmin(req *structures.AdminAddAmountToUserWalletRequest) error {

}

func (q *Query) SendAmountToMasterDistributorWalletFromAdmin() {
	
}