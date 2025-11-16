package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) AddNewBank(req *structures.BankModel) error {
	query := `INSERT INTO banks (bank_name, ifsc_code) VALUES ($1, $2)`
	_, err := q.Pool.Exec(context.Background(), query, req.BankName, req.IFSCCode)
	return err
}

func (q *Query) GetBanks() (*[]structures.BankModel, error) {
	query := `SELECT bank_name, ifsc_code FROM banks`
	rows, err := q.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banks []structures.BankModel
	for rows.Next() {
		var bank structures.BankModel
		err := rows.Scan(&bank.BankName, &bank.IFSCCode)
		if err != nil {
			return nil, err
		}
		banks = append(banks, bank)
	}
	return &banks, nil
}