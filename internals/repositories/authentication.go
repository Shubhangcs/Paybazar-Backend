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
}

func NewAuthRepository(query *queries.Query, jwtUtils *pkg.JwtUtils, passwordUtils *pkg.PasswordUtils, jsonUtils *pkg.JsonUtils) *authRepository {
	return &authRepository{
		query:         query,
		jwtUtils:      jwtUtils,
		passwordUtils: passwordUtils,
		jsonUtils:     jsonUtils,
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
		return "", fmt.Errorf("failed to create user: %w", err)
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
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	if err := e.Validate(res); err != nil {
		return "", fmt.Errorf("invalid response from database: %w", err)
	}
	token, err := ar.jwtUtils.GenerateToken(res, time.Hour*24*365)
	return token, err
}

func (ar *authRepository) RegisterDistributor() {
	
}
