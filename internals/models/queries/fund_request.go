package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/jackc/pgx/v5"
)

func (q *Query) GetFundRequestsById(requesterId string) (*[]structures.GetFundRequestModel, error) {
	var fundRequests []structures.GetFundRequestModel

	const query = `
		SELECT
			request_unique_id,
			requester_unique_id,
			requester_name,
			requester_type,
			amount,
			bank_name,
			utr_number,
			request_date,
			request_status,
			remarks
		FROM 
			fund_requests
		WHERE 
			requester_id = $1
		ORDER BY 
			created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, requesterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fr structures.GetFundRequestModel
		if err := rows.Scan(
			&fr.RequestUniqueId,   // request_unique_id
			&fr.RequesterUniqueId, // requester_unique_id
			&fr.RequesterName,     // requester_name
			&fr.RequesterType,     // requester_type
			&fr.Amount,            // amount
			&fr.BankName,          // bank_name
			&fr.UTRNumber,         // utr_number
			&fr.RequestDate,
			&fr.RequestStatus,     // request_status
			&fr.Remarks,           // remarks
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &fundRequests, nil
}

func (q *Query) GetAllFundRequests(adminId string) (*[]structures.GetFundRequestModel, error) {
	var fundRequests []structures.GetFundRequestModel

	const query = `
		SELECT
			request_id,
			request_unique_id,
			requester_id,
			requester_unique_id,
			requester_name,
			requester_type,
			amount,
			bank_name,
			utr_number,
			request_date,
			request_status,
			remarks
		FROM 
			fund_requests
		WHERE 
			admin_id = $1
		ORDER BY 
			created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, adminId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fr structures.GetFundRequestModel
		if err := rows.Scan(
			&fr.RequestId,         // request_id
			&fr.RequestUniqueId,   // request_unique_id
			&fr.RequesterId,       // requester_id
			&fr.RequesterUniqueId, // requester_unique_id
			&fr.RequesterName,     // requester_name
			&fr.RequesterType,     // requester_type
			&fr.Amount,            // amount
			&fr.BankName,       // bank_branch
			&fr.UTRNumber,         // utr_number
			&fr.RequestDate,
			&fr.RequestStatus,     // request_status
			&fr.Remarks,           // remarks
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &fundRequests, nil
}

func (q *Query) RejectFundRequest(requestId string) error {
	const query = `
		UPDATE fund_requests
		SET 
			request_status = 'REJECTED',
			updated_at = NOW()
		WHERE 
			request_id = $1
			AND request_status = 'PENDING';
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		requestId,
	)
	return err
}

func (q *Query) CreateFundRequest(req *structures.CreateFundRequestModel) error {
	const query = `
		INSERT INTO fund_requests (
			admin_id,
			requester_id,
			requester_unique_id,
			requester_name,
			requester_type,
			amount,
			bank_name,
			request_date,
			utr_number,
			remarks,
			request_status
		)
		VALUES (
			$1,  -- admin_id
			$2,  -- requester_id
			$3,  -- requester_unique_id
			$4,  -- requester_name
			$5,  -- requester_type
			$6,  -- amount
			$7,  -- bank_name
			$8, -- request date
			$9, -- utr_number
			$10, -- remarks
			'PENDING'  -- default status for new requests
		);
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		req.AdminID,           // $1
		req.RequesterID,       // $2
		req.RequesterUniqueID, // $3
		req.RequesterName,     // $4
		req.RequesterType,     // $5
		req.Amount,            // $6
		req.BankName,          // $7
		req.RequestDate,
		req.UTRNumber,         // $11
		req.Remarks,           // $12
	)

	return err
}


func (q *Query) AcceptFundRequest(req *structures.AcceptFundRequestModel) error {
	ctx := context.Background()

	tx, err := q.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // safe no-op if already committed
	}()

	// 1) Lock and read fund_request (ensure it's PENDING)
	var (
		requesterID       string
		requesterType     string
		requesterUniqueID string
		requesterName     string
		amountStr         string // numeric as text
	)
	err = tx.QueryRow(ctx, `
		SELECT requester_id, requester_type, requester_unique_id, requester_name, amount
		FROM fund_requests
		WHERE request_id = $1 AND admin_id = $2 AND request_status = 'PENDING'
		FOR UPDATE
	`, req.RequestID, req.AdminID).Scan(&requesterID, &requesterType, &requesterUniqueID, &requesterName, &amountStr)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("fund request not found or not pending for given request_id and admin_id")
		}
		return fmt.Errorf("select fund_request: %w", err)
	}

	// 2) Lock admin row and get admin name + ensure balance exists
	var adminBalanceStr string
	var adminName string
	err = tx.QueryRow(ctx, `
		SELECT admin_wallet_balance::text, admin_name
		FROM admins
		WHERE admin_id = $1
		FOR UPDATE
	`, req.AdminID).Scan(&adminBalanceStr, &adminName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("admin not found")
		}
		return fmt.Errorf("select admin: %w", err)
	}

	// 3) Ensure admin has sufficient balance
	var sufficient bool
	err = tx.QueryRow(ctx, `
		SELECT (admin_wallet_balance >= $1::numeric)
		FROM admins
		WHERE admin_id = $2
	`, amountStr, req.AdminID).Scan(&sufficient)
	if err != nil {
		return fmt.Errorf("check admin balance: %w", err)
	}
	if !sufficient {
		return fmt.Errorf("admin has insufficient wallet balance")
	}

	// 4) Lock requester row and read balance for update (just to ensure existence)
	var requesterBalanceStr string
	switch requesterType {
	case "USER":
		err = tx.QueryRow(ctx, `
			SELECT user_wallet_balance::text
			FROM users
			WHERE user_id = $1
			FOR UPDATE
		`, requesterID).Scan(&requesterBalanceStr)
	case "DISTRIBUTOR":
		err = tx.QueryRow(ctx, `
			SELECT distributor_wallet_balance::text
			FROM distributors
			WHERE distributor_id = $1
			FOR UPDATE
		`, requesterID).Scan(&requesterBalanceStr)
	case "MASTER_DISTRIBUTOR":
		err = tx.QueryRow(ctx, `
			SELECT master_distributor_wallet_balance::text
			FROM master_distributors
			WHERE master_distributor_id = $1
			FOR UPDATE
		`, requesterID).Scan(&requesterBalanceStr)
	default:
		return fmt.Errorf("unknown requester_type: %s", requesterType)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("requester not found")
		}
		return fmt.Errorf("select requester wallet: %w", err)
	}

	// 5) Update admin wallet: subtract amount
	_, err = tx.Exec(ctx, `
		UPDATE admins
		SET admin_wallet_balance = admin_wallet_balance - $1::numeric,
		    updated_at = NOW()
		WHERE admin_id = $2
	`, amountStr, req.AdminID)
	if err != nil {
		return fmt.Errorf("update admin wallet: %w", err)
	}

	// 6) Update requester wallet: add amount
	switch requesterType {
	case "USER":
		_, err = tx.Exec(ctx, `
			UPDATE users
			SET user_wallet_balance = user_wallet_balance + $1::numeric,
			    updated_at = NOW()
			WHERE user_id = $2
		`, amountStr, requesterID)
	case "DISTRIBUTOR":
		_, err = tx.Exec(ctx, `
			UPDATE distributors
			SET distributor_wallet_balance = distributor_wallet_balance + $1::numeric,
			    updated_at = NOW()
			WHERE distributor_id = $2
		`, amountStr, requesterID)
	case "MASTER_DISTRIBUTOR":
		_, err = tx.Exec(ctx, `
			UPDATE master_distributors
			SET master_distributor_wallet_balance = master_distributor_wallet_balance + $1::numeric,
			    updated_at = NOW()
			WHERE master_distributor_id = $2
		`, amountStr, requesterID)
	}
	if err != nil {
		return fmt.Errorf("update requester wallet: %w", err)
	}

	// 7) Mark fund_request as APPROVED (ensure we update the same row)
	res, err := tx.Exec(ctx, `
		UPDATE fund_requests
		SET request_status = 'APPROVED', updated_at = NOW()
		WHERE request_id = $1
	`, req.RequestID)
	if err != nil {
		return fmt.Errorf("update fund_request status: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("failed to update fund_request status")
	}

	// 8) Insert two transactions (each with its own transaction_id)
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
			gen_random_uuid(),         -- unique transaction_id for admin debit
			$2::uuid,                  -- transactor_id (admin)
			$3::uuid,                  -- receiver_id (requester)
			$4,                        -- transactor_name (admin_name)
			$5,                        -- receiver_name (requester_name)
			'ADMIN',                   -- transactor_type
			$6,                        -- receiver_type (USER|DISTRIBUTOR|MASTER_DISTRIBUTOR)
			'DEBIT',                   -- transaction_type
			'FUND_REQUEST',            -- transaction_service
			$7::numeric,               -- amount
			'SUCCESS',                 -- transaction_status
			('Fund request approved - admin debit | request_id=' || $1::text),
			NOW()
		), (
			gen_random_uuid(),         -- unique transaction_id for requester credit
			$3::uuid,                  -- transactor_id (requester)
			$2::uuid,                  -- receiver_id (admin)
			$5,                        -- transactor_name (requester_name)
			$4,                        -- receiver_name (admin_name)
			$6,                        -- transactor_type (USER|DISTRIBUTOR|MASTER_DISTRIBUTOR)
			'ADMIN',                   -- receiver_type
			'CREDIT',                  -- transaction_type
			'FUND_REQUEST',
			$7::numeric,
			'SUCCESS',
			('Fund request approved - requester credit | request_id=' || $1::text),
			NOW()
		);
	`, req.RequestID, req.AdminID, requesterID, adminName, requesterName, requesterType, amountStr)
	if err != nil {
		return fmt.Errorf("insert transactions: %w", err)
	}

	// 9) Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}