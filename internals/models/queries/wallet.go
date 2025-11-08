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

func (q *Query) AdminWalletTopup(req *structures.AdminWalletTopupRequest) (*structures.AdminWalletTopupResponse, error) {
	var res structures.AdminWalletTopupResponse

	query := `
	WITH upd AS (
		UPDATE admin_wallets
		SET balance = balance + $2
		WHERE admin_id = $1
		RETURNING admin_id, balance
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
		RETURNING transaction_id, admin_id
	)
	SELECT 
		ins.admin_id::TEXT AS admin_id,
		ins.transaction_id::TEXT AS transaction_id,
		upd.balance AS latest_balance
	FROM ins
	JOIN upd ON upd.admin_id = ins.admin_id;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminId,
		req.Amount,
		req.Remarks,
	).Scan(
		&res.AdminId,
		&res.TransactionId,
		&res.Balance,
	)

	return &res, err
}

// Master Distributor Wallet Function

func (q *Query) GetMasterDistributorWalletBalance(masterDistributorId string) (string, error) {
	var balance string
	query := `SELECT balance FROM master_distributor_wallets WHERE master_distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, masterDistributorId).Scan(&balance)
	return balance, err
}

func (q *Query) MasterDistributorWalletTopup(req *structures.MasterDistributorWalletTopupRequest) (*structures.MasterDistributorWalletTopupResponse, error) {
	var res structures.MasterDistributorWalletTopupResponse

	query := `
	WITH sel_admin AS (
		SELECT a.admin_id, a.admin_unique_id
		FROM admins a
		WHERE a.admin_id = $1
	),
	sel_md AS (
		SELECT m.master_distributor_id, m.master_distributor_unique_id
		FROM master_distributors m
		WHERE m.master_distributor_id = $2
	),
	-- Deduct from admin if balance is sufficient
	deduct_admin AS (
		UPDATE admin_wallets aw
		SET balance = aw.balance - $3
		FROM sel_admin sa
		WHERE aw.admin_id = sa.admin_id
		  AND aw.balance >= $3
		RETURNING aw.admin_id, sa.admin_unique_id
	),
	-- Record admin debit transaction (reference → MD unique_id)
	admin_tx AS (
		INSERT INTO admin_wallet_transactions (
			admin_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks
		)
		SELECT
			d.admin_id,
			$3,
			'DEBIT',
			'MD',
			sm.master_distributor_unique_id,  -- use unique_id here
			$4
		FROM deduct_admin d
		JOIN sel_md sm ON TRUE
		RETURNING 1
	),
	-- Credit master distributor wallet
	credit_md AS (
		UPDATE master_distributor_wallets mw
		SET balance = balance + $3
		FROM sel_md sm
		WHERE mw.master_distributor_id = sm.master_distributor_id
		RETURNING mw.master_distributor_id, mw.balance
	),
	-- Record master distributor credit transaction (reference → Admin unique_id)
	md_tx AS (
		INSERT INTO master_distributor_wallet_transactions (
			master_distributor_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks
		)
		SELECT
			c.master_distributor_id,
			$3,
			'CREDIT',
			'ADMIN',
			sa.admin_unique_id,  -- use admin unique_id here
			$4
		FROM credit_md c
		JOIN sel_admin sa ON TRUE
		RETURNING transaction_id, master_distributor_id
	)
	SELECT 
		mt.master_distributor_id::TEXT AS master_distributor_id,
		mt.transaction_id::TEXT AS transaction_id,
		c.balance AS balance
	FROM md_tx mt
	JOIN credit_md c ON c.master_distributor_id = mt.master_distributor_id;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminId,
		req.MasterDistributorId,
		req.Amount,
		req.Remarks,
	).Scan(
		&res.MasterDistributorId,
		&res.TransactionId,
		&res.Balance,
	)

	return &res, err
}

// Distributor Wallet Function

func (q *Query) GetDistributorWalletBalance(distributorId string) (string, error) {
	var balance string
	query := `SELECT balance FROM distributor_wallets WHERE distributor_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, distributorId).Scan(&balance)
	return balance, err
}

func (q *Query) DistributorWalletTopup(req *structures.DistributorWalletTopupRequest) (*structures.DistributorWalletTopupResponse, error) {
	var res structures.DistributorWalletTopupResponse

	query := `
	WITH sel_admin AS (
		SELECT a.admin_id, a.admin_unique_id
		FROM admins a
		WHERE a.admin_id = $1
	),
	sel_distributor AS (
		SELECT d.distributor_id, d.distributor_unique_id
		FROM distributors d
		WHERE d.distributor_id = $2
	),
	-- Deduct from admin if balance is sufficient
	deduct_admin AS (
		UPDATE admin_wallets aw
		SET balance = aw.balance - $3
		FROM sel_admin sa
		WHERE aw.admin_id = sa.admin_id
		  AND aw.balance >= $3
		RETURNING aw.admin_id, sa.admin_unique_id
	),
	-- Record admin debit transaction
	admin_tx AS (
		INSERT INTO admin_wallet_transactions (
			admin_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks
		)
		SELECT
			d.admin_id,
			$3,
			'DEBIT',
			'DISTRIBUTOR',
			sd.distributor_unique_id::UUID,  -- reference unique_id of distributor
			$4
		FROM deduct_admin d
		JOIN sel_distributor sd ON TRUE
		RETURNING 1
	),
	-- Credit distributor wallet
	credit_distributor AS (
		UPDATE distributor_wallets dw
		SET balance = balance + $3
		FROM sel_distributor sd
		WHERE dw.distributor_id = sd.distributor_id
		RETURNING dw.distributor_id, dw.balance
	),
	-- Record distributor credit transaction
	distributor_tx AS (
		INSERT INTO distributor_wallet_transactions (
			distributor_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks
		)
		SELECT
			c.distributor_id,
			$3,
			'CREDIT',
			'ADMIN',
			sa.admin_unique_id::UUID,  -- reference unique_id of admin
			$4
		FROM credit_distributor c
		JOIN sel_admin sa ON TRUE
		RETURNING transaction_id, distributor_id
	)
	SELECT 
		dt.distributor_id::TEXT AS distributor_id,
		dt.transaction_id::TEXT AS transaction_id,
		c.balance AS balance
	FROM distributor_tx dt
	JOIN credit_distributor c ON c.distributor_id = dt.distributor_id;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminId,
		req.DistributorId,
		req.Amount,
		req.Remarks,
	).Scan(
		&res.DistributorId,
		&res.TransactionId,
		&res.Balance,
	)

	return &res, err
}

// User Wallet Function

func (q *Query) GetUserWalletBalance(userId string) (string, error) {
	var balance string
	query := `SELECT balance FROM users WHERE user_id=$1`
	err := q.Pool.QueryRow(context.Background(), query, userId).Scan(&balance)
	return balance, err
}

func (q *Query) UserWalletTopup(req *structures.UserWalletTopupRequest) (*structures.UserWalletTopupResponse, error) {
	var res structures.UserWalletTopupResponse

	query := `
	WITH sel_admin AS (
		SELECT a.admin_id, a.admin_unique_id
		FROM admins a
		WHERE a.admin_id = $1
	),
	sel_user AS (
		SELECT u.user_id, u.user_unique_id
		FROM users u
		WHERE u.user_id = $2
	),
	-- Deduct from admin if balance is sufficient
	deduct_admin AS (
		UPDATE admin_wallets aw
		SET balance = aw.balance - $3
		FROM sel_admin sa
		WHERE aw.admin_id = sa.admin_id
		  AND aw.balance >= $3
		RETURNING aw.admin_id, sa.admin_unique_id
	),
	-- Record admin debit transaction (reference → user unique_id)
	admin_tx AS (
		INSERT INTO admin_wallet_transactions (
			admin_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks
		)
		SELECT
			d.admin_id,
			$3,
			'DEBIT',
			'USER',
			su.user_unique_id,  -- reference user unique_id
			$4
		FROM deduct_admin d
		JOIN sel_user su ON TRUE
		RETURNING 1
	),
	-- Credit user wallet
	credit_user AS (
		UPDATE user_wallets uw
		SET balance = balance + $3
		FROM sel_user su
		WHERE uw.user_id = su.user_id
		RETURNING uw.user_id, uw.balance
	),
	-- Record user credit transaction (reference → admin unique_id)
	user_tx AS (
		INSERT INTO user_wallet_transactions (
			user_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks
		)
		SELECT
			c.user_id,
			$3,
			'CREDIT',
			'ADMIN',
			sa.admin_unique_id,  -- reference admin unique_id
			$4
		FROM credit_user c
		JOIN sel_admin sa ON TRUE
		RETURNING transaction_id, user_id
	)
	SELECT 
		ut.user_id::TEXT AS user_id,
		ut.transaction_id::TEXT AS transaction_id,
		c.balance AS balance
	FROM user_tx ut
	JOIN credit_user c ON c.user_id = ut.user_id;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminId,
		req.UserId,
		req.Amount,
		req.Remarks,
	).Scan(
		&res.UserId,
		&res.TransactionId,
		&res.Balance,
	)

	return &res, err
}
