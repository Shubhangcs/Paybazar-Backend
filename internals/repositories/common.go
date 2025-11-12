package repositories

import (
	"log"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type commonRepo struct {
	query *queries.Query
}

func NewCommonRepository(q *queries.Query) *commonRepo {
	return &commonRepo{query: q}
}

// Master Distributors by Admin ID
func (r *commonRepo) GetAllMasterDistributorsByAdminID(adminId string) (*[]structures.MasterDistributorGetResponse, error) {
	res, err := r.query.GetAllMasterDistributorsByID(adminId)
	if err != nil {
		log.Println("DB error while fetching master distributors by admin ID:", err)
		return nil, echo.NewHTTPError(500, "Failed to fetch master distributors")
	}
	if res == nil {
		// return empty slice instead of nil to avoid nil pointer issues in callers
		empty := []structures.MasterDistributorGetResponse{}
		return &empty, nil
	}
	return res, nil
}

// Distributors by Master Distributor ID
func (r *commonRepo) GetAllDistributorsByMasterDistributorID(masterDistributorId string) (*[]structures.DistributorGetResponse, error) {
	res, err := r.query.GetAllDistributorsByMasterDistributorID(masterDistributorId)
	if err != nil {
		log.Println("DB error while fetching distributors by master distributor ID:", err)
		return nil, echo.NewHTTPError(500, "Failed to fetch distributors")
	}
	if res == nil {
		empty := []structures.DistributorGetResponse{}
		return &empty, nil
	}
	return res, nil
}

// Users by Distributor ID
func (r *commonRepo) GetAllUsersByDistributorID(distributorId string) (*[]structures.UserGetResponse, error) {
	res, err := r.query.GetAllUsersByDistributorID(distributorId)
	if err != nil {
		log.Println("DB error while fetching users by distributor ID:", err)
		return nil, echo.NewHTTPError(500, "Failed to fetch users")
	}
	if res == nil {
		empty := []structures.UserGetResponse{}
		return &empty, nil
	}
	return res, nil
}
