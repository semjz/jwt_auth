package database

import (
	"context"
	"fmt"
	"log"

	"auth_project/cmd/jwt_auth/ent"
	"auth_project/config"

	_ "github.com/lib/pq"
)

func DBConnect() *ent.Client {
	// Load the environment variables
	host := config.LoadEnvVariable("DB_HOST")
	port := config.LoadEnvVariable("DB_PORT")
	user := config.LoadEnvVariable("DB_USER")
	db_name := config.LoadEnvVariable("DB_NAME")
	pass := config.LoadEnvVariable("DB_PASS")

	// Correctly format the connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, db_name, pass,
	)

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	// defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
