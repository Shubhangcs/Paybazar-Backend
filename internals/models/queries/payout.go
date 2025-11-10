package queries

import (
	"context"
	"errors"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/jackc/pgx/v5"
)

func (q *Query) CheckUserBalance(userId string, amount string) (bool, error) {
	var hasBalance bool
	query := `SELECT 
  CASE 
    WHEN balance >= $2::numeric THEN TRUE 
    ELSE FALSE 
  END AS has_sufficient_balance
FROM user_wallets
WHERE user_id = $1;
`
	err := q.Pool.QueryRow(context.Background(), query, userId, amount).Scan(&hasBalance)
	return hasBalance, err
}

func (q *Query) CheckPayoutLimit(userId string, amount string) (bool, error) {
	var hasPayoutLimit bool
	query := `
	SELECT 
  CASE 
    WHEN COALESCE(SUM(amount), 0) + $2::numeric > 25000 THEN FALSE
    ELSE TRUE
  END AS can_transact
FROM payout_service
WHERE user_id = $1
  AND DATE(created_at) = CURRENT_DATE;
`
	err := q.Pool.QueryRow(context.Background(), query, userId, amount).Scan(&hasPayoutLimit)
	return hasPayoutLimit, err
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
		transaction_id::TEXT,
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

var ErrPayoutSuccessNotApplied = errors.New("payout success not applied (missing/invalid transaction or zero rows affected)")

func (q *Query) PayoutSuccess(req *structures.PayoutApiSuccessResponse) error {
	const sql = `
WITH sel AS (
  SELECT 
    p.transaction_id,
    p.user_id,
    p.amount::numeric AS amt,
    p.commision::numeric AS commission
  FROM payout_service p
  WHERE p.transaction_id = $1
  FOR UPDATE
),
-- update payout with operator_transaction_id and order_id and mark SUCCESS
upd_payout AS (
  UPDATE payout_service p
  SET operator_transaction_id = $2,
      order_id               = $3,
      transaction_status     = 'SUCCESS',
      updated_at             = NOW()
  FROM sel s
  WHERE p.transaction_id = s.transaction_id
  RETURNING 1
),
u AS (
  SELECT
    s.transaction_id,
    s.user_id,
    s.amt,
    s.commission,
    usr.user_unique_id,
    usr.admin_id,
    usr.master_distributor_id,
    usr.distributor_id
  FROM sel s
  JOIN users usr ON usr.user_id = s.user_id
),
splits AS (
  SELECT
    u.*,
    ROUND(u.commission * 0.50, 2) AS user_share,
    ROUND(u.commission * 0.20, 2) AS distributor_share,
    ROUND(u.commission * 0.05, 2) AS md_share,
    (u.commission - (ROUND(u.commission * 0.50, 2)
                   + ROUND(u.commission * 0.20, 2)
                   + ROUND(u.commission * 0.05, 2))) AS admin_share
  FROM u
),

-- 1) Deduct payout amount from USER wallet
user_deduct AS (
  UPDATE user_wallets uw
  SET balance = uw.balance - s.amt
  FROM splits s
  WHERE uw.user_id = s.user_id
  RETURNING 1
),
-- User DEBIT tx for payout amount
user_debit_tx AS (
  INSERT INTO user_wallet_transactions (
    user_id, amount, transaction_type, transaction_service, reference_id, remarks
  )
  SELECT
    s.user_id,
    s.amt,
    'DEBIT',
    'PAYOUT',
    s.user_unique_id,
    'SUCCESS'
  FROM splits s
  RETURNING 1
),

-- 2) Credit USER with 50% commission
user_comm_credit AS (
  UPDATE user_wallets uw
  SET balance = uw.balance + s.user_share
  FROM splits s
  WHERE uw.user_id = s.user_id
  RETURNING 1
),
user_comm_tx AS (
  INSERT INTO user_wallet_transactions (
    user_id, amount, transaction_type, transaction_service, reference_id, remarks
  )
  SELECT
    s.user_id,
    s.user_share,
    'CREDIT',
    'PAYOUT',
    s.user_unique_id,
    'COMMISSION'
  FROM splits s
  RETURNING 1
),

-- 3) Credit ADMIN share
adm_upd AS (
  UPDATE admin_wallets aw
  SET balance = aw.balance + s.admin_share
  FROM splits s
  WHERE aw.admin_id = s.admin_id
  RETURNING 1
),
adm_tx AS (
  INSERT INTO admin_wallet_transactions (
    admin_id, amount, transaction_type, transaction_service, reference_id, remarks
  )
  SELECT
    s.admin_id,
    s.admin_share,
    'CREDIT',
    'PAYOUT',
    s.user_unique_id,
    'COMMISSION'
  FROM splits s
  RETURNING 1
),

-- 4) Credit MASTER DISTRIBUTOR share
md_upd AS (
  UPDATE master_distributor_wallets mw
  SET balance = mw.balance + s.md_share
  FROM splits s
  WHERE mw.master_distributor_id = s.master_distributor_id
  RETURNING 1
),
md_tx AS (
  INSERT INTO master_distributor_wallet_transactions (
    master_distributor_id, amount, transaction_type, transaction_service, reference_id, remarks
  )
  SELECT
    s.master_distributor_id,
    s.md_share,
    'CREDIT',
    'PAYOUT',
    s.user_unique_id,
    'COMMISSION'
  FROM splits s
  RETURNING 1
),

-- 5) Credit DISTRIBUTOR share
dist_upd AS (
  UPDATE distributor_wallets dw
  SET balance = dw.balance + s.distributor_share
  FROM splits s
  WHERE dw.distributor_id = s.distributor_id
  RETURNING 1
),
dist_tx AS (
  INSERT INTO distributor_wallet_transactions (
    distributor_id, amount, transaction_type, transaction_service, reference_id, remarks
  )
  SELECT
    s.distributor_id,
    s.distributor_share,
    'CREDIT',
    'PAYOUT',
    s.user_unique_id,
    'COMMISSION'
  FROM splits s
  RETURNING 1
)

-- If ANY of the above CTEs returns 0 rows, this JOIN chain produces no row -> ErrNoRows
SELECT 1
FROM sel
JOIN upd_payout           ON TRUE
JOIN user_deduct          ON TRUE
JOIN user_debit_tx        ON TRUE
JOIN user_comm_credit     ON TRUE
JOIN user_comm_tx         ON TRUE
JOIN adm_upd              ON TRUE
JOIN adm_tx               ON TRUE
JOIN md_upd               ON TRUE
JOIN md_tx                ON TRUE
JOIN dist_upd             ON TRUE
JOIN dist_tx              ON TRUE;
`
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var ok int
	err := q.Pool.QueryRow(
		ctx,
		sql,
		req.PartnerRequestID,      // $1 transaction_id
		req.OperatorTransactionID, // $2
		req.OrderID,               // $3
	).Scan(&ok)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPayoutSuccessNotApplied
		}
		return err
	}
	return nil
}

func (q *Query) PayoutFailure(transactionId string) error {
	const query = `
		UPDATE payout_service
		SET 
			transaction_status = 'FAILED',
			updated_at = NOW()
		WHERE transaction_id = $1
		  AND transaction_status = 'PENDING';
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := q.Pool.Exec(ctx, query, transactionId)
	return err
}
