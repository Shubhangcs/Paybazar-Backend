package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) AddNewBeneficiary(req *structures.BeneficiaryModel) error {
	query := `INSERT INTO beneficiaries (mobile_number, bank_name, ifsc_code, account_number, beneficiary_name, beneficiary_phone) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := q.Pool.Exec(context.Background(), query, req.MobileNumber, req.BankName, req.IFSCCode, req.AccountNumber, req.BeneficiaryName, req.BeneficiaryPhone)
	return err
}

func (q *Query) GetBeneficiaries(mobileNumber string) (*[]structures.BeneficiaryModel, error) {
	query := `SELECT beneficiary_id ,mobile_number ,bank_name, ifsc_code, account_number, beneficiary_name, beneficiary_phone, beneficiary_verified FROM beneficiaries WHERE mobile_number = $1`
	rows, err := q.Pool.Query(context.Background(), query , mobileNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beneficiaries []structures.BeneficiaryModel
	for rows.Next() {
		var beneficiary structures.BeneficiaryModel
		err := rows.Scan(&beneficiary.BeneficiaryID,&beneficiary.MobileNumber,&beneficiary.BankName, &beneficiary.IFSCCode, &beneficiary.AccountNumber, &beneficiary.BeneficiaryName, &beneficiary.BeneficiaryPhone, &beneficiary.BeneficiaryVerified)
		if err != nil {
			return nil, err
		}
		beneficiaries = append(beneficiaries, beneficiary)
	}
	return &beneficiaries, nil
}

func (q *Query) VerifyBenificary(beneficiaryId string) error {
	query := `UPDATE beneficiaries SET beneficiary_verified = TRUE WHERE beneficiary_id = $1`
	_, err := q.Pool.Exec(context.Background(), query, beneficiaryId)	
	return err
}

func (q *Query) DeleteBeneficiary(beneficiaryId string) error {
	query := `DELETE FROM beneficiaries WHERE beneficiary_id=$1`
	_, err := q.Pool.Exec(context.Background(),query,beneficiaryId)
	return err
}
