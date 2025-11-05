package main

import (
	"github.com/Srujankm12/paybazar-api/internals/handlers"
	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/repositories"
	"github.com/Srujankm12/paybazar-api/pkg"
	"github.com/labstack/echo/v4"
)

type Routes struct {
	Query *queries.Query
	PasswordUtils *pkg.PasswordUtils
	JwtUtils *pkg.JwtUtils
	JsonUtils *pkg.JsonUtils
}

func newRoutes(query *queries.Query) *Routes {
	var passwordUtils = &pkg.PasswordUtils{}
	var jwtUtils = &pkg.JwtUtils{}
	var jsonUtils = &pkg.JsonUtils{}
	return &Routes{
		Query: query,
		PasswordUtils: passwordUtils,
		JwtUtils: jwtUtils,
		JsonUtils: jsonUtils,
	}
}

func (r *Routes) AuthRouter(router *echo.Echo) {
	var repo = repositories.NewAuthRepository(
		r.Query,
		r.JwtUtils,
		r.PasswordUtils,
		r.JsonUtils,
	)
	var handler = handlers.NewAuthHandler(repo)

	router.POST("/register" , handler.RegisterAdminRequest)
}