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
	JsonUtils     *pkg.JsonUtils
	TwillioUtils  *pkg.TwillioUtils
}

func newRoutes(query *queries.Query) *Routes {
	var passwordUtils = &pkg.PasswordUtils{}
	var jwtUtils = &pkg.JwtUtils{}
	var jsonUtils = &pkg.JsonUtils{}
	var twillioUtils = &pkg.TwillioUtils{}
	return &Routes{
		Query:         query,
		PasswordUtils: passwordUtils,
		JwtUtils:      jwtUtils,
		JsonUtils:     jsonUtils,
		TwillioUtils:  twillioUtils,
	}
}

func (r *Routes) AuthRouter(router *echo.Echo) {
	var repo = repositories.NewAuthRepository(
		r.Query,
		r.JwtUtils,
		r.PasswordUtils,
		r.JsonUtils,
		r.TwillioUtils,
	)
	var handler = handlers.NewAuthHandler(repo)

	router.POST("/admin/register", handler.RegisterAdminRequest)
	router.POST("/md/register", handler.RegisterMasterDistributorRequest)
	router.POST("/distributor/register", handler.RegisterDistributorRequest)
	router.POST("/user/register", handler.RegisterUserRequest)
	router.POST("/admin/login", handler.LoginAdminRequest)
	router.POST("/md/login", handler.LoginMasterDistributorRequest)
	router.POST("/distributor/login", handler.LoginDistributorRequest)
	router.POST("/user/sendotp", handler.LoginUserSendOTPRequest)
	router.POST("/user/validateotp", handler.LoginUserValidateOTPRequest)
}
