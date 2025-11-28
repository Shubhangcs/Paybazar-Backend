package queries

import (
	"context"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/jackc/pgx/v5"
)

// Admin Wallet Functions

func (q *Query) GetAdminWalletBalance(adminId string) (string, error) {
	var balance string
	query := `SELECT admin_wallet_balance FROM admins WHERE admin_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, adminId).Scan(&balance)
	return balance, err
}

func (q *Query) AdminWalletTopup(req *structures.AdminWalletTopupRequest) error {
	query1 := `UPDATE admins SET admin_wallet_balance = admin_wallet_balance + $2::numeric WHERE admin_id=$1 RETURNING admin_name`
	query2 := `INSERT INTO 
	transactions(
		transactor_id,
		receiver_id,
		transactor_name,
		receiver_name,
		transactor_type,
		receiver_type,
		transaction_type,
		transaction_service,
		amount,
		transaction_status,
		remarks
	) VALUES(
		$1,
		$1,
		$2,
		$2,
		'ADMIN',
		'ADMIN',
		'CREDIT',
		'TOPUP',
		$3,
		'SUCCESS',
		$4
	)
	`
	var adminName string
	tx, err := q.Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(context.Background()) }()

	if err := tx.QueryRow(context.Background(), query1, req.AdminId, req.Amount).Scan(&adminName); err != nil {
		return err
	}

	if _, err := tx.Exec(context.Background(), query2, req.AdminId, adminName, req.Amount, req.Remarks); err != nil {
		return err
	}
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

// Master Distributor Wallet Function

func (q *Query) GetMasterDistributorWalletBalance(masterDistributorId string) (string, error) {
	var balance string
	query := `SELECT master_distributor_wallet_balance FROM master_distributors WHERE master_distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, masterDistributorId).Scan(&balance)
	return balance, err
}

// Distributor Wallet Function

func (q *Query) GetDistributorWalletBalance(distributorId string) (string, error) {
	var balance string
	query := `SELECT distributor_wallet_balance FROM distributors WHERE distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, distributorId).Scan(&balance)
	return balance, err
}

// User Wallet Function

func (q *Query) GetUserWalletBalance(userId string) (string, error) {
	var balance string
	query := `SELECT user_wallet_balance FROM users WHERE user_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, userId).Scan(&balance)
	return balance, err
}

func (q *Query) GetTransactions(userId string) (*[]structures.WalletTransaction, error) {
	const query = `
		SELECT
			transaction_id,
			transactor_id,
			receiver_id,
			transactor_name,
			receiver_name,
			transactor_type,
			receiver_type,
			transaction_type,
			transaction_service,
			amount,
			transaction_status,
			remarks
		FROM 
			transactions
		WHERE 
			transactor_id = $1
			OR receiver_id = $1
		ORDER BY 
			created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []structures.WalletTransaction
	for rows.Next() {
		var t structures.WalletTransaction
		if err := rows.Scan(
			&t.TransactionID,
			&t.TransactorID,
			&t.ReceiverID,
			&t.TransactorName,
			&t.ReceiverName,
			&t.TransactorType,
			&t.ReceiverType,
			&t.TransactionType,
			&t.TransactionService,
			&t.Amount,
			&t.TransactionStatus,
			&t.Remarks,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (q *Query) UserRefund(req *structures.RefundRequest) error {
	updateUserWalletBalanceQuery := `
		UPDATE users
		SET user_wallet_balance = user_wallet_balance - $1::NUMERIC
		WHERE user_phone = $2 AND user_wallet_balance >= $1::NUMERIC;
	`
	updateAdminWalletBalanceQuery := `
		UPDATE admins
		SET admin_wallet_balance = admin_wallet_balance + $1::NUMERIC
		WHERE admin_id = $2;
	`

	ctx := context.Background()

	tx, err := q.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Deduct from user wallet
	cmdTag, err := tx.Exec(ctx, updateUserWalletBalanceQuery, req.Amount, req.PhoneNumber)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		// No row updated => user not found or insufficient balance
		return fmt.Errorf("user refund failed: insufficient balance or user not found")
	}

	// 2. Credit admin wallet
	if _, err := tx.Exec(ctx, updateAdminWalletBalanceQuery, req.Amount, req.AdminID); err != nil {
		return err
	}

	// 3. Commit
	return tx.Commit(ctx)
}

func (q *Query) MasterDistributorRefund(req *structures.RefundRequest) error {
	updateMdWalletBalanceQuery := `
		UPDATE master_distributors
		SET master_distributor_wallet_balance = master_distributor_wallet_balance - $1::NUMERIC
		WHERE master_distributor_phone = $2 AND master_distributor_wallet_balance >= $1::NUMERIC;
	`
	updateAdminWalletBalanceQuery := `
		UPDATE admins
		SET admin_wallet_balance = admin_wallet_balance + $1::NUMERIC
		WHERE admin_id = $2;
	`

	ctx := context.Background()

	tx, err := q.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Deduct from MD wallet
	cmdTag, err := tx.Exec(ctx, updateMdWalletBalanceQuery, req.Amount, req.PhoneNumber)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("master distributor refund failed: insufficient balance or MD not found")
	}

	// 2. Credit admin wallet
	if _, err := tx.Exec(ctx, updateAdminWalletBalanceQuery, req.Amount, req.AdminID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (q *Query) DistributorRefund(req *structures.RefundRequest) error {
	updateDistributorWalletBalanceQuery := `
		UPDATE distributors
		SET distributor_wallet_balance = distributor_wallet_balance - $1::NUMERIC
		WHERE distributor_phone = $2 AND distributor_wallet_balance >= $1::NUMERIC;
	`
	updateAdminWalletBalanceQuery := `
		UPDATE admins
		SET admin_wallet_balance = admin_wallet_balance + $1::NUMERIC
		WHERE admin_id = $2;
	`

	ctx := context.Background()

	tx, err := q.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Deduct from distributor wallet
	cmdTag, err := tx.Exec(ctx, updateDistributorWalletBalanceQuery, req.Amount, req.PhoneNumber)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("distributor refund failed: insufficient balance or distributor not found")
	}

	// 2. Credit admin wallet
	if _, err := tx.Exec(ctx, updateAdminWalletBalanceQuery, req.Amount, req.AdminID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
