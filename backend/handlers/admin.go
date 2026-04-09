package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Seeyko/casamia-api/services"
	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	db    *services.Database
	img   *ImageHandler
}

func NewAdminHandler(db *services.Database, img *ImageHandler) *AdminHandler {
	return &AdminHandler{db: db, img: img}
}

// === NEWS CRUD ===

func (h *AdminHandler) ListNews(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.DB.Query(`SELECT id, title, content, image_path, published, created_at, updated_at FROM news ORDER BY created_at DESC`)
	if err != nil {
		jsonError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var news []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content, imagePath string
		var published bool
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&id, &title, &content, &imagePath, &published, &createdAt, &updatedAt); err != nil {
			continue
		}
		news = append(news, map[string]interface{}{
			"id": id, "title": title, "content": content, "image_path": imagePath,
			"published": published, "created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if news == nil {
		news = []map[string]interface{}{}
	}
	jsonResponse(w, news)
}

func (h *AdminHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, "invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	published := r.FormValue("published") == "true"

	var imagePath string
	if _, _, err := r.FormFile("image"); err == nil {
		filename, err := h.img.SaveUpload(r, "image", 5<<20)
		if err != nil {
			jsonError(w, "image upload failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		imagePath = filename
	}

	var id int
	err := h.db.DB.QueryRow(
		`INSERT INTO news (title, content, image_path, published) VALUES ($1, $2, $3, $4) RETURNING id`,
		title, content, imagePath, published,
	).Scan(&id)
	if err != nil {
		jsonError(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, map[string]interface{}{"id": id, "title": title, "image_path": imagePath})
}

func (h *AdminHandler) UpdateNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, "invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	published := r.FormValue("published") == "true"

	// Check for new image
	if _, _, err := r.FormFile("image"); err == nil {
		// Delete old image
		var oldImage string
		h.db.DB.QueryRow(`SELECT image_path FROM news WHERE id = $1`, id).Scan(&oldImage)
		h.img.DeleteFile(oldImage)

		filename, err := h.img.SaveUpload(r, "image", 5<<20)
		if err != nil {
			jsonError(w, "image upload failed", http.StatusBadRequest)
			return
		}
		h.db.DB.Exec(`UPDATE news SET title=$1, content=$2, published=$3, image_path=$4, updated_at=CURRENT_TIMESTAMP WHERE id=$5`,
			title, content, published, filename, id)
	} else {
		h.db.DB.Exec(`UPDATE news SET title=$1, content=$2, published=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$4`,
			title, content, published, id)
	}

	jsonResponse(w, map[string]string{"status": "updated"})
}

func (h *AdminHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var imagePath string
	h.db.DB.QueryRow(`SELECT image_path FROM news WHERE id = $1`, id).Scan(&imagePath)
	h.img.DeleteFile(imagePath)

	h.db.DB.Exec(`DELETE FROM news WHERE id = $1`, id)
	jsonResponse(w, map[string]string{"status": "deleted"})
}

// === MENU CATEGORIES CRUD ===

func (h *AdminHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.DB.Query(`SELECT id, name, section, sort_order FROM menu_categories ORDER BY section, sort_order, id`)
	if err != nil {
		jsonError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cats []map[string]interface{}
	for rows.Next() {
		var id, sortOrder int
		var name, section string
		if err := rows.Scan(&id, &name, &section, &sortOrder); err != nil {
			continue
		}
		cats = append(cats, map[string]interface{}{"id": id, "name": name, "section": section, "sort_order": sortOrder})
	}
	if cats == nil {
		cats = []map[string]interface{}{}
	}
	jsonResponse(w, cats)
}

func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name      string `json:"name"`
		Section   string `json:"section"`
		SortOrder int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid json", http.StatusBadRequest)
		return
	}

	var id int
	err := h.db.DB.QueryRow(
		`INSERT INTO menu_categories (name, section, sort_order) VALUES ($1, $2, $3) RETURNING id`,
		body.Name, body.Section, body.SortOrder,
	).Scan(&id)
	if err != nil {
		jsonError(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, map[string]interface{}{"id": id})
}

func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body struct {
		Name      string `json:"name"`
		Section   string `json:"section"`
		SortOrder int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid json", http.StatusBadRequest)
		return
	}

	h.db.DB.Exec(`UPDATE menu_categories SET name=$1, section=$2, sort_order=$3 WHERE id=$4`,
		body.Name, body.Section, body.SortOrder, id)
	jsonResponse(w, map[string]string{"status": "updated"})
}

func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Delete associated item images
	rows, _ := h.db.DB.Query(`SELECT image_path FROM menu_items WHERE category_id = $1 AND image_path != ''`, id)
	if rows != nil {
		for rows.Next() {
			var img string
			rows.Scan(&img)
			h.img.DeleteFile(img)
		}
		rows.Close()
	}

	h.db.DB.Exec(`DELETE FROM menu_categories WHERE id = $1`, id)
	jsonResponse(w, map[string]string{"status": "deleted"})
}

// === MENU ITEMS CRUD ===

func (h *AdminHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	query := `SELECT id, category_id, name, description, price, image_path, sort_order, available, badge, note FROM menu_items`
	var rows interface {
		Next() bool
		Scan(dest ...interface{}) error
		Close() error
	}
	var err error

	if categoryID != "" {
		query += ` WHERE category_id = $1 ORDER BY sort_order, id`
		rows, err = h.db.DB.Query(query, categoryID)
	} else {
		query += ` ORDER BY category_id, sort_order, id`
		rows, err = h.db.DB.Query(query)
	}
	if err != nil {
		jsonError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, catID, sortOrder int
		var name, desc, price, imagePath, badge, note string
		var available bool
		if err := rows.Scan(&id, &catID, &name, &desc, &price, &imagePath, &sortOrder, &available, &badge, &note); err != nil {
			continue
		}
		items = append(items, map[string]interface{}{
			"id": id, "category_id": catID, "name": name, "description": desc,
			"price": price, "image_path": imagePath, "sort_order": sortOrder,
			"available": available, "badge": badge, "note": note,
		})
	}
	if items == nil {
		items = []map[string]interface{}{}
	}
	jsonResponse(w, items)
}

func (h *AdminHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, "invalid form data", http.StatusBadRequest)
		return
	}

	categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
	name := r.FormValue("name")
	description := r.FormValue("description")
	price := r.FormValue("price")
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))
	badge := r.FormValue("badge")
	note := r.FormValue("note")

	var imagePath string
	if _, _, err := r.FormFile("image"); err == nil {
		filename, err := h.img.SaveUpload(r, "image", 5<<20)
		if err != nil {
			jsonError(w, "image upload failed", http.StatusBadRequest)
			return
		}
		imagePath = filename
	}

	var id int
	err := h.db.DB.QueryRow(
		`INSERT INTO menu_items (category_id, name, description, price, image_path, sort_order, badge, note)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		categoryID, name, description, price, imagePath, sortOrder, badge, note,
	).Scan(&id)
	if err != nil {
		jsonError(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, map[string]interface{}{"id": id})
}

func (h *AdminHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, "invalid form data", http.StatusBadRequest)
		return
	}

	categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
	name := r.FormValue("name")
	description := r.FormValue("description")
	price := r.FormValue("price")
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))
	available := r.FormValue("available") != "false"
	badge := r.FormValue("badge")
	note := r.FormValue("note")

	if _, _, err := r.FormFile("image"); err == nil {
		var oldImage string
		h.db.DB.QueryRow(`SELECT image_path FROM menu_items WHERE id = $1`, id).Scan(&oldImage)
		h.img.DeleteFile(oldImage)

		filename, err := h.img.SaveUpload(r, "image", 5<<20)
		if err != nil {
			jsonError(w, "image upload failed", http.StatusBadRequest)
			return
		}
		h.db.DB.Exec(`UPDATE menu_items SET category_id=$1, name=$2, description=$3, price=$4, image_path=$5, sort_order=$6, available=$7, badge=$8, note=$9 WHERE id=$10`,
			categoryID, name, description, price, filename, sortOrder, available, badge, note, id)
	} else {
		h.db.DB.Exec(`UPDATE menu_items SET category_id=$1, name=$2, description=$3, price=$4, sort_order=$5, available=$6, badge=$7, note=$8 WHERE id=$9`,
			categoryID, name, description, price, sortOrder, available, badge, note, id)
	}

	jsonResponse(w, map[string]string{"status": "updated"})
}

func (h *AdminHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var imagePath string
	h.db.DB.QueryRow(`SELECT image_path FROM menu_items WHERE id = $1`, id).Scan(&imagePath)
	h.img.DeleteFile(imagePath)

	h.db.DB.Exec(`DELETE FROM menu_items WHERE id = $1`, id)
	jsonResponse(w, map[string]string{"status": "deleted"})
}

// === LOCATIONS ===

func (h *AdminHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body struct {
		Name         string          `json:"name"`
		Address      string          `json:"address"`
		Phone        string          `json:"phone"`
		OpeningHours json.RawMessage `json:"opening_hours"`
		OrderMethod  string          `json:"order_method"`
		OrderInfo    string          `json:"order_info"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid json", http.StatusBadRequest)
		return
	}

	h.db.DB.Exec(`UPDATE locations SET name=$1, address=$2, phone=$3, opening_hours=$4, order_method=$5, order_info=$6, updated_at=CURRENT_TIMESTAMP WHERE id=$7`,
		body.Name, body.Address, body.Phone, body.OpeningHours, body.OrderMethod, body.OrderInfo, id)
	jsonResponse(w, map[string]string{"status": "updated"})
}

// === Helpers ===

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
