package handlers

import (
	"encoding/json"
	"net/http"
	"secure-api/models"
	"secure-api/utils"
	"time"

	"gorm.io/gorm"
)

func Register(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }


        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        if input.Username == "" || input.Password == "" {
            http.Error(w, "Username and password required", http.StatusBadRequest)
            return
        }

        if len(input.Password) < 8 {
            http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
            return
        }


        hashed, err := utils.HashPassword(input.Password)
        if err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        user := models.User{
            Username: input.Username,
            Password: hashed,
            Role:     "user",
        }

        if err := db.Create(&user).Error; err != nil {
            http.Error(w, "User already exists", http.StatusConflict)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "User created successfully",
        })
    }
}

func Login(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        var user models.User
        if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        if !utils.CheckPassword(user.Password, input.Password) {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        accessToken, refreshToken, err := utils.GenerateJWT(user.ID, user.Role)
        if err != nil {
            http.Error(w, "Could not generate tokens", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:     "refresh_token",
            Value:    refreshToken,
            HttpOnly: true,
            Secure:   true, // HTTPS only
            SameSite: http.SameSiteStrictMode,
            Path:     "/refresh",
            MaxAge:   int(7 * 24 * time.Hour),
        })

        json.NewEncoder(w).Encode(map[string]string{
            "access_token": accessToken,
            "role":         user.Role,
        })
    }
}