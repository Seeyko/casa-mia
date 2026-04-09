package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Seeyko/casamia-api/services"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db  *services.Database
	jwt *services.JWTService
}

func NewAuthHandler(db *services.Database, jwt *services.JWTService) *AuthHandler {
	return &AuthHandler{db: db, jwt: jwt}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if body.Username == "" || body.Password == "" {
		jsonError(w, "username and password required", http.StatusBadRequest)
		return
	}

	var id int
	var passwordHash string
	err := h.db.DB.QueryRow(`SELECT id, password_hash FROM admin_users WHERE username = $1`, body.Username).Scan(&id, &passwordHash)
	if err != nil {
		log.Printf("[AUTH] Login failed - user not found: %s", body.Username)
		time.Sleep(200 * time.Millisecond)
		jsonError(w, "identifiants incorrects", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.Password)); err != nil {
		log.Printf("[AUTH] Login failed - wrong password: %s", body.Username)
		time.Sleep(200 * time.Millisecond)
		jsonError(w, "identifiants incorrects", http.StatusUnauthorized)
		return
	}

	token, expiresAt, err := h.jwt.GenerateToken(id, body.Username)
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Printf("[AUTH] Login success: %s", body.Username)
	jsonResponse(w, map[string]interface{}{
		"token":      token,
		"expires_at": expiresAt,
		"username":   body.Username,
	})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// User is already authenticated via middleware — get user from context
	claims := r.Context().Value("claims").(*services.Claims)

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(body.NewPassword) < 6 {
		jsonError(w, "le mot de passe doit faire au moins 6 caracteres", http.StatusBadRequest)
		return
	}

	// Verify current password
	var passwordHash string
	err := h.db.DB.QueryRow(`SELECT password_hash FROM admin_users WHERE id = $1`, claims.UserID).Scan(&passwordHash)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.CurrentPassword)); err != nil {
		jsonError(w, "mot de passe actuel incorrect", http.StatusUnauthorized)
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.db.DB.Exec(`UPDATE admin_users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, string(newHash), claims.UserID)

	log.Printf("[AUTH] Password changed: %s", claims.Username)
	jsonResponse(w, map[string]string{"status": "password_changed"})
}

func (h *AuthHandler) RequestReset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	resetToken := hex.EncodeToString(tokenBytes)
	expires := time.Now().Add(1 * time.Hour)

	result, _ := h.db.DB.Exec(`UPDATE admin_users SET reset_token = $1, reset_expires = $2 WHERE username = $3`,
		resetToken, expires, body.Username)

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("[AUTH] Reset token generated for: %s — token: %s (expires %s)", body.Username, resetToken, expires.Format("15:04"))
	}

	// Always return success (don't reveal if user exists)
	jsonResponse(w, map[string]string{
		"status":  "ok",
		"message": "Si le compte existe, un token de reinitialisation a ete genere. Consultez les logs du serveur.",
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if body.Token == "" || len(body.NewPassword) < 6 {
		jsonError(w, "token et nouveau mot de passe (min 6 car.) requis", http.StatusBadRequest)
		return
	}

	// Find user with valid reset token
	var id int
	var username string
	err := h.db.DB.QueryRow(
		`SELECT id, username FROM admin_users WHERE reset_token = $1 AND reset_expires > CURRENT_TIMESTAMP`,
		body.Token,
	).Scan(&id, &username)
	if err != nil {
		time.Sleep(200 * time.Millisecond)
		jsonError(w, "token invalide ou expire", http.StatusBadRequest)
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Update password and clear reset token
	h.db.DB.Exec(`UPDATE admin_users SET password_hash = $1, reset_token = '', reset_expires = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		string(newHash), id)

	log.Printf("[AUTH] Password reset via token: %s", username)
	jsonResponse(w, map[string]string{"status": "password_reset", "message": "Mot de passe reinitialise avec succes"})
}
