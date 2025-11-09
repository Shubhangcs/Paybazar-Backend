package queries

import (
	"context"
	"time"

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
			remarks
		)
		SELECT
			admin_id,
			$2,
			'CREDIT',
			'SELF',
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

func (q *Query) GetAdminWalletTransactions(adminId string) (*[]structures.AdminWalletTransactions, error) {
	const query = `
	SELECT 
		transaction_id::TEXT,
		admin_id::TEXT,
		amount,
		transaction_type,
		transaction_service,
		reference_id,
		remarks,
		created_at::TEXT
	FROM admin_wallet_transactions
	WHERE admin_id = $1
	ORDER BY created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, adminId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []structures.AdminWalletTransactions

	for rows.Next() {
		var tx structures.AdminWalletTransactions
		err := rows.Scan(
			&tx.TransactionId,
			&tx.AdminId,
			&tx.Amount,
			&tx.TransactionType,
			&tx.TransactionService,
			&tx.ReferenceId,
			&tx.Remarks,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &transactions, nil
}

// Master Distributor Wallet Function

func (q *Query) GetMasterDistributorWalletBalance(masterDistributorId string) (string, error) {
	var balance string
	query := `SELECT balance FROM master_distributor_wallets WHERE master_distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, masterDistributorId).Scan(&balance)
	return balance, err
}

func (q *Query) GetMasterDistributorWalletTransactions(masterDistributorId string) (*[]structures.MasterDistributorWalletTransactions, error) {
	const query = `
	SELECT 
		transaction_id::TEXT,
		master_distributor_id::TEXT,
		amount,
		transaction_type,
		transaction_service,
		reference_id,
		remarks,
		created_at
	FROM master_distributor_wallet_transactions
	WHERE master_distributor_id = $1
	ORDER BY created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, masterDistributorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []structures.MasterDistributorWalletTransactions

	for rows.Next() {
		var tx structures.MasterDistributorWalletTransactions
		err := rows.Scan(
			&tx.TransactionId,
			&tx.MasterDistributorId,
			&tx.Amount,
			&tx.TransactionType,
			&tx.TransactionService,
			&tx.ReferenceId,
			&tx.Remarks,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &transactions, nil
}


// Distributor Wallet Function

func (q *Query) GetDistributorWalletBalance(distributorId string) (string, error) {
	var balance string
	query := `SELECT balance FROM distributor_wallets WHERE distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, distributorId).Scan(&balance)
	return balance, err
}

func (q *Query) GetDistributorWalletTransactions(distributorId string) (*[]structures.DistributorWalletTransactions, error) {
	const query = `
	SELECT 
		transaction_id::TEXT,
		distributor_id::TEXT,
		amount::TEXT,
		transaction_type,
		transaction_service,
		reference_id,
		remarks,
		TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
	FROM distributor_wallet_transactions
	WHERE distributor_id = $1
	ORDER BY created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, distributorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []structures.DistributorWalletTransactions

	for rows.Next() {
		var tx structures.DistributorWalletTransactions
		err := rows.Scan(
			&tx.TransactionId,
			&tx.DistributorId,
			&tx.Amount,
			&tx.TransactionType,
			&tx.TransactionService,
			&tx.ReferenceId,
			&tx.Remarks,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &transactions, nil
}


// User Wallet Function

func (q *Query) GetUserWalletBalance(userId string) (string, error) {
	var balance string
	query := `SELECT balance FROM user_wallets WHERE user_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, userId).Scan(&balance)
	return balance, err
}

func (q *Query) GetUserWalletTransactions(userId string) (*[]structures.UserWalletTransactions, error) {
	const query = `
	SELECT 
		transaction_id::TEXT,
		user_id::TEXT,
		amount::TEXT,
		transaction_type,
		transaction_service,
		reference_id,
		remarks,
		TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
	FROM user_wallet_transactions
	WHERE user_id = $1
	ORDER BY created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []structures.UserWalletTransactions

	for rows.Next() {
		var tx structures.UserWalletTransactions
		err := rows.Scan(
			&tx.TransactionId,
			&tx.UserId,
			&tx.Amount,
			&tx.TransactionType,
			&tx.TransactionService,
			&tx.ReferenceId,
			&tx.Remarks,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &transactions, nil
}

