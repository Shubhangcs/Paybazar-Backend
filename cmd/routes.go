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
	rg.GET("/wallet/get/transactions/:admin_id", walletHandler.GetAdminWalletTransactionsRequest)
	rg.POST("/wallet/topup", walletHandler.AdminWalletTopupRequest)
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
	rg.POST("/create/distributor", authHandler.RegisterDistributorRequest)

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
	rg.GET("/wallet/get/transactions/:master_distributor_id", walletHandler.GetMasterDistributorWalletTransactionsRequest)
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
	rg.POST("/create/user", authHandler.RegisterUserRequest)

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
	rg.GET("/wallet/get/transactions/:distributor_id", walletHandler.GetDistributorWalletTransactionsRequest)
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
	rg.GET("/wallet/get/transactions/:user_id", walletHandler.GetUserWalletTransactionsRequest)

	// Payout Request
	var payoutRepo = repositories.NewPayoutRepository(
		r.Query,
	)
	var payoutHandler = handlers.NewPayoutHandler(payoutRepo)
	rg.POST("/payout", payoutHandler.PayoutRequest)
}
