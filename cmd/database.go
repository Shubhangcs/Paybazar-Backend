package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool struct {
	Pool *pgxpool.Pool
}

func newDatabasePoolConnection() *ConnectionPool {
	var databaseUrl string = os.Getenv("DATABASE_URL")

	// Checking whether database url exists in .env file
	if databaseUrl == "" {
		log.Fatalln("failed to connect to database no database url in .env file")
	}

	// Configuring the database pool connection
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("failed to configure the database pool connection: %v", err)
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	// Creating a pool of connections
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Returning the database connection pool
	return &ConnectionPool{
		Pool: pool,
	}
}

func (conn *ConnectionPool) CheckDatabaseWorkingStatus() {
	// Ping the database for checking the connectivity
	if err := conn.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to communicate with database: %v",err)
	}
	log.Println("database connection is working properly")
}

func (conn *ConnectionPool) CloseConnection() {
	// Close database connection
	conn.Pool.Close()
}