package main

import (
	"github.com/Srujankm12/paybazar-api/internals/handlers"
	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/repositories"
	"github.com/Srujankm12/paybazar-api/pkg"
	"github.com/labstack/echo/v4"
)

type Routes struct {
	Query         *queries.Query
	PasswordUtils *pkg.PasswordUtils
	JwtUtils      *pkg.JwtUtils
	TwillioUtils  *pkg.TwillioUtils
}

func newRoutes(query *queries.Query) *Routes {
	var passwordUtils = &pkg.PasswordUtils{}
	var jwtUtils = &pkg.JwtUtils{}
	var twillioUtils = &pkg.TwillioUtils{}
	return &Routes{
		Query:         query,
		PasswordUtils: passwordUtils,
		JwtUtils:      jwtUtils,
		TwillioUtils:  twillioUtils,
	}
}

func (r *Routes) AdminRoutes(rg *echo.Group) {
	// Authentication
	var authRepo = repositories.NewAuthRepository(
		r.Query,
		r.JwtUtils,
		r.PasswordUtils,
		r.TwillioUtils,
	)
	var authHandler = handlers.NewAuthHandler(authRepo)
	rg.POST("/register", authHandler.RegisterAdminRequest)
	rg.POST("/login", authHandler.LoginAdminRequest)
	rg.POST("/create/md", authHandler.RegisterMasterDistributorRequest)
	rg.POST("/create/distributor", authHandler.RegisterDistributorRequest)
	rg.POST("/create/user", authHandler.RegisterUserRequest)

	// Fund Request
	var fundRequestRepo = repositories.NewFundRequestRepository(
		r.Query,
	)
	var fundRequestHandler = handlers.NewFundRequestHandler(fundRequestRepo)
	rg.GET("/get/fund/requests/:admin_id", fundRequestHandler.GetAllFundRequests)
	rg.GET("/reject/fund/request/:request_id", fundRequestHandler.RejectFundRequest)
	rg.POST("/accept/fund/request", fundRequestHandler.AcceptFundRequest)

	// Wallet Request
	var walletRepo = repositories.NewWalletRepository(
		r.Query,
	)
	var walletHandler = handlers.NewWalletHandler(walletRepo)
	rg.GET("/wallet/get/balance/:admin_id", walletHandler.GetAdminWalletBalanceRequest)
	rg.GET("/wallet/get/transactions/:id", walletHandler.GetTransactionsRequest)
	rg.POST("/wallet/topup", walletHandler.AdminWalletTopupRequest)

	// Common Request
	var commonRepo = repositories.NewCommonRepository(
		r.Query,
	)
	var commonHandler = handlers.NewCommonHandler(commonRepo)
	rg.GET("/get/md/:admin_id", commonHandler.GetAllMasterDistributorsByAdminID)
	rg.GET("/get/distributor/:master_distributor_id", commonHandler.GetAllDistributorsByMasterDistributorID)
	rg.GET("/get/user/:distributor_id" , commonHandler.GetAllUsersByDistributorID)
	rg.GET("/get/distributor/:admin_id" , commonHandler.GetAllDistributorsByAdminID)
	rg.GET("/get/user/:admin_id" , commonHandler.GetAllUsersByAdminID)
}

func (r *Routes) MasterDistributorRoutes(rg *echo.Group) {
	// Authentication
	var authRepo = repositories.NewAuthRepository(
		r.Query,
		r.JwtUtils,
		r.PasswordUtils,
		r.TwillioUtils,
	)
	var authHandler = handlers.NewAuthHandler(authRepo)
	rg.POST("/login", authHandler.LoginMasterDistributorRequest)

	// Fund Request
	var fundRequestRepo = repositories.NewFundRequestRepository(
		r.Query,
	)
	var fundRequestHandler = handlers.NewFundRequestHandler(fundRequestRepo)
	rg.POST("/create/fund/request", fundRequestHandler.CreateFundRequest)
	rg.GET("/get/fund/request/:requester_id", fundRequestHandler.GetFundRequestsById)

	// Wallet Request
	var walletRepo = repositories.NewWalletRepository(
		r.Query,
	)
	var walletHandler = handlers.NewWalletHandler(walletRepo)
	rg.GET("/wallet/get/balance/:master_distributor_id", walletHandler.GetMasterDistributorWalletBalanceRequest)
	rg.GET("/wallet/get/transactions/:id", walletHandler.GetTransactionsRequest)
}

func (r *Routes) DistributorRoutes(rg *echo.Group) {
	// Authentication
	var authRepo = repositories.NewAuthRepository(
		r.Query,
		r.JwtUtils,
		r.PasswordUtils,
		r.TwillioUtils,
	)
	var authHandler = handlers.NewAuthHandler(authRepo)
	rg.POST("/login", authHandler.LoginDistributorRequest)

	// Fund Request
	var fundRequestRepo = repositories.NewFundRequestRepository(
		r.Query,
	)
	var fundRequestHandler = handlers.NewFundRequestHandler(fundRequestRepo)
	rg.POST("/create/fund/request", fundRequestHandler.CreateFundRequest)
	rg.GET("/get/fund/request/:requester_id", fundRequestHandler.GetFundRequestsById)

	// Wallet Request
	var walletRepo = repositories.NewWalletRepository(
		r.Query,
	)
	var walletHandler = handlers.NewWalletHandler(walletRepo)
	rg.GET("/wallet/get/balance/:distributor_id", walletHandler.GetDistributorWalletBalanceRequest)
	rg.GET("/wallet/get/transactions/:id", walletHandler.GetTransactionsRequest)
}

func (r *Routes) UserRoutes(rg *echo.Group) {
	// Authentication
	var authRepo = repositories.NewAuthRepository(
		r.Query,
		r.JwtUtils,
		r.PasswordUtils,
		r.TwillioUtils,
	)
	var authHandler = handlers.NewAuthHandler(authRepo)
	rg.POST("/login/send/otp", authHandler.LoginUserSendOTPRequest)
	rg.POST("/login/validate/otp", authHandler.LoginUserValidateOTPRequest)
	rg.POST("/set/mpin" , authHandler.SetMpinRequest)

	// Fund Request
	var fundRequestRepo = repositories.NewFundRequestRepository(
		r.Query,
	)
	var fundRequestHandler = handlers.NewFundRequestHandler(fundRequestRepo)
	rg.POST("/create/fund/request", fundRequestHandler.CreateFundRequest)
	rg.GET("/get/fund/request/:requester_id", fundRequestHandler.GetFundRequestsById)

	// Wallet Request
	var walletRepo = repositories.NewWalletRepository(
		r.Query,
	)
	var walletHandler = handlers.NewWalletHandler(walletRepo)
	rg.GET("/wallet/get/balance/:user_id", walletHandler.GetUserWalletBalanceRequest)
	rg.GET("/wallet/get/transactions/:id", walletHandler.GetTransactionsRequest)

	// Payout Request
	var payoutRepo = repositories.NewPayoutRepository(
		r.Query,
	)
	var payoutHandler = handlers.NewPayoutHandler(payoutRepo)
	rg.POST("/payout", payoutHandler.PayoutRequest)
}
