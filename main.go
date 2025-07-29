package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"secure-api/config"
	"secure-api/handlers"
	"secure-api/middleware"
	"secure-api/models"
	"secure-api/utils"
)

var db *gorm.DB

func main() {
    cfg := config.LoadConfig()
    utils.InitJWT(cfg.JWTSecret)

    var err error
    db, err = gorm.Open(sqlite.Open(cfg.DBFile), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to SQLite database")
    }

    err = db.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("Failed to migrate database")
    }

    r := mux.NewRouter()
    r.HandleFunc("/register", handlers.Register(db)).Methods("POST")
    r.HandleFunc("/login", handlers.Login(db)).Methods("POST")
    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware(""))
    api.HandleFunc("/profile", handlers.GetProfile(db)).Methods("GET")
    r.Use(middleware.SecurityHeader)
    r.Use(middleware.CORS(cfg.AllowedOrigins))

    log.Printf("Server starting on port %s (DB: %s)", cfg.Port, cfg.DBFile)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}