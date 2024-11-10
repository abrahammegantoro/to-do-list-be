package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abrahammegantoro/to-do-list-be/internal/repository/psql"
	"github.com/abrahammegantoro/to-do-list-be/internal/rest"
	"github.com/abrahammegantoro/to-do-list-be/internal/rest/middlewares"
	"github.com/abrahammegantoro/to-do-list-be/todo"
	"github.com/abrahammegantoro/to-do-list-be/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	e := echo.New()
	// e.Use(middlewares.CORS)
	e.Use(middleware.CORS())

	userRepo := psql.NewUserRepository(conn)
	todoRepo := psql.NewTodoRepository(conn)

	userService := user.NewUserService(userRepo)
	todoService := todo.NewTodoService(todoRepo)

	api := e.Group("/api/v1")

	rest.NewUserHandler(api, userService)

	todoApi := api.Group("/todos")
	todoApi.Use(middlewares.AuthMiddleware(userRepo))

	rest.NewTodoHandler(todoApi, todoService)

	e.Logger.Fatal(e.Start(":8080"))
}
