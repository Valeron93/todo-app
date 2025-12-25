package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Valeron93/todo-app/internal/assets"
	"github.com/Valeron93/todo-app/internal/controllers/auth"
	"github.com/Valeron93/todo-app/internal/controllers/todo"
	"github.com/Valeron93/todo-app/internal/migrations"
	"github.com/Valeron93/todo-app/internal/model"
	"github.com/Valeron93/todo-app/internal/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("ERROR: failed to close db: %v", err)
		}
	}()

	if err := migrations.RunMigrations(db); err != nil {
		log.Panic(err)
	}

	todoRepo := model.NewTodoRepoSql(db)
	authController := auth.New(model.NewUserRepoSql(db), model.NewSessionManagerSql(db))
	todoController := todo.New(todoRepo)
	views := views.NewViewHandler(todoRepo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(authController.InjectSessionMiddleware)

	r.Post("/api/register", authController.HandleRegister)
	r.Post("/api/login", authController.HandleLogin)

	r.Get("/register", views.HandleRegisterPage)
	r.Get("/login", views.HandleLoginPage)

	r.Handle("/static/*", assets.StaticHandler)

	// protected pages
	r.Group(func(r chi.Router) {
		r.Use(authController.AuthRedirectMiddleware)
		r.Get("/", views.HandleIndexPage)
	})

	// protected API endpoints
	r.Group(func(r chi.Router) {
		r.Use(authController.AuthMiddleware)

		r.Post("/api/todo", todoController.HandlePostTodo)
		r.Delete("/api/todo/{id}", todoController.HandleDeleteTodo)
		r.Post("/api/logout", authController.HandleLogout)
	})

	const addr = ":3000"
	log.Print("listening on " + addr)
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Panic(err)
	}
}
