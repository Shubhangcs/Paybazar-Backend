package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

// Admin Wallet Functions

func (q *Query) GetAdminWalletBalance(adminId string) (string, error) {
	var balance string
	query := `SELECT balance FROM admin_wallets WHERE admin_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, adminId).Scan(&balance)
	return balance, err
}

func (q *Query) AdminWalletTopup(req *structures.AdminWalletTopupRequest) error {
	query := `
	WITH upd AS (
		UPDATE admin_wallets
		SET balance = balance + $2
		WHERE admin_id = $1
		RETURNING admin_id
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
			'CREDIT',
			'SELF',
			NULL,
			$3
		FROM upd
	)
	SELECT 1;
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		req.AdminId,
		req.Amount,
		req.Remarks,
	)

	return err
}

// Master Distributor Wallet Function

func (q *Query) GetMasterDistributorWalletBalance(masterDistributorId string) (string, error) {
	var balance string
	query := `SELECT balance FROM master_distributor_wallets WHERE master_distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, masterDistributorId).Scan(&balance)
	return balance, err
}
// Distributor Wallet Function

func (q *Query) GetDistributorWalletBalance(distributorId string) (string, error) {
	var balance string
	query := `SELECT balance FROM distributor_wallets WHERE distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, distributorId).Scan(&balance)
	return balance, err
}
// User Wallet Function

func (q *Query) GetUserWalletBalance(userId string) (string, error) {
	var balance string
	query := `SELECT balance FROM users WHERE user_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, userId).Scan(&balance)
	return balance, err
}
