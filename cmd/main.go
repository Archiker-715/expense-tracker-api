package main

import (
	"log"

	"github.com/Archiker-715/expense-tracker-api/internal/handlers"
	"github.com/Archiker-715/expense-tracker-api/internal/middleware"
	"github.com/Archiker-715/expense-tracker-api/internal/repository/pg"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// TODO: add routing and middleware for service's apipaths
// init DB
// dockerize that all
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment variables")
	}

	pg.Connect()

	authRepo := pg.NewAuthRepository(pg.DB)
	authHandler := handlers.NewAuthHadler(authRepo)

	expenseRepo := pg.NewExpenseRepository(pg.DB)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo)

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	authRouter := api.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("", authHandler.SignIn).Methods("POST")
	authRouter.HandleFunc("", authHandler.SignUp).Methods("POST")

	expenseRouter := api.PathPrefix("/expenses").Subrouter()
	expenseRouter.HandleFunc("", expenseHandler.GetExpenses).Methods("GET")
	expenseRouter.HandleFunc("", expenseHandler.CreateExpense).Methods("POST")
	expenseRouter.HandleFunc("", expenseHandler.UpdateExpense).Methods("PUT")
	expenseRouter.HandleFunc("", expenseHandler.DeleteExpense).Methods("DELETE")
	expenseRouter.Use(middleware.AuthMiddleware)
}
