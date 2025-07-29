package handlers

import (
	"encoding/json"
	"net/http"
	"secure-api/models"

	"gorm.io/gorm"
	// "github.com/gorilla/mux"
)

func GetProfile(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID").(uint)
        var user models.User
        if err := db.Select("id, username, role, created_at").First(&user, userID).Error; err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        json.NewEncoder(w).Encode(user)
    }
}