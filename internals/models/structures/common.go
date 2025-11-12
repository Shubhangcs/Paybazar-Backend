package structures

type User struct {
	UserID           string `json:"user_id"`
	UserUniqueID     string `json:"user_unique_id"`
	UserName         string `json:"user_name"`
	UserEmail        string `json:"user_email"`
	UserPassword     string `json:"user_password"`
	UserPhone        string `json:"user_phone"`
	UserAadharNumber string `json:"user_aadhar_number"`
	UserPanNumber    string `json:"user_pan_number"`
	UserCity         string `json:"user_city"`
	UserState        string `json:"user_state"`
	UserAddress      string `json:"user_address"`
	UserPincode      string `json:"user_pincode"`
	UserDateOfBirth  string `json:"user_date_of_birth"`
	UserGender       string `json:"user_gender"`
	UserKYCStatus    bool   `json:"user_kyc_status"`
}

type Admin struct {
	AdminID       string `json:"admin_id"`
	AdminUniqueID string `json:"admin_unique_id"`
	AdminName     string `json:"admin_name"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	AdminPhone    string `json:"admin_phone"`
}

type MasterDistributor struct {
	MasterDistributorID       string `json:"master_distributor_id"`
	MasterDistributorUniqueID string `json:"master_distributor_unique_id"`
	MasterDistributorName     string `json:"master_distributor_name"`
	MasterDistributorEmail    string `json:"master_distributor_email"`
	MasterDistributorPassword string `json:"master_distributor_password"`
	MasterDistributorPhone    string `json:"master_distributor_phone"`
}

type Distributor struct {
	DistributorID       string `json:"distributor_id"`
	DistributorUniqueID string `json:"distributor_unique_id"`
	DistributorName     string `json:"distributor_name"`
	DistributorEmail    string `json:"distributor_email"`
	DistributorPassword string `json:"distributor_password"`
	DistributorPhone    string `json:"distributor_phone"`
}

type MasterDistributorGetResponse struct {
	MasterDistributorID            string `json:"master_distributor_id"`
	MasterDistributorUniqueID      string `json:"master_distributor_unique_id"`
	MasterDistributorName          string `json:"master_distributor_name"`
	MasterDistributorEmail         string `json:"master_distributor_email"`
	MasterDistributorPhone         string `json:"master_distributor_phone"`
	MasterDistributorWalletBalance string `json:"master_distributor_wallet_balance"`
}

type DistributorGetResponse struct {
	DistributorID            string `json:"distributor_id"`
	DistributorUniqueID      string `json:"distributor_unique_id"`
	DistributorName          string `json:"distributor_name"`
	DistributorEmail         string `json:"distributor_email"`
	DistributorPhone         string `json:"distributor_phone"`
	DistributorWalletBalance string `json:"distributor_wallet_balance"`
}

type UserGetResponse struct {
	UserID            string `json:"user_id"`
	UserUniqueID      string `json:"user_unique_id"`
	UserName          string `json:"user_name"`
	UserEmail         string `json:"user_email"`
	UserPhone         string `json:"user_phone"`
	UserWalletBalance string `json:"user_wallet_balance"`
}

type CommonResponse struct {
	Message string      `json:"msg"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}
