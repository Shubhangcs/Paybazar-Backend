package interfaces

import "github.com/Srujankm12/paybazar-api/internals/models/structures"

type CommonInterface interface {
	GetAllMasterDistributorsByAdminID(adminId string) (*[]structures.MasterDistributorGetResponse, error)
	GetAllDistributorsByMasterDistributorID(masterDistributorId string) (*[]structures.DistributorGetResponse, error)
	GetAllUsersByDistributorID(distributorId string) (*[]structures.UserGetResponse, error)
	GetAllDistributorsByAdminID(adminId string) (*[]structures.DistributorGetResponse, error)
	GetAllUsersByAdminID(adminId string) (*[]structures.UserGetResponse, error)
}
