package queries

import (
	"context"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) PayoutRequestInitilizationRequest(req *structures.PayoutInitilizationRequest) (string, error) {
	var transactionId string

	const query = `
	WITH sel_user AS (
		SELECT u.user_id, uw.balance::numeric AS wallet_balance
		FROM users u
		JOIN user_wallets uw ON uw.user_id = u.user_id
		WHERE u.user_id = $1
	),
	check_balance AS (
		SELECT wallet_balance >= $7::numeric AS sufficient
		FROM sel_user
	),
	deduct_user AS (
		UPDATE user_wallets uw
		SET balance = uw.balance - $7::numeric
		FROM sel_user su
		WHERE uw.user_id = su.user_id
		  AND su.wallet_balance >= $7::numeric
		RETURNING uw.user_id
	),
	insert_payout AS (
		INSERT INTO payout_service (
			user_id,
			mobile_number,
			account_number,
			ifsc_code,
			bank_name,
			beneficiary_name,
			amount,
			transfer_type,
			transaction_status,
			remarks
		)
		SELECT
			su.user_id,
			$2, $3, $4, $5, $6,
			$7, $8, $9, $10
		FROM sel_user su
		JOIN deduct_user du ON du.user_id = su.user_id
		RETURNING transaction_id
	)
	SELECT transaction_id::TEXT
	FROM insert_payout;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := q.Pool.QueryRow(
		ctx,
		query,
		req.UserID,            // $1
		req.MobileNumber,      // $2
		req.AccountNumber,     // $3
		req.IFSCCode,          // $4
		req.BankName,          // $5
		req.BeneficiaryName,   // $6
		req.Amount,            // $7 (numeric)
		req.TransferType,      // $8 ('IMPS' | 'NEFT')
		req.TransactionStatus, // $9 (likely 'PENDING')
		req.Remarks,           // $10
	).Scan(&transactionId)

	return transactionId, err
}

func (q *Query) PayoutFinalSuccessRequest(req structures.PayoutFinalSuccessRequest) error {
	const query = `
	WITH sel AS (
		SELECT 
			ps.user_id,
			ps.amount::numeric AS amount
		FROM payout_service ps
		WHERE ps.transaction_id = $1
		  AND ps.transaction_status = 'PENDING'
		FOR UPDATE
	),
	deduct_user AS (
		UPDATE user_wallets uw
		SET balance = uw.balance - s.amount
		FROM sel s
		WHERE uw.user_id = s.user_id
		  AND uw.balance >= s.amount
		RETURNING uw.user_id
	),
	upd AS (
		UPDATE payout_service ps
		SET 
			operator_transaction_id = $2,
			order_id = $3,
			transaction_status = $4, -- e.g. 'SUCCESS'
			updated_at = NOW()
		WHERE ps.transaction_id = $1
		  AND EXISTS (SELECT 1 FROM deduct_user)
		RETURNING ps.transaction_id::TEXT
	)
	SELECT transaction_id FROM upd;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var tid string
	err := q.Pool.QueryRow(
		ctx,
		query,
		req.TransactionID,         // $1
		req.OperatorTransactionID, // $2
		req.OrderID,               // $3
		req.TransactionStatus,     // $4 ('SUCCESS' | 'FAILED' | etc.)
	).Scan(&tid)

	return err
}
