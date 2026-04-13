package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Seeyko/casamia-api/models"
	"github.com/Seeyko/casamia-api/services"
)

type MenuHandler struct {
	db *services.Database
}

func NewMenuHandler(db *services.Database) *MenuHandler {
	return &MenuHandler{db: db}
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	section := r.URL.Query().Get("section")

	var catRows interface{ Close() error }
	var err error

	query := `SELECT id, name, section, sort_order FROM menu_categories`
	if section != "" {
		query += ` WHERE section = $1`
		query += ` ORDER BY sort_order, id`
		rows, e := h.db.DB.Query(query, section)
		catRows = rows
		err = e
		if err == nil {
			defer catRows.Close()
			categories := scanCategories(rows)
			h.fillItems(categories)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.MenuResponse{Categories: categories})
		}
	} else {
		query += ` ORDER BY sort_order, id`
		rows, e := h.db.DB.Query(query)
		catRows = rows
		err = e
		if err == nil {
			defer catRows.Close()
			categories := scanCategories(rows)
			h.fillItems(categories)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.MenuResponse{Categories: categories})
		}
	}

	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
	}
}

func scanCategories(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
}) []models.MenuCategory {
	var categories []models.MenuCategory
	for rows.Next() {
		var c models.MenuCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Section, &c.SortOrder); err != nil {
			continue
		}
		categories = append(categories, c)
	}
	if categories == nil {
		categories = []models.MenuCategory{}
	}
	return categories
}

func (h *MenuHandler) fillItems(categories []models.MenuCategory) {
	for i := range categories {
		rows, err := h.db.DB.Query(
			`SELECT id, category_id, name, description, price, image_path, sort_order, available, badge, note
			 FROM menu_items WHERE category_id = $1 AND available = true ORDER BY sort_order, id`,
			categories[i].ID,
		)
		if err != nil {
			categories[i].Items = []models.MenuItem{}
			continue
		}

		var items []models.MenuItem
		for rows.Next() {
			var item models.MenuItem
			if err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.ImagePath, &item.SortOrder, &item.Available, &item.Badge, &item.Note); err != nil {
				continue
			}
			items = append(items, item)
		}
		rows.Close()

		if items == nil {
			items = []models.MenuItem{}
		}
		categories[i].Items = items
	}
}
