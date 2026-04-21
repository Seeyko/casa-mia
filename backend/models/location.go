package models

import (
	"encoding/json"
	"time"
)

type OpeningSlot struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type DayHours struct {
	Slots []OpeningSlot `json:"slots"`
}

// OpeningHours maps day names (lundi, mardi, ...) to their hours
type OpeningHours map[string]*DayHours

type Location struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Slug           string          `json:"slug"`
	Address        string          `json:"address"`
	Phone          string          `json:"phone"`
	OpeningHours   json.RawMessage `json:"opening_hours"`
	OrderMethod    string          `json:"order_method"`
	OrderInfo      string          `json:"order_info"`
	ClosureStart   *string         `json:"closure_start"`   // YYYY-MM-DD or null
	ClosureEnd     *string         `json:"closure_end"`     // YYYY-MM-DD or null
	ClosureMessage string          `json:"closure_message"` // optional message to show during closure
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type LocationStatus struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	IsOpen         bool   `json:"is_open"`
	NextChange     string `json:"next_change"`     // e.g. "Ouvre mardi à 9h" or "Ferme à 21h"
	IsOnVacation   bool   `json:"is_on_vacation"`  // true if today falls inside closure period
	ClosureMessage string `json:"closure_message"` // optional admin message ("Vacances d'été", etc.)
	ClosureUntil   string `json:"closure_until"`   // formatted end date (e.g. "31/04") when on vacation
}
