package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Valeron93/todo-app/internal/controllers/auth"
	"github.com/Valeron93/todo-app/internal/controllers/todo"
	"github.com/Valeron93/todo-app/internal/migrations"
	"github.com/Valeron93/todo-app/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	if err := migrations.RunMigrations(db); err != nil {
		log.Panic(err)
	}

	todoRepo := model.NewTodoRepoSql(db)
	authController := auth.New(model.NewUserRepoSql(db), model.NewSessionManagerSql(db))
	todoController := todo.New(todoRepo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.HandleFunc("/register", authController.HandleRegister)
	r.HandleFunc("/login", authController.HandleLogin)

	// protected endpoints
	r.Group(func(r chi.Router) {
		r.Use(authController.AuthMiddleware)

		r.Get("/", todoController.HandleIndex)
		r.Post("/api/todo", todoController.HandlePostTodo)
		r.Delete("/api/todo/{id}", todoController.HandleDeleteTodo)
		r.Post("/logout", authController.HandleLogout)
	})

	const addr = ":3000"
	log.Print("listening on " + addr)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Panic(err)
	}
}
