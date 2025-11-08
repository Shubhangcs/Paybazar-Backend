package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) GetAllTransactions(req *structures.TransactionRequest) (*[]structures.TransactionResponse, error) {
	var transactions []structures.TransactionResponse

	query := `
	SELECT *
	FROM (
		SELECT
			transaction_id,
			admin_id AS transactor_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks,
			created_at,
			'ADMIN' AS transactor_type
		FROM admin_wallet_transactions
		WHERE $2 = 'ADMIN'

		UNION ALL

		SELECT
			transaction_id,
			master_distributor_id AS transactor_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks,
			created_at,
			'MASTER_DISTRIBUTOR' AS transactor_type
		FROM master_distributor_wallet_transactions
		WHERE $2 = 'MASTER_DISTRIBUTOR'

		UNION ALL

		SELECT
			transaction_id,
			distributor_id AS transactor_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks,
			created_at,
			'DISTRIBUTOR' AS transactor_type
		FROM distributor_wallet_transactions
		WHERE $2 = 'DISTRIBUTOR'

		UNION ALL

		SELECT
			transaction_id,
			user_id AS transactor_id,
			amount,
			transaction_type,
			transaction_service,
			reference_id,
			remarks,
			created_at,
			'USER' AS transactor_type
		FROM user_wallet_transactions
		WHERE $2 = 'USER'
	) AS all_transactions
	WHERE transactor_id = $1
	ORDER BY created_at DESC;
	`

	rows, err := q.Pool.Query(
		context.Background(),
		query,
		req.TransactorId,
		req.TransactorType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tr structures.TransactionResponse
		if err := rows.Scan(
			&tr.TransactionId,
			&tr.TransactorId,
			&tr.Amount,
			&tr.TransactionType,
			&tr.TransactionService,
			&tr.ReferenceId,
			&tr.Remarks,
			&tr.CreatedAt,
			&tr.TransactorType, // new: included from SELECT
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, tr)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &transactions, nil
}
