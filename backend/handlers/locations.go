package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Seeyko/casamia-api/models"
	"github.com/Seeyko/casamia-api/services"
)

type LocationHandler struct {
	db *services.Database
}

func NewLocationHandler(db *services.Database) *LocationHandler {
	return &LocationHandler{db: db}
}

func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.DB.Query(`SELECT id, name, slug, address, phone, opening_hours, order_method, order_info, created_at, updated_at FROM locations ORDER BY id`)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		var orderMethod, orderInfo string
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Slug, &loc.Address, &loc.Phone, &loc.OpeningHours, &orderMethod, &orderInfo, &loc.CreatedAt, &loc.UpdatedAt); err != nil {
			continue
		}
		locations = append(locations, loc)
	}

	if locations == nil {
		locations = []models.Location{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (h *LocationHandler) Status(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.DB.Query(`SELECT id, name, slug, opening_hours FROM locations ORDER BY id`)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	loc, _ := time.LoadLocation("Europe/Paris")
	now := time.Now().In(loc)

	var statuses []models.LocationStatus
	for rows.Next() {
		var id int
		var name, slug string
		var hoursJSON json.RawMessage
		if err := rows.Scan(&id, &name, &slug, &hoursJSON); err != nil {
			continue
		}

		isOpen, nextChange := computeStatus(now, hoursJSON)
		statuses = append(statuses, models.LocationStatus{
			ID:         id,
			Name:       name,
			Slug:       slug,
			IsOpen:     isOpen,
			NextChange: nextChange,
		})
	}

	if statuses == nil {
		statuses = []models.LocationStatus{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

var dayNames = []string{"dimanche", "lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi"}

func computeStatus(now time.Time, hoursJSON json.RawMessage) (bool, string) {
	var hours models.OpeningHours
	if err := json.Unmarshal(hoursJSON, &hours); err != nil {
		return false, "Horaires indisponibles"
	}

	dayName := dayNames[now.Weekday()]
	dayHours, exists := hours[dayName]

	// Check if currently open
	if exists && dayHours != nil {
		for _, slot := range dayHours.Slots {
			openTime, err1 := parseTime(slot.Open)
			closeTime, err2 := parseTime(slot.Close)
			if err1 != nil || err2 != nil {
				continue
			}

			nowMinutes := now.Hour()*60 + now.Minute()
			if nowMinutes >= openTime && nowMinutes < closeTime {
				closeH := closeTime / 60
				closeM := closeTime % 60
				return true, fmt.Sprintf("Ferme a %dh%02d", closeH, closeM)
			}
		}
	}

	// Currently closed — find next opening
	nextOpen := findNextOpening(now, hours)
	return false, nextOpen
}

func findNextOpening(now time.Time, hours models.OpeningHours) string {
	dayName := dayNames[now.Weekday()]
	nowMinutes := now.Hour()*60 + now.Minute()

	// Check later today
	if dh, ok := hours[dayName]; ok && dh != nil {
		for _, slot := range dh.Slots {
			openTime, err := parseTime(slot.Open)
			if err != nil {
				continue
			}
			if openTime > nowMinutes {
				h := openTime / 60
				m := openTime % 60
				return fmt.Sprintf("Ouvre a %dh%02d", h, m)
			}
		}
	}

	// Check next days
	for i := 1; i <= 7; i++ {
		futureDay := now.AddDate(0, 0, i)
		futureDayName := dayNames[futureDay.Weekday()]
		if dh, ok := hours[futureDayName]; ok && dh != nil && len(dh.Slots) > 0 {
			openTime, err := parseTime(dh.Slots[0].Open)
			if err != nil {
				continue
			}
			h := openTime / 60
			m := openTime % 60
			capitalDay := capitalize(futureDayName)
			return fmt.Sprintf("Ouvre %s a %dh%02d", capitalDay, h, m)
		}
	}

	return "Horaires indisponibles"
}

func parseTime(s string) (int, error) {
	// Parse "9h", "9h00", "13h30", "18h00", "21h30"
	var h, m int
	n, _ := fmt.Sscanf(s, "%dh%d", &h, &m)
	if n >= 1 {
		return h*60 + m, nil
	}
	return 0, fmt.Errorf("invalid time: %s", s)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
