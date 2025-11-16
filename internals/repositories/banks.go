package repositories

import (
	"fmt"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)


type bankRepo struct {
	query *queries.Query
}

func NewBankRepo(query *queries.Query) *bankRepo {
	return &bankRepo{query: query}
}

func (r *bankRepo) GetBanks(e echo.Context) (*[]structures.BankModel, error) {
	res , err := r.query.GetBanks()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bank details")
	}
	return res , nil
}

func (r *bankRepo) AddNewBank(e echo.Context) error {
	req := &structures.BankModel{}
	if err := e.Bind(req); err != nil {
		return fmt.Errorf("invalid request format")
	}
	err := r.query.AddNewBank(req)
	if err != nil {
		return fmt.Errorf("failed to add bank")
	}
	return nil
}
