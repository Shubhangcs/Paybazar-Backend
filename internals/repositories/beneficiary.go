package repositories

import (
	"fmt"
	"log"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type beneficiaryRepo struct {
	query *queries.Query
}

func NewBeneficiaryRepo(query *queries.Query) *beneficiaryRepo {
	return &beneficiaryRepo{query: query}
}

func (r *beneficiaryRepo) GetBeneficiaries(e echo.Context) (*[]structures.BeneficiaryModel, error) {
	var phone = e.Param("phone")
	if phone == "" {
		return nil, fmt.Errorf("phone number not found")
	}
	res, err := r.query.GetBeneficiaries(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch benificiries")
	}
	return res, nil
}

func (r *beneficiaryRepo) AddNewBeneficiary(e echo.Context) error {
	req := &structures.BeneficiaryModel{}
	if err := e.Bind(req); err != nil {
		return err
	}
	err := r.query.AddNewBeneficiary(req)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to add benificary")
	}
	return nil
}

func (r *beneficiaryRepo) VerifyBeneficiary(e echo.Context) error {
	beneficiaryId := e.Param("ben_id")
	if beneficiaryId == "" {
		return fmt.Errorf("beneficiary id not found")
	}
	if err := r.query.VerifyBenificary(beneficiaryId); err != nil {
		return fmt.Errorf("failed to verify beneficiary")
	}
	return nil
}
