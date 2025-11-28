package queries

import (
	"context"
	"fmt"

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
			$1, $2, $3, $4
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
			master_distributor_password,
			master_distributor_aadhar_number,
			master_distributor_pan_number,
			master_distributor_date_of_birth,
			master_distributor_gender,
			master_distributor_city,
			master_distributor_state,
			master_distributor_address,
			master_distributor_pincode,
			business_name,
			business_type,
			gst_number
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16
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
		req.MasterDistributorAadharNumber,
		req.MasterDistributorPanNumber,
		req.MasterDistributorDateOfBirth,
		req.MasterDistributorGender,
		req.MasterDistributorCity,
		req.MasterDistributorState,
		req.MasterDistributorAddress,
		req.MasterDistributorPincode,
		req.BusinessName,
		req.BusinessType,
		req.GSTNumber,
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
			distributor_password,
			distributor_aadhar_number,
			distributor_pan_number,
			distributor_date_of_birth,
			distributor_gender,
			distributor_city,
			distributor_state,
			distributor_address,
			distributor_pincode,
			business_name,
			business_type,
			gst_number
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, $17
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
		req.MasterDistributorID,      // $1
		req.AdminID,                  // $2
		req.DistributorName,          // $3
		req.DistributorPhone,         // $4
		req.DistributorEmail,         // $5
		req.DistributorPassword,      // $6
		req.DistributorAadharNumber,  // $7
		req.DistributorPanNumber,     // $8
		req.DistributorDateOfBirth,   // $9
		req.DistributorGender,        // $10
		req.DistributorCity,          // $11
		req.DistributorState,         // $12
		req.DistributorAddress,       // $13
		req.DistributorPincode,       // $14
		req.BusinessName,
		req.BusinessType,
		req.GSTNumber,
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
			user_password,
			user_aadhar_number,
			user_pan_number,
			user_date_of_birth,
			user_gender,
			user_city,
			user_state,
			user_address,
			user_pincode,
			business_name,
			business_type,
			gst_number
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15,
			$16, $17, $18
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
		req.AdminID,              // $1
		req.MasterDistributorID,  // $2
		req.DistributorID,        // $3
		req.UserName,             // $4
		req.UserPhone,            // $5
		req.UserEmail,            // $6
		req.UserPassword,         // $7
		req.UserAadharNumber,     // $8
		req.UserPanNumber,        // $9
		req.UserDateOfBirth,      // $10
		req.UserGender,           // $11
		req.UserCity,             // $12
		req.UserState,            // $13
		req.UserAddress,          // $14
		req.UserPincode,          // $15
		req.BusinessName,
		req.BusinessType,
		req.GSTNumber,
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
	var mpin any

	query := `
	WITH validated AS (
		SELECT 
			u.user_id::TEXT AS user_id,
			u.user_unique_id,
			u.user_name,
			u.admin_id::TEXT AS admin_id,
			u.master_distributor_id::TEXT AS master_distributor_id,
			u.user_mpin::TEXT as user_mpin,
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
		(SELECT user_mpin FROM validated),
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
		&mpin,
		&res.DistributorID,
	)

	if mpin != "" {
		res.IsMpinSet = true
	}

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

func (q *Query) SetMpin(userId string, mpin string) (*structures.UserAuthResponse, error) {
	var res structures.UserAuthResponse
	var mpins string
	query := `UPDATE users SET user_mpin=$1 WHERE user_id=$2 RETURNING admin_id, master_distributor_id, distributor_id, user_id, user_unique_id, user_mpin;`
	err := q.Pool.QueryRow(context.Background(), query, mpin, userId).Scan(
		&res.AdminID,
		&res.MasterDistributorID,
		&res.DistributorID,
		&res.UserID,
		&res.UserUniqueID,
		&mpins,
	)
	if mpins == "" {
		return nil, fmt.Errorf("mpin not set")
	}
	res.IsMpinSet = true
	return &res, err
}

func (q *Query) VerifyMPIN(userId string, mpin string) (bool, error) {
	var isValidMpin bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1 AND user_mpin=$2) as user_exists`
	err := q.Pool.QueryRow(context.Background(), query, userId, mpin).Scan(&isValidMpin)
	return isValidMpin, err
}

func (q *Query) UpdateProfile(req *structures.UpdateUserProfile) error {
	query := `
		UPDATE users
		SET user_name=$1 , user_email=$2 , user_phone=$3,
		user_aadhar_number=$4, user_pan_number=$5,
		user_city=$6, user_state=$7, user_address=$8,
		user_pincode=$9, user_date_of_birth=$10,
		user_gender=$11 WHERE user_id=$12;
	`
	_, err := q.Pool.Exec(
		context.Background(),
		query,
		req.UserName,
		req.UserEmail,
		req.UserPhone,
		req.UserAadharNumber,
		req.UserPanNumber,
		req.UserCity,
		req.UserState,
		req.UserAddress,
		req.UserPincode,
		req.UserDateOfBirth,
		req.UserGender,
		req.UserID,
	)
	return err
}

func (q *Query) FetchProfileDetails(userId string) (*structures.GetUserProfile, error) {
	var res structures.GetUserProfile
	query := `
		SELECT user_id , user_unique_id ,user_name , user_email , user_phone,
		user_aadhar_number, user_pan_number,
		user_city, user_state, user_address,
		user_pincode, user_date_of_birth,
		user_gender FROM users  WHERE user_id=$1;
	`
	err := q.Pool.QueryRow(context.Background(), query, userId).Scan(
		&res.UserID,
		&res.UserUniqueID,
		&res.UserName,
		&res.UserEmail,
		&res.UserPhone,
		&res.UserAadharNumber,
		&res.UserPanNumber,
		&res.UserCity,
		&res.UserState,
		&res.UserAddress,
		&res.UserPincode,
		&res.UserDateOfBirth,
		&res.UserGender,
	)
	return &res, err
}

// UpdateMasterDistributorProfile updates the master distributor's profile details.
func (q *Query) UpdateMasterDistributorProfile(req *structures.UpdateMasterDistributorProfile) error {
	query := `
		UPDATE master_distributors
		SET master_distributor_name=$1,
			master_distributor_email=$2,
			master_distributor_phone=$3,
			master_distributor_aadhar_number=$4,
			master_distributor_pan_number=$5,
			master_distributor_city=$6,
			master_distributor_state=$7,
			master_distributor_address=$8,
			master_distributor_pincode=$9,
			master_distributor_date_of_birth=$10,
			master_distributor_gender=$11
		WHERE master_distributor_id=$12;
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		req.MasterDistributorName,
		req.MasterDistributorEmail,
		req.MasterDistributorPhone,
		req.MasterDistributorAadharNumber,
		req.MasterDistributorPanNumber,
		req.MasterDistributorCity,
		req.MasterDistributorState,
		req.MasterDistributorAddress,
		req.MasterDistributorPincode,
		req.MasterDistributorDateOfBirth,
		req.MasterDistributorGender,
		req.MasterDistributorID,
	)

	return err
}

// FetchMasterDistributorProfileDetails fetches the master distributor's profile details.
func (q *Query) FetchMasterDistributorProfileDetails(masterDistributorID string) (*structures.GetMasterDistributorProfile, error) {
	var res structures.GetMasterDistributorProfile

	query := `
		SELECT 
			master_distributor_id,
			master_distributor_unique_id,
			master_distributor_name,
			master_distributor_email,
			master_distributor_phone,
			master_distributor_aadhar_number,
			master_distributor_pan_number,
			master_distributor_city,
			master_distributor_state,
			master_distributor_address,
			master_distributor_pincode,
			master_distributor_date_of_birth,
			master_distributor_gender
		FROM master_distributors
		WHERE master_distributor_id=$1;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		masterDistributorID,
	).Scan(
		&res.MasterDistributorID,
		&res.MasterDistributorUniqueID,
		&res.MasterDistributorName,
		&res.MasterDistributorEmail,
		&res.MasterDistributorPhone,
		&res.MasterDistributorAadharNumber,
		&res.MasterDistributorPanNumber,
		&res.MasterDistributorCity,
		&res.MasterDistributorState,
		&res.MasterDistributorAddress,
		&res.MasterDistributorPincode,
		&res.MasterDistributorDateOfBirth,
		&res.MasterDistributorGender,
	)

	return &res, err
}

// UpdateDistributorProfile updates the distributor's profile details.
func (q *Query) UpdateDistributorProfile(req *structures.UpdateDistributorProfile) error {
	query := `
		UPDATE distributors
		SET distributor_name=$1,
			distributor_email=$2,
			distributor_phone=$3,
			distributor_aadhar_number=$4,
			distributor_pan_number=$5,
			distributor_city=$6,
			distributor_state=$7,
			distributor_address=$8,
			distributor_pincode=$9,
			distributor_date_of_birth=$10,
			distributor_gender=$11
		WHERE distributor_id=$12;
	`

	_, err := q.Pool.Exec(
		context.Background(),
		query,
		req.DistributorName,
		req.DistributorEmail,
		req.DistributorPhone,
		req.DistributorAadharNumber,
		req.DistributorPanNumber,
		req.DistributorCity,
		req.DistributorState,
		req.DistributorAddress,
		req.DistributorPincode,
		req.DistributorDateOfBirth,
		req.DistributorGender,
		req.DistributorID,
	)

	return err
}

// FetchDistributorProfileDetails fetches the distributor's profile details.
func (q *Query) FetchDistributorProfileDetails(distributorID string) (*structures.GetDistributorProfile, error) {
	var res structures.GetDistributorProfile

	query := `
		SELECT 
			distributor_id,
			distributor_unique_id,
			distributor_name,
			distributor_email,
			distributor_phone,
			distributor_aadhar_number,
			distributor_pan_number,
			distributor_city,
			distributor_state,
			distributor_address,
			distributor_pincode,
			distributor_date_of_birth,
			distributor_gender
		FROM distributors
		WHERE distributor_id=$1;
	`

	err := q.Pool.QueryRow(
		context.Background(),
		query,
		distributorID,
	).Scan(
		&res.DistributorID,
		&res.DistributorUniqueID,
		&res.DistributorName,
		&res.DistributorEmail,
		&res.DistributorPhone,
		&res.DistributorAadharNumber,
		&res.DistributorPanNumber,
		&res.DistributorCity,
		&res.DistributorState,
		&res.DistributorAddress,
		&res.DistributorPincode,
		&res.DistributorDateOfBirth,
		&res.DistributorGender,
	)

	return &res, err
}

