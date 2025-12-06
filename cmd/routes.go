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
	rg.POST("/user/wallet/refund", walletHandler.UserRefundRequest)
	rg.POST("/md/wallet/refund", walletHandler.MasterDistributorRefundRequest)
	rg.POST("/distributor/wallet/refund", walletHandler.DistributorRefundRequest)

	// Common Request
	var commonRepo = repositories.NewCommonRepository(
		r.Query,
	)
	var commonHandler = handlers.NewCommonHandler(commonRepo)
	rg.GET("/get/md/:admin_id", commonHandler.GetAllMasterDistributorsByAdminID)
	rg.GET("/get/distributors/:master_distributor_id", commonHandler.GetAllDistributorsByMasterDistributorID)
	rg.GET("/get/users/:distributor_id", commonHandler.GetAllUsersByDistributorID)
	rg.GET("/get/distributor/:admin_id", commonHandler.GetAllDistributorsByAdminID)
	rg.GET("/get/user/:admin_id", commonHandler.GetAllUsersByAdminID)
	rg.GET("/get/user/phone/:phone", commonHandler.GetUserByPhone)
	rg.GET("/get/md/phone/:phone", commonHandler.GetMasterDistributorByPhone)
	rg.GET("/get/distributor/phone/:phone", commonHandler.GetDistributorByPhone)

	// Ticket Requests
	var ticketRepo = repositories.NewTicketRepo(r.Query)
	var ticketHan = handlers.NewTicketHandler(ticketRepo)
	rg.GET("/get/tickets/:admin_id", ticketHan.GetAllTickets)
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
	rg.GET("/get/profile/:master_distributor_id", authHandler.GetMasterDistributorProfileRequest)
	rg.POST("/update/profile", authHandler.UpdateMasterDistributorProfileRequest)

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
	rg.GET("/get/profile/:distributor_id", authHandler.GetDistributorProfileRequest)
	rg.POST("/update/profile", authHandler.UpdateDistributorProfileRequest)

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
	rg.POST("/set/mpin", authHandler.SetMpinRequest)
	rg.POST("/verify/mpin", authHandler.VerifyMPINRequest)
	rg.GET("/get/profile/:user_id", authHandler.GetUserProfileRequest)
	rg.POST("/update/profile", authHandler.UpdateUserProfileRequest)

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
		r.JwtUtils,
	)
	var payoutHandler = handlers.NewPayoutHandler(payoutRepo)
	rg.POST("/payout", payoutHandler.PayoutRequest)
	rg.GET("/payout/get/transactions/:user_id", payoutHandler.GetPayoutTransactionRequest)
	rg.GET("/payout/:user_id/:phone/:account_number/:ifsc", payoutHandler.VerifyPayoutAccountNumber)

	// Bank Requests
	var bankRepo = repositories.NewBankRepo(r.Query)
	var bankHandler = handlers.NewBankHandler(bankRepo)
	rg.GET("/get/banks", bankHandler.GetAllBanks)
	rg.POST("/add/bank", bankHandler.CreateBank)

	// Beneficary Requests
	var benRepo = repositories.NewBeneficiaryRepo(r.Query)
	var benHandler = handlers.NewBeneficiaryHandler(benRepo)
	rg.GET("/get/beneficiaries/:phone", benHandler.GetBeneficiaries)
	rg.GET("/verify/beneficiaries/:ben_id", benHandler.VerifyBeneficiary)
	rg.POST("/add/beneficiary", benHandler.AddNewBeneficiary)
	rg.GET("/delete/beneficiary/:ben_id", benHandler.DeleteBeneficiary)

	// Ticket Requests
	var ticketRepo = repositories.NewTicketRepo(r.Query)
	var ticketHan = handlers.NewTicketHandler(ticketRepo)
	rg.POST("/add/ticket", ticketHan.AddNewTicket)
}
