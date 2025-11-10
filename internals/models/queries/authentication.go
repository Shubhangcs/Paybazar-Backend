package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) CreateAdmin(req *structures.AdminRegisterRequest) (*structures.AdminAuthResponse, error) {
	var res structures.AdminAuthResponse

	query := `
		INSERT INTO admins (
			admin_name,
			admin_phone,
			admin_email,
			admin_password
		)
		VALUES (
			$1,  -- admin_name
			$2,  -- admin_phone
			$3,  -- admin_email
			$4   -- hashed password passed directly
		)
		RETURNING admin_id, admin_unique_id, admin_name;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminName,
		req.AdminPhone,
		req.AdminEmail,
		req.AdminPassword,
	).Scan(&res.AdminID, &res.AdminUniqueID, &res.AdminName)

	return &res, err
}

func (q *Query) CreateMasterDistributor(req *structures.MasterDistributorRegisterRequest) (*structures.MasterDistributorAuthResponse, error) {
	var res structures.MasterDistributorAuthResponse

	query := `
		INSERT INTO master_distributors (
			admin_id,
			master_distributor_name,
			master_distributor_phone,
			master_distributor_email,
			master_distributor_password
		)
		VALUES (
			$1,  -- admin_id
			$2,  -- master_distributor_name
			$3,  -- master_distributor_phone
			$4,  -- master_distributor_email
			$5   -- hashed password passed directly
		)
		RETURNING 
			master_distributor_id,
			master_distributor_unique_id,
			master_distributor_name,
			admin_id;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminID,
		req.MasterDistributorName,
		req.MasterDistributorPhoneNumber,
		req.MasterDistributorEmail,
		req.MasterDistributorPassword,
	).Scan(
		&res.MasterDistributorID,
		&res.MasterDistributorUniqueID,
		&res.MasterDistributorName,
		&res.AdminID,
	)

	return &res, err
}

func (q *Query) CreateDistributor(req *structures.DistributorRegisterRequest) (*structures.DistributorAuthResponse, error) {
	var res structures.DistributorAuthResponse

	query := `
		INSERT INTO distributors (
			master_distributor_id,
			admin_id,
			distributor_name,
			distributor_phone,
			distributor_email,
			distributor_password
		)
		VALUES (
			$1,  -- master_distributor_id
			$2,  -- admin_id
			$3,  -- distributor_name
			$4,  -- distributor_phone
			$5,  -- distributor_email
			$6   -- hashed password passed directly
		)
		RETURNING 
			distributor_id,
			distributor_unique_id,
			master_distributor_id,
			admin_id,
			distributor_name;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.MasterDistributorID, // $1
		req.AdminID,             // $2
		req.DistributorName,     // $3
		req.DistributorPhone,    // $4
		req.DistributorEmail,    // $5
		req.DistributorPassword, // $6
	).Scan(
		&res.DistributorID,
		&res.DistributorUniqueID,
		&res.MasterDistributorID,
		&res.AdminID,
		&res.DistributorName,
	)

	return &res, err
}

func (q *Query) CreateUser(req *structures.UserRegistrationRequest) (*structures.UserAuthResponse, error) {
	var res structures.UserAuthResponse

	query := `
		INSERT INTO users (
			admin_id,
			master_distributor_id,
			distributor_id,
			user_name,
			user_phone,
			user_email,
			user_password
		)
		VALUES (
			$1,  -- admin_id
			$2,  -- master_distributor_id
			$3,  -- distributor_id
			$4,  -- user_name
			$5,  -- user_phone
			$6,  -- user_email
			$7   -- hashed password passed directly
		)
		RETURNING 
			user_id,
			user_unique_id,
			user_name,
			distributor_id,
			master_distributor_id,
			admin_id;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.AdminID,             // $1
		req.MasterDistributorID, // $2
		req.DistributorID,       // $3
		req.UserName,            // $4
		req.UserPhone,           // $5
		req.UserEmail,           // $6
		req.UserPassword,        // $7
	).Scan(
		&res.UserID,
		&res.UserUniqueID,
		&res.UserName,
		&res.DistributorID,
		&res.MasterDistributorID,
		&res.AdminID,
	)

	return &res, err
}

func (q *Query) GetAdminPassword(email string) (string, error) {
	var password string
	query := `SELECT admin_password FROM admins WHERE admin_email=$1`
	err := q.Pool.QueryRow(context.Background(), query, email).Scan(&password)
	return password, err
}

func (q *Query) GetMasterDistributorPassword(email string) (string, error) {
	var password string
	query := `SELECT master_distributor_password FROM master_distributors WHERE master_distributor_email=$1`
	err := q.Pool.QueryRow(context.Background(), query, email).Scan(&password)
	return password, err
}

func (q *Query) GetDistributorPassword(email string) (string, error) {
	var password string
	query := `SELECT distributor_password FROM distributors WHERE distributor_email=$1`
	err := q.Pool.QueryRow(context.Background(), query, email).Scan(&password)
	return password, err
}

func (q *Query) GenerateOTPForUser(phone string) (string, error) {
	var otp string
	query := `
		INSERT INTO otps (phone)
		VALUES ($1)
		RETURNING otp;
	`
	err := q.Pool.QueryRow(context.Background(), query, phone).Scan(&otp)
	return otp, err
}

func (q *Query) ValidateOTP(req *structures.UserLoginRequest) (*structures.UserAuthResponse, error) {
	var res structures.UserAuthResponse

	query := `
	WITH validated AS (
		SELECT 
			u.user_id::TEXT AS user_id,
			u.user_unique_id,
			u.user_name,
			u.admin_id::TEXT AS admin_id,
			u.master_distributor_id::TEXT AS master_distributor_id,
			u.distributor_id::TEXT AS distributor_id,
			o.phone
		FROM otps o
		JOIN users u ON u.user_phone = o.phone
		WHERE o.phone = $1 AND o.otp = $2
	)
	DELETE FROM otps 
	WHERE phone = (SELECT phone FROM validated)
	RETURNING 
		(SELECT user_id FROM validated),
		(SELECT user_unique_id FROM validated),
		(SELECT user_name FROM validated),
		(SELECT admin_id FROM validated),
		(SELECT master_distributor_id FROM validated),
		(SELECT distributor_id FROM validated);
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.Phone,
		req.OTP,
	).Scan(
		&res.UserID,
		&res.UserUniqueID,
		&res.UserName,
		&res.AdminID,
		&res.MasterDistributorID,
		&res.DistributorID,
	)

	return &res, err
}

func (q *Query) LoginAdmin(req *structures.AdminLoginRequest) (*structures.AdminAuthResponse, error) {
	var res structures.AdminAuthResponse
	query := `
		SELECT 
  			admin_id::TEXT AS admin_id,
  			admin_unique_id,
  			admin_name
		FROM admins
		WHERE admin_email = $1;`

	err := q.Pool.QueryRow(context.Background(), query, req.AdminEmail).Scan(&res.AdminID, &res.AdminUniqueID, &res.AdminName)
	return &res, err
}

func (q *Query) LoginMasterDistributor(req *structures.MasterDistributorLoginRequest) (*structures.MasterDistributorAuthResponse, error) {
	var res structures.MasterDistributorAuthResponse

	query := `
	SELECT 
		master_distributor_id::TEXT AS master_distributor_id,
		master_distributor_unique_id,
		master_distributor_name,
		admin_id::TEXT AS admin_id
	FROM master_distributors
	WHERE master_distributor_email = $1;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.MasterDistributorEmail,
	).Scan(
		&res.MasterDistributorID,
		&res.MasterDistributorUniqueID,
		&res.MasterDistributorName,
		&res.AdminID,
	)

	return &res, err
}

func (q *Query) LoginDistributor(req *structures.DistributorLoginRequest) (*structures.DistributorAuthResponse, error) {
	var res structures.DistributorAuthResponse

	query := `
	SELECT 
		distributor_id::TEXT AS distributor_id,
		distributor_unique_id,
		distributor_name,
		admin_id::TEXT AS admin_id,
		master_distributor_id::TEXT AS master_distributor_id
	FROM distributors
	WHERE distributor_email = $1;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		req.DistributorEmail,
	).Scan(
		&res.DistributorID,
		&res.DistributorUniqueID,
		&res.DistributorName,
		&res.AdminID,
		&res.MasterDistributorID,
	)

	return &res, err
}

func (q *Query) CheckUserExistViaPhone(phone string) (bool, error) {
	var isUserExist bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE user_phone=$1) as user_exists"
	err := q.Pool.QueryRow(context.Background(), query, phone).Scan(&isUserExist)
	return isUserExist, err
}
