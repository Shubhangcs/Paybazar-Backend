package main

import (
	"log"
	"os"

	// "github.com/Shubhangcs/South_canara_agromart_main_http_server/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/labstack/echo/v4"
)

func start() {
	// Creating a new echo router
	var router *echo.Echo = echo.New()

	// Connecting to database
	var conn *ConnectionPool = newDatabasePoolConnection()

	// Checking database connection status
	conn.CheckDatabaseWorkingStatus()

	// Closing the database connection once the function ends
	defer conn.CloseConnection()

	// Initilize database
	query := queries.NewQuery(conn.Pool)
	query.InitializeDatabase()

	// Add JSON valodator
	router.Validator = newValidator()

	// middlewares
	var middlewares *Middlewares = newMiddleware()
	router.Use(middlewares.LoggerMiddleware)

	// routes
	var routes *Routes = newRoutes(query)
	routes.AuthRouter(router)

	// Starting the Server
	var serverPort string = os.Getenv("SERVER_PORT")
	if serverPort == "" {
		log.Fatalln("failed to start server, server port not declared in .env file")
	}
	log.Fatal(router.Start(serverPort))
}
