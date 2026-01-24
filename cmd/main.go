package main

import (
	"log"
	"net/http"

	"github.com/Archiker-715/expense-tracker-api/internal/handlers"
	"github.com/Archiker-715/expense-tracker-api/internal/middleware"
	"github.com/Archiker-715/expense-tracker-api/internal/repository/pg"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// dockerize that all
func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Println(".env not found, using environment variables")
	// }

	pg.Connect()

	authRepo := pg.NewAuthRepository(pg.DB)
	authHandler := handlers.NewAuthHadler(authRepo)

	expenseRepo := pg.NewExpenseRepository(pg.DB)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo)

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	authRouter := api.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signin", authHandler.SignIn).Methods("POST")
	authRouter.HandleFunc("/signup", authHandler.SignUp).Methods("POST")

	expenseRouter := api.PathPrefix("/expenses").Subrouter()
	expenseRouter.HandleFunc("", expenseHandler.GetExpenses).Methods("GET")
	expenseRouter.HandleFunc("", expenseHandler.CreateExpense).Methods("POST")
	expenseRouter.HandleFunc("/{id}", expenseHandler.UpdateExpense).Methods("PUT")
	expenseRouter.HandleFunc("/{id}", expenseHandler.DeleteExpense).Methods("DELETE")
	expenseRouter.Use(middleware.AuthMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("NOT FOUND: %s %s", r.Method, r.URL.Path)
		http.Error(w, "not found", http.StatusNotFound)
	})

	log.Println("Server starting :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
