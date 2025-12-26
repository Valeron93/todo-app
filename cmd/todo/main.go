package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Valeron93/todo-app/internal/assets"
	"github.com/Valeron93/todo-app/internal/controller"
	"github.com/Valeron93/todo-app/internal/middleware"
	"github.com/Valeron93/todo-app/internal/migrations"
	"github.com/Valeron93/todo-app/internal/model"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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
	userRepo := model.NewUserRepoSql(db)
	sessionRepo := model.NewSessionManagerSql(db)

	authController := controller.NewAuth(userRepo, sessionRepo)
	todoController := controller.NewTodo(todoRepo)
	authMiddleware := middleware.NewAuthMiddleware(sessionRepo)

	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)

	r.Use(authMiddleware.InjectSession)

	r.Post("/api/register", authController.HandleRegister)
	r.Post("/api/login", authController.HandleLogin)

	r.Get("/register", authController.HandleRegisterPage)
	r.Get("/login", authController.HandleLoginPage)

	r.Handle("/static/*", assets.StaticHandler)

	// protected pages
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.UnauthorizedRedirect("/login"))
		r.Get("/", todoController.HandleTodoListPage)
	})

	// protected API endpoints
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Unauthorized401)

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
