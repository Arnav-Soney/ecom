package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	repo "github.com/Arnav-Soney/ecom/internal/adapters/postgresql/sqlc"
	"github.com/Arnav-Soney/ecom/internal/orders"
	"github.com/Arnav-Soney/ecom/internal/products"
)

// mount
func (app application) mount() http.Handler {
	// packages for handling hhtp server are : gorilla mux, chi router, echo framework, fiber
	r := chi.NewRouter()

	// A good base middleware stack

	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analystics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All good"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)

	orderService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(orderService)
	r.Post("/orders", ordersHandler.PlaceOrder)

	return r
}

// run
func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at address %s", app.config.addr)

	return server.ListenAndServe()
}

type application struct {
	config Config
	// logger
	// db driver dependency injection
	db *pgx.Conn
}

type Config struct {
	addr string // address of server
	db   dbConfig
}

type dbConfig struct {
	dsn string // domain string
}
