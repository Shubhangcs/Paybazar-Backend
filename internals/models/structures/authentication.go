package structures

// User Authentication Models

type UserLoginRequest struct {
	Phone string `json:"phone" validate:"required,phoneIN"`
	OTP   string `json:"otp" validate:"required"`
}

type UserRegistrationRequest struct {
	AdminID             string `json:"admin_id" validate:"required,uuid4"`
	MasterDistributorID string `json:"master_distributor_id" validate:"required,uuid4"`
	DistributorID       string `json:"distributor_id" validate:"required,uuid4"`
	UserName            string `json:"user_name" validate:"required,min=2,max=50"`
	UserEmail           string `json:"user_email" validate:"required,email"`
	UserPassword        string `json:"user_password" validate:"required,passwordStrong"`
	UserPhone           string `json:"user_phone" validate:"required,phoneIN"`
}

type UserAuthResponse struct {
	UserID              string `json:"user_id" validate:"required,uuid4"`
	UserUniqueID        string `json:"user_unique_id" validate:"required"`
	UserName            string `json:"user_name" validate:"required,min=2,max=50"`
	AdminID             string `json:"admin_id" validate:"required,uuid4"`
	MasterDistributorID string `json:"master_distributor_id" validate:"required,uuid4"`
	DistributorID       string `json:"distributor_id" validate:"required,uuid4"`
}

// Admin Authentication Models

type AdminLoginRequest struct {
	AdminEmail    string `json:"admin_email" validate:"required,email"`
	AdminPassword string `json:"admin_password" validate:"required,passwordStrong"`
}

type AdminRegisterRequest struct {
	AdminName     string `json:"admin_name" validate:"required,min=2,max=50"`
	AdminEmail    string `json:"admin_email" validate:"required,email"`
	AdminPassword string `json:"admin_password" validate:"required,passwordStrong"`
	AdminPhone    string `json:"admin_phone" validate:"required,phoneIN"`
}

type AdminAuthResponse struct {
	AdminID       string `json:"admin_id" validate:"required,uuid4"`
	AdminUniqueID string `json:"admin_unique_id" validate:"required"`
	AdminName     string `json:"admin_name" validate:"required,min=2,max=50"`
}

// Master Distributor Authentication Models

type MasterDistributorLoginRequest struct {
	MasterDistributorEmail    string `json:"master_distributor_email" validate:"required,email"`
	MasterDistributorPassword string `json:"master_distributor_password" validate:"required,passwordStrong"`
}

type MasterDistributorRegisterRequest struct {
	AdminID                      string `json:"admin_id" validate:"required,uuid4"`
	MasterDistributorName        string `json:"master_distributor_name" validate:"required,min=2,max=50"`
	MasterDistributorEmail       string `json:"master_distributor_email" validate:"required,email"`
	MasterDistributorPassword    string `json:"master_distributor_password" validate:"required,passwordStrong"`
	MasterDistributorPhoneNumber string `json:"master_distributor_phone_number" validate:"required,phoneIN"`
}

type MasterDistributorAuthResponse struct {
	MasterDistributorID       string `json:"master_distributor_id" validate:"required,uuid4"`
	MasterDistributorUniqueID string `json:"master_distributor_unique_id" validate:"required"`
	MasterDistributorName     string `json:"master_distributor_name" validate:"required,min=2,max=50"`
	AdminID                   string `json:"admin_id" validate:"required,uuid4"`
}

// Distributor Authentication Models

type DistributorLoginRequest struct {
	DistributorEmail    string `json:"distributor_email" validate:"required,email"`
	DistributorPassword string `json:"distributor_password" validate:"required,passwordStrong"`
}

type DistributorRegisterRequest struct {
	AdminID             string `json:"admin_id" validate:"required,uuid4"`
	MasterDistributorID string `json:"master_distributor_id" validate:"required,uuid4"`
	DistributorName     string `json:"distributor_name" validate:"required,min=2,max-50"`
	DistributorEmail    string `json:"distributor_email" validate:"required,email"`
	DistributorPassword string `json:"distributor_password" validate:"required,passwordStrong"`
	DistributorPhone    string `json:"distributor_phone" validate:"required,phoneIN"`
}

type DistributorAuthResponse struct {
	DistributorID       string `json:"distributor_id" validate:"required,uuid4"`
	DistributorUniqueID string `json:"distributor_unique_id" validate:"required"`
	DistributorName     string `json:"distributor_name" validate:"required,min=2,max-50"`
	MasterDistributorID string `json:"master_distributor_id" validate:"required,uuid4"`
	AdminID             string `json:"admin_id" validate:"required,uuid4"`
}

// Final Auth Response

type AuthResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data"`
}
