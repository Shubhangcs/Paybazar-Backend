package repositories

import (
	"log"
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
	twillioUtils  *pkg.TwillioUtils
}

func NewAuthRepository(query *queries.Query, jwtUtils *pkg.JwtUtils, passwordUtils *pkg.PasswordUtils, twillioUtils *pkg.TwillioUtils) *authRepository {
	return &authRepository{
		query:         query,
		jwtUtils:      jwtUtils,
		passwordUtils: passwordUtils,
		twillioUtils:  twillioUtils,
	}
}

// ----------------------------
// Helper methods
// ----------------------------

func (ar *authRepository) bindAndValidate(e echo.Context, v interface{}) error {
	if err := e.Bind(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request format")
	}
	if err := e.Validate(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request data")
	}
	return nil
}

func (ar *authRepository) validateDBResponse(e echo.Context, res interface{}) error {
	if err := e.Validate(res); err != nil {
		return echo.NewHTTPError(500, "Invalid response from database")
	}
	return nil
}

func (ar *authRepository) generateTokenFor(res interface{}, duration time.Duration) (string, error) {
	token, err := ar.jwtUtils.GenerateToken(res, duration)
	if err != nil {
		log.Println("Token generation failed:", err)
		return "", echo.NewHTTPError(500, "Failed to generate token")
	}
	return token, nil
}

func (ar *authRepository) hashPassword(raw string) (string, error) {
	hashPassword, err := ar.passwordUtils.HashPassword(raw)
	if err != nil {
		log.Println("Password hashing failed:", err)
		return "", echo.NewHTTPError(500, "Failed to secure password")
	}
	return hashPassword, nil
}

func (ar *authRepository) verifyPasswordAndLogin(
	e echo.Context,
	getDBPass func() (string, error),
	inputPassword string,
	loginFunc func() (interface{}, error),
	tokenDuration time.Duration,
) (string, error) {
	dbPass, err := getDBPass()
	if err != nil {
		log.Println("DB password retrieval failed:", err)
		return "", echo.NewHTTPError(401, "Invalid credentials")
	}
	if err := ar.passwordUtils.VerifyPassword(dbPass, inputPassword); err != nil {
		return "", echo.NewHTTPError(401, "Invalid password")
	}
	res, err := loginFunc()
	if err != nil {
		log.Println("Login query failed:", err)
		return "", echo.NewHTTPError(401, "Failed to log in")
	}
	if err := ar.validateDBResponse(e, res); err != nil {
		return "", err
	}
	token, err := ar.generateTokenFor(res, tokenDuration)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ----------------------------
// Public methods
// ----------------------------

func (ar *authRepository) RegisterAdmin(e echo.Context) (string, error) {
	var req structures.AdminRegisterRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	hashPassword, err := ar.hashPassword(req.AdminPassword)
	if err != nil {
		return "", err
	}
	req.AdminPassword = hashPassword
	res, err := ar.query.CreateAdmin(&req)
	if err != nil {
		log.Println("DB create admin error:", err)
		return "", echo.NewHTTPError(500, "Failed to create admin")
	}
	if err := ar.validateDBResponse(e, res); err != nil {
		return "", err
	}
	return ar.generateTokenFor(res, time.Hour*24)
}

func (ar *authRepository) RegisterMasterDistributor(e echo.Context) (string, error) {
	var req structures.MasterDistributorRegisterRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	hashPassword, err := ar.hashPassword(req.MasterDistributorPassword)
	if err != nil {
		return "", err
	}
	req.MasterDistributorPassword = hashPassword
	res, err := ar.query.CreateMasterDistributor(&req)
	if err != nil {
		log.Println("DB create master distributor error:", err)
		return "", echo.NewHTTPError(500, "Failed to create master distributor")
	}
	if err := ar.validateDBResponse(e, res); err != nil {
		return "", err
	}
	return ar.generateTokenFor(res, time.Hour*24)
}

func (ar *authRepository) RegisterDistributor(e echo.Context) (string, error) {
	var req structures.DistributorRegisterRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	hashPassword, err := ar.hashPassword(req.DistributorPassword)
	if err != nil {
		return "", err
	}
	req.DistributorPassword = hashPassword
	res, err := ar.query.CreateDistributor(&req)
	if err != nil {
		log.Println("DB create distributor error:", err)
		return "", echo.NewHTTPError(500, "Failed to create distributor")
	}
	if err := ar.validateDBResponse(e, res); err != nil {
		return "", err
	}
	return ar.generateTokenFor(res, time.Hour*24*365)
}

func (ar *authRepository) RegisterUser(e echo.Context) (string, error) {
	var req structures.UserRegistrationRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	hashPassword, err := ar.hashPassword(req.UserPassword)
	if err != nil {
		return "", err
	}
	req.UserPassword = hashPassword
	res, err := ar.query.CreateUser(&req)
	if err != nil {
		log.Println("DB create user error:", err)
		return "", echo.NewHTTPError(500, "Failed to create user")
	}
	if err := ar.validateDBResponse(e, res); err != nil {
		return "", err
	}
	return ar.generateTokenFor(res, time.Hour*24*365)
}

func (ar *authRepository) LoginAdmin(e echo.Context) (string, error) {
	var req structures.AdminLoginRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	return ar.verifyPasswordAndLogin(
		e,
		func() (string, error) { return ar.query.GetAdminPassword(req.AdminEmail) },
		req.AdminPassword,
		func() (interface{}, error) { return ar.query.LoginAdmin(&req) },
		time.Hour*24*365,
	)
}

func (ar *authRepository) LoginMasterDistributor(e echo.Context) (string, error) {
	var req structures.MasterDistributorLoginRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	return ar.verifyPasswordAndLogin(
		e,
		func() (string, error) { return ar.query.GetMasterDistributorPassword(req.MasterDistributorEmail) },
		req.MasterDistributorPassword,
		func() (interface{}, error) { return ar.query.LoginMasterDistributor(&req) },
		time.Hour*24*365,
	)
}

func (ar *authRepository) LoginDistributor(e echo.Context) (string, error) {
	var req structures.DistributorLoginRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	return ar.verifyPasswordAndLogin(
		e,
		func() (string, error) { return ar.query.GetDistributorPassword(req.DistributorEmail) },
		req.DistributorPassword,
		func() (interface{}, error) { return ar.query.LoginDistributor(&req) },
		time.Hour*24*365,
	)
}

func (ar *authRepository) LoginUserSendOTP(e echo.Context) (string, error) {
	var req structures.UserLoginRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	exists, err := ar.query.CheckUserExistViaPhone(req.Phone)
	if err != nil {
		log.Println("DB user existence check error:", err)
		return "", echo.NewHTTPError(404, "Failed to find user")
	}
	if !exists {
		return "", echo.NewHTTPError(404, "User not found")
	}
	otp, err := ar.query.GenerateOTPForUser(req.Phone)
	if err != nil {
		log.Println("DB OTP generation error:", err)
		return "", echo.NewHTTPError(500, "Failed to generate OTP")
	}
	if err := ar.twillioUtils.SendOTP(req.Phone, otp); err != nil {
		log.Println("Twillio send error:", err)
		return "", echo.NewHTTPError(500, "Failed to send OTP")
	}
	return "OTP sent successfully", nil
}

func (ar *authRepository) LoginUserValidateOTP(e echo.Context) (string, error) {
	var req structures.UserLoginRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	res, err := ar.query.ValidateOTP(&req)
	if err != nil {
		log.Println("DB OTP validation error:", err)
		return "", echo.NewHTTPError(401, "Invalid or expired OTP")
	}
	if err := ar.validateDBResponse(e, res); err != nil {
		return "", err
	}
	return ar.generateTokenFor(res, time.Hour*24*365)
}

func (ar *authRepository) SetUserMpin(e echo.Context) (string, error) {
	var req structures.UserMpinRequest
	if err := ar.bindAndValidate(e, &req); err != nil {
		return "", err
	}
	err := ar.query.SetMpin(req.UserID, req.UserMPIN)
	if err != nil {
		log.Println("DB MPIN Setting error:", err)
		return "", echo.NewHTTPError(401, "Failed to Set MPIN")
	}
	return "MPIN Set Successfull", nil
}
