package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/jackc/pgx/v5"
)

func (q *Query) CheckUserBalance(userId string, amount string, commission string) (bool, error) {
	var hasBalance bool
	query := `SELECT 
    CASE 
        WHEN user_wallet_balance >= ($2::numeric + $3::numeric) THEN TRUE
        ELSE FALSE
    END AS has_sufficient_balance
FROM 
    users
WHERE 
    user_id = $1;
`
	err := q.Pool.QueryRow(context.Background(), query, userId, amount, commission).Scan(&hasBalance)
	return hasBalance, err
}

// func (q *Query) CheckPayoutLimit(userId string, amount string) (bool, error) {
// 	var hasPayoutLimit bool
// 	query := `
// 	SELECT
//     CASE
//         WHEN COALESCE(SUM(amount), 0) + $2::numeric <= 25000 THEN TRUE
//         ELSE FALSE
//     END AS within_limit
// FROM
//     payout_service
// WHERE
//     user_id = $1
//     AND transaction_status = 'SUCCESS'
//     AND created_at::date = CURRENT_DATE;
// `
// 	err := q.Pool.QueryRow(context.Background(), query, userId, amount).Scan(&hasPayoutLimit)
// 	return hasPayoutLimit, err
// }

func (q *Query) CheckMpin(userID string, mpin string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1 AND user_mpin = $2)`
	var hasMpin bool
	err := q.Pool.QueryRow(context.Background(), query, userID, mpin).Scan(&hasMpin)
	return hasMpin, err
}

func (q *Query) InitilizePayoutRequest(req *structures.PayoutInitilizationRequest) (*structures.PayoutApiRequest, error) {
	const query = `
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
		remarks,
		commision
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,UPPER($8),'PENDING',$9,$10)
	RETURNING 
		payout_transaction_id::TEXT,
		mobile_number AS mobile_no,
		account_number AS account_no,
		ifsc_code AS ifsc,
		bank_name,
		beneficiary_name AS benificiary_name,
		amount::TEXT,
		CASE 
			WHEN UPPER(transfer_type) = 'IMPS' THEN '5'
			WHEN UPPER(transfer_type) = 'NEFT' THEN '6'
		END AS transfer_type;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var res structures.PayoutApiRequest
	err := q.Pool.QueryRow(
		ctx,
		query,
		req.UserID,
		req.MobileNumber,
		req.AccountNumber,
		req.IFSCCode,
		req.BankName,
		req.BeneficiaryName,
		req.Amount,
		req.TransferType,
		req.Remarks,
		req.Commission,
	).Scan(
		&res.PartnerRequestID,
		&res.MobileNumber,
		&res.AccountNumber,
		&res.IFSCCode,
		&res.BankName,
		&res.BeneficiaryName,
		&res.Amount,
		&res.TransferType,
	)

	return &res, err
}

func (q *Query) PayoutSuccess(req *structures.PayoutApiSuccessResponse) error {
	ctx := context.Background()

	tx, err := q.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1) Lock payout_service row and read required fields (only if PENDING)
	var (
		userID        string
		amountStr     string
		commissionStr string
		currentStatus string
	)
	err = tx.QueryRow(ctx, `
		SELECT user_id, amount::text, commision::text, transaction_status
		FROM payout_service
		WHERE payout_transaction_id = $1
		FOR UPDATE
	`, req.PartnerRequestID).Scan(&userID, &amountStr, &commissionStr, &currentStatus)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("payout not found")
		}
		return fmt.Errorf("select payout_service: %w", err)
	}

	if currentStatus != "PENDING" {
		return fmt.Errorf("payout is not pending (status=%s)", currentStatus)
	}

	// 2) Lock user row and get wallet, name, and hierarchy ids
	var (
		userBalanceStr      string
		userName            string
		distributorID       *string
		masterDistributorID *string
		adminID             *string
	)
	err = tx.QueryRow(ctx, `
		SELECT user_wallet_balance::text, user_name, distributor_id, master_distributor_id, admin_id
		FROM users
		WHERE user_id = $1
		FOR UPDATE
	`, userID).Scan(&userBalanceStr, &userName, &distributorID, &masterDistributorID, &adminID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("select user: %w", err)
	}

	// 3) Ensure user has sufficient balance for the payout amount
	var sufficient bool
	err = tx.QueryRow(ctx, `
		SELECT (user_wallet_balance >= $1::numeric)
		FROM users
		WHERE user_id = $2
	`, amountStr, userID).Scan(&sufficient)
	if err != nil {
		return fmt.Errorf("check user balance: %w", err)
	}
	if !sufficient {
		return fmt.Errorf("user has insufficient wallet balance")
	}

	// 4) Deduct amount from user wallet
	_, err = tx.Exec(ctx, `
		UPDATE users
		SET user_wallet_balance = user_wallet_balance - $1::numeric - $3::numeric,
		    updated_at = NOW()
		WHERE user_id = $2
	`, amountStr, userID, commissionStr)
	if err != nil {
		return fmt.Errorf("deduct user wallet: %w", err)
	}

	// 5) Insert payout debit transaction (transactor and receiver both user; receiver_name = 'Payout')
	_, err = tx.Exec(ctx, `
		INSERT INTO transactions (
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
			remarks,
			created_at
		) VALUES (
			gen_random_uuid(),
			$1::uuid,
			$1::uuid,
			$2,
			'Payout',
			'USER',
			'USER',
			'DEBIT',
			'PAYOUT',
			$3::numeric + $5::numeric,
			'SUCCESS',
			('Payout processed | payout_id=' || $4::text),
			NOW()
		)
	`, userID, userName, amountStr, req.PartnerRequestID, commissionStr)
	if err != nil {
		return fmt.Errorf("insert payout transaction: %w", err)
	}

	// 6) Calculate commission shares and credit them
	//    admin: 25%, distributor: 20%, master distributor: 5%, retailer(user): 50%

	// a) Admin share (if admin_id not null)
	if adminID != nil && *adminID != "" {
		// credit admin wallet
		_, err = tx.Exec(ctx, `
			WITH share AS (SELECT ($1::numeric * 0.2917)::numeric AS amt)
			UPDATE admins
			SET admin_wallet_balance = admin_wallet_balance + (SELECT amt FROM share),
			    updated_at = NOW()
			WHERE admin_id = $2
		`, commissionStr, *adminID)
		if err != nil {
			return fmt.Errorf("credit admin commission: %w", err)
		}

		// fetch admin name (for transaction row)
		var adminName string
		err = tx.QueryRow(ctx, `SELECT admin_name FROM admins WHERE admin_id = $1`, *adminID).Scan(&adminName)
		if err != nil {
			return fmt.Errorf("fetch admin name: %w", err)
		}

		// insert admin commission transaction
		_, err = tx.Exec(ctx, `
			INSERT INTO transactions (
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
				remarks,
				created_at
			) VALUES (
				gen_random_uuid(),
				$1::uuid,
				$1::uuid,
				$2,
				$2,
				'ADMIN',
				'ADMIN',
				'CREDIT',
				'PAYOUT',
				($3::numeric * 0.2917)::numeric,
				'SUCCESS',
				('Payout commission credited to admin | payout_id=' || $4::text),
				NOW()
			)
		`, *adminID, adminName, commissionStr, req.PartnerRequestID)
		if err != nil {
			return fmt.Errorf("insert admin commission transaction: %w", err)
		}
	}

	// b) Distributor share (if distributor exists)
	if distributorID != nil && *distributorID != "" {
		_, err = tx.Exec(ctx, `
			WITH share AS (SELECT ($1::numeric * 0.1667)::numeric AS amt)
			UPDATE distributors
			SET distributor_wallet_balance = distributor_wallet_balance + (SELECT amt FROM share),
			    updated_at = NOW()
			WHERE distributor_id = $2
		`, commissionStr, *distributorID)
		if err != nil {
			return fmt.Errorf("credit distributor commission: %w", err)
		}

		var distName string
		err = tx.QueryRow(ctx, `SELECT distributor_name FROM distributors WHERE distributor_id = $1`, *distributorID).Scan(&distName)
		if err != nil {
			return fmt.Errorf("fetch distributor name: %w", err)
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO transactions (
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
				remarks,
				created_at
			) VALUES (
				gen_random_uuid(),
				$1::uuid,
				$1::uuid,
				$2,
				$2,
				'DISTRIBUTOR',
				'DISTRIBUTOR',
				'CREDIT',
				'PAYOUT',
				($3::numeric * 0.1667)::numeric,
				'SUCCESS',
				('Payout commission credited to distributor | payout_id=' || $4::text),
				NOW()
			)
		`, *distributorID, distName, commissionStr, req.PartnerRequestID)
		if err != nil {
			return fmt.Errorf("insert distributor commission transaction: %w", err)
		}
	}

	// c) Master distributor share (if exists)
	if masterDistributorID != nil && *masterDistributorID != "" {
		_, err = tx.Exec(ctx, `
			WITH share AS (SELECT ($1::numeric * 0.0417)::numeric AS amt)
			UPDATE master_distributors
			SET master_distributor_wallet_balance = master_distributor_wallet_balance + (SELECT amt FROM share),
			    updated_at = NOW()
			WHERE master_distributor_id = $2
		`, commissionStr, *masterDistributorID)
		if err != nil {
			return fmt.Errorf("credit master distributor commission: %w", err)
		}

		var mdName string
		err = tx.QueryRow(ctx, `SELECT master_distributor_name FROM master_distributors WHERE master_distributor_id = $1`, *masterDistributorID).Scan(&mdName)
		if err != nil {
			return fmt.Errorf("fetch master distributor name: %w", err)
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO transactions (
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
				remarks,
				created_at
			) VALUES (
				gen_random_uuid(),
				$1::uuid,
				$1::uuid,
				$2,
				$2,
				'MASTER_DISTRIBUTOR',
				'MASTER_DISTRIBUTOR',
				'CREDIT',
				'PAYOUT',
				($3::numeric * 0.0417)::numeric,
				'SUCCESS',
				('Payout commission credited to master distributor | payout_id=' || $4::text),
				NOW()
			)
		`, *masterDistributorID, mdName, commissionStr, req.PartnerRequestID)
		if err != nil {
			return fmt.Errorf("insert master distributor commission transaction: %w", err)
		}
	}

	// d) Retailer (user) share (50%) — credit back to user
	_, err = tx.Exec(ctx, `
		WITH share AS (SELECT ($1::numeric * 0.50)::numeric AS amt)
		UPDATE users
		SET user_wallet_balance = user_wallet_balance + (SELECT amt FROM share),
		    updated_at = NOW()
		WHERE user_id = $2
	`, commissionStr, userID)
	if err != nil {
		return fmt.Errorf("credit retailer commission: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO transactions (
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
			remarks,
			created_at
		) VALUES (
			gen_random_uuid(),
			$1::uuid,
			$1::uuid,
			$2,
			$2,
			'USER',
			'USER',
			'CREDIT',
			'PAYOUT',
			($3::numeric * 0.50)::numeric,
			'SUCCESS',
			('Payout commission credited to retailer (user) | payout_id=' || $4::text),
			NOW()
		)
	`, userID, userName, commissionStr, req.PartnerRequestID)
	if err != nil {
		return fmt.Errorf("insert retailer commission transaction: %w", err)
	}

	// 7) Finally update payout_service to SUCCESS and set operator_transaction_id, order_id
	_, err = tx.Exec(ctx, `
		UPDATE payout_service
		SET transaction_status = 'SUCCESS',
		    operator_transaction_id = $2,
		    order_id = $3,
		    updated_at = NOW()
		WHERE payout_transaction_id = $1
	`, req.PartnerRequestID, req.OperatorTransactionID, req.OrderID)
	if err != nil {
		return fmt.Errorf("update payout_service status: %w", err)
	}

	// 8) Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (q *Query) PayoutFailure(req *structures.PayoutApiFailureResponse) error {
	ctx := context.Background()

	tx, err := q.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1) Lock payout_service row and read required fields (only if PENDING)
	var (
		userID        string
		amountStr     string
		currentStatus string
	)
	err = tx.QueryRow(ctx, `
		SELECT user_id, amount::text, transaction_status
		FROM payout_service
		WHERE payout_transaction_id = $1
		FOR UPDATE
	`, req.PayoutTransactionID).Scan(&userID, &amountStr, &currentStatus)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("payout not found")
		}
		return fmt.Errorf("select payout_service: %w", err)
	}

	if currentStatus != "PENDING" {
		return fmt.Errorf("payout is not pending (status=%s)", currentStatus)
	}

	// 2) Update payout_service to FAILED and set operator_transaction_id, order_id, remarks
	_, err = tx.Exec(ctx, `
		UPDATE payout_service
		SET transaction_status = 'FAILED',
		    operator_transaction_id = $2,
		    order_id = $3,
		    remarks = $4,
		    updated_at = NOW()
		WHERE payout_transaction_id = $1
	`, req.PayoutTransactionID, "INVALID", "INVALID", "INVALID")
	if err != nil {
		return fmt.Errorf("update payout_service to FAILED: %w", err)
	}

	// 3) Insert a failed transaction record for audit (transactor & receiver are user)
	//    receiver_name set to 'Payout' to indicate payout service
	//    We record the attempted amount and include provider remarks in transaction remarks.
	_, err = tx.Exec(ctx, `
		INSERT INTO transactions (
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
			remarks,
			created_at
		) VALUES (
			gen_random_uuid(),
			$1::uuid,
			$1::uuid,
			(SELECT user_name FROM users WHERE user_id = $1),
			'Payout',
			'USER',
			'USER',
			'DEBIT',
			'PAYOUT',
			$2::numeric,
			'FAILED',
			('Payout failed | payout_id=' || $3::text || ' | provider_remarks=' || COALESCE($4, '') ),
			NOW()
		)
	`, userID, amountStr, req.PayoutTransactionID, "INVALID")
	if err != nil {
		return fmt.Errorf("insert failed payout transaction: %w", err)
	}

	// 4) Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (q *Query) GetPayoutTransactions(userId string) (*[]structures.GetPayoutLogs, error) {
	query := `
		SELECT operator_transaction_id, mobile_number,
		bank_name, beneficiary_name, amount, commision,
		transfer_type, transaction_status, created_at::text
		FROM payout_service
		WHERE user_id=$1;
	`
	res, err := q.Pool.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var payoutTransactions []structures.GetPayoutLogs
	for res.Next() {
		var payoutTransaction structures.GetPayoutLogs
		if err := res.Scan(
			&payoutTransaction.TransactionID,
			&payoutTransaction.PhoneNumber,
			&payoutTransaction.BankName,
			&payoutTransaction.BeneficiaryName,
			&payoutTransaction.Amount,
			&payoutTransaction.Commission,
			&payoutTransaction.TransferType,
			&payoutTransaction.TransactionStatus,
			&payoutTransaction.TransactionDateAndTime,
		); err != nil {
			return nil, err
		}

		payoutTransactions = append(payoutTransactions, payoutTransaction)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return &payoutTransactions, nil
}
