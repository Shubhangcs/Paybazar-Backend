package repositories

import (
	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

type commonRepo struct {
	query *queries.Query
}

func NewCommonRepository(q *queries.Query) *commonRepo {
	return &commonRepo{query: q}
}

// Master Distributors by Admin ID
func (r *commonRepo) GetAllMasterDistributorsByAdminID(adminId string) (*[]structures.MasterDistributorGetResponse, error) {
	return r.query.GetAllMasterDistributorsByID(adminId)
}

// Distributors by Master Distributor ID
func (r *commonRepo) GetAllDistributorsByMasterDistributorID(masterDistributorId string) (*[]structures.DistributorGetResponse, error) {
	return r.query.GetAllDistributorsByMasterDistributorID(masterDistributorId)
}

// Users by Distributor ID
func (r *commonRepo) GetAllUsersByDistributorID(distributorId string) (*[]structures.UserGetResponse, error) {
	return r.query.GetAllUsersByDistributorID(distributorId)
}
