package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Seeyko/casamia-api/models"
	"github.com/Seeyko/casamia-api/services"
)

type NewsHandler struct {
	db *services.Database
}

func NewNewsHandler(db *services.Database) *NewsHandler {
	return &NewsHandler{db: db}
}

func (h *NewsHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.DB.Query(`SELECT id, title, content, image_path, published, created_at, updated_at FROM news WHERE published = true ORDER BY created_at DESC LIMIT 5`)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.ImagePath, &n.Published, &n.CreatedAt, &n.UpdatedAt); err != nil {
			continue
		}
		news = append(news, n)
	}

	if news == nil {
		news = []models.News{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
