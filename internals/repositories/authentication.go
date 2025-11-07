package repositories

import (
	"fmt"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/Srujankm12/paybazar-api/pkg"
	"github.com/labstack/echo/v4"
)

type authRepository struct {
	query         *queries.Query
	jwtUtils      *pkg.JwtUtils
	passwordUtils *pkg.PasswordUtils
	jsonUtils     *pkg.JsonUtils
	twillioUtils  *pkg.TwillioUtils
}

func NewAuthRepository(query *queries.Query, jwtUtils *pkg.JwtUtils, passwordUtils *pkg.PasswordUtils, jsonUtils *pkg.JsonUtils, twillioUtils *pkg.TwillioUtils) *authRepository {
	return &authRepository{
		query:         query,
		jwtUtils:      jwtUtils,
		passwordUtils: passwordUtils,
		jsonUtils:     jsonUtils,
		twillioUtils:  twillioUtils,
	}
}

func (ar *authRepository) RegisterAdmin(e echo.Context) (string, error) {
	var req structures.AdminRegisterRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	hashPassword, err := ar.passwordUtils.HashPassword(req.AdminPassword)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	req.AdminPassword = hashPassword
	res, err := ar.query.CreateAdmin(&req)
	if err != nil {
		return "", fmt.Errorf("failed to create admin: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) RegisterMasterDistributor(e echo.Context) (string, error) {
	var req structures.MasterDistributorRegisterRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	hashPassword, err := ar.passwordUtils.HashPassword(req.MasterDistributorPassword)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	req.MasterDistributorPassword = hashPassword
	res, err := ar.query.CreateMasterDistributor(&req)
	if err != nil {
		return "", fmt.Errorf("failed to create master distributor: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) RegisterDistributor(e echo.Context) (string, error) {
	var req structures.DistributorRegisterRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	hashPassword, err := ar.passwordUtils.HashPassword(req.DistributorPassword)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	req.DistributorPassword = hashPassword
	res, err := ar.query.CreateDistributor(&req)
	if err != nil {
		return "", fmt.Errorf("failed to create distributor: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) RegisterUser(e echo.Context) (string, error) {
	var req structures.UserRegistrationRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	hashPassword, err := ar.passwordUtils.HashPassword(req.UserPassword)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	req.UserPassword = hashPassword
	res, err := ar.query.CreateUser(&req)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) LoginAdmin(e echo.Context) (string, error) {
	var req structures.AdminLoginRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	dbPass, err := ar.query.GetAdminPassword(req.AdminPassword)
	if err != nil {
		return "", fmt.Errorf("failed to retrive password from database: %w", err)
	}
	if err := ar.passwordUtils.VerifyPassword(dbPass, req.AdminPassword); err != nil {
		return "", fmt.Errorf("incorrect password: %w", err)
	}
	res, err := ar.query.LoginAdmin(&req)
	if err != nil {
		return "", fmt.Errorf("failed to login admin: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) LoginMasterDistributor(e echo.Context) (string, error) {
	var req structures.MasterDistributorLoginRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	dbPass, err := ar.query.GetMasterDistributorPassword(req.MasterDistributorPassword)
	if err != nil {
		return "", fmt.Errorf("failed to retrive password from database: %w", err)
	}
	if err := ar.passwordUtils.VerifyPassword(dbPass, req.MasterDistributorPassword); err != nil {
		return "", fmt.Errorf("incorrect password: %w", err)
	}
	res, err := ar.query.LoginMasterDistributor(&req)
	if err != nil {
		return "", fmt.Errorf("failed to login master distributor: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) LoginDistributor(e echo.Context) (string, error) {
	var req structures.DistributorLoginRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	dbPass, err := ar.query.GetDistributorPassword(req.DistributorPassword)
	if err != nil {
		return "", fmt.Errorf("failed to retrive password from database: %w", err)
	}
	if err := ar.passwordUtils.VerifyPassword(dbPass, req.DistributorPassword); err != nil {
		return "", fmt.Errorf("incorrect password: %w", err)
	}
	res, err := ar.query.LoginDistributor(&req)
	if err != nil {
		return "", fmt.Errorf("failed to login distributor: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) LoginUserSendOTP(e echo.Context) (string, error) {
	var req structures.UserLoginRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	exists, err := ar.query.CheckUserExistViaPhone(req.Phone)
	if err != nil {
		return "", fmt.Errorf("failed to find user in database: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("invalid phone number")
	}
	otp, err := ar.query.GenerateOTPForUser(req.Phone)
	if err != nil {
		return "", fmt.Errorf("failed to generate otp in database: %w", err)
	}
	if err := ar.twillioUtils.SendOTP(req.Phone, otp); err != nil {
		return "", fmt.Errorf("failed to send otp from twillio: %w", err)
	}
	return "OTP sent successfully", err
}

func (ar *authRepository) LoginUserValidateOTP(e echo.Context) (string, error) {
	var req structures.UserLoginRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request body: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}
	res, err := ar.query.ValidateOTP(&req)
	if err != nil {
		return "", fmt.Errorf("failed to validate OTP: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("failed to validate response: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}
