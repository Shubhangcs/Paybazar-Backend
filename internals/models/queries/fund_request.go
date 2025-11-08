package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) GetFundRequestsByID(requesterId string) (*[]structures.FundRequest, error) {
	var fundRequest structures.FundRequest
	var fundRequests []structures.FundRequest

	query := `
		SELECT admin_id, requester_id, requester_type, amount, bank_name, account_number, ifsc_code, bank_branch, utr_number, remarks, request_status FROM fund_requests WHERE requester_id=$1
	`
	res, err := q.Pool.Query(context.Background(), query, requesterId)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		if err := res.Scan(
			&fundRequest.AdminId,
			&fundRequest.RequesterId,
			&fundRequest.RequesterType,
			&fundRequest.Amount,
			&fundRequest.BankName,
			&fundRequest.IFSCCode,
			&fundRequest.BankBranch,
			&fundRequest.UTRNumber,
			&fundRequest.Remarks,
			&fundRequest.RequestStatus,
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fundRequest)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}
	return &fundRequests, nil
}

func (q *Query) RejectFundRequest(requestId string) error {

	query := `
	UPDATE fund_requests
	SET 
		request_status = 'REJECTED',
		updated_at = NOW()
	WHERE request_id = $1
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		requestId,
	)
	return err
}

func (q *Query) CreateFundRequest(req *structures.FundRequest) error {
	query := `
	INSERT INTO fund_requests (
		requester_id,
		requester_type,
		amount,
		bank_name,
		account_number,
		ifsc_code,
		bank_branch,
		utr_number,
		request_status,
		remarks
	)
	VALUES (
		$1,        -- requester_id (UUID)
		$2,        -- requester_type ('USER' | 'DISTRIBUTOR' | 'MASTER_DISTRIBUTOR')
		$3,        -- amount
		$4,        -- bank_name
		$5,        -- account_number
		$6,        -- ifsc_code
		$7,        -- bank_branch
		$8,        -- utr_number
		'PENDING', -- default status
		$9         -- remarks
	);
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		req.RequesterId,
		req.RequesterType,
		req.Amount,
		req.BankName,
		req.AccountNumber,
		req.IFSCCode,
		req.BankBranch,
		req.UTRNumber,
		req.Remarks,
	)

	return err
}

func (q *Query) AcceptFundRequest(req *structures.AcceptFundRequest) error {
	query := `
	WITH
	-- Admin and target selectors gated by RequesterType
	sel_admin AS (
		SELECT a.admin_id, a.admin_unique_id
		FROM admins a
		WHERE a.admin_id = $1
	),
	sel_user AS (
		SELECT u.user_id, u.user_unique_id
		FROM users u
		WHERE u.user_id = $2 AND $3 = 'USER'
	),
	sel_distributor AS (
		SELECT d.distributor_id, d.distributor_unique_id
		FROM distributors d
		WHERE d.distributor_id = $2 AND $3 = 'DISTRIBUTOR'
	),
	sel_md AS (
		SELECT m.master_distributor_id, m.master_distributor_unique_id
		FROM master_distributors m
		WHERE m.master_distributor_id = $2 AND $3 = 'MASTER_DISTRIBUTOR'
	),

	-- Deduct from admin only if sufficient balance
	deduct_admin AS (
		UPDATE admin_wallets aw
		SET balance = aw.balance - $4::numeric
		FROM sel_admin sa
		WHERE aw.admin_id = sa.admin_id
		  AND aw.balance >= $4::numeric
		RETURNING aw.admin_id
	),

	-- Admin DEBIT transaction (reference_id = requester's unique_id)
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
			sa.admin_id,
			$4::numeric,
			'DEBIT',
			CASE 
				WHEN $3 = 'USER' THEN 'USER'
				WHEN $3 = 'DISTRIBUTOR' THEN 'DISTRIBUTOR'
				WHEN $3 = 'MASTER_DISTRIBUTOR' THEN 'MD'
			END,
			COALESCE(
				(SELECT su.user_unique_id FROM sel_user su),
				(SELECT sd.distributor_unique_id FROM sel_distributor sd),
				(SELECT sm.master_distributor_unique_id FROM sel_md sm)
			),
			$5
		FROM deduct_admin da
		JOIN sel_admin sa ON TRUE
		RETURNING 1
	),

	-- Credit target wallet (only one branch matches)
	credit_user AS (
		UPDATE user_wallets uw
		SET balance = uw.balance + $4::numeric
		FROM sel_user su
		WHERE uw.user_id = su.user_id
		RETURNING uw.user_id
	),
	credit_distributor AS (
		UPDATE distributor_wallets dw
		SET balance = dw.balance + $4::numeric
		FROM sel_distributor sd
		WHERE dw.distributor_id = sd.distributor_id
		RETURNING dw.distributor_id
	),
	credit_md AS (
		UPDATE master_distributor_wallets mw
		SET balance = mw.balance + $4::numeric
		FROM sel_md sm
		WHERE mw.master_distributor_id = sm.master_distributor_id
		RETURNING mw.master_distributor_id
	),

	-- Target CREDIT transaction (reference_id = admin unique_id)
	user_tx AS (
		INSERT INTO user_wallet_transactions (
			user_id, amount, transaction_type, transaction_service, reference_id, remarks
		)
		SELECT
			cu.user_id,
			$4::numeric,
			'CREDIT',
			'ADMIN',
			(SELECT sa.admin_unique_id FROM sel_admin sa),
			$5
		FROM credit_user cu
		RETURNING 1
	),
	distributor_tx AS (
		INSERT INTO distributor_wallet_transactions (
			distributor_id, amount, transaction_type, transaction_service, reference_id, remarks
		)
		SELECT
			cd.distributor_id,
			$4::numeric,
			'CREDIT',
			'ADMIN',
			(SELECT sa.admin_unique_id FROM sel_admin sa),
			$5
		FROM credit_distributor cd
		RETURNING 1
	),
	md_tx AS (
		INSERT INTO master_distributor_wallet_transactions (
			master_distributor_id, amount, transaction_type, transaction_service, reference_id, remarks
		)
		SELECT
			cm.master_distributor_id,
			$4::numeric,
			'CREDIT',
			'ADMIN',
			(SELECT sa.admin_unique_id FROM sel_admin sa),
			$5
		FROM credit_md cm
		RETURNING 1
	),

	-- Accept the most recent matching PENDING fund request
	-- This runs ONLY if admin deduction succeeded (joined via deduct_admin)
	upd_fund_req AS (
		UPDATE fund_requests fr
		SET request_status = 'ACCEPTED',
			updated_at = NOW()
		WHERE fr.request_id = (
			SELECT fr2.request_id
			FROM fund_requests fr2
			WHERE fr2.requester_id = $2
			  AND fr2.requester_type = $3
			  AND fr2.amount = $4::numeric
			  AND fr2.request_status = 'PENDING'
			ORDER BY fr2.created_at DESC
			LIMIT 1
		)
		AND EXISTS (SELECT 1 FROM deduct_admin)  -- ensure debit happened
		RETURNING 1
	)

	SELECT 1;
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		req.AdminId,       // $1
		req.RequesterId,   // $2
		req.RequesterType, // $3: 'USER' | 'DISTRIBUTOR' | 'MASTER_DISTRIBUTOR'
		req.Amount,        // $4 (cast to numeric in SQL)
		req.Remarks,       // $5
	)
	return err
}
