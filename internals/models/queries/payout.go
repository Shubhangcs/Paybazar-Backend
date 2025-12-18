package queries

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
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

func (q *Query) CheckMpin(userID string, mpin string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1 AND user_mpin = $2)`
	var hasMpin bool
	err := q.Pool.QueryRow(context.Background(), query, userID, mpin).Scan(&hasMpin)
	return hasMpin, err
}

func (q *Query) InitilizePayoutRequest(req *structures.PayoutInitilizationRequest) (*structures.PayoutApiRequest, error) {

	const getUserDetailsQuery = `
		SELECT admin_id, master_distributor_id, distributor_id
		FROM users
		WHERE user_id=$1;
	`

	const checkCommisionExists = `
		SELECT EXISTS(SELECT 1 FROM commisions WHERE distributor_id=$1) AS is_user_commision;
	`

	const getCommisions = `
		SELECT admin_commision::TEXT, master_distributor_commision::TEXT,
		distributor_commision::TEXT, user_commision::TEXT, commision::TEXT
		FROM commisions
		WHERE distributor_id=$1;
	`

	const insertPayoutTransactionQuery = `
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
	VALUES ($1,$2,$3,$4,$5,$6,$7,UPPER($8),'PENDING',$9,($7::NUMERIC * ($10::NUMERIC / 100::NUMERIC )))
	RETURNING 
		payout_transaction_id::TEXT,
		mobile_number AS mobile_no,
		account_number AS account_no,
		ifsc_code AS ifsc,
		bank_name,
		beneficiary_name AS benificiary_name,
		amount::TEXT,
		commision,
		CASE 
			WHEN UPPER(transfer_type) = 'IMPS' THEN '5'
			WHEN UPPER(transfer_type) = 'NEFT' THEN '6'
		END AS transfer_type;
	`

	const updateUserWalletBalance = `
		UPDATE users SET
		user_wallet_balance = user_wallet_balance - ($1::NUMERIC + ($2::NUMERIC * $3::NUMERIC))
		WHERE user_id=$4;
	`

	const updateDistributorWalletBalance = `
		UPDATE distributors SET
		distributor_wallet_balance = distributor_wallet_balance + ($1::NUMERIC * $2::NUMERIC)
		WHERE distributor_id=$3;
	`

	const updateMasterDistributorWalletBalance = `
		UPDATE master_distributors SET
		master_distributor_wallet_balance = master_distributor_wallet_balance + ($1::NUMERIC * $2::NUMERIC)
		WHERE master_distributor_id=$3;
	`

	const updateAdminWalletBalance = `
		UPDATE admins SET
		admin_wallet_balance = admin_wallet_balance + ($1::NUMERIC * $2::NUMERIC)
		WHERE admin_id=$3;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := q.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var UserDetails struct {
		adminID             string
		masterDistributorID string
		distributorID       string
	}

	if err := tx.QueryRow(ctx, getUserDetailsQuery, req.UserID).Scan(
		&UserDetails.adminID,
		&UserDetails.masterDistributorID,
		&UserDetails.distributorID,
	); err != nil {
		return nil, err
	}

	var hasCommision bool
	if err := tx.QueryRow(ctx, checkCommisionExists, UserDetails.distributorID).Scan(&hasCommision); err != nil {
		return nil, err
	}

	var Commision struct {
		AdminCommision             string
		MasterDistributorCommision string
		DistributorCommision       string
		RetailerCommision          string
		TotalCommision             string
	}
	if !hasCommision {
		Commision.TotalCommision = "1.2"
		Commision.AdminCommision = "0.2917"
		Commision.MasterDistributorCommision = "0.0417"
		Commision.DistributorCommision = "0.1667"
		Commision.RetailerCommision = "0.50"
	} else {
		if err := tx.QueryRow(ctx, getCommisions, UserDetails.distributorID).Scan(
			&Commision.AdminCommision,
			&Commision.MasterDistributorCommision,
			&Commision.DistributorCommision,
			&Commision.RetailerCommision,
			&Commision.TotalCommision,
		); err != nil {
			return nil, err
		}
	}

	var commision string

	var res structures.PayoutApiRequest
	if err := tx.QueryRow(
		ctx,
		insertPayoutTransactionQuery,
		req.UserID,
		req.MobileNumber,
		req.AccountNumber,
		req.IFSCCode,
		req.BankName,
		req.BeneficiaryName,
		req.Amount,
		req.TransferType,
		req.Remarks,
		Commision.TotalCommision,
	).Scan(
		&res.PartnerRequestID,
		&res.MobileNumber,
		&res.AccountNumber,
		&res.IFSCCode,
		&res.BankName,
		&res.BeneficiaryName,
		&res.Amount,
		&res.TransferType,
		&commision,
	); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, updateUserWalletBalance, req.Amount, commision, Commision.RetailerCommision, req.UserID); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, updateAdminWalletBalance, commision, Commision.AdminCommision, UserDetails.adminID); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, updateMasterDistributorWalletBalance, commision, Commision.MasterDistributorCommision, UserDetails.masterDistributorID); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, updateDistributorWalletBalance, commision, Commision.DistributorCommision, UserDetails.distributorID); err != nil {
		return nil, err
	}

	return &res, tx.Commit(ctx)
}
func (q *Query) FinalPayout(req *structures.PayoutFinal) error {
	query := `
		UPDATE payout_service SET
		operator_transaction_id=$1, order_id=$2,
		transaction_status=$3
		WHERE payout_transaction_id=$4;
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if req.Status == 1 {
		if _, err := q.Pool.Exec(ctx, query, req.OpertaorTransactionID, req.OrderID, "SUCCESS", req.PartnerRequestID); err != nil {
			return err
		}
	}
	if req.Status == 2 {
		if _, err := q.Pool.Exec(ctx, query, req.OpertaorTransactionID, req.OrderID, "PENDING", req.PartnerRequestID); err != nil {
			return err
		}
	}
	if req.Status == 3 {
		if _, err := q.Pool.Exec(ctx, query, req.OpertaorTransactionID, req.OrderID, "FAILED", req.PartnerRequestID); err != nil {
			return err
		}
	}
	return nil
}

func (q *Query) GetPayoutTransactions(userId string) (*[]structures.GetPayoutLogs, error) {
	query := `
		SELECT payout_transaction_id,operator_transaction_id, mobile_number,
		bank_name, beneficiary_name, amount, commision,
		transfer_type, transaction_status, created_at::text, account_number
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
			&payoutTransaction.PayoutTransactionID,
			&payoutTransaction.TransactionID,
			&payoutTransaction.PhoneNumber,
			&payoutTransaction.BankName,
			&payoutTransaction.BeneficiaryName,
			&payoutTransaction.Amount,
			&payoutTransaction.Commission,
			&payoutTransaction.TransferType,
			&payoutTransaction.TransactionStatus,
			&payoutTransaction.TransactionDateAndTime,
			&payoutTransaction.AccountNumber,
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

func (q *Query) DeductUserBalanceForVerification(userId string) error {
	ctx := context.Background()
	tx, err := q.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		WITH admin_user AS (
			SELECT admin_id 
			FROM users 
			WHERE user_id=$1
			LIMIT 1
		),
		deduct AS (
			UPDATE users
			SET user_wallet_balance = user_wallet_balance - 3
			WHERE user_id = $1 AND user_wallet_balance >= 3
			RETURNING 3 AS amount_deducted
		)
		UPDATE admins
		SET admin_wallet_balance = admin_wallet_balance + (SELECT amount_deducted FROM deduct)
		WHERE admin_id = (SELECT admin_id FROM admin_user);
	`

	// Execute the combined transaction query
	cmdTag, err := tx.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("verification wallet transaction failed: %w", err)
	}

	// If nothing was deducted, it means user had insufficient balance
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("insufficient balance")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (q *Query) PayoutTransactionRefund(req *structures.PayoutRefund) error {
	getUserDetailsQuery := `
		SELECT user_id, amount, commision FROM payout_service WHERE payout_transaction_id=$1;
	`
	getAllIdsFromUser := `
		SELECT master_distributor_id, distributor_id, admin_id FROM users WHERE user_id=$1;
	`
	adminCutAmount := `
		WITH admin_commision AS(
			SELECT ($1::numeric * 0.2917)::numeric AS amt
		)
		UPDATE admins SET admin_wallet_balance = admin_wallet_balance - (SELECT amt FROM admin_commision)
		WHERE admin_id=$2 AND admin_wallet_balance >= (SELECT amt FROM admin_commision);
	`
	masterDistributorCutAmount := `
		WITH master_distributor_commision AS(
			SELECT ($1::NUMERIC *  0.0417)::numeric as amt
		)
		UPDATE master_distributors SET master_distributor_wallet_balance = master_distributor_wallet_balance - (SELECT amt FROM master_distributor_commision)
		WHERE master_distributor_id=$2 AND master_distributor_wallet_balance >= (SELECT amt FROM master_distributor_commision);
	`
	distributorCutAmount := `
		WITH distributor_commision AS(
			SELECT ($1::NUMERIC * 0.1667)::numeric AS amt
		)
		UPDATE distributors SET distributor_wallet_balance = distributor_wallet_balance - (SELECT amt FROM distributor_commision)
		WHERE distributor_id=$2 AND distributor_wallet_balance >= (SELECT amt FROM distributor_commision);
	`
	addAmountToUser := `
		WITH user_commision_cut AS(
			SELECT ($1::numeric * 0.5)::numeric - $1 AS amt
		)
		UPDATE users SET user_wallet_balance = user_wallet_balance + (SELECT amt FROM user_commision_cut) + $3::numeric
		WHERE user_id=$2;
	`
	updateTransactionStatus := `
		UPDATE payout_service SET transaction_status = 'REFUND'
		WHERE payout_transaction_id=$1;
	`

	tx, err := q.Pool.Begin(context.Background())
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to refund database error")
	}
	defer func() { _ = tx.Rollback(context.Background()) }()

	var transactionDetails struct {
		UserID     string
		Amount     string
		Commission string
	}
	if err := tx.QueryRow(context.Background(), getUserDetailsQuery, req.TransactionID).Scan(
		&transactionDetails.UserID,
		&transactionDetails.Amount,
		&transactionDetails.Commission,
	); err != nil {
		log.Println(err)
		return fmt.Errorf("failed to execuite transaction")
	}

	var usersDetails struct {
		MasterDistributorID string
		DistributorID       string
		AdminID             string
	}

	if err := tx.QueryRow(context.Background(), getAllIdsFromUser, transactionDetails.UserID).Scan(
		&usersDetails.MasterDistributorID,
		&usersDetails.DistributorID,
		&usersDetails.AdminID,
	); err != nil {
		log.Println(err)
		return fmt.Errorf("failed to get user details")
	}

	aCut, err := tx.Exec(context.Background(), adminCutAmount, transactionDetails.Commission, usersDetails.AdminID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to execuite admin deduct transaction")
	}

	if aCut.RowsAffected() == 0 {
		return fmt.Errorf("insufficient balance in admin")
	}

	mCut, err := tx.Exec(context.Background(), masterDistributorCutAmount, transactionDetails.Commission, usersDetails.MasterDistributorID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to execuite md deduct transaction")
	}

	if mCut.RowsAffected() == 0 {
		return fmt.Errorf("insufficient balance in md")
	}

	dCut, err := tx.Exec(context.Background(), distributorCutAmount, transactionDetails.Commission, usersDetails.DistributorID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to execuite md deduct transaction")
	}

	if dCut.RowsAffected() == 0 {
		return fmt.Errorf("insufficient balance in distributor")
	}

	_, err = tx.Exec(context.Background(), addAmountToUser, transactionDetails.Commission, transactionDetails.UserID, transactionDetails.Amount)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to add amount to user")
	}

	_, err = tx.Exec(context.Background(), updateTransactionStatus, req.TransactionID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to update transaction status")
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Println(err)
		return fmt.Errorf("failed to commit transaction")
	}
	return nil
}

func (q *Query) UpdatePayoutTransaction(req *structures.UpdatePayoutTransaction) error {
	query := `
		UPDATE payout_service SET
		operator_transaction_id=$1,transaction_status=$2
		WHERE payout_transaction_id=$3;
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := q.Pool.Exec(ctx, query, req.OperatorTransactionID, req.Status, req.PayoutTransactionID); err != nil {
		return err
	}
	return nil
}
