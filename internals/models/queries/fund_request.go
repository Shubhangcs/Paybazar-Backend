package queries

import (
	"context"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) GetFundRequestsById(requesterId string) (*[]structures.FundRequest, error) {
	var fundRequests []structures.FundRequest

	query := `
	SELECT 
		admin_id,
		request_id,
		requester_id,
		requester_type,
		amount,
		bank_name,
		account_number,
		ifsc_code,
		bank_branch,
		utr_number,
		remarks,
		request_status
	FROM fund_requests
	WHERE requester_id = $1;
	`

	rows, err := q.Pool.Query(context.Background(), query, requesterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fr structures.FundRequest
		if err := rows.Scan(
			&fr.AdminId,
			&fr.RequestId,
			&fr.RequesterId,
			&fr.RequesterType,
			&fr.Amount,
			&fr.BankName,
			&fr.AccountNumber,
			&fr.IFSCCode,
			&fr.BankBranch,
			&fr.UTRNumber,
			&fr.Remarks,
			&fr.RequestStatus,
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fr)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &fundRequests, nil
}

func (q *Query) GetAllFundRequests(adminId string) (*[]structures.FundRequest, error) {
	var fundRequests []structures.FundRequest

	query := `
	SELECT 
		admin_id,
		request_id,
		requester_id,
		requester_type,
		amount,
		bank_name,
		account_number,
		ifsc_code,
		bank_branch,
		utr_number,
		remarks,
		request_status
	FROM fund_requests
	WHERE admin_id=$1;
	`

	rows, err := q.Pool.Query(context.Background(), query,adminId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fr structures.FundRequest
		if err := rows.Scan(
			&fr.AdminId,
			&fr.RequestId,
			&fr.RequesterId,
			&fr.RequesterType,
			&fr.Amount,
			&fr.BankName,
			&fr.AccountNumber,
			&fr.IFSCCode,
			&fr.BankBranch,
			&fr.UTRNumber,
			&fr.Remarks,
			&fr.RequestStatus,
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fr)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
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
		admin_id,
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
		$10,
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
		req.AdminId,
	)

	return err
}

func (q *Query) AcceptFundRequest(req *structures.AcceptFundRequest) error {
	const sql = `
	WITH sel_admin AS (
		SELECT a.admin_id, a.admin_unique_id
		FROM admins a
		WHERE a.admin_id = $1
	),
	sel_req AS (
		SELECT fr.request_id, fr.requester_id, fr.requester_type, fr.amount, fr.remarks
		FROM fund_requests fr
		WHERE fr.request_id = $2
		  AND fr.request_status = 'PENDING'
		FOR UPDATE
	),
	sel_user AS (
		SELECT u.user_id, u.user_unique_id
		FROM users u
		JOIN sel_req r ON r.requester_type = 'USER' AND u.user_id = r.requester_id
	),
	sel_distributor AS (
		SELECT d.distributor_id, d.distributor_unique_id
		FROM distributors d
		JOIN sel_req r ON r.requester_type = 'DISTRIBUTOR' AND d.distributor_id = r.requester_id
	),
	sel_md AS (
		SELECT m.master_distributor_id, m.master_distributor_unique_id
		FROM master_distributors m
		JOIN sel_req r ON r.requester_type = 'MASTER_DISTRIBUTOR' AND m.master_distributor_id = r.requester_id
	),
	deduct_admin AS (
		UPDATE admin_wallets aw
		SET balance = aw.balance - r.amount
		FROM sel_admin sa, sel_req r
		WHERE aw.admin_id = sa.admin_id
		  AND aw.balance >= r.amount
		RETURNING aw.admin_id
	),
	admin_tx AS (
		INSERT INTO admin_wallet_transactions (
			admin_id, amount, transaction_type, transaction_service, reference_id, remarks
		)
		SELECT
			sa.admin_id,
			r.amount,
			'DEBIT',
			CASE r.requester_type
				WHEN 'USER' THEN 'USER'
				WHEN 'DISTRIBUTOR' THEN 'DISTRIBUTOR'
				WHEN 'MASTER_DISTRIBUTOR' THEN 'MD'
			END,
			COALESCE(
				(SELECT su.user_unique_id FROM sel_user su),
				(SELECT sd.distributor_unique_id FROM sel_distributor sd),
				(SELECT sm.master_distributor_unique_id FROM sel_md sm)
			),
			r.remarks
		FROM sel_admin sa, sel_req r
		JOIN deduct_admin d ON TRUE
		RETURNING 1
	),
	credit_user AS (
		UPDATE user_wallets uw
		SET balance = uw.balance + r.amount
		FROM sel_user su, sel_req r
		WHERE uw.user_id = su.user_id
		RETURNING uw.user_id
	),
	credit_distributor AS (
		UPDATE distributor_wallets dw
		SET balance = dw.balance + r.amount
		FROM sel_distributor sd, sel_req r
		WHERE dw.distributor_id = sd.distributor_id
		RETURNING dw.distributor_id
	),
	credit_md AS (
		UPDATE master_distributor_wallets mw
		SET balance = mw.balance + r.amount
		FROM sel_md sm, sel_req r
		WHERE mw.master_distributor_id = sm.master_distributor_id
		RETURNING mw.master_distributor_id
	),
	user_tx AS (
		INSERT INTO user_wallet_transactions (
			user_id, amount, transaction_type, transaction_service, reference_id, remarks
		)
		SELECT
			cu.user_id,
			r.amount,
			'CREDIT',
			'ADMIN',
			sa.admin_unique_id,
			r.remarks
		FROM credit_user cu, sel_admin sa, sel_req r
		RETURNING 1
	),
	distributor_tx AS (
		INSERT INTO distributor_wallet_transactions (
			distributor_id, amount, transaction_type, transaction_service, reference_id, remarks
		)
		SELECT
			cd.distributor_id,
			r.amount,
			'CREDIT',
			'ADMIN',
			sa.admin_unique_id,
			r.remarks
		FROM credit_distributor cd, sel_admin sa, sel_req r
		RETURNING 1
	),
	md_tx AS (
		INSERT INTO master_distributor_wallet_transactions (
			master_distributor_id, amount, transaction_type, transaction_service, reference_id, remarks
		)
		SELECT
			cm.master_distributor_id,
			r.amount,
			'CREDIT',
			'ADMIN',
			sa.admin_unique_id,
			r.remarks
		FROM credit_md cm, sel_admin sa, sel_req r
		RETURNING 1
	),
	upd_fund_req AS (
		UPDATE fund_requests fr
		SET request_status = 'APPROVED',
			updated_at = NOW()
		FROM sel_req r
		WHERE fr.request_id = r.request_id
		  AND EXISTS (SELECT 1 FROM deduct_admin)
		RETURNING 1
	)
	SELECT 1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := q.Pool.Exec(
		ctx,
		sql,
		req.AdminId,   // $1
		req.RequestId, // $2
	)
	return err
}
